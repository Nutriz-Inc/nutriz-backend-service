package repositories

import (
	"context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
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
) (*[]entities.Address, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select("a.*").
		From("address", "a").
		PaginationPaged(1, entities.MAX_ADDRESS_QUANTITY_PER_USER).
		OrderBy(q.OrderBy{Column: "a.created_at"}).
		WhereAnd(q.Where{Column: "a.removed_at", Type: "IS NULL"}).
		WhereAnd(q.Where{Column: "a.id_user", Type: "=", Val: userId})

	return utils.ListQuery[entities.Address](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(entities.MAX_ADDRESS_QUANTITY_PER_USER),
		false,
	)
}
