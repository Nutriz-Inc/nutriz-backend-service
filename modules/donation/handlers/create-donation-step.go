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

type HandlerCreateDonationStep struct {
	db                       *fluxgo.Database
	config                   *config.Env
	donationRepo             *repositories.DonationRepository
	donationStepRepo         *repositories.DonationStepRepository
	donationStepTimelineRepo *repositories.DonationStepTimelineRepository
	userRepo                 *repositories.UserRepository
	addressRepo              *repositories.AddressRepository
}

func HandlerCreateDonationStepStart(
	db *fluxgo.Database,
	config *config.Env,
	donationRepo *repositories.DonationRepository,
	donationStepRepo *repositories.DonationStepRepository,
	donationStepTimelineRepo *repositories.DonationStepTimelineRepository,
	userRepo *repositories.UserRepository,
	addressRepo *repositories.AddressRepository,
) *HandlerCreateDonationStep {
	return &HandlerCreateDonationStep{
		db,
		config,
		donationRepo,
		donationStepRepo,
		donationStepTimelineRepo,
		userRepo,
		addressRepo,
	}
}

func (h *HandlerCreateDonationStep) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.CreateDonationStepReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 201}, nil
}

func (h *HandlerCreateDonationStep) Execute(ctx c.Context, data *dto.CreateDonationStepReq) (*dto.CreateDonationStepRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type != entities.EnumUserTypeAdmin {
		return nil, utils.ErrorForbidden("User does not have permission to create donation step", "user.forbidden")
	}

	donation, err := h.donationRepo.GetDonationById(ctx, data.IdDonation)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation")
	}
	if donation == nil {
		return nil, fluxgo.ErrorNotFound("Donation not found")
	}
	if !donation.IsActive {
		return nil, fluxgo.ErrorBadRequest("Donation is not active", "donation.inactive")
	}

	donationSteps, _, err := h.donationStepRepo.GetDonationStepsByIdDonation(ctx, data.IdDonation)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation steps")
	}

	isFirstStep := donationSteps == nil || len(*donationSteps) == 0

	if isFirstStep && data.Name == entities.EnumDonationStepCollectMilk {
		donationUser, err := h.userRepo.GetUserById(ctx, donation.CreatedBy)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to get donation user")
		}
		if donationUser == nil {
			return nil, fluxgo.ErrorNotFound("Donation user not found")
		}
		if !donationUser.IsBloodExamValid() {
			return nil, fluxgo.ErrorBadRequest(
				"Donation can only start at the milk collection step if the user has a valid blood exam",
				"donation_step.blood_exam_invalid",
			)
		}
	}

	if !isFirstStep {
		if len(*donationSteps) == entities.NUMBER_OF_DONATION_STEPS {
			return nil, fluxgo.ErrorBadRequest("Donation already has the maximum number of steps", "donation_step.max_steps")
		}

		for _, step := range *donationSteps {
			if step.Status != entities.EnumDonationStepStatusDone {
				return nil, fluxgo.ErrorBadRequest(
					"Previous donation step is not completed",
					"previous_step.incomplete",
				)
			}
			if step.Name == data.Name {
				return nil, fluxgo.ErrorBadRequest(
					"Donation step with the same name already exists",
					"donation_step.duplicate_name",
				)
			}
		}
	}

	var setDateTime *time.Time

	if data.SetDate != nil {
		if !utils.IsFutureDate(*data.SetDate) {
			return nil, fluxgo.ErrorBadRequest("Set date must be in the future", "set_date.invalid")
		}

		setDateTime, err = utils.StringToTime(*data.SetDate)
		if err != nil {
			return nil, fluxgo.ErrorBadRequest("Invalid set date format", "donation_step.invalid_set_date_format")
		}
	}

	idDonationStep := utils.IdGenerate(utils.DonationStepEntity)

	err = h.db.RunTransaction(ctx, func(ctx c.Context, tx *sqlx.Tx) error {
		var idAddress *string

		if data.HasAddress() {
			var handleErr *fluxgo.GlobalError
			idAddress, handleErr = h.handleAddress(ctx, tx, data, donation)
			if handleErr != nil {
				return &utils.TxError{Err: handleErr}
			}
		}

		err := h.donationStepRepo.CreateDonationStepTx(ctx, tx, &repositories.CreateDonationStepRepositoryReq{
			IdDonationStep: idDonationStep,
			IdDonation:     data.IdDonation,
			IdUser:         data.ActionBy,
			IdAddress:      idAddress,
			Name:           data.Name,
			Description:    data.Description,
			Status:         entities.EnumDonationStepStatusPending,
			SetDate:        setDateTime,
		})
		if err != nil {
			return fmt.Errorf("error to create donation step: %w", err)
		}

		idDonationStepTimeline := utils.IdGenerate(utils.DonationStepTimelineEntity)
		err = h.donationStepTimelineRepo.CreateDonationStepTimelineTx(ctx, tx, &repositories.CreateDonationStepTimelineRepositoryReq{
			IdDonationStepTimeline: idDonationStepTimeline,
			IdDonationStep:         idDonationStep,
			IdAddress:              idAddress,
			Description:            data.Description,
			Status:                 entities.EnumDonationStepStatusPending,
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

	donationStep, err := h.donationStepRepo.GetDonationStepById(ctx, idDonationStep)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation step")
	}
	if donationStep == nil {
		return nil, fluxgo.ErrorNotFound("Donation step not found")
	}

	return &dto.CreateDonationStepRes{
		DonationStep: *donationStep,
	}, nil
}

func (h *HandlerCreateDonationStep) handleAddress(ctx c.Context, tx *sqlx.Tx, data *dto.CreateDonationStepReq, donation *entities.Donation) (*string, *fluxgo.GlobalError) {
	var idAddress *string

	if data.IdAddress != nil {
		address, err := h.addressRepo.GetAddressById(ctx, *data.IdAddress)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to get address")
		}
		if address == nil {
			return nil, fluxgo.ErrorNotFound("Address not found")
		}

		if address.IdUser != nil && *address.IdUser != donation.CreatedBy {
			return nil, utils.ErrorForbidden("Address does not belong to the user", "address.forbidden")
		}

		idAddress = data.IdAddress
	} else {
		address, err := h.addressRepo.GetAddressByZipcode(ctx, data.Address.ZipCode)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to get address")
		}

		if address == nil || (address.IdUser != nil && *address.IdUser != donation.CreatedBy) {
			var handleErr *fluxgo.GlobalError

			idAddress, handleErr = h.createAddress(ctx, data.Address, tx)
			if handleErr != nil {
				return nil, handleErr
			}
		} else {
			idAddress = &address.IdAddress
		}
	}

	return idAddress, nil
}

func (h *HandlerCreateDonationStep) createAddress(ctx c.Context, data *dto.AddressCreateBase, tx *sqlx.Tx) (*string, *fluxgo.GlobalError) {
	addressData, err := utils.GetAddressByZipCodeOptionalCoordinates(ctx, data.ZipCode, h.config)
	if err != nil {
		return nil, fluxgo.ErrorInternalError(err.Error())
	}

	idAddress := utils.IdGenerate(utils.AddressEntity)

	repoData := &repositories.CreateAddressRepositoryReq{
		IdAddress:    idAddress,
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
