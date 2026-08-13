package handlers

import (
	c "context"
	"fmt"
	dto "nutriz-backend-service/modules/donation/dtos"
	sharedDto "nutriz-backend-service/shared/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerCreateDonation struct {
	donationRepo *repositories.DonationRepository
	userRepo     *repositories.UserRepository
}

func HandlerCreateDonationStart(
	donationRepo *repositories.DonationRepository,
	userRepo *repositories.UserRepository,
) *HandlerCreateDonation {
	return &HandlerCreateDonation{
		donationRepo,
		userRepo,
	}
}

func (h *HandlerCreateDonation) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.CreateDonationReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 201}, nil
}

func (h *HandlerCreateDonation) Execute(ctx c.Context, data *dto.CreateDonationReq) (*dto.CreateDonationRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type != entities.EnumUserTypeCommon {
		return nil, utils.ErrorForbidden("User does not have permission to create donation", "user.forbidden")
	}

	_, totalDonationsActive, err := h.donationRepo.ListDonationByFilters(ctx, &dto.ListDonationReq{
		IsActive: utils.BoolPtr(true),
		ActionBy: &data.ActionBy,
		PaginationReq: utils.PaginationReq{
			PageSize: 1,
			Page:     1,
		},
	})
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donations active by user id")
	}
	if totalDonationsActive >= entities.MAX_DONATION_ACTIVE_QUANTITY_PER_USER {
		return nil, fluxgo.ErrorBadRequest(fmt.Sprintf("User can have up to %d active donations", entities.MAX_DONATION_ACTIVE_QUANTITY_PER_USER), "donation.max_active_quantity_reached")
	}

	idDonation := utils.IdGenerate(utils.DonationEntity)
	repoData := &repositories.CreateDonationRepositoryReq{
		IdDonation: idDonation,
		IdUser:     user.IdUser,
		IsActive:   true, // default value for new donation
	}

	err = h.donationRepo.CreateDonation(ctx, repoData)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to create donation")
	}

	donation, err := h.donationRepo.GetDonationById(ctx, idDonation)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation")
	}
	if donation == nil {
		return nil, fluxgo.ErrorNotFound("Donation not found")
	}

	return &dto.CreateDonationRes{
		DonationOut: sharedDto.NewDonationOut(*donation),
	}, nil
}
