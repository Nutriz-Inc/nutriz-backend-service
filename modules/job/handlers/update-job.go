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

type HandlerUpdateJob struct {
	jobRepo  *repositories.JobRepository
	userRepo *repositories.UserRepository
}

func HandlerUpdateJobStart(
	jobRepo *repositories.JobRepository,
	userRepo *repositories.UserRepository,
) *HandlerUpdateJob {
	return &HandlerUpdateJob{
		jobRepo:  jobRepo,
		userRepo: userRepo,
	}
}

func (h *HandlerUpdateJob) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.UpdateJobReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerUpdateJob) Execute(ctx c.Context, data *dto.UpdateJobReq) (*dto.UpdateJobRes, *fluxgo.GlobalError) {
	actor, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if actor == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}

	if actor.Type != entities.EnumUserTypeAdmin && actor.Type != entities.EnumUserTypeNurse {
		return nil, utils.ErrorForbidden(
			"Only admins and nurses can update jobs",
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

	if job.Status != entities.EnumJobStatusPending {
		return nil, fluxgo.ErrorBadRequest(
			"Only pending jobs can be updated",
			"job.not_pending",
		)
	}

	fieldsToUpdate := 0
	repoData := &repositories.UpdateJobRepositoryReq{
		IdJob:     data.Id,
		UpdatedBy: data.ActionBy,
	}

	switch actor.Type {
	case entities.EnumUserTypeAdmin:
		if err := h.handleAdmUpdate(ctx, data, job, repoData, &fieldsToUpdate); err != nil {
			return nil, err
		}
	case entities.EnumUserTypeNurse:
		if err := h.handleNurseUpdate(ctx, data, job, repoData, &fieldsToUpdate); err != nil {
			return nil, err
		}
	}

	if fieldsToUpdate == 0 {
		return nil, fluxgo.ErrorBadRequest(
			"At least one field must be different to update",
			"job.no_fields_to_update",
		)
	}

	err = h.jobRepo.UpdateJob(ctx, repoData)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to update job")
	}

	updatedJob, err := h.jobRepo.GetJobById(ctx, data.Id)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get updated job")
	}
	if updatedJob == nil {
		return nil, fluxgo.ErrorNotFound("Updated job not found")
	}

	return &dto.UpdateJobRes{
		JobOut: entities.NewJobOut(*updatedJob),
	}, nil
}

func (h *HandlerUpdateJob) handleAdmUpdate(ctx c.Context, data *dto.UpdateJobReq, job *entities.Job, repoData *repositories.UpdateJobRepositoryReq, fieldsToUpdate *int) *fluxgo.GlobalError {
	validator := data.ValidateUpdateJobOptionalFields()

	if validator.HasStatus || validator.HasUserFeedback {
		return utils.ErrorForbidden(
			"Admins can only update id_user, description and date_set",
			"job.invalid_fields_for_admin",
		)
	}

	if validator.HasIdUser && *data.IdUser != job.IdUser {
		assignee, err := h.userRepo.GetUserById(ctx, *data.IdUser)
		if err != nil {
			return fluxgo.ErrorInternalError("Error to get assignee user")
		}
		if assignee == nil {
			return fluxgo.ErrorNotFound("Assignee user not found")
		}
		if assignee.Type != entities.EnumUserTypeNurse {
			return utils.ErrorForbidden(
				"Jobs can only be assigned to nurses",
				"job.invalid_assignee",
			)
		}
		repoData.IdUser = data.IdUser
		(*fieldsToUpdate)++
	}

	if validator.HasDescription && *data.Description != job.Description {
		repoData.Description = data.Description
		(*fieldsToUpdate)++
	}

	if validator.HasDateSet && (job.DateSet == nil || *data.DateSet != job.DateSet.Format("2006-01-02T15:04:05Z07:00")) {
		if !utils.IsFutureDate(*data.DateSet) {
			return fluxgo.ErrorBadRequest(
				"Set date cannot be in the past",
				"job.invalid_set_date",
			)
		}
		setDateTime, err := utils.StringToTime(*data.DateSet)
		if err != nil {
			return fluxgo.ErrorBadRequest(
				"Invalid set date format",
				"job.invalid_set_date_format",
			)
		}
		repoData.DateSet = setDateTime
		(*fieldsToUpdate)++
	}

	return nil
}

func (h *HandlerUpdateJob) handleNurseUpdate(ctx c.Context, data *dto.UpdateJobReq, job *entities.Job, repoData *repositories.UpdateJobRepositoryReq, fieldsToUpdate *int) *fluxgo.GlobalError {
	validator := data.ValidateUpdateJobOptionalFields()

	if validator.HasIdUser || validator.HasDescription || validator.HasDateSet {
		return utils.ErrorForbidden(
			"Nurses can only update status and user_feedback",
			"job.invalid_fields_for_nurse",
		)
	}

	if validator.HasStatus && *data.Status != job.Status {
		if !validator.HasUserFeedback {
			return fluxgo.ErrorBadRequest(
				"User feedback is required when updating status",
				"job.feedback_required",
			)
		}
		repoData.Status = data.Status
		(*fieldsToUpdate)++
	}

	if validator.HasUserFeedback && (job.UserFeedback == nil || *data.UserFeedback != *job.UserFeedback) {
		repoData.UserFeedback = data.UserFeedback
		(*fieldsToUpdate)++
	}

	return nil
}
