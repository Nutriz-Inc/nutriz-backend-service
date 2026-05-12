package handlers

import (
	"context"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/gofiber/fiber/v2"
)

type HandlerCreateAddress struct {
	addressRepo *repositories.AddressRepository
	userRepo    *repositories.UserRepository
}

func HandlerCreateAddressStart(
	consentRepo *repositories.AddressRepository,
	userRepo *repositories.UserRepository,
) *HandlerCreateAddress {
	return &HandlerCreateAddress{
		addressRepo: consentRepo,
		userRepo:    userRepo,
	}
}

func (h *HandlerCreateAddress) HandleHttp(c *fiber.Ctx, income interface{}) (*fluxgo.GlobalResponse, *fluxgo.GlobalError) {
	resp, err := h.Execute(c.UserContext(), income.(*dto.CreateConsentReq))
	if err != nil {
		return nil, err
	}
	return &fluxgo.GlobalResponse{Content: resp, Status: 200}, nil
}

func (h *HandlerCreateAddress) Execute(ctx context.Context, data *dto.CreateConsentReq) (*dto.CreateConsentRes, *fluxgo.GlobalError) {
	user, err := h.userRepo.GetUserById(ctx, data.ActionBy)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get user")
	}
	if user == nil {
		return nil, fluxgo.ErrorNotFound("User not found")
	}

	provider, err := 

	idAddress := utils.IdGenerate(utils.AddressEntity)

	repoData := &repositories.CreateConsentRepositoryReq{
		TermsVersion: data.TermsVersion,
		Ip:           data.IpAddress,
		UserAgent:    data.UserAgent,
		IdUser:       user.IdUser,
		IdAddress:    idAddress,
	}

	err = h.addressRepo.CreateAddress(ctx, repoData)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to create consent")
	}

	address, err := h.addressRepo.GetAddressById(ctx, idAddress)
	if err != nil {
		return nil, fluxgo.ErrorInternalError("Error to get consent log")
	}
	if address == nil {
		return nil, fluxgo.ErrorNotFound("Consent log not found")
	}

	return &dto.CreateConsentRes{
		Address: *consentLog,
	}, nil
}
