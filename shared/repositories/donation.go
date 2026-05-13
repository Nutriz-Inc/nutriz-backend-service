package repositories

import (
	"context"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
)

type DonationRepository struct {
	fluxgo.Repository[entities.Donation]
}

func DonationRepositoryStart(db *fluxgo.Database) *DonationRepository {
	return &DonationRepository{*fluxgo.NewRepository[entities.Donation](db)}
}

func (r *DonationRepository) ListDonationByFilters(
	ctx context.Context,
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

	qb.WhereAnd(q.Where{
		Column: "d.created_by",
		Type:   "=",
		Val:    filter.ActionBy,
	})

	return utils.ListQuery[entities.Donation](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(filter.PageSize),
		true,
	)
}

func (r *DonationRepository) GetDonationById(ctx context.Context, id string) (*entities.Donation, error) {
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
	IdDonation  string
	IdUser      string
	IsActive    bool
	Description string
}

func (r *DonationRepository) CreateDonation(
	ctx context.Context,
	data *CreateDonationRepositoryReq,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		INSERT INTO donation (
			id_donation,
			is_active,
			description,
			created_by,
			created_at
		) VALUES (
		 	:id_donation,
			:is_active,
			:description,
			:id_user,
			now()
		)
	`

	params := map[string]any{
		"id_donation": data.IdDonation,
		"id_user":     data.IdUser,
		"is_active":   data.IsActive,
		"description": data.Description,
	}

	return utils.Insert(
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		query,
		params,
	)
}
