package repositories

import (
	"context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"strings"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
	"github.com/jmoiron/sqlx"
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
		true,
	)
}

func (r *AddressRepository) GetAddressById(ctx context.Context, id string) (*entities.Address, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.Address](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "address" WHERE id_address = $1 AND removed_at IS NULL`,
		id,
	)
}

func (r *AddressRepository) GetAddressByZipcodeAndIdUser(ctx context.Context, zipcode string, userId string) (*entities.Address, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.Address](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "address" WHERE zipcode = $1 AND id_user = $2 AND removed_at IS NULL`,
		zipcode,
		userId,
	)
}

type CreateAddressRepositoryReq struct {
	IdAddress    string
	IdUser       string
	Zipcode      string
	Street       string
	Number       *string
	City         string
	State        string
	Neighborhood string
	Complement   *string
	Latitude     *float64
	Longitude    *float64
}

func (r *AddressRepository) createAddress(
	ctx context.Context,
	exec sqlx.ExtContext,
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

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *AddressRepository) CreateAddressTx(
	ctx context.Context,
	tx *sqlx.Tx,
	data *CreateAddressRepositoryReq,
) error {
	return r.createAddress(
		ctx,
		tx,
		data,
	)
}

func (r *AddressRepository) CreateAddress(
	ctx context.Context,
	data *CreateAddressRepositoryReq,
) error {
	return r.createAddress(
		ctx,
		r.DB.WriteDB(),
		data,
	)
}

type UpdateAddressRepositoryReq struct {
	IdAddress    string
	IdUser       string
	Zipcode      *string
	Street       *string
	Number       *string
	City         *string
	State        *string
	Neighborhood *string
	Complement   *string
	Latitude     *float64
	Longitude    *float64
}

func (r *AddressRepository) UpdateAddress(ctx context.Context, data UpdateAddressRepositoryReq) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	sets := []string{}
	params := map[string]any{
		"id_address": data.IdAddress,
		"updated_by": data.IdUser,
	}

	if data.Zipcode != nil {
		sets = append(sets, "zipcode = :zipcode")
		params["zipcode"] = data.Zipcode
	}
	if data.Street != nil {
		sets = append(sets, "street = :street")
		params["street"] = data.Street
	}
	if data.Number != nil {
		sets = append(sets, "number = :number")
		params["number"] = data.Number
	}
	if data.City != nil {
		sets = append(sets, "city = :city")
		params["city"] = data.City
	}
	if data.State != nil {
		sets = append(sets, "state = :state")
		params["state"] = data.State
	}
	if data.Neighborhood != nil {
		sets = append(sets, "neighborhood = :neighborhood")
		params["neighborhood"] = data.Neighborhood
	}
	if data.Complement != nil {
		sets = append(sets, "complement = :complement")
		params["complement"] = data.Complement
	}
	if data.Latitude != nil {
		sets = append(sets, "latitude = :latitude")
		params["latitude"] = data.Latitude
	}
	if data.Longitude != nil {
		sets = append(sets, "longitude = :longitude")
		params["longitude"] = data.Longitude
	}

	if len(sets) == 0 {
		return nil
	}

	query := `
		UPDATE address
		SET ` + strings.Join(sets, ", ") + `,
		    updated_at = now()
		WHERE id_address = :id_address
		  AND removed_at IS NULL
		  AND id_user = :updated_by
	`

	return utils.Update(
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		query,
		params,
	)
}

func (r *AddressRepository) RemoveAddress(ctx context.Context, id, actionBy string) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		UPDATE address
		SET removed_at = now()
		WHERE id_address = $1 AND id_user = $2 AND removed_at IS NULL
	`

	return utils.Delete(
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		query,
		id,
		actionBy,
	)
}
