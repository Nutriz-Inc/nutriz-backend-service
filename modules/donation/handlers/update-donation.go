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

type HandlerUpdateDonation struct {
	donationRepo *repositories.DonationRepository
	userRepo     *repositories.UserRepository
}

func HandlerUpdateDonationStart(
	donationRepo *repositories.DonationRepository,
	userRepo *repositories.UserRepository,
) *HandlerUpdateDonation {
	return &HandlerUpdateDonation{
		donationRepo,
		userRepo,
	}
}

func (h *HandlerUpdateDonation) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.UpdateDonationReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerUpdateDonation) Execute(ctx c.Context, data *dto.UpdateDonationReq) (*dto.UpdateDonationRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type == entities.EnumUserTypeNurse {
		return nil, utils.ErrorForbidden("User does not have permission to update donation", "user.forbidden")
	}

	donation, err := h.donationRepo.GetDonationById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation")
	}
	if donation == nil {
		return nil, fluxgo.ErrorNotFound("Donation not found")
	}

	fieldsToUpdate := 0
	req := repositories.UpdateDonationRepositoryReq{
		IdDonation: data.Id,
		IdUser:     data.ActionBy,
	}
	validator := data.ValidateUpdateDonationOptionalFields()

	if user.Type == entities.EnumUserTypeCommon {
		if donation.CreatedBy != user.IdUser {
			return nil, utils.ErrorForbidden("You don't have permission to access this resource", "donation.forbidden")
		}

		if validator.HasUserFeedback {
			req.UserFeedback = data.UserFeedback
			fieldsToUpdate++
		}
	}

	if user.Type == entities.EnumUserTypeAdmin {
		if validator.HasQuantityDonated {
			req.QuantityDonated = data.QuantityDonated
			fieldsToUpdate++
		}

		if validator.HasIsActive {
			req.IsActive = data.IsActive
			fieldsToUpdate++
		}
	}

	if fieldsToUpdate == 0 {
		return nil, fluxgo.ErrorBadRequest("At least one field must be sent to update", "no_fields_to_update")
	}

	err = h.donationRepo.UpdateDonation(ctx, &req)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to update donation")
	}

	donation, err = h.donationRepo.GetDonationById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation")
	}
	if donation == nil {
		return nil, fluxgo.ErrorNotFound("Donation not found")
	}

	return &dto.UpdateDonationRes{
		Donation: *donation,
	}, nil
}
