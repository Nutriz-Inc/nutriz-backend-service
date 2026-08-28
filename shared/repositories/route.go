package repositories

import (
	c "context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	"github.com/jmoiron/sqlx"
)

type RouteRepository struct {
	fluxgo.Repository[entities.Route]
}

func RouteRepositoryStart(db *fluxgo.Database) *RouteRepository {
	return &RouteRepository{*fluxgo.NewRepository[entities.Route](db)}
}

func (r *RouteRepository) GetRouteById(ctx c.Context, id string) (*entities.Route, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.Route](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "route" WHERE id_route = $1 AND removed_at IS NULL`,
		id,
	)
}

type CreateRouteRepositoryReq struct {
	IdRoute      string
	IdDriver     string
	IdUser       string
	Name         string
	Description  string
	City         *string
	Neighborhood *string
	Status       entities.EnumRouteStatus
	DateSet      time.Time
}

func (r *RouteRepository) createRoute(
	ctx c.Context,
	exec sqlx.ExtContext,
	data *CreateRouteRepositoryReq,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		INSERT INTO route (
			id_route,
			id_driver,
			name,
			description,
			city,
			neighborhood,
			status,
			date_set,
			created_at,
			created_by
		) VALUES (
			:id_route,
			:id_driver,
			:name,
			:description,
			:city,
			:neighborhood,
			:status,
			:date_set,
			now(),
			:id_user
		)
	`

	params := map[string]any{
		"id_route":     data.IdRoute,
		"id_driver":    data.IdDriver,
		"id_user":      data.IdUser,
		"name":         data.Name,
		"description":  data.Description,
		"city":         data.City,
		"neighborhood": data.Neighborhood,
		"status":       data.Status,
		"date_set":     data.DateSet,
	}

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *RouteRepository) CreateRouteTx(
	ctx c.Context,
	tx *sqlx.Tx,
	data *CreateRouteRepositoryReq,
) error {
	return r.createRoute(
		ctx,
		tx,
		data,
	)
}

func (r *RouteRepository) CreateRoute(
	ctx c.Context,
	data *CreateRouteRepositoryReq,
) error {
	return r.createRoute(
		ctx,
		r.DB.WriteDB(),
		data,
	)
}
