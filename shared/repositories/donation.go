package repositories

import (
	"context"
	"database/sql"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/entities"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
	"go.opentelemetry.io/otel/codes"
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
) ([]entities.Donation, int, error) {
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

	query, args := qb.ToSelectSql()
	resp := make([]entities.Donation, 0, filter.PageSize)

	err := r.DB.ReadOnlyDB().SelectContext(ctx, &resp, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, 0, err
	}

	queryTotal, argsTotal := qb.ToSelectTotalSql()

	var total int
	err = r.DB.ReadOnlyDB().GetContext(ctx, &total, queryTotal, argsTotal...)

	if err == sql.ErrNoRows {
		return nil, total, nil
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, total, err
	}

	return resp, total, nil
}
