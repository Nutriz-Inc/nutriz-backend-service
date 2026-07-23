package handlers

import (
	c "context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"

	dto "nutriz-backend-service/modules/job/dtos"

	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerGetJob struct {
	jobRepo  *repositories.JobRepository
	userRepo *repositories.UserRepository
}

func HandlerGetJobStart(jobRepo *repositories.JobRepository, userRepo *repositories.UserRepository) *HandlerGetJob {
	return &HandlerGetJob{
		jobRepo,
		userRepo,
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
	user, err := h.userRepo.GetUserById(ctx, filters.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type == entities.EnumUserTypeCommon {
		return nil, utils.ErrorForbidden("User does not have permission to get job", "user.forbidden")
	}

	job, err := h.jobRepo.GetJobInfoById(ctx, filters.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get job")
	}
	if job == nil {
		return nil, fluxgo.ErrorNotFound("Job not found")
	}

	if user.Type == entities.EnumUserTypeNurse {
		if job.IdUser != filters.ActionBy {
			return nil, utils.ErrorForbidden("You don't have permission to access this resource", "job.forbidden")
		}
	}

	return &dto.GetJobRes{
		JobInfoRes: *job,
	}, nil
}
