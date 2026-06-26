package handlers

import (
	"context"
	c "context"
	"nutriz-backend-service/config"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerListDonationPoints struct {
	donationPointRepo *repositories.DonationPointRepository
	config            *config.Env
}

func HandlerListDonationPointsStart(donationPointRepo *repositories.DonationPointRepository, config *config.Env) *HandlerListDonationPoints {
	return &HandlerListDonationPoints{donationPointRepo, config}
}

func (h *HandlerListDonationPoints) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.ListDonationPointsReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerListDonationPoints) Execute(ctx c.Context, filters *dto.ListDonationPointsReq) (*dto.ListDonationPointsRes, *fluxgo.GlobalError) {
	if filters.ZipCode != nil {
		coordinates, err := h.getCoordinatesByZipcode(ctx, *filters.ZipCode)
		if err != nil {
			return nil, err
		}

		filters.Latitude = &coordinates.Latitude
		filters.Longitude = &coordinates.Longitude
	}

	donationPoints, total, err := h.donationPointRepo.ListDonationPointsByFilters(ctx, filters)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to list donation points")
	}
	if donationPoints == nil || len(*donationPoints) == 0 {
		return nil, fluxgo.ErrorNotFound("Donation points not found")
	}

	return &dto.ListDonationPointsRes{
		Data: *donationPoints,
		PaginationRes: utils.PaginationRes{
			Page:     filters.Page,
			PageSize: filters.PageSize,
			Total:    total,
		},
	}, nil
}

func (h *HandlerListDonationPoints) getCoordinatesByZipcode(ctx context.Context, zipCode string) (*entities.Coordinates, *fluxgo.GlobalError) {
	addressData, err := utils.GetAddressByZipCode(ctx, zipCode, h.config)
	if err != nil {
		return nil, fluxgo.ErrorInternalError(err.Error())
	}

	if addressData.Latitude == nil || addressData.Longitude == nil {
		return nil, fluxgo.ErrorInternalError("Coordinates not found for the given zipcode")
	}

	return &entities.Coordinates{
		Latitude:  *addressData.Latitude,
		Longitude: *addressData.Longitude,
	}, nil
}
