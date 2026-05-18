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

type HandlerUpdateDonationStep struct {
	donationRepo             *repositories.DonationRepository
	donationStepRepo         *repositories.DonationStepRepository
	DonationStepTimelineRepo *repositories.DonationStepTimelineRepository
	userRepo                 *repositories.UserRepository
}

func HandlerUpdateDonationStepStart(
	donationRepo *repositories.DonationRepository,
	donationStepRepo *repositories.DonationStepRepository,
	donationStepTimelineRepo *repositories.DonationStepTimelineRepository,
	userRepo *repositories.UserRepository,
) *HandlerUpdateDonationStep {
	return &HandlerUpdateDonationStep{
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

func (h *HandlerUpdateDonationStep) Execute(ctx context.Context, data *dto.UpdateDonationStepReq) (*dto.UpdateDonationStepRes, *fluxgo.GlobalError) {
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
		IdDonation: data.Id,
		IdUser:     data.ActionBy,
	}
	validator := data.ValidateUpdateDonationStepOptionalFields()

	if fieldsToUpdate == 0 {
		return nil, fluxgo.ErrorBadRequest("At least one field must be sent to update", "no_fields_to_update")
	}

	err = h.donationStepRepo.UpdateDonationStep(ctx, req)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to update donation step")
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
