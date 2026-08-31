package repositories

import (
	c "context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
	"github.com/jmoiron/sqlx"
)

type BottleRepository struct {
	fluxgo.Repository[entities.Bottle]
}

func BottleRepositoryStart(db *fluxgo.Database) *BottleRepository {
	return &BottleRepository{*fluxgo.NewRepository[entities.Bottle](db)}
}

func (r *BottleRepository) GetBottlesByIdDonation(
	ctx c.Context,
	idDonation string,
) (*[]entities.Bottle, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select("b.*").
		From("bottle", "b").
		OrderBy(q.OrderBy{Column: "b.created_at"}).
		WhereAnd(q.Where{Column: "b.id_donation", Type: "=", Val: idDonation})

	return utils.ListQuery[entities.Bottle](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		nil,
		false,
	)
}

func (r *BottleRepository) GetBottleById(ctx c.Context, id string) (*entities.Bottle, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.Bottle](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "bottle" WHERE id_bottle = $1`,
		id,
	)
}

type CreateBottleRepositoryReq struct {
	IdBottle          string
	IdDonation        string
	IdUser            string
	QuantityDonatedMl *float64
	Discarded         *bool
	Description       *string
}

func (r *BottleRepository) createBottle(
	ctx c.Context,
	exec sqlx.ExtContext,
	data *CreateBottleRepositoryReq,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		INSERT INTO bottle (
			id_bottle,
			id_donation,
			quantity_donated_ml,
			discarded,
			description,
			created_at,
			created_by
		) VALUES (
			:id_bottle,
			:id_donation,
			:quantity_donated_ml,
			:discarded,
			:description,
			now(),
			:id_user
		)
	`

	params := map[string]any{
		"id_bottle":           data.IdBottle,
		"id_donation":         data.IdDonation,
		"id_user":             data.IdUser,
		"quantity_donated_ml": data.QuantityDonatedMl,
		"discarded":           data.Discarded,
		"description":         data.Description,
	}

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *BottleRepository) CreateBottleTx(
	ctx c.Context,
	tx *sqlx.Tx,
	data *CreateBottleRepositoryReq,
) error {
	return r.createBottle(
		ctx,
		tx,
		data,
	)
}

func (r *BottleRepository) CreateBottle(
	ctx c.Context,
	data *CreateBottleRepositoryReq,
) error {
	return r.createBottle(
		ctx,
		r.DB.WriteDB(),
		data,
	)
}

func (r *BottleRepository) deleteBottlesByIdDonation(
	ctx c.Context,
	exec sqlx.ExtContext,
	idDonation string,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `DELETE FROM bottle WHERE id_donation = :id_donation`

	params := map[string]any{
		"id_donation": idDonation,
	}

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *BottleRepository) DeleteBottlesByIdDonationTx(
	ctx c.Context,
	tx *sqlx.Tx,
	idDonation string,
) error {
	return r.deleteBottlesByIdDonation(ctx, tx, idDonation)
}
