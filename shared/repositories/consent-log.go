package repositories

import (
	c "context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
)

type ConsentLogRepository struct {
	fluxgo.Repository[entities.ConsentLog]
}

func ConsentLogRepositoryStart(db *fluxgo.Database) *ConsentLogRepository {
	return &ConsentLogRepository{*fluxgo.NewRepository[entities.ConsentLog](db)}
}

type CreateConsentRepositoryReq struct {
	TermsVersion string
	Ip           string
	UserAgent    string
	IdUser       string
	IdConsentLog string
}

func (r *ConsentLogRepository) CreateConsentLog(
	ctx c.Context,
	data *CreateConsentRepositoryReq,
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
		"id_consent_log": data.IdConsentLog,
		"id_user":        data.IdUser,
		"terms_version":  data.TermsVersion,
		"ip_address":     data.Ip,
		"user_agent":     data.UserAgent,
	}

	return utils.Insert(
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		query,
		params,
	)
}

func (r *ConsentLogRepository) GetConsentLogById(ctx c.Context, id string) (*entities.ConsentLog, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.ConsentLog](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "consent_log" WHERE id_consent_log = $1`,
		id,
	)
}
