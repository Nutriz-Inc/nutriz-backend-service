package repositories

import (
	c "context"
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
	ctx c.Context,
	filter *dto.ListJobsReq,
) (*[]dto.JobRes, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select(
			"j.*",
			"un.name AS user_nurse_name",
			"uc.name AS user_common_name",
			`a.id_address         AS "address.id_address"`,
			`a.id_user            AS "address.id_user"`,
			`a.id_donation_point  AS "address.id_donation_point"`,
			`a.zipcode            AS "address.zipcode"`,
			`a.street             AS "address.street"`,
			`a.number             AS "address.number"`,
			`a.city               AS "address.city"`,
			`a.state              AS "address.state"`,
			`a.neighborhood       AS "address.neighborhood"`,
			`a.complement         AS "address.complement"`,
			`a.latitude           AS "address.latitude"`,
			`a.longitude          AS "address.longitude"`,
			`a.created_at         AS "address.created_at"`,
			`a.updated_at         AS "address.updated_at"`,
			`a.removed_at         AS "address.removed_at"`,
		).
		From("job", "j").
		Join(q.Join{
			Table: "user",
			As:    "un",
			On:    "un.id_user = j.id_user AND un.removed_at IS NULL",
			Type:  q.LeftJoin,
		}).
		Join(q.Join{
			Table: "donation_step",
			As:    "ds",
			On:    "ds.id_donation_step = j.id_step",
			Type:  q.LeftJoin,
		}).
		Join(q.Join{
			Table: "donation",
			As:    "d",
			On:    "d.id_donation = ds.id_donation",
			Type:  q.LeftJoin,
		}).
		Join(q.Join{
			Table: "user",
			As:    "uc",
			On:    "uc.id_user = d.created_by AND uc.removed_at IS NULL",
			Type:  q.LeftJoin,
		}).
		Join(q.Join{
			Table: "address",
			As:    "a",
			On:    "a.id_address = ds.id_address AND a.removed_at IS NULL",
			Type:  q.LeftJoin,
		}).
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

	if filter.IdUserCommon != nil {
		qb.WhereAnd(q.Where{
			Column: "d.created_by",
			Type:   "=",
			Val:    *filter.IdUserCommon,
		})
	}

	if filter.IdUserNurse != nil {
		qb.WhereAnd(q.Where{
			Column: "j.id_user",
			Type:   "=",
			Val:    *filter.IdUserNurse,
		})
	}

	if filter.UserCommonName != nil {
		qb.WhereAnd(q.Where{
			Column: "uc.name",
			Type:   "ILIKE",
			Val:    "%" + *filter.UserCommonName + "%",
		})
	}

	if filter.UserNurseName != nil {
		qb.WhereAnd(q.Where{
			Column: "un.name",
			Type:   "ILIKE",
			Val:    "%" + *filter.UserNurseName + "%",
		})
	}

	return utils.ListQuery[dto.JobRes](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(filter.PageSize),
		true,
	)
}

func (r *JobRepository) GetJobById(ctx c.Context, id string) (*entities.Job, error) {
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
	ctx c.Context,
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

func (r *JobRepository) RemoveJob(ctx c.Context, id, actionBy string) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		UPDATE job
		SET removed_at = now(),
		    removed_by = $2			
		WHERE id_job = $1 AND removed_at IS NULL
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

type UpdateJobRepositoryReq struct {
	IdJob        string
	IdUser       *string
	Status       *entities.EnumJobStatus
	Description  *string
	DateSet      *time.Time
	UserFeedback *string
	UpdatedBy    string
}

func (r *JobRepository) updateJob(ctx c.Context, exec sqlx.ExtContext, data *UpdateJobRepositoryReq) error {
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
	ctx c.Context,
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
	ctx c.Context,
	data *UpdateJobRepositoryReq,
) error {
	return r.updateJob(
		ctx,
		r.DB.WriteDB(),
		data,
	)
}
