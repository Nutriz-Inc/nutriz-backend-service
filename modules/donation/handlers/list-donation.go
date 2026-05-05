package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerListDonationPoints struct {
	donationPointsRepo *repositories.DonationPointRepository
}

func HandlerListDonationPointsStart(donationPointsRepo *repositories.DonationPointRepository) *HandlerListDonationPoints {
	return &HandlerListDonationPoints{donationPointsRepo}
}

func (h *HandlerListDonationPoints) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.ListDonationPointsReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerListDonationPoints) Execute(ctx c.Context, filters *dto.ListDonationPointsReq) (*dto.ListDonationPointsRes, *fluxgo.GlobalError) {
	donationPoints, total, err := h.donationPointsRepo.ListDonationPointsByFilters(ctx, filters)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to list donation points")
	}

	return &dto.ListDonationPointsRes{
		Data: donationPoints,
		PaginationRes: utils.PaginationRes{
			Page:     filters.Page,
			PageSize: filters.PageSize,
			Total:    total,
		},
	}, nil
}
