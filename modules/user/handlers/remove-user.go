package handlers

import (
	c "context"
	"fmt"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type HandlerRemoveUser struct {
	db           *fluxgo.Database
	userRepo     *repositories.UserRepository
	donationRepo *repositories.DonationRepository
}

func HandlerRemoveUserStart(
	db *fluxgo.Database,
	userRepo *repositories.UserRepository,
	donationRepo *repositories.DonationRepository,
) *HandlerRemoveUser {
	return &HandlerRemoveUser{
		db,
		userRepo,
		donationRepo,
	}
}

func (h *HandlerRemoveUser) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.RemoveUserReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerRemoveUser) Execute(ctx c.Context, data *dto.RemoveUserReq) (*dto.RemoveUserRes, *fluxgo.GlobalError) {
	userAction, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user action")
	}
	if userAction == nil {
		return nil, fluxgo.ErrorNotFound("User action not found")
	}

	user, err := h.userRepo.GetUserById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}

	if !userAction.CanRemoveUser(*user) {
		return nil, utils.ErrorForbidden("You don't have permission to access this resource", "user.forbidden")
	}

	err = h.db.RunTransaction(ctx, func(ctx c.Context, tx *sqlx.Tx) error {
		err := h.userRepo.RemoveUserTx(ctx, tx, data.Id, data.ActionBy)
		if err != nil {
			return fmt.Errorf("error to remove user: %w", err)
		}

		err = h.donationRepo.DisableUserDonationsTx(ctx, tx, data.Id, data.ActionBy)
		if err != nil {
			return fmt.Errorf("error to create donation step timeline: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fluxgo.ErrorInternalError(err.Error())
	}

	user, err = h.userRepo.GetUserById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}

	return &dto.RemoveUserRes{
		Success: user == nil,
	}, nil
}
