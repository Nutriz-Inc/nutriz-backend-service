package repositories

import (
	"context"
	"database/sql"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/entities"

	q "github.com/MMortari/go-query-builder"

	fluxgo "github.com/MMortari/FluxGo"
	"go.opentelemetry.io/otel/codes"
)

type DonationPointRepository struct {
	fluxgo.Repository[entities.DonationPoint]
}

func DonationPointRepositoryStart(db *fluxgo.Database) *DonationPointRepository {
	return &DonationPointRepository{*fluxgo.NewRepository[entities.DonationPoint](db)}
}

func (r *DonationPointRepository) ListDonationPointsByFilters(
	ctx context.Context,
	filter *dto.ListDonationPointsReq,
) ([]entities.DonationPoint, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select("dp.*").
		From("donation_point", "dp").
		OrderBy(q.OrderBy{Column: "dp.created_at"}).
		PaginationPaged(filter.Page, filter.PageSize).
		WhereAnd(q.Where{Column: "dp.removed_at", Type: "IS NULL"})

	if filter.Name != nil {
		qb.WhereAnd(q.Where{
			Column: "dp.name",
			Type:   "ILIKE",
			Val:    "%" + *filter.Name + "%",
		})
	}
	if filter.Cnpj != nil {
		qb.WhereAnd(q.Where{
			Column: "dp.cnpj",
			Type:   "=",
			Val:    *filter.Cnpj,
		})
	}
	if filter.HasHome != nil {
		qb.WhereAnd(q.Where{
			Column: "dp.has_home",
			Type:   "=",
			Val:    *filter.HasHome,
		})
	}

	query, args := qb.ToSelectSql()
	resp := make([]entities.DonationPoint, 0, filter.PageSize)

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
