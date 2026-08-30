package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerListDonationStepTimeline struct {
	donationStepTimelineRepo *repositories.DonationStepTimelineRepository
	donationStepRepo         *repositories.DonationStepRepository
	donationRepo             *repositories.DonationRepository
	userRepo                 *repositories.UserRepository
}

func HandlerListDonationStepTimelineStart(donationStepTimelineRepo *repositories.DonationStepTimelineRepository, donationStepRepo *repositories.DonationStepRepository, donationRepo *repositories.DonationRepository, userRepo *repositories.UserRepository) *HandlerListDonationStepTimeline {
	return &HandlerListDonationStepTimeline{
		donationStepTimelineRepo,
		donationStepRepo,
		donationRepo,
		userRepo,
	}
}

func (h *HandlerListDonationStepTimeline) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.ListDonationStepTimelineReq))
	if err != nil {
		return nil, err
	}

	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerListDonationStepTimeline) Execute(ctx c.Context, filters *dto.ListDonationStepTimelineReq) (*dto.ListDonationStepTimelineRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, filters.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if !user.Action().CanViewDonationStepTimeline {
		return nil, utils.ErrorForbidden("User does not have permission to list donation step timeline", "user.forbidden")
	}

	donationStep, err := h.donationStepRepo.GetDonationStepById(ctx, filters.IdDonationStep)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation step")
	}
	if donationStep == nil {
		return nil, fluxgo.ErrorNotFound("Donation step not found")
	}

	donation, err := h.donationRepo.GetDonationById(ctx, donationStep.IdDonation)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation")
	}
	if donation == nil {
		return nil, fluxgo.ErrorNotFound("Donation not found")
	}

	if user.Type == entities.EnumUserTypeCommon && donation.CreatedBy != filters.ActionBy {
		return nil, utils.ErrorForbidden("User does not have permission to list donation step timeline", "user.forbidden")
	}

	donationStepTimeline, _, err := h.donationStepTimelineRepo.GetDonationStepsTimelineByIdDonationStep(ctx, filters.IdDonationStep)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to list donation step timeline")
	}
	if donationStepTimeline == nil || len(*donationStepTimeline) == 0 {
		return nil, fluxgo.ErrorNotFound("Donation step timeline not found")
	}

	return &dto.ListDonationStepTimelineRes{
		Data: *donationStepTimeline,
	}, nil
}
