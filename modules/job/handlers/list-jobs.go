package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/job/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerListJobs struct {
	jobRepo  *repositories.JobRepository
	userRepo *repositories.UserRepository
}

func HandlerListJobsStart(jobRepo *repositories.JobRepository, userRepo *repositories.UserRepository) *HandlerListJobs {
	return &HandlerListJobs{
		jobRepo,
		userRepo,
	}
}

func (h *HandlerListJobs) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.ListJobsReq))
	if err != nil {
		return nil, err
	}

	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerListJobs) Execute(ctx c.Context, filters *dto.ListJobsReq) (*dto.ListJobsRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, filters.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type == entities.EnumUserTypeDonor {
		return nil, utils.ErrorForbidden("User does not have permission to list jobs", "user.forbidden")
	}

	jobs, total, err := h.jobRepo.ListJobsByFilters(ctx, filters)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to list jobs")
	}
	if jobs == nil || len(*jobs) == 0 {
		return nil, fluxgo.ErrorNotFound("Jobs not found")
	}

	return &dto.ListJobsRes{
		Data: *jobs,
		PaginationRes: utils.PaginationRes{
			Page:     filters.Page,
			PageSize: filters.PageSize,
			Total:    total,
		},
	}, nil
}
