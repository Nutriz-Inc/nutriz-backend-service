package handlers

import (
	c "context"
	"errors"
	"fmt"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type HandlerUpdateDonation struct {
	db           *fluxgo.Database
	donationRepo *repositories.DonationRepository
	userRepo     *repositories.UserRepository
}

func HandlerUpdateDonationStart(
	db *fluxgo.Database,
	donationRepo *repositories.DonationRepository,
	userRepo *repositories.UserRepository,
) *HandlerUpdateDonation {
	return &HandlerUpdateDonation{
		db,
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
	if !user.Action().CanUpdateDonation {
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

		if validator.HasFeedback {
			req.UserFeedback = data.UserFeedback
			req.ScoreFeedback = data.ScoreFeedback
			fieldsToUpdate = fieldsToUpdate + 2
		}
	}

	if user.Type == entities.EnumUserTypeAdmin {
		if validator.HasIsActive {
			req.IsActive = data.IsActive
			fieldsToUpdate++
		}
	}

	if fieldsToUpdate == 0 {
		return nil, fluxgo.ErrorBadRequest("At least one field must be sent to update", "no_fields_to_update")
	}

	err = h.db.RunTransaction(ctx, func(ctx c.Context, tx *sqlx.Tx) error {
		err := h.donationRepo.UpdateDonationTx(ctx, tx, &req)
		if err != nil {
			return fmt.Errorf("error to update donation: %w", err)
		}

		return nil
	})
	if err != nil {
		var txErr *utils.TxError
		if errors.As(err, &txErr) {
			return nil, txErr.Err
		}
		return nil, fluxgo.ErrorInternalError(err.Error())
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
