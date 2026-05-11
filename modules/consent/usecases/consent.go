package usecases

import (
	"context"
	"nutriz-backend-service/modules/consent/dtos"
	"nutriz-backend-service/shared/repositories"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

	idConsentLog := uuid.New().String()
	//idConsentLog := utils.IdGenerate(utils.ConsentEntity)

	err := uc.consentRepo.CreateConsent(ctx, *data, idUser, idConsentLog)
	if err != nil {

		log.Ctx(ctx).Error().Err(err).
			Str("id_user", idUser).
			Str("terms_version", data.TermsVersion).
			Msg("failed to create consent")

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
