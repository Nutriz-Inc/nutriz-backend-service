package repositories

import (
	"context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"strings"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
	"github.com/jmoiron/sqlx"
)

type DonationStepRepository struct {
	fluxgo.Repository[entities.DonationStep]
}

func DonationStepRepositoryStart(db *fluxgo.Database) *DonationStepRepository {
	return &DonationStepRepository{*fluxgo.NewRepository[entities.DonationStep](db)}
}

func (r *DonationStepRepository) GetDonationStepsByIdDonation(
	ctx context.Context,
	idDonation string,
) (*[]entities.DonationStep, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select("ds.*").
		From("donation_step", "ds").
		OrderBy(q.OrderBy{Column: "ds.created_at"}).
		PaginationPaged(1, entities.NUMBER_OF_DONATION_STEPS).
		WhereAnd(q.Where{Column: "ds.id_donation", Type: "=", Val: idDonation})

	return utils.ListQuery[entities.DonationStep](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(entities.NUMBER_OF_DONATION_STEPS),
		false,
	)
}
func (r *DonationStepRepository) GetDonationStepById(ctx context.Context, id string) (*entities.DonationStep, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.DonationStep](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "donation_step" WHERE id_donation_step = $1`,
		id,
	)
}

type CreateDonationStepRepositoryReq struct {
	IdDonationStep string
	IdDonation     string
	IdUser         string
	Name           entities.EnumDonationSteps
	Description    string
	Status         entities.EnumDonationStepStatus
	SetDate        *time.Time
}

func (r *DonationStepRepository) createDonationStep(
	ctx context.Context,
	exec sqlx.ExtContext,
	data *CreateDonationStepRepositoryReq,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		INSERT INTO donation_step (
			id_donation_step,
			id_donation,
			name,
			description,
			status,
			set_date,
			created_at,
			created_by
		) VALUES (
			:id_donation_step,
			:id_donation,
			:name,
			:description,
			:status,
			:set_date,
			now(),
			:id_user
		)
	`

	params := map[string]any{
		"id_donation_step": data.IdDonationStep,
		"id_donation":      data.IdDonation,
		"id_user":          data.IdUser,
		"name":             data.Name,
		"description":      data.Description,
		"status":           data.Status,
		"set_date":         data.SetDate,
	}

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *DonationStepRepository) CreateDonationStepTx(
	ctx context.Context,
	tx *sqlx.Tx,
	data *CreateDonationStepRepositoryReq,
) error {
	return r.createDonationStep(
		ctx,
		tx,
		data,
	)
}

func (r *DonationStepRepository) CreateDonationStep(
	ctx context.Context,
	data *CreateDonationStepRepositoryReq,
) error {
	return r.createDonationStep(
		ctx,
		r.DB.WriteDB(),
		data,
	)
}

type UpdateDonationStepRepositoryReq struct {
	IdDonationStep string
	IdUser         string
}

func (r *DonationStepRepository) UpdateDonationStep(ctx context.Context, data UpdateDonationStepRepositoryReq) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	sets := []string{}
	params := map[string]any{
		"id_donation_step": data.IdDonationStep,
		"updated_by":       data.IdUser,
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

	return utils.Update(
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		query,
		params,
	)
}
