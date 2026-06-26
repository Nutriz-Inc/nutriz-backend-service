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
) (*[]entities.Donation, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select("d.*").
		From("donation", "d").
		OrderBy(q.OrderBy{Column: "d.created_at"}).
		PaginationPaged(filter.Page, filter.PageSize).
		WhereAnd(q.Where{Column: "d.removed_at", Type: "IS NULL"})

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

	return utils.ListQuery[entities.Donation](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(filter.PageSize),
		true,
	)
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
}

func (r *DonationRepository) UpdateDonation(ctx c.Context, data UpdateDonationRepositoryReq) error {
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
