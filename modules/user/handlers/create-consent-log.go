package handlers

import (
	c "context"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerCreateConsentLog struct {
	consentLogRepo *repositories.ConsentLogRepository
	userRepo       *repositories.UserRepository
}

func HandlerCreateConsentLogStart(
	consentRepo *repositories.ConsentLogRepository,
	userRepo *repositories.UserRepository,
) *HandlerCreateConsentLog {
	return &HandlerCreateConsentLog{
		consentLogRepo: consentRepo,
		userRepo:       userRepo,
	}
}

func (h *HandlerCreateConsentLog) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.CreateConsentReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 201}, nil
}

func (h *HandlerCreateConsentLog) Execute(ctx c.Context, data *dto.CreateConsentReq) (*dto.CreateConsentRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}
	if user.Type != entities.EnumUserTypeCommon {
		return nil, utils.ErrorForbidden("User does not have permission to create consent log", "user.forbidden")
	}

	idConsentLog := utils.IdGenerate(utils.ConsentLogEntity)

	repoData := &repositories.CreateConsentRepositoryReq{
		TermsVersion: data.TermsVersion,
		Ip:           data.IpAddress,
		UserAgent:    data.UserAgent,
		IdUser:       user.IdUser,
		IdConsentLog: idConsentLog,
	}

	err = h.consentLogRepo.CreateConsentLog(ctx, repoData)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to create consent")
	}

	consentLog, err := h.consentLogRepo.GetConsentLogById(ctx, idConsentLog)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get consent log")
	}
	if consentLog == nil {
		return nil, fluxgo.ErrorNotFound("Consent log not found")
	}

	return &dto.CreateConsentRes{
		ConsentLogOut: entities.NewConsentLogOut(*consentLog),
	}, nil
}
