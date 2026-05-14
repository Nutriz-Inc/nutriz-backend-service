package repositories

import (
	"context"
	dto "nutriz-backend-service/modules/job/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
)

type JobRepository struct {
	fluxgo.Repository[entities.Job]
}

func JobRepositoryStart(db *fluxgo.Database) *JobRepository {
	return &JobRepository{*fluxgo.NewRepository[entities.Job](db)}
}

func (r *JobRepository) ListJobsByFilters(
	ctx context.Context,
	filter *dto.ListJobsReq,
) (*[]entities.Job, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select("j.*").
		From("job", "j").
		OrderBy(q.OrderBy{Column: "j.date_set"}).
		PaginationPaged(filter.Page, filter.PageSize).
		WhereAnd(q.Where{Column: "j.removed_at", Type: "IS NULL"})

	if filter.DateSet != nil {
		qb.WhereAnd(q.Where{
			Column: "DATE(j.date_set)",
			Type:   "=",
			Val:    *filter.DateSet,
		})
	}

	if filter.ActionBy != nil {
		qb.WhereAnd(q.Where{
			Column: "j.id_user",
			Type:   "=",
			Val:    *filter.ActionBy,
		})
	}

	return utils.ListQuery[entities.Job](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(filter.PageSize),
		true,
	)
}

func (r *JobRepository) GetJobById(ctx context.Context, id string) (*entities.Job, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.Job](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "job" WHERE id_job = $1 AND removed_at IS NULL`,
		id,
	)
}
