package handlers

import (
	"context"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerCreateDonationStep struct {
	donationRepo     *repositories.DonationRepository
	donationStepRepo *repositories.DonationStepRepository
	userRepo         *repositories.UserRepository
}

func HandlerCreateDonationStepStart(
	donationRepo *repositories.DonationRepository,
	donationStepRepo *repositories.DonationStepRepository,
	userRepo *repositories.UserRepository,
) *HandlerCreateDonationStep {
	return &HandlerCreateDonationStep{
		donationRepo,
		donationStepRepo,
		userRepo,
	}
}

func (h *HandlerCreateDonationStep) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.CreateDonationStepReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerCreateDonationStep) Execute(ctx context.Context, data *dto.CreateDonationStepReq) (*dto.CreateDonationStepRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
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

	if donationSteps != nil && len(*donationSteps) >= 0 {
		for _, step := range *donationSteps {
			if step.Status != entities.EnumDonationStepStatusDone && step.Status != entities.EnumDonationStepStatusWarn {
				return nil, fluxgo.ErrorBadRequest("Previous donation step is not completed", "previous_step.incomplete")
			}
		}
	}

	if data.SetDate != nil {
		if !utils.IsFutureDate(*data.SetDate) {
			return nil, fluxgo.ErrorBadRequest("Set date must be in the future", "set_date.invalid")
		}
	}

	idDonation := utils.IdGenerate(utils.DonationStepEntity)

	reqData := repositories.CreateDonationStepRepositoryReq{
		IdDonationStep: idDonation,
		IdDonation:     data.IdDonation,
		IdUser:         data.ActionBy,
		Name:           data.Name,
		Description:    data.Description,
		Status:         entities.EnumDonationStepStatusPending,
		SetDate:        data.SetDate,
	}

	err = h.donationStepRepo.CreateDonationStep(ctx, &reqData)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to create donation step")
	}

	donationStep, err := h.donationStepRepo.GetDonationStepById(ctx, idDonation)
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
