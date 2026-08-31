package repositories

import (
	c "context"
	dto "nutriz-backend-service/modules/route/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"strings"
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

	if filter.Name != nil {
		qb.WhereAnd(q.Where{
			Column: "r.name",
			Type:   "ILIKE",
			Val:    "%" + *filter.Name + "%",
		})
	}

	if filter.City != nil {
		qb.WhereAnd(q.Where{
			Column: "r.city",
			Type:   "ILIKE",
			Val:    "%" + *filter.City + "%",
		})
	}

	if filter.Neighborhood != nil {
		qb.WhereAnd(q.Where{
			Column: "r.neighborhood",
			Type:   "ILIKE",
			Val:    "%" + *filter.Neighborhood + "%",
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

type UpdateRouteRepositoryReq struct {
	IdRoute      string
	UpdatedBy    string
	Name         *string
	City         *string
	Neighborhood *string
	Status       *entities.EnumRouteStatus
	Description  *string
	DateSet      *time.Time
	SetDateStart bool
	SetDateEnd   bool
	Mileage      *float64
	UserFeedback *string
}

func (r *RouteRepository) updateRoute(
	ctx c.Context,
	exec sqlx.ExtContext,
	data *UpdateRouteRepositoryReq,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	sets := []string{}
	params := map[string]any{
		"id_route":   data.IdRoute,
		"updated_by": data.UpdatedBy,
	}

	if data.Name != nil {
		sets = append(sets, "name = :name")
		params["name"] = *data.Name
	}
	if data.City != nil {
		sets = append(sets, "city = :city")
		params["city"] = *data.City
	}
	if data.Neighborhood != nil {
		sets = append(sets, "neighborhood = :neighborhood")
		params["neighborhood"] = *data.Neighborhood
	}
	if data.Status != nil {
		sets = append(sets, "status = :status")
		params["status"] = *data.Status
	}
	if data.Description != nil {
		sets = append(sets, "description = :description")
		params["description"] = *data.Description
	}
	if data.DateSet != nil {
		sets = append(sets, "date_set = :date_set")
		params["date_set"] = *data.DateSet
	}
	if data.Mileage != nil {
		sets = append(sets, "mileage = :mileage")
		params["mileage"] = *data.Mileage
	}
	if data.UserFeedback != nil {
		sets = append(sets, "user_feedback = :user_feedback")
		params["user_feedback"] = *data.UserFeedback
	}
	if data.SetDateStart {
		sets = append(sets, "date_start = now()")
	}
	if data.SetDateEnd {
		sets = append(sets, "date_end = now()")
	}

	if len(sets) == 0 {
		return nil
	}

	query := `
		UPDATE route
		SET ` + strings.Join(sets, ", ") + `,
		    updated_at = now(),
			updated_by = :updated_by
		WHERE id_route = :id_route AND removed_at IS NULL
	`

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *RouteRepository) UpdateRouteTx(
	ctx c.Context,
	tx *sqlx.Tx,
	data *UpdateRouteRepositoryReq,
) error {
	return r.updateRoute(
		ctx,
		tx,
		data,
	)
}

func (r *RouteRepository) UpdateRoute(
	ctx c.Context,
	data *UpdateRouteRepositoryReq,
) error {
	return r.updateRoute(
		ctx,
		r.DB.WriteDB(),
		data,
	)
}

func (r *RouteRepository) touchRoute(
	ctx c.Context,
	exec sqlx.ExtContext,
	idRoute string,
	updatedBy string,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		UPDATE route
		SET updated_at = now(),
		    updated_by = :updated_by
		WHERE id_route = :id_route AND removed_at IS NULL
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

func (r *RouteRepository) TouchRouteTx(
	ctx c.Context,
	tx *sqlx.Tx,
	idRoute string,
	updatedBy string,
) error {
	return r.touchRoute(ctx, tx, idRoute, updatedBy)
}
