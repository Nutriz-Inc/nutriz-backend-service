package repositories

import (
	c "context"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"strings"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
	"github.com/jmoiron/sqlx"
)

type DonationRepository struct {
	fluxgo.Repository[entities.Donation]
}

func DonationRepositoryStart(db *fluxgo.Database) *DonationRepository {
	return &DonationRepository{*fluxgo.NewRepository[entities.Donation](db)}
}

func (r *DonationRepository) ListDonationByFilters(
	ctx c.Context,
	filter *dto.ListDonationReq,
) (*[]dto.DonationRes, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	currentStepNameSubquery := `(
		SELECT ds.name
		FROM "donation_step" ds
		WHERE ds.id_donation = d.id_donation
		ORDER BY ds.created_at DESC
		LIMIT 1
	)`

	currentStepStatusSubquery := `(
		SELECT ds.status
		FROM "donation_step" ds
		WHERE ds.id_donation = d.id_donation
		ORDER BY ds.created_at DESC
		LIMIT 1
	)`

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select(
			"d.*",
			"COALESCE("+currentStepNameSubquery+"::text, '') AS current_step",
			"COALESCE("+currentStepStatusSubquery+" = '"+string(entities.EnumDonationStepStatusFailed)+"', false) AS has_error",
		).
		From("donation", "d").
		Join(q.Join{
			Table: "user",
			As:    "u",
			On:    "u.id_user = d.created_by AND u.removed_at IS NULL",
			Type:  q.LeftJoin,
		}).
		OrderBy(q.OrderBy{Column: "d.created_at", Type: "DESC"}).
		PaginationPaged(filter.Page, filter.PageSize).
		WhereAnd(q.Where{Column: "d.removed_at", Type: "IS NULL"})

	if filter.ActionBy == nil {
		qb.Select(
			"u.name AS user_name",
			"u.cpf AS user_document",
		)
	}

	if filter.IsActive != nil {
		qb.WhereAnd(q.Where{
			Column: "d.is_active",
			Type:   "=",
			Val:    *filter.IsActive,
		})
	}

	if filter.ActionBy != nil {
		qb.WhereAnd(q.Where{
			Column: "d.created_by",
			Type:   "=",
			Val:    *filter.ActionBy,
		})
	}

	if filter.IdUserCommon != nil {
		qb.WhereAnd(q.Where{
			Column: "d.created_by",
			Type:   "=",
			Val:    *filter.IdUserCommon,
		})
	}

	if filter.UserDocument != nil {
		qb.WhereAnd(q.Where{
			Column: "u.cpf",
			Type:   "=",
			Val:    *filter.UserDocument,
		})
	}

	if filter.UserName != nil {
		qb.WhereAnd(q.Where{
			Column: "u.name",
			Type:   "ILIKE",
			Val:    "%" + *filter.UserName + "%",
		})
	}

	if filter.CurrentStep != nil {
		qb.WhereAnd(q.Where{
			Column: currentStepNameSubquery,
			Type:   "=",
			Val:    *filter.CurrentStep,
		})
	}

	rows, total, err := utils.ListQuery[donationRow](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(filter.PageSize),
		true,
	)
	if err != nil || rows == nil {
		return nil, total, err
	}

	donations := make([]dto.DonationRes, len(*rows))
	for i, row := range *rows {
		donations[i] = dto.DonationRes{
			DonationOut:  entities.NewDonationOut(row.Donation),
			UserName:     row.UserName,
			UserDocument: row.UserDocument,
			CurrentStep:  row.CurrentStep,
			HasError:     row.HasError,
		}
	}

	return &donations, total, nil
}

// donationRow is a scan-only shape used to map SQL columns (including the
// computed current_step/has_error and the joined user fields) via sqlx.
// It never leaves this file - callers only see dto.DonationRes.
type donationRow struct {
	entities.Donation
	UserName     *string                    `db:"user_name"`
	UserDocument *string                    `db:"user_document"`
	CurrentStep  entities.EnumDonationSteps `db:"current_step"`
	HasError     bool                       `db:"has_error"`
}

func (r *DonationRepository) GetDonationById(ctx c.Context, id string) (*entities.Donation, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.Donation](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "donation" WHERE id_donation = $1 AND removed_at IS NULL`,
		id,
	)
}

type CreateDonationRepositoryReq struct {
	IdDonation string
	IdUser     string
	IsActive   bool
}

func (r *DonationRepository) CreateDonation(
	ctx c.Context,
	data *CreateDonationRepositoryReq,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		INSERT INTO donation (
			id_donation,
			is_active,
			created_by,
			created_at
		) VALUES (
		 	:id_donation,
			:is_active,
			:id_user,
			now()
		)
	`

	params := map[string]any{
		"id_donation": data.IdDonation,
		"id_user":     data.IdUser,
		"is_active":   data.IsActive,
	}

	return utils.Insert(
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		query,
		params,
	)
}

type UpdateDonationRepositoryReq struct {
	IdDonation      string
	IdUser          string
	IsActive        *bool
	QuantityDonated *float64
	UserFeedback    *string
	ScoreFeedback   *int16
}

func (r *DonationRepository) updateDonation(ctx c.Context, exec sqlx.ExtContext, data *UpdateDonationRepositoryReq) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	sets := []string{}
	params := map[string]any{
		"id_donation": data.IdDonation,
		"updated_by":  data.IdUser,
	}

	if data.IsActive != nil {
		sets = append(sets, "is_active = :is_active")
		params["is_active"] = data.IsActive
	}
	if data.QuantityDonated != nil {
		sets = append(sets, "quantity_donated = :quantity_donated")
		params["quantity_donated"] = data.QuantityDonated
	}
	if data.UserFeedback != nil {
		sets = append(sets, "user_feedback = :user_feedback")
		params["user_feedback"] = data.UserFeedback
	}
	if data.ScoreFeedback != nil {
		sets = append(sets, "score_feedback = :score_feedback")
		params["score_feedback"] = data.ScoreFeedback
	}

	if len(sets) == 0 {
		return nil
	}

	query := `
		UPDATE donation
		SET ` + strings.Join(sets, ", ") + `,
		    updated_at = now(),
			updated_by = :updated_by
		WHERE id_donation = :id_donation
		  AND removed_at IS NULL
	`

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *DonationRepository) UpdateDonationTx(
	ctx c.Context,
	tx *sqlx.Tx,
	data *UpdateDonationRepositoryReq,
) error {
	return r.updateDonation(
		ctx,
		tx,
		data,
	)
}

func (r *DonationRepository) UpdateDonation(
	ctx c.Context,
	data *UpdateDonationRepositoryReq,
) error {
	return r.updateDonation(
		ctx,
		r.DB.WriteDB(),
		data,
	)
}

func (r *DonationRepository) disableUserDonations(ctx c.Context, exec sqlx.ExtContext, idUser, actionBy string) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		UPDATE "donation"
		set is_active = false, updated_at = now(), updated_by = :action_by
		WHERE created_by = :id_user AND is_active = true AND removed_at IS NULL
	`

	params := map[string]any{
		"id_user":   idUser,
		"action_by": actionBy,
	}

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *DonationRepository) DisableUserDonationsTx(
	ctx c.Context,
	tx *sqlx.Tx,
	id, actionBy string,
) error {
	return r.disableUserDonations(
		ctx,
		tx,
		id,
		actionBy,
	)
}

func (r *DonationRepository) DisableUserDonations(
	ctx c.Context,
	id, actionBy string,
) error {
	return r.disableUserDonations(
		ctx,
		r.DB.WriteDB(),
		id,
		actionBy,
	)
}

func (r *DonationRepository) CountCompletedDonations(ctx c.Context, idUser string) (int64, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	result, err := utils.Get[int64](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`
		SELECT COUNT(*)
		FROM (
			SELECT d.id_donation
			FROM donation d
			JOIN donation_step ds
				ON ds.id_donation = d.id_donation
			WHERE d.is_active = false
				AND d.removed_at IS NULL
				AND d.created_by = $1
			GROUP BY d.id_donation
			HAVING COUNT(*) = 4
				AND COUNT(*) FILTER (WHERE ds.status = 'done') = 4
		) completed
		`,
		idUser,
	)

	if err != nil {
		return 0, err
	}

	if result == nil {
		return 0, nil
	}

	return *result, nil
}
