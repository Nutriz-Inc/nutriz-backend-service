package repositories

import (
	"context"
	dto "nutriz-backend-service/modules/job/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"strings"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
	"github.com/jmoiron/sqlx"
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

	if filter.IdStep != nil {
		qb.WhereAnd(q.Where{
			Column: "j.id_step",
			Type:   "=",
			Val:    *filter.IdStep,
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
	Status       entities.EnumJobStatus
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
			status,
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
			:status,
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
		"status":        data.Status,
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

type UpdateJobRepositoryReq struct {
	IdJob        string
	IdUser       *string
	Status       *entities.EnumJobStatus
	Description  *string
	DateSet      *time.Time
	UserFeedback *string
	UpdatedBy    string
}

func (r *JobRepository) updateJob(ctx context.Context, exec sqlx.ExtContext, data *UpdateJobRepositoryReq) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	sets := []string{}
	params := map[string]any{
		"id_job":     data.IdJob,
		"updated_by": data.UpdatedBy,
	}

	if data.Status != nil {
		sets = append(sets, "status = :status")
		params["status"] = *data.Status
	}
	if data.DateSet != nil {
		sets = append(sets, "date_set = :date_set")
		params["date_set"] = *data.DateSet
	}
	if data.UserFeedback != nil {
		sets = append(sets, "user_feedback = :user_feedback")
		params["user_feedback"] = *data.UserFeedback
	}
	if data.Description != nil {
		sets = append(sets, "description = :description")
		params["description"] = *data.Description
	}
	if data.IdUser != nil {
		sets = append(sets, "id_user = :id_user")
		params["id_user"] = *data.IdUser
	}

	query := `
		UPDATE job
		SET ` + strings.Join(sets, ", ") + `,
		    updated_at = now(),
			updated_by = :updated_by
		WHERE id_job = :id_job AND removed_at IS NULL
	`

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *JobRepository) UpdateJobTx(
	ctx context.Context,
	tx *sqlx.Tx,
	data *UpdateJobRepositoryReq,
) error {
	return r.updateJob(
		ctx,
		tx,
		data,
	)
}

func (r *JobRepository) UpdateJob(
	ctx context.Context,
	data *UpdateJobRepositoryReq,
) error {
	return r.updateJob(
		ctx,
		r.DB.WriteDB(),
		data,
	)
}
