package handlers

import (
	c "context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	donationDto "nutriz-backend-service/modules/donation/dtos"
	donationHandler "nutriz-backend-service/modules/donation/handlers"
	dto "nutriz-backend-service/modules/user/dtos"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerGetUser struct {
	userRepo           *repositories.UserRepository
	donationRepo       *repositories.DonationRepository
	addressRepo        *repositories.AddressRepository
	babyRepo           *repositories.UserBabyRepository
	getDonationHandler *donationHandler.HandlerGetDonation
}

func HandlerGetUserStart(userRepo *repositories.UserRepository, donationRepo *repositories.DonationRepository, addressRepo *repositories.AddressRepository, babyRepo *repositories.UserBabyRepository, getDonationHandler *donationHandler.HandlerGetDonation) *HandlerGetUser {
	return &HandlerGetUser{
		userRepo,
		donationRepo,
		addressRepo,
		babyRepo,
		getDonationHandler,
	}
}

func (h *HandlerGetUser) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.GetUserReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerGetUser) Execute(ctx c.Context, filters *dto.GetUserReq) (*dto.GetUserRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, filters.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}

	userAction, err := h.userRepo.GetUserById(ctx, filters.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user action")
	}
	if userAction == nil {
		return nil, fluxgo.ErrorNotFound("User action not found")
	}

	if userAction.Type == entities.EnumUserTypeCommon && user.IdUser != filters.ActionBy {
		return nil, utils.ErrorForbidden("You don't have permission to access this resource", "user.forbidden")
	}

	var addresses *[]entities.Address
	var babies *[]entities.UserBaby
	var donationsCompleted *int64
	var currentDonation *donationDto.GetDonationRes

	isCommonUser := user.Type == entities.EnumUserTypeCommon

	if filters.ShowAddress && isCommonUser {
		addresses, _, err = h.addressRepo.GetAddressesByUserId(ctx, filters.Id)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to get addresses")
		}
	}

	if filters.ShowBaby && isCommonUser {
		babies, _, err = h.babyRepo.GetUserBabiesByUserId(ctx, filters.Id)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to get babies")
		}
	}

	if filters.ShowDonationsCompleted && isCommonUser {
		count, err := h.donationRepo.CountCompletedDonations(ctx, filters.Id)
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to count completed donations")
		}
		donationsCompleted = &count
	}

	if filters.ShowCurrentDonation && isCommonUser {
		latestDonation, _, err := h.donationRepo.ListDonationByFilters(ctx, &donationDto.ListDonationReq{
			IsActive: utils.BoolPtr(true),
			ActionBy: &filters.ActionBy,
			PaginationReq: utils.PaginationReq{
				PageSize: 1,
				Page:     1,
			},
		})
		if err != nil {
			return nil, fluxgo.ErrorInternalError("Error to get current donation")
		}

		if len(*latestDonation) > 0 {
			var handlerErr *fluxgo.GlobalError

			currentDonation, handlerErr = h.getDonationHandler.Execute(ctx, &donationDto.GetDonationReq{
				ActionBy: filters.ActionBy,
				GetReq: utils.GetReq{
					Id: (*latestDonation)[0].IdDonation,
				},
			})
			if handlerErr != nil {
				return nil, fluxgo.ErrorInternalError("Error to get current donation")
			}
		}
	}

	var addressesOut *[]entities.AddressOut
	if addresses != nil {
		out := entities.NewAddressesOut(*addresses)
		addressesOut = &out
	}

	var babiesOut *[]entities.UserBabyOut
	if babies != nil {
		out := entities.NewUserBabiesOut(*babies)
		babiesOut = &out
	}

	return &dto.GetUserRes{
		UserOut:            entities.NewUserOut(*user),
		Addresses:          addressesOut,
		Babies:             babiesOut,
		DonationsCompleted: donationsCompleted,
		CurrentDonation:    currentDonation,
	}, nil
}
