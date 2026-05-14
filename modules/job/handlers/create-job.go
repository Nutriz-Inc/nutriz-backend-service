package handlers

import (
	"context"
	"time"

	dto "nutriz-backend-service/modules/job/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

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

func (h *HandlerCreateJob) Execute(ctx context.Context, data *dto.CreateJobReq) (*dto.CreateJobRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}

	if user.Type != "nurse" {
		return nil, utils.ErrorForbidden(
			"Only nurses can create jobs",
			"job.forbidden",
		)
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

	if donationStep.Status != entities.EnumDonationStepStatusPending {
		return nil, fluxgo.ErrorBadRequest(
			"Donation step is not pending",
			"STEP_NOT_PENDING",
		)
	}

	idJob := utils.IdGenerate(utils.JobEntity)
	now := time.Now()

	repoData := &repositories.CreateJobRepositoryReq{
		IdJob:        idJob,
		IdUser:       user.IdUser,
		IdStep:       data.IdStep,
		Name:         data.Name,
		Description:  data.Description,
		DateSet:      data.DateSet,
		UserFeedback: data.UserFeedback,
		CreatedBy:    user.IdUser,
	}

	err = h.jobRepo.CreateJob(ctx, repoData)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to create job")
	}

	return &dto.CreateJobRes{
		Job: entities.Job{
			IdJob:        idJob,
			IdUser:       user.IdUser,
			IdStep:       data.IdStep,
			Name:         data.Name,
			Description:  data.Description,
			DateSet:      data.DateSet,
			UserFeedback: data.UserFeedback,
			CreatedAt:    now,
			CreatedBy:    user.IdUser,
		},
	}, nil
}
