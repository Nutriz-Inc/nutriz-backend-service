package usecases

import (
	"context"
	"log"
	"nutriz-backend-service/modules/user/consent/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
)

type CreateConsentUseCase struct {
	consentRepo repositories.IConsentRepository
}

func NewCreateConsentUseCase(
	consentRepo repositories.IConsentRepository,
) *CreateConsentUseCase {
	return &CreateConsentUseCase{
		consentRepo: consentRepo,
	}
}

func (uc *CreateConsentUseCase) Execute(
	ctx context.Context,
	idUser string,
	data *dtos.CreateConsentReq,
) (*dtos.CreateConsentRes, *fluxgo.GlobalError) {

	if !utils.IsValidIP(data.IpAddress) {
		return nil, fluxgo.ErrorBadRequest(
			"invalid ip address",
			"INVALID_IP",
		)
	}

	idConsentLog := utils.IdGenerate(utils.ConsentLogEntity)

	err := uc.consentRepo.CreateConsent(ctx, *data, idUser, idConsentLog)
	if err != nil {
		log.Printf("failed to create consent - user: %s, terms_version: %s, error: %v",
			idUser, data.TermsVersion, err)

		return nil, fluxgo.ErrorInternalError("Error to create consent")
	}

	return &dtos.CreateConsentRes{
		IdConsentLog: idConsentLog,
		IdUser:       idUser,
		TermsVersion: data.TermsVersion,
		IpAddress:    data.IpAddress,
		UserAgent:    data.UserAgent,
	}, nil
}
