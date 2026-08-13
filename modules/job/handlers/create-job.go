package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/job/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerCreateJob struct {
	jobRepo          *repositories.JobRepository
	userRepo         *repositories.UserRepository
	donationStepRepo *repositories.DonationStepRepository
	donationRepo     *repositories.DonationRepository
}

func HandlerCreateJobStart(
	jobRepo *repositories.JobRepository,
	userRepo *repositories.UserRepository,
	donationStepRepo *repositories.DonationStepRepository,
	donationRepo *repositories.DonationRepository,
) *HandlerCreateJob {
	return &HandlerCreateJob{
		jobRepo:          jobRepo,
		userRepo:         userRepo,
		donationStepRepo: donationStepRepo,
		donationRepo:     donationRepo,
	}
}

func (h *HandlerCreateJob) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.CreateJobReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 201}, nil
}

func (h *HandlerCreateJob) Execute(ctx c.Context, data *dto.CreateJobReq) (*dto.CreateJobRes, *fluxgo.GlobalError) {
	creator, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if creator == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}

	if creator.Type != entities.EnumUserTypeAdmin {
		return nil, utils.ErrorForbidden(
			"Only adms can create jobs",
			"job.forbidden",
		)
	}

	assignee, err := h.userRepo.GetUserById(ctx, data.IdUser)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get assignee user")
	}
	if assignee == nil {
		return nil, fluxgo.ErrorNotFound("Assignee user not found")
	}

	if assignee.Type != entities.EnumUserTypeNurse {
		return nil, utils.ErrorForbidden(
			"Jobs can only be assigned to nurses",
			"job.invalid_assignee",
		)
	}

	var setDateTime *time.Time
	if data.DateSet != nil {
		if !utils.IsFutureDate(*data.DateSet) {
			return nil, fluxgo.ErrorBadRequest(
				"Set date cannot be in the past",
				"job.invalid_set_date",
			)
		}

		setDateTime, err = utils.StringToTime(*data.DateSet)
		if err != nil {
			return nil, fluxgo.ErrorBadRequest(
				"Invalid set date format",
				"job.invalid_set_date_format",
			)
		}
	}

	donationStep, err := h.donationStepRepo.GetDonationStepById(ctx, data.IdStep)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation step")
	}
	if donationStep == nil {
		return nil, fluxgo.ErrorNotFound("Donation step not found")
	}

	donation, err := h.donationRepo.GetDonationById(ctx, donationStep.IdDonation)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get donation")
	}
	if donation == nil {
		return nil, fluxgo.ErrorNotFound("Donation not found")
	}

	if !donation.IsActive {
		return nil, fluxgo.ErrorBadRequest(
			"Donation is not active",
			"DONATION_NOT_ACTIVE",
		)
	}

	idJob := utils.IdGenerate(utils.JobEntity)
	initalStatus := entities.EnumJobStatusPending // All new jobs start with pending status, and it's the nurse's responsibility to update it to done or failed when they complete the job

	repoData := &repositories.CreateJobRepositoryReq{
		IdJob:       idJob,
		IdUser:      data.IdUser,
		IdStep:      data.IdStep,
		Status:      initalStatus,
		Name:        data.Name,
		Description: data.Description,
		DateSet:     setDateTime,
		CreatedBy:   data.ActionBy,
	}

	err = h.jobRepo.CreateJob(ctx, repoData)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to create job")
	}

	job, err := h.jobRepo.GetJobById(ctx, idJob)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get created job")
	}
	if job == nil {
		return nil, fluxgo.ErrorNotFound("Created job not found")
	}

	return &dto.CreateJobRes{
		JobOut: entities.NewJobOut(*job),
	}, nil
}
