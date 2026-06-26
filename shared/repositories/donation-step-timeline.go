package repositories

import (
	c "context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
	"github.com/jmoiron/sqlx"
)

type DonationStepTimelineRepository struct {
	fluxgo.Repository[entities.DonationStepTimeline]
}

func DonationStepTimelineRepositoryStart(db *fluxgo.Database) *DonationStepTimelineRepository {
	return &DonationStepTimelineRepository{*fluxgo.NewRepository[entities.DonationStepTimeline](db)}
}

type CreateDonationStepTimelineRepositoryReq struct {
	IdDonationStepTimeline string
	IdDonationStep         string
	Description            string
	Status                 entities.EnumDonationStepStatus
	SetDate                *time.Time
	IdUser                 string
}

func (r *DonationStepTimelineRepository) createDonationStepTimeline(
	ctx c.Context,
	exec sqlx.ExtContext,
	data *CreateDonationStepTimelineRepositoryReq,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		INSERT INTO donation_step_timeline (
			id_donation_step_timeline,
			id_donation_step,
			description,
			status,
			set_date,
			created_by,
			created_at
		) VALUES (
		 	:id_donation_step_timeline,
			:id_donation_step,
			:description,
			:status,
			:set_date,
			:action_by,
			NOW()
		)
	`

	params := map[string]any{
		"id_donation_step_timeline": data.IdDonationStepTimeline,
		"id_donation_step":          data.IdDonationStep,
		"description":               data.Description,
		"status":                    data.Status,
		"set_date":                  data.SetDate,
		"action_by":                 data.IdUser,
	}

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *DonationStepTimelineRepository) CreateDonationStepTimelineTx(
	ctx c.Context,
	tx *sqlx.Tx,
	data *CreateDonationStepTimelineRepositoryReq,
) error {
	return r.createDonationStepTimeline(
		ctx,
		tx,
		data,
	)
}

func (r *DonationStepTimelineRepository) CreateDonationStepTimeline(
	ctx c.Context,
	data *CreateDonationStepTimelineRepositoryReq,
) error {
	return r.createDonationStepTimeline(
		ctx,
		r.DB.WriteDB(),
		data,
	)
}

func (r *DonationStepTimelineRepository) GetDonationStepsTimelineByIdDonationStep(
	ctx c.Context,
	idDonationStep string,
) (*[]entities.DonationStepTimeline, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select("dtl.*").
		From("donation_step_timeline", "dtl").
		OrderBy(q.OrderBy{Column: "dtl.created_at"}).
		PaginationPaged(1, 50).
		WhereAnd(q.Where{Column: "dtl.id_donation_step", Type: "=", Val: idDonationStep})

	return utils.ListQuery[entities.DonationStepTimeline](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(50),
		false,
	)
}
