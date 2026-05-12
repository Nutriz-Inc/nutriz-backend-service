package repositories

import (
	"context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
	"github.com/paemuri/brdoc/v2"
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

func (r *AddressRepository) GetAddressById(ctx context.Context, id string) (*entities.Address, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.Address](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "address" WHERE id_address = $1`,
		id,
	)
}

type CreateAddressRepositoryReq struct {
	IdAddress    string
	IdUser       string
	Zipcode      string
	Street       string
	Number       *string
	City         string
	State        brdoc.FederativeUnit
	Neighborhood string
	Complement   *string
	Latitude     *float64
	Longitude    *float64
}

func (r *AddressRepository) CreateAddress(
	ctx context.Context,
	data *CreateAddressRepositoryReq,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		INSERT INTO address (
			id_address,
			id_user,
			zipcode,
			street,
			number,
			city,
			state,
			neighborhood,
			complement,
			latitude,
			longitude,
			created_at
		) VALUES (
		 	:id_address,
		 	:id_user,
		 	:zipcode,
		 	:street,
		 	:number,
		 	:city,
		 	:state,
		 	:neighborhood,
		 	:complement,
		 	:latitude,
		 	:longitude,
			now()
		)
	`

	params := map[string]any{
		"id_address":   data.IdAddress,
		"id_user":      data.IdUser,
		"zipcode":      data.Zipcode,
		"street":       data.Street,
		"number":       data.Number,
		"city":         data.City,
		"state":        data.State,
		"neighborhood": data.Neighborhood,
		"complement":   data.Complement,
		"latitude":     data.Latitude,
		"longitude":    data.Longitude,
	}

	return utils.Insert(
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		query,
		params,
	)
}
