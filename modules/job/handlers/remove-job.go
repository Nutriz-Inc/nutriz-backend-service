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

type HandlerRemoveJob struct {
	jobRepo  *repositories.JobRepository
	userRepo *repositories.UserRepository
}

func HandlerRemoveJobStart(
	jobRepo *repositories.JobRepository,
	userRepo *repositories.UserRepository,
) *HandlerRemoveJob {
	return &HandlerRemoveJob{
		jobRepo:  jobRepo,
		userRepo: userRepo,
	}
}

func (h *HandlerRemoveJob) HandleHttp(c *fiber.Ctx, income any) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.RemoveJobReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerRemoveJob) Execute(ctx c.Context, data *dto.RemoveJobReq) (*dto.RemoveJobRes, *fluxgo.GlobalError) {
	actor, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if actor == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}

	if actor.Type != entities.EnumUserTypeAdmin {
		return nil, utils.ErrorForbidden(
			"Only admins can remove jobs",
			"job.forbidden",
		)
	}

	job, err := h.jobRepo.GetJobById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get job")
	}
	if job == nil {
		return nil, fluxgo.ErrorNotFound("Job not found")
	}

	err = h.jobRepo.RemoveJob(ctx, data.Id, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to remove job")
	}

	removed, err := h.jobRepo.GetJobById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get job")
	}

	return &dto.RemoveJobRes{
		Success: removed == nil,
	}, nil
}
