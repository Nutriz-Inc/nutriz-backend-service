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
			"ds.id_address AS id_address_ref",
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

	if filter.Status != nil {
		qb.WhereAnd(q.Where{
			Column: "j.status",
			Type:   "=",
			Val:    *filter.Status,
		})
	}

	rows, total, err := utils.ListQuery[jobRow](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(filter.PageSize),
		true,
	)
	if err != nil || rows == nil {
		return nil, total, err
	}

	jobs, err := r.attachAddresses(ctx, span, *rows)
	if err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

type jobRow struct {
	dto.JobRes
	IdAddressRef *string `db:"id_address_ref"`
}

func (r *JobRepository) attachAddresses(ctx c.Context, span fluxgo.Span, rows []jobRow) (*[]dto.JobRes, error) {
	addressIds := make([]string, 0, len(rows))
	seenAddressIds := make(map[string]bool, len(rows))
	for _, row := range rows {
		if row.IdAddressRef == nil || seenAddressIds[*row.IdAddressRef] {
			continue
		}
		seenAddressIds[*row.IdAddressRef] = true
		addressIds = append(addressIds, *row.IdAddressRef)
	}

	addressById := make(map[string]entities.Address, len(addressIds))
	if len(addressIds) > 0 {
		addressQb := q.NewQueryBuilder(q.SetOtelSpan(span)).
			Select("*").
			From("address").
			WhereAnd(q.Where{Column: "id_address", Type: "IN", Val: addressIds}).
			WhereAnd(q.Where{Column: "removed_at", Type: "IS NULL"})

		addresses, _, err := utils.ListQuery[entities.Address](
			ctx,
			r.DB.ReadOnlyDB(),
			span,
			addressQb,
			nil,
			false,
		)
		if err != nil {
			return nil, err
		}
		for _, address := range *addresses {
			addressById[address.IdAddress] = address
		}
	}

	jobs := make([]dto.JobRes, len(rows))
	for i, row := range rows {
		jobs[i] = row.JobRes
		if row.IdAddressRef == nil {
			continue
		}
		if address, ok := addressById[*row.IdAddressRef]; ok {
			jobs[i].Address = &address
		}
	}

	return &jobs, nil
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

func (r *JobRepository) GetJobInfoById(ctx c.Context, id string) (*dto.JobInfoRes, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[dto.JobInfoRes](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT
			j.*,
			d.created_by AS id_user_common,
			ds.id_address AS id_address
		FROM "job" j
		LEFT JOIN donation_step ds ON ds.id_donation_step = j.id_step
		LEFT JOIN donation d ON d.id_donation = ds.id_donation
		WHERE j.id_job = $1 AND j.removed_at IS NULL`,
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
