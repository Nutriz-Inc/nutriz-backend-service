package handlers

import (
	c "context"
	"fmt"
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
	donationRepo             *repositories.DonationRepository
	donationStepRepo         *repositories.DonationStepRepository
	donationStepTimelineRepo *repositories.DonationStepTimelineRepository
	userRepo                 *repositories.UserRepository
}

func HandlerUpdateDonationStepStart(
	db *fluxgo.Database,
	donationRepo *repositories.DonationRepository,
	donationStepRepo *repositories.DonationStepRepository,
	donationStepTimelineRepo *repositories.DonationStepTimelineRepository,
	userRepo *repositories.UserRepository,
) *HandlerUpdateDonationStep {
	return &HandlerUpdateDonationStep{
		db,
		donationRepo,
		donationStepRepo,
		donationStepTimelineRepo,
		userRepo,
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
	if validator.HasStatus {
		if *data.Status == entities.EnumDonationStepStatusDone || *data.Status == entities.EnumDonationStepStatusFailed {
			if validator.HasSetDate {
				return nil, fluxgo.ErrorBadRequest(
					"You can't set a set date if the status is done or failed",
					"donation_step.invalid_status",
				)
			}

			req.IsComplete = true
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

	if fieldsToUpdate == 0 {
		return nil, fluxgo.ErrorBadRequest("At least one field must be sent to update", "no_fields_to_update")
	}

	err = h.db.RunTransaction(ctx, func(ctx c.Context, tx *sqlx.Tx) error {
		err := h.donationStepRepo.UpdateDonationStepTx(ctx, tx, &req)
		if err != nil {
			return fmt.Errorf("error to update donation step: %w", err)
		}

		if req.IsComplete {
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
