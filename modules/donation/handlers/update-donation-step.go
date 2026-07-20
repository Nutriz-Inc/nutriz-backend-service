package handlers

import (
	c "context"
	"errors"
	"fmt"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type HandlerUpdateDonationStep struct {
	db                       *fluxgo.Database
	config                   *config.Env
	donationRepo             *repositories.DonationRepository
	donationStepRepo         *repositories.DonationStepRepository
	donationStepTimelineRepo *repositories.DonationStepTimelineRepository
	userRepo                 *repositories.UserRepository
	addressRepo              *repositories.AddressRepository
}

func HandlerUpdateDonationStepStart(
	db *fluxgo.Database,
	config *config.Env,
	donationRepo *repositories.DonationRepository,
	donationStepRepo *repositories.DonationStepRepository,
	donationStepTimelineRepo *repositories.DonationStepTimelineRepository,
	userRepo *repositories.UserRepository,
	addressRepo *repositories.AddressRepository,
) *HandlerUpdateDonationStep {
	return &HandlerUpdateDonationStep{
		db,
		config,
		donationRepo,
		donationStepRepo,
		donationStepTimelineRepo,
		userRepo,
		addressRepo,
	}
}

func (h *HandlerUpdateDonationStep) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.UpdateDonationStepReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerUpdateDonationStep) Execute(ctx c.Context, data *dto.UpdateDonationStepReq) (*dto.UpdateDonationStepRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type != entities.EnumUserTypeAdmin {
		return nil, utils.ErrorForbidden("User does not have permission to update donation step", "user.forbidden")
	}

	donationStep, err := h.donationStepRepo.GetDonationStepById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation step")
	}
	if donationStep == nil {
		return nil, fluxgo.ErrorNotFound("Donation step not found")
	}

	if !donationStep.CanUpdate() {
		return nil, utils.ErrorForbidden("Donation step cannot be updated because it's already done or failed", "donation_step.cannot_update")
	}

	donation, err := h.donationRepo.GetDonationById(ctx, donationStep.IdDonation)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation")
	}
	if donation == nil {
		return nil, fluxgo.ErrorNotFound("Donation not found")
	}

	if !donation.IsActive {
		return nil, utils.ErrorForbidden("Donation is not active", "donation.inactive")
	}

	fieldsToUpdate := 0
	req := repositories.UpdateDonationStepRepositoryReq{
		IdDonationStep: data.Id,
		IdUser:         data.ActionBy,
		Description:    data.Description,
		IsComplete:     false,
	}

	validator := data.ValidateUpdateDonationStepOptionalFields()

	var setDateTime *time.Time
	var isFailed bool
	var isDone bool

	if validator.HasStatus {
		if *data.Status == entities.EnumDonationStepStatusDone || *data.Status == entities.EnumDonationStepStatusFailed {
			if validator.HasSetDate {
				return nil, fluxgo.ErrorBadRequest(
					"You can't set a set date if the status is done or failed",
					"donation_step.invalid_status",
				)
			}

			req.IsComplete = true

			if *data.Status == entities.EnumDonationStepStatusDone {
				isDone = true
			}

			if *data.Status == entities.EnumDonationStepStatusFailed {
				isFailed = true
			}
		}

		req.Status = data.Status
		fieldsToUpdate++
	}

	if validator.HasSetDate {
		if !utils.IsFutureDate(*data.SetDate) {
			return nil, fluxgo.ErrorBadRequest(
				"Set date cannot be in the past",
				"donation_step.invalid_set_date",
			)
		}

		setDateTime, err = utils.StringToTime(*data.SetDate)
		if err != nil {
			return nil, fluxgo.ErrorBadRequest(
				"Invalid set date format",
				"donation_step.invalid_set_date_format",
			)
		}

		req.SetDate = setDateTime
		fieldsToUpdate++
	}

	if validator.HasAddress {
		fieldsToUpdate++
	}

	if fieldsToUpdate == 0 {
		return nil, fluxgo.ErrorBadRequest("At least one field must be sent to update", "no_fields_to_update")
	}

	err = h.db.RunTransaction(ctx, func(ctx c.Context, tx *sqlx.Tx) error {
		var idAddress *string

		if validator.HasAddress {
			var handleErr *fluxgo.GlobalError
			idAddress, handleErr = h.handleAddress(ctx, tx, &dto.CreateDonationStepReq{
				ActionBy:  data.ActionBy,
				IdAddress: data.IdAddress,
				Address:   data.Address,
			})
			if handleErr != nil {
				return &utils.TxError{Err: handleErr}
			}

			req.IdAddress = idAddress
		}

		err := h.donationStepRepo.UpdateDonationStepTx(ctx, tx, &req)
		if err != nil {
			return fmt.Errorf("error to update donation step: %w", err)
		}

		if (isDone && donationStep.Name == entities.EnumDonationStepMilkAnalysis) || isFailed {
			err = h.donationRepo.UpdateDonationTx(ctx, tx, &repositories.UpdateDonationRepositoryReq{
				IdDonation: donation.IdDonation,
				IdUser:     data.ActionBy,
				IsActive:   utils.BoolPtr(false),
			})
			if err != nil {
				return fmt.Errorf("error to update donation step: %w", err)
			}
		}

		idDonationStepTimeline := utils.IdGenerate(utils.DonationStepTimelineEntity)
		status := donationStep.Status

		if req.Status != nil && *req.Status != donationStep.Status {
			status = *req.Status
		}

		err = h.donationStepTimelineRepo.CreateDonationStepTimelineTx(ctx, tx, &repositories.CreateDonationStepTimelineRepositoryReq{
			IdDonationStepTimeline: idDonationStepTimeline,
			IdDonationStep:         data.Id,
			IdAddress:              idAddress,
			Description:            data.Description,
			Status:                 status,
			SetDate:                setDateTime,
			IdUser:                 data.ActionBy,
		})
		if err != nil {
			return fmt.Errorf("error to create donation step timeline: %w", err)
		}

		return nil
	})
	if err != nil {
		var txErr *utils.TxError
		if errors.As(err, &txErr) {
			return nil, txErr.Err
		}
		return nil, fluxgo.ErrorInternalError(err.Error())
	}

	donationStep, err = h.donationStepRepo.GetDonationStepById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation step")
	}
	if donationStep == nil {
		return nil, fluxgo.ErrorNotFound("Donation step not found")
	}

	return &dto.UpdateDonationStepRes{
		DonationStep: *donationStep,
	}, nil
}

func (h *HandlerUpdateDonationStep) handleAddress(ctx c.Context, tx *sqlx.Tx, data *dto.CreateDonationStepReq) (*string, *fluxgo.GlobalError) {
	var idAddress *string

	if data.IdAddress != nil {
		address, err := h.addressRepo.GetAddressById(ctx, *data.IdAddress)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to get address")
		}
		if address == nil {
			return nil, fluxgo.ErrorNotFound("Address not found")
		}

		if address.IdUser != nil && *address.IdUser != data.ActionBy {
			return nil, utils.ErrorForbidden("Address does not belong to the user", "address.forbidden")
		}

		idAddress = data.IdAddress
	} else {
		address, err := h.addressRepo.GetAddressByZipcode(ctx, data.Address.ZipCode)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to get address")
		}

		if address == nil || (address.IdUser != nil && *address.IdUser != data.ActionBy) {
			var handleErr *fluxgo.GlobalError

			idAddress, handleErr = h.createAddress(ctx, data.Address, data.ActionBy, tx)
			if handleErr != nil {
				return nil, handleErr
			}
		} else {
			idAddress = &address.IdAddress
		}
	}

	return idAddress, nil
}

func (h *HandlerUpdateDonationStep) createAddress(ctx c.Context, data *dto.AddressCreateBase, actionBy string, tx *sqlx.Tx) (*string, *fluxgo.GlobalError) {
	addressData, err := utils.GetAddressByZipCode(ctx, data.ZipCode, h.config)
	if err != nil {
		return nil, fluxgo.ErrorInternalError(err.Error())
	}

	idAddress := utils.IdGenerate(utils.AddressEntity)

	repoData := &repositories.CreateAddressRepositoryReq{
		IdAddress:    idAddress,
		IdUser:       actionBy,
		Zipcode:      data.ZipCode,
		Street:       addressData.Street,
		Number:       data.Number,
		City:         addressData.City,
		State:        addressData.State,
		Neighborhood: addressData.Neighborhood,
		Complement:   data.Complement,
		Latitude:     addressData.Latitude,
		Longitude:    addressData.Longitude,
	}

	err = h.addressRepo.CreateAddressTx(ctx, tx, repoData)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to create address")
	}

	return &idAddress, nil
}
