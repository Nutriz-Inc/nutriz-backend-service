package repositories

import (
	"context"
	"database/sql"
	"nutriz-backend-service/shared/entities"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
	"go.opentelemetry.io/otel/codes"
)

type AddressRepository struct {
	fluxgo.Repository[entities.Address]
}

func AddressRepositoryStart(db *fluxgo.Database) *AddressRepository {
	return &AddressRepository{*fluxgo.NewRepository[entities.Address](db)}
}

func (r *AddressRepository) GetAddressesByUserId(
	ctx context.Context,
	userId string,
) (*[]entities.Address, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select("a.*").
		From("address", "a").
		PaginationPaged(1, entities.MAX_ADDRESS_QUANTITY_PER_USER).
		OrderBy(q.OrderBy{Column: "a.created_at"}).
		WhereAnd(q.Where{Column: "a.removed_at", Type: "IS NULL"}).
		WhereAnd(q.Where{Column: "a.id_user", Type: "=", Val: userId})

	query, args := qb.ToSelectSql()
	resp := make([]entities.Address, 0, entities.MAX_ADDRESS_QUANTITY_PER_USER)

	err := r.DB.ReadOnlyDB().SelectContext(ctx, &resp, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	queryTotal, argsTotal := qb.ToSelectTotalSql()

	var total int
	err = r.DB.ReadOnlyDB().GetContext(ctx, &total, queryTotal, argsTotal...)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &resp, nil
}
