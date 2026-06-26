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

type HandlerCreateDonationStep struct {
	db                       *fluxgo.Database
	donationRepo             *repositories.DonationRepository
	donationStepRepo         *repositories.DonationStepRepository
	donationStepTimelineRepo *repositories.DonationStepTimelineRepository
	userRepo                 *repositories.UserRepository
}

func HandlerCreateDonationStepStart(
	db *fluxgo.Database,
	donationRepo *repositories.DonationRepository,
	donationStepRepo *repositories.DonationStepRepository,
	donationStepTimelineRepo *repositories.DonationStepTimelineRepository,
	userRepo *repositories.UserRepository,
) *HandlerCreateDonationStep {
	return &HandlerCreateDonationStep{
		db,
		donationRepo,
		donationStepRepo,
		donationStepTimelineRepo,
		userRepo,
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

	if donationSteps != nil && len(*donationSteps) > 0 {
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
		err := h.donationStepRepo.CreateDonationStepTx(ctx, tx, &repositories.CreateDonationStepRepositoryReq{
			IdDonationStep: idDonationStep,
			IdDonation:     data.IdDonation,
			IdUser:         data.ActionBy,
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
