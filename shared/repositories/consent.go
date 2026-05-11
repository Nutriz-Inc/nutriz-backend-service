package repositories

import (
	"context"
	"nutriz-backend-service/modules/user/consent/dtos"
	"nutriz-backend-service/shared/entities"

	fluxgo "github.com/MMortari/FluxGo"
)

type IConsentRepository interface {
	CreateConsent(
		ctx context.Context,
		data dtos.CreateConsentReq,
		idUser string,
		idConsentLog string,
	) error
}

type ConsentRepository struct {
	*fluxgo.Repository[entities.ConsentLog]
}

func ConsentRepositoryStart(
	db *fluxgo.Database,
) *ConsentRepository {
	return &ConsentRepository{
		fluxgo.NewRepository[entities.ConsentLog](db),
	}
}

func (r *ConsentRepository) CreateConsent(
	ctx context.Context,
	data dtos.CreateConsentReq,
	idUser string,
	idConsentLog string,
) error {

	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		INSERT INTO consent_log (
			id_consent_log,
			id_user,
			terms_version,
			accepted_at,
			ip_address,
			user_agent
		) VALUES (
			:id_consent_log,
			:id_user,
			:terms_version,
			now(),
			:ip_address,
			:user_agent
		)
	`

	params := map[string]any{
		"id_consent_log": idConsentLog,
		"id_user":        idUser,
		"terms_version":  data.TermsVersion,
		"ip_address":     data.IpAddress,
		"user_agent":     data.UserAgent,
	}

	_, err := r.DB.WriteDB().NamedExecContext(ctx, query, params)
	if err != nil {
		span.SetError(err)
		return err
	}

	return nil
}
