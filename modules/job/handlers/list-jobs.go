package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/job/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerListJobs struct {
	jobRepo *repositories.JobRepository
}

func HandlerListJobsStart(jobRepo *repositories.JobRepository) *HandlerListJobs {
	return &HandlerListJobs{jobRepo}
}

func (h *HandlerListJobs) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.ListJobsReq))
	if err != nil {
		return nil, err
	}

	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerListJobs) Execute(ctx c.Context, filters *dto.ListJobsReq) (*dto.ListJobsRes, *fluxgo.GlobalError) {
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
