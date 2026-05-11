package handlers

import (
	c "context"
	"nutriz-backend-service/shared/repositories"

	dto "nutriz-backend-service/modules/job/dtos"

	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerGetJob struct {
	jobRepo *repositories.JobRepository
}

func HandlerGetJobStart(jobRepo *repositories.JobRepository) *HandlerGetJob {
	return &HandlerGetJob{
		jobRepo: jobRepo,
	}
}

func (h *HandlerGetJob) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.GetJobReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerGetJob) Execute(ctx c.Context, filters *dto.GetJobReq) (*dto.GetJobRes, *fluxgo.GlobalError) {
	job, err := h.jobRepo.GetJobById(ctx, filters.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get job")
	}
	if job == nil {
		return nil, fluxgo.ErrorNotFound("Job not found")
	}
	if job.IdUser != filters.ActionBy {
		return nil, utils.ErrorForbidden("You don't have permission to access this resource", "job.forbidden")
	}

	return &dto.GetJobRes{
		Job: *job,
	}, nil
}
