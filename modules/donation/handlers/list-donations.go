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

type HandlerListDonations struct {
	donationRepo *repositories.DonationRepository
	userRepo     *repositories.UserRepository
}

func HandlerListDonationsStart(donationRepo *repositories.DonationRepository, userRepo *repositories.UserRepository) *HandlerListDonations {
	return &HandlerListDonations{donationRepo, userRepo}
}

func (h *HandlerListDonations) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.ListDonationReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerListDonations) Execute(ctx c.Context, filters *dto.ListDonationReq) (*dto.ListDonationRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, *filters.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type != entities.EnumUserTypeCommon {
		filters.ActionBy = nil
	} else {
		filters.UserName = nil
		filters.UserDocument = nil
	}

	donations, total, err := h.donationRepo.ListDonationByFilters(ctx, filters)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to list donations")
	}
	if donations == nil || len(*donations) == 0 {
		return nil, fluxgo.ErrorNotFound("Donations not found")
	}

	donationsFormatted := h.mapResponse(donations, user.Type)

	return &dto.ListDonationRes{
		Data: *donationsFormatted,
		PaginationRes: utils.PaginationRes{
			Page:     filters.Page,
			PageSize: filters.PageSize,
			Total:    total,
		},
	}, nil
}

func (h *HandlerListDonations) mapResponse(donations *[]dto.DonationRes, userType entities.EnumUserType) *[]dto.DonationRes {
	if userType == entities.EnumUserTypeCommon {
		for i := range *donations {
			(*donations)[i].UserName = nil
			(*donations)[i].UserDocument = nil
		}
	}

	return donations
}
