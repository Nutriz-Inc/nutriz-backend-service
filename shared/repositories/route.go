package repositories

import (
	c "context"
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
	"github.com/jmoiron/sqlx"
)

type RouteRepository struct {
	fluxgo.Repository[entities.Route]
}

func RouteRepositoryStart(db *fluxgo.Database) *RouteRepository {
	return &RouteRepository{*fluxgo.NewRepository[entities.Route](db)}
}

func (r *RouteRepository) ListRoutesByFilters(
	ctx c.Context,
	filter *dto.ListRoutesReq,
) (*[]dto.RouteRes, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select(
			"r.*",
			"ud.name AS driver_name",
		).
		From("route", "r").
		Join(q.Join{
			Table: "user",
			As:    "ud",
			On:    "ud.id_user = r.id_driver AND ud.removed_at IS NULL",
			Type:  q.LeftJoin,
		}).
		OrderBy(q.OrderBy{Column: "r.date_set", Type: "DESC"}).
		PaginationPaged(filter.Page, filter.PageSize).
		WhereAnd(q.Where{Column: "r.removed_at", Type: "IS NULL"})

	if filter.IdDriver != nil {
		qb.WhereAnd(q.Where{
			Column: "r.id_driver",
			Type:   "=",
			Val:    *filter.IdDriver,
		})
	}

	if filter.DriverName != nil {
		qb.WhereAnd(q.Where{
			Column: "ud.name",
			Type:   "ILIKE",
			Val:    "%" + *filter.DriverName + "%",
		})
	}

	if filter.Status != nil {
		qb.WhereAnd(q.Where{
			Column: "r.status",
			Type:   "=",
			Val:    *filter.Status,
		})
	}

	if filter.DateSet != nil {
		qb.WhereAnd(q.Where{
			Column: "DATE(r.date_set)",
			Type:   "=",
			Val:    *filter.DateSet,
		})
	}

	return utils.ListQuery[dto.RouteRes](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(filter.PageSize),
		true,
	)
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
