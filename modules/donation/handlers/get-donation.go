package handlers

import (
	c "context"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	dto "nutriz-backend-service/modules/donation/dtos"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerGetDonation struct {
	donationRepo     *repositories.DonationRepository
	donationStepRepo *repositories.DonationStepRepository
}

func HandlerGetDonationStart(donationRepo *repositories.DonationRepository, donationStepRepo *repositories.DonationStepRepository) *HandlerGetDonation {
	return &HandlerGetDonation{
		donationRepo:     donationRepo,
		donationStepRepo: donationStepRepo,
	}
}

func (h *HandlerGetDonation) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.GetDonationReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerGetDonation) Execute(ctx c.Context, filters *dto.GetDonationReq) (*dto.GetDonationRes, *fluxgo.GlobalError) {
	donation, err := h.donationRepo.GetDonationById(ctx, filters.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation")
	}
	if donation == nil {
		return nil, fluxgo.ErrorNotFound("Donation not found")
	}
	if donation.CreatedBy != filters.ActionBy {
		return nil, utils.ErrorForbidden("You don't have permission to access this resource", "donation.forbidden")
	}

	donationSteps, _, err := h.donationStepRepo.GetDonationStepsByIdDonation(ctx, filters.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation steps")
	}

	return &dto.GetDonationRes{
		Donation: *donation,
		Steps:    donationSteps,
	}, nil
}
