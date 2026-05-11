package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerListDonations struct {
	donationRepo *repositories.DonationRepository
}

func HandlerListDonationsStart(donationRepo *repositories.DonationRepository) *HandlerListDonations {
	return &HandlerListDonations{donationRepo}
}

func (h *HandlerListDonations) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.ListDonationReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerListDonations) Execute(ctx c.Context, filters *dto.ListDonationReq) (*dto.ListDonationRes, *fluxgo.GlobalError) {
	donations, total, err := h.donationRepo.ListDonationByFilters(ctx, filters)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to list donations")
	}
	if donations == nil || len(*donations) == 0 {
		return nil, fluxgo.ErrorNotFound("Donations not found")
	}

	return &dto.ListDonationRes{
		Data: *donations,
		PaginationRes: utils.PaginationRes{
			Page:     filters.Page,
			PageSize: filters.PageSize,
			Total:    total,
		},
	}, nil
}
