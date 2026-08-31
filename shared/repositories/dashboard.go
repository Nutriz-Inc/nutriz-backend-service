package repositories

import (
	c "context"
	"math"
	dto "nutriz-backend-service/modules/dashboard/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
)

type DashboardRepository struct {
	DB *fluxgo.Database
}

func DashboardRepositoryStart(db *fluxgo.Database) *DashboardRepository {
	return &DashboardRepository{DB: db}
}

func (r *DashboardRepository) StartSpan(ctx c.Context, name string) (c.Context, fluxgo.Span) {
	return r.DB.StartSpan(ctx, "repository/dashboard/"+name)
}

const dateRangeFilter = `($1::timestamp IS NULL OR d.created_at >= $1) AND ($2::timestamp IS NULL OR d.created_at < $2)`

const nonDiscardedBottles = `
	FROM bottle b
	JOIN donation d ON d.id_donation = b.id_donation AND d.removed_at IS NULL
	WHERE b.discarded IS NOT TRUE
	  AND ` + dateRangeFilter

func (r *DashboardRepository) GetTotalMilkCollected(ctx c.Context, start, end *time.Time) (float64, error) {
	ctx, span := r.StartSpan(ctx, "GetTotalMilkCollected")
	defer span.End()

	result, err := utils.Get[float64](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT COALESCE(SUM(b.quantity_donated_ml), 0)`+nonDiscardedBottles,
		start,
		end,
	)
	if err != nil || result == nil {
		return 0, err
	}

	return *result, nil
}

func (r *DashboardRepository) GetMilkCollectedByMonth(ctx c.Context, start, end *time.Time) ([]dto.MilkCollectedByMonth, error) {
	ctx, span := r.StartSpan(ctx, "GetMilkCollectedByMonth")
	defer span.End()

	result, err := utils.List[dto.MilkCollectedByMonth](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`
		SELECT
			TO_CHAR(d.created_at, 'YYYY-MM') AS month,
			COALESCE(SUM(b.quantity_donated_ml), 0) AS total
		`+nonDiscardedBottles+`
		GROUP BY TO_CHAR(d.created_at, 'YYYY-MM')
		ORDER BY TO_CHAR(d.created_at, 'YYYY-MM')
		`,
		start,
		end,
	)
	if err != nil || result == nil {
		return nil, err
	}

	return *result, nil
}

func (r *DashboardRepository) GetFeedbackByScore(ctx c.Context, start, end *time.Time) ([]dto.FeedbackScoreCount, error) {
	ctx, span := r.StartSpan(ctx, "GetFeedbackByScore")
	defer span.End()

	rows, err := utils.List[dto.FeedbackScoreCount](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`
		SELECT
			d.score_feedback AS score,
			COUNT(*) AS count
		FROM donation d
		WHERE d.removed_at IS NULL
		  AND d.score_feedback IS NOT NULL
		  AND `+dateRangeFilter+`
		GROUP BY d.score_feedback
		ORDER BY d.score_feedback
		`,
		start,
		end,
	)
	if err != nil || rows == nil {
		return nil, err
	}

	countByScore := make(map[int16]int64, len(*rows))
	for _, row := range *rows {
		countByScore[row.Score] = row.Count
	}

	feedback := make([]dto.FeedbackScoreCount, 0, 5)
	for score := int16(1); score <= 5; score++ {
		feedback = append(feedback, dto.FeedbackScoreCount{
			Score: score,
			Count: countByScore[score],
		})
	}

	return feedback, nil
}

type avgServiceTimeRow struct {
	AvgHours *float64 `db:"avg_hours"`
}

func (r *DashboardRepository) GetAverageServiceTimeHours(ctx c.Context, start, end *time.Time) (*float64, error) {
	ctx, span := r.StartSpan(ctx, "GetAverageServiceTimeHours")
	defer span.End()

	result, err := utils.Get[avgServiceTimeRow](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`
		SELECT AVG(EXTRACT(EPOCH FROM (first_collect.completed_at - u.created_at)) / 3600) AS avg_hours
		FROM "user" u
		JOIN (
			SELECT d.created_by AS id_user, MIN(ds.completed_at) AS completed_at
			FROM donation_step ds
			JOIN donation d ON d.id_donation = ds.id_donation AND d.removed_at IS NULL
			WHERE ds.name = $3
			  AND ds.status = $4
			  AND ds.completed_at IS NOT NULL
			GROUP BY d.created_by
		) first_collect ON first_collect.id_user = u.id_user
		WHERE u.removed_at IS NULL
		  AND ($1::timestamp IS NULL OR first_collect.completed_at >= $1)
		  AND ($2::timestamp IS NULL OR first_collect.completed_at < $2)
		`,
		start,
		end,
		entities.EnumDonationStepCollectMilk,
		entities.EnumDonationStepStatusDone,
	)
	if err != nil || result == nil {
		return nil, err
	}

	return result.AvgHours, nil
}

func (r *DashboardRepository) GetDonationsWithErrorCount(ctx c.Context, start, end *time.Time) (int64, error) {
	ctx, span := r.StartSpan(ctx, "GetDonationsWithErrorCount")
	defer span.End()

	result, err := utils.Get[int64](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`
		SELECT COUNT(*)
		FROM donation d
		WHERE d.removed_at IS NULL
		  AND `+dateRangeFilter+`
		  AND (
			SELECT ds.status
			FROM donation_step ds
			WHERE ds.id_donation = d.id_donation
			ORDER BY ds.created_at DESC
			LIMIT 1
		  ) = $3
		`,
		start,
		end,
		entities.EnumDonationStepStatusFailed,
	)
	if err != nil || result == nil {
		return 0, err
	}

	return *result, nil
}

type activeDonationStepCountRow struct {
	Step  string `db:"step"`
	Count int64  `db:"count"`
}

func (r *DashboardRepository) GetActiveDonationsByCurrentStep(ctx c.Context, start, end *time.Time) ([]dto.ActiveDonationsByStep, error) {
	ctx, span := r.StartSpan(ctx, "GetActiveDonationsByCurrentStep")
	defer span.End()

	rows, err := utils.List[activeDonationStepCountRow](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`
		SELECT
			COALESCE(current_step.name::text, '') AS step,
			COUNT(*) AS count
		FROM donation d
		LEFT JOIN LATERAL (
			SELECT ds.name
			FROM donation_step ds
			WHERE ds.id_donation = d.id_donation
			ORDER BY ds.created_at DESC
			LIMIT 1
		) current_step ON true
		WHERE d.removed_at IS NULL
		  AND d.is_active = true
		  AND `+dateRangeFilter+`
		GROUP BY current_step.name
		`,
		start,
		end,
	)
	if err != nil || rows == nil {
		return nil, err
	}

	countByStep := make(map[string]int64, len(*rows))
	var total int64
	for _, row := range *rows {
		countByStep[row.Step] = row.Count
		total += row.Count
	}

	steps := []entities.EnumDonationSteps{
		entities.EnumDonationStepBloodTest,
		entities.EnumDonationStepDeliverMilkingKit,
		entities.EnumDonationStepCollectMilk,
		entities.EnumDonationStepMilkAnalysis,
	}

	result := make([]dto.ActiveDonationsByStep, 0, len(steps))
	for _, step := range steps {
		count := countByStep[string(step)]

		var percentage float64
		if total > 0 {
			percentage = math.Round(float64(count)*10000/float64(total)) / 100
		}

		result = append(result, dto.ActiveDonationsByStep{
			Step:       step,
			Count:      count,
			Percentage: percentage,
		})
	}

	return result, nil
}

func (r *DashboardRepository) GetDonorRecurrenceRate(ctx c.Context, start, end *time.Time) (float64, error) {
	ctx, span := r.StartSpan(ctx, "GetDonorRecurrenceRate")
	defer span.End()

	result, err := utils.Get[float64](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`
		WITH donor_counts AS (
			SELECT d.created_by AS id_user, COUNT(*) AS total
			FROM donation d
			JOIN "user" u ON u.id_user = d.created_by AND u.removed_at IS NULL
			WHERE d.removed_at IS NULL
			  AND `+dateRangeFilter+`
			GROUP BY d.created_by
		)
		SELECT
			(CASE WHEN COUNT(*) = 0 THEN 0
			ELSE ROUND((COUNT(*) FILTER (WHERE total > 1))::numeric * 100.0 / COUNT(*), 2)
			END)::float8
		FROM donor_counts
		`,
		start,
		end,
	)
	if err != nil || result == nil {
		return 0, err
	}

	return *result, nil
}
