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
	bottleRepo   *repositories.BottleRepository
	userRepo     *repositories.UserRepository
}

func HandlerUpdateDonationStart(
	db *fluxgo.Database,
	donationRepo *repositories.DonationRepository,
	bottleRepo *repositories.BottleRepository,
	userRepo *repositories.UserRepository,
) *HandlerUpdateDonation {
	return &HandlerUpdateDonation{
		db,
		donationRepo,
		bottleRepo,
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

		if validator.HasFeedback {
			req.UserFeedback = data.UserFeedback
			req.ScoreFeedback = data.ScoreFeedback
			fieldsToUpdate = fieldsToUpdate + 2
		}
	}

	if validator.HasBottles && user.Type != entities.EnumUserTypeAdmin {
		return nil, utils.ErrorForbidden("Only admins can update donation bottles", "donation.bottles_forbidden")
	}

	if user.Type == entities.EnumUserTypeAdmin {
		if validator.HasIsActive {
			req.IsActive = data.IsActive
			fieldsToUpdate++
		}

		if validator.HasBottles {
			fieldsToUpdate = fieldsToUpdate + len(*data.Bottles)
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

		if validator.HasBottles {
			if err := h.bottleRepo.DeleteBottlesByIdDonationTx(ctx, tx, data.Id); err != nil {
				return fmt.Errorf("error to remove donation bottles: %w", err)
			}

			milkDonated := 0.0
			for _, bottle := range *data.Bottles {
				err := h.bottleRepo.CreateBottleTx(ctx, tx, &repositories.CreateBottleRepositoryReq{
					IdBottle:          utils.IdGenerate(utils.BottleEntity),
					IdDonation:        data.Id,
					IdUser:            data.ActionBy,
					QuantityDonatedMl: bottle.QuantityDonatedMl,
					Discarded:         bottle.Discarded,
					Description:       bottle.Description,
				})
				if err != nil {
					return fmt.Errorf("error to create donation bottle: %w", err)
				}

				if bottle.QuantityDonatedMl != nil && (bottle.Discarded == nil || !*bottle.Discarded) {
					milkDonated = milkDonated + *bottle.QuantityDonatedMl
				}
			}

			err = h.userRepo.UpdateUserTx(ctx, tx, &repositories.UpdateUserRepositoryReq{
				IdUser:      donation.CreatedBy,
				ActionBy:    data.ActionBy,
				MilkDonated: &milkDonated,
			})
			if err != nil {
				return fmt.Errorf("error to update user milk donated: %w", err)
			}
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
