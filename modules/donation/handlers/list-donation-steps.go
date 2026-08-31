package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerListDonationSteps struct {
	donationStepRepo *repositories.DonationStepRepository
	userRepo         *repositories.UserRepository
}

func HandlerListDonationStepsStart(
	donationStepRepo *repositories.DonationStepRepository,
	userRepo *repositories.UserRepository,
) *HandlerListDonationSteps {
	return &HandlerListDonationSteps{
		donationStepRepo,
		userRepo,
	}
}

func (h *HandlerListDonationSteps) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.ListDonationStepsReq))
	if err != nil {
		return nil, err
	}

	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerListDonationSteps) Execute(ctx c.Context, filters *dto.ListDonationStepsReq) (*dto.ListDonationStepsRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, filters.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if !user.Action().CanListDonationStep {
		return nil, utils.ErrorForbidden("User does not have permission to list donation steps", "user.forbidden")
	}

	rows, total, err := h.donationStepRepo.ListDonationStepsByFilters(ctx, filters)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to list donation steps")
	}

	steps := make([]dto.DonationStepRes, 0, len(*rows))
	for _, row := range *rows {
		steps = append(steps, dto.DonationStepRes{
			DonationStep: row.DonationStep,
			Address:      row.Address(),
		})
	}

	return &dto.ListDonationStepsRes{
		Data: steps,
		PaginationRes: utils.PaginationRes{
			Page:     filters.Page,
			PageSize: filters.PageSize,
			Total:    total,
		},
	}, nil
}
