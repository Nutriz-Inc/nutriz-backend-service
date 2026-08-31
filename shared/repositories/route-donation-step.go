package repositories

import (
	c "context"
	"time"

	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/jmoiron/sqlx"
)

type RouteDonationStepRepository struct {
	fluxgo.Repository[entities.RouteDonationStep]
}

func RouteDonationStepRepositoryStart(db *fluxgo.Database) *RouteDonationStepRepository {
	return &RouteDonationStepRepository{*fluxgo.NewRepository[entities.RouteDonationStep](db)}
}

func (r *RouteDonationStepRepository) GetRouteDonationStepsByIdRoute(
	ctx c.Context,
	idRoute string,
) (*[]entities.RouteDonationStep, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.List[entities.RouteDonationStep](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "route_donation_step"
		 WHERE id_route = $1 AND removed_at IS NULL
		 ORDER BY stop_order ASC`,
		idRoute,
	)
}

type CreateRouteDonationStepRepositoryReq struct {
	IdRouteDonationStep string
	IdRoute             string
	IdDonationStep      string
	IdUser              string
	StopOrder           int16
}

func (r *RouteDonationStepRepository) createRouteDonationStep(
	ctx c.Context,
	exec sqlx.ExtContext,
	data *CreateRouteDonationStepRepositoryReq,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		INSERT INTO route_donation_step (
			id_route_donation_step,
			id_route,
			id_donation_step,
			stop_order,
			created_at,
			created_by
		) VALUES (
			:id_route_donation_step,
			:id_route,
			:id_donation_step,
			:stop_order,
			now(),
			:id_user
		)
	`

	params := map[string]any{
		"id_route_donation_step": data.IdRouteDonationStep,
		"id_route":               data.IdRoute,
		"id_donation_step":       data.IdDonationStep,
		"id_user":                data.IdUser,
		"stop_order":             data.StopOrder,
	}

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *RouteDonationStepRepository) CreateRouteDonationStepTx(
	ctx c.Context,
	tx *sqlx.Tx,
	data *CreateRouteDonationStepRepositoryReq,
) error {
	return r.createRouteDonationStep(
		ctx,
		tx,
		data,
	)
}

func (r *RouteDonationStepRepository) CreateRouteDonationStep(
	ctx c.Context,
	data *CreateRouteDonationStepRepositoryReq,
) error {
	return r.createRouteDonationStep(
		ctx,
		r.DB.WriteDB(),
		data,
	)
}

func (r *RouteDonationStepRepository) GetRouteDonationStepById(
	ctx c.Context,
	id string,
) (*entities.RouteDonationStep, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.RouteDonationStep](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "route_donation_step" WHERE id_route_donation_step = $1 AND removed_at IS NULL`,
		id,
	)
}

func (r *RouteDonationStepRepository) removeRouteDonationStep(
	ctx c.Context,
	exec sqlx.ExtContext,
	id string,
	removedBy string,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		UPDATE route_donation_step
		SET removed_at = now(),
		    removed_by = :removed_by
		WHERE id_route_donation_step = :id_route_donation_step AND removed_at IS NULL
	`

	params := map[string]any{
		"id_route_donation_step": id,
		"removed_by":             removedBy,
	}

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *RouteDonationStepRepository) RemoveRouteDonationStepTx(
	ctx c.Context,
	tx *sqlx.Tx,
	id string,
	removedBy string,
) error {
	return r.removeRouteDonationStep(ctx, tx, id, removedBy)
}

func (r *RouteDonationStepRepository) RemoveRouteDonationStep(
	ctx c.Context,
	id string,
	removedBy string,
) error {
	return r.removeRouteDonationStep(ctx, r.DB.WriteDB(), id, removedBy)
}

type RouteDonationStepWithLocation struct {
	entities.RouteDonationStep
	Latitude  *float64 `db:"latitude" json:"latitude"`
	Longitude *float64 `db:"longitude" json:"longitude"`
}

func (r *RouteDonationStepRepository) GetRouteDonationStepsWithLocationByIdRoute(
	ctx c.Context,
	idRoute string,
) (*[]RouteDonationStepWithLocation, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.List[RouteDonationStepWithLocation](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT
			rds.*,
			a.latitude AS latitude,
			a.longitude AS longitude
		 FROM route_donation_step rds
		 INNER JOIN donation_step ds ON ds.id_donation_step = rds.id_donation_step
		 LEFT JOIN address a ON a.id_address = ds.id_address AND a.removed_at IS NULL
		 WHERE rds.id_route = $1 AND rds.removed_at IS NULL
		 ORDER BY rds.stop_order ASC`,
		idRoute,
	)
}

func (r *RouteDonationStepRepository) updateStopOrder(
	ctx c.Context,
	exec sqlx.ExtContext,
	id string,
	stopOrder int16,
	updatedBy string,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		UPDATE route_donation_step
		SET stop_order = :stop_order,
		    updated_at = now(),
		    updated_by = :updated_by
		WHERE id_route_donation_step = :id_route_donation_step AND removed_at IS NULL
	`

	params := map[string]any{
		"id_route_donation_step": id,
		"stop_order":             stopOrder,
		"updated_by":             updatedBy,
	}

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *RouteDonationStepRepository) UpdateStopOrderTx(
	ctx c.Context,
	tx *sqlx.Tx,
	id string,
	stopOrder int16,
	updatedBy string,
) error {
	return r.updateStopOrder(ctx, tx, id, stopOrder, updatedBy)
}

func (r *RouteDonationStepRepository) updateDateStart(
	ctx c.Context,
	exec sqlx.ExtContext,
	id string,
	updatedBy string,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		UPDATE route_donation_step
		SET date_start = now(),
		    updated_at = now(),
		    updated_by = :updated_by
		WHERE id_route_donation_step = :id_route_donation_step AND removed_at IS NULL
	`

	params := map[string]any{
		"id_route_donation_step": id,
		"updated_by":             updatedBy,
	}

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *RouteDonationStepRepository) UpdateDateStartTx(
	ctx c.Context,
	tx *sqlx.Tx,
	id string,
	updatedBy string,
) error {
	return r.updateDateStart(ctx, tx, id, updatedBy)
}

func (r *RouteDonationStepRepository) setStartedStepsDateEnd(
	ctx c.Context,
	exec sqlx.ExtContext,
	idRoute string,
	updatedBy string,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		UPDATE route_donation_step
		SET date_end = now(),
		    updated_at = now(),
		    updated_by = :updated_by
		WHERE id_route = :id_route
		  AND date_start IS NOT NULL
		  AND date_end IS NULL
		  AND removed_at IS NULL
	`

	params := map[string]any{
		"id_route":   idRoute,
		"updated_by": updatedBy,
	}

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *RouteDonationStepRepository) SetStartedStepsDateEndTx(
	ctx c.Context,
	tx *sqlx.Tx,
	idRoute string,
	updatedBy string,
) error {
	return r.setStartedStepsDateEnd(ctx, tx, idRoute, updatedBy)
}

type RouteDonationStepWithAddress struct {
	entities.RouteDonationStep
	AddrIdAddress       *string    `db:"addr_id_address"`
	AddrIdUser          *string    `db:"addr_id_user"`
	AddrIdDonationPoint *string    `db:"addr_id_donation_point"`
	AddrZipcode         *string    `db:"addr_zipcode"`
	AddrStreet          *string    `db:"addr_street"`
	AddrNumber          *string    `db:"addr_number"`
	AddrCity            *string    `db:"addr_city"`
	AddrState           *string    `db:"addr_state"`
	AddrNeighborhood    *string    `db:"addr_neighborhood"`
	AddrComplement      *string    `db:"addr_complement"`
	AddrLatitude        *float64   `db:"addr_latitude"`
	AddrLongitude       *float64   `db:"addr_longitude"`
	AddrCreatedAt       *time.Time `db:"addr_created_at"`
	AddrUpdatedAt       *time.Time `db:"addr_updated_at"`
}

func (s RouteDonationStepWithAddress) Address() *entities.Address {
	if s.AddrIdAddress == nil {
		return nil
	}

	return &entities.Address{
		IdAddress:       *s.AddrIdAddress,
		IdUser:          s.AddrIdUser,
		IdDonationPoint: s.AddrIdDonationPoint,
		Zipcode:         utils.DerefString(s.AddrZipcode),
		Street:          utils.DerefString(s.AddrStreet),
		Number:          s.AddrNumber,
		City:            utils.DerefString(s.AddrCity),
		State:           utils.DerefString(s.AddrState),
		Neighborhood:    utils.DerefString(s.AddrNeighborhood),
		Complement:      s.AddrComplement,
		Latitude:        s.AddrLatitude,
		Longitude:       s.AddrLongitude,
		CreatedAt:       utils.DerefTime(s.AddrCreatedAt),
		UpdatedAt:       s.AddrUpdatedAt,
	}
}

func (r *RouteDonationStepRepository) GetRouteDonationStepsWithAddressByIdRoute(
	ctx c.Context,
	idRoute string,
) (*[]RouteDonationStepWithAddress, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.List[RouteDonationStepWithAddress](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT
			rds.*,
			a.id_address AS addr_id_address,
			a.id_user AS addr_id_user,
			a.id_donation_point AS addr_id_donation_point,
			a.zipcode AS addr_zipcode,
			a.street AS addr_street,
			a.number AS addr_number,
			a.city AS addr_city,
			a.state AS addr_state,
			a.neighborhood AS addr_neighborhood,
			a.complement AS addr_complement,
			a.latitude AS addr_latitude,
			a.longitude AS addr_longitude,
			a.created_at AS addr_created_at,
			a.updated_at AS addr_updated_at
		 FROM route_donation_step rds
		 INNER JOIN donation_step ds ON ds.id_donation_step = rds.id_donation_step
		 LEFT JOIN address a ON a.id_address = ds.id_address AND a.removed_at IS NULL
		 WHERE rds.id_route = $1 AND rds.removed_at IS NULL
		 ORDER BY rds.stop_order ASC`,
		idRoute,
	)
}
