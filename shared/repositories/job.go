package repositories

import (
	"context"
	dto "nutriz-backend-service/modules/job/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"time"

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
		WhereAnd(q.Where{Column: "j.removed_at", Type: "IS NULL"}).
		WhereAnd(q.Where{Column: "j.id_user", Type: "=", Val: filter.ActionBy})

	if filter.DateSet != nil {
		qb.WhereAnd(q.Where{
			Column: "DATE(j.date_set)",
			Type:   "=",
			Val:    *filter.DateSet,
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

type CreateJobRepositoryReq struct {
	IdJob        string
	IdUser       string
	IdStep       string
	Name         string
	Description  string
	DateSet      *time.Time
	UserFeedback *string
	CreatedBy    string
}

func (r *JobRepository) CreateJob(
	ctx context.Context,
	data *CreateJobRepositoryReq,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		INSERT INTO job (
			id_job,
			id_user,
			id_step,
			name,
			description,
			date_set,
			user_feedback,
			created_at,
			created_by
		) VALUES (
			:id_job,
			:id_user,
			:id_step,
			:name,
			:description,
			:date_set,
			:user_feedback,
			now(),
			:created_by
		)
	`

	params := map[string]any{
		"id_job":        data.IdJob,
		"id_user":       data.IdUser,
		"id_step":       data.IdStep,
		"name":          data.Name,
		"description":   data.Description,
		"date_set":      data.DateSet,
		"user_feedback": data.UserFeedback,
		"created_by":    data.CreatedBy,
	}

	return utils.Insert(
		ctx,
		r.DB.WriteDB(),
		span,
		query,
		params,
	)
}
