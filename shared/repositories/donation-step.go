package repositories

import (
	"context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
)

type DonationStepRepository struct {
	fluxgo.Repository[entities.DonationStep]
}

func DonationStepRepositoryStart(db *fluxgo.Database) *DonationStepRepository {
	return &DonationStepRepository{*fluxgo.NewRepository[entities.DonationStep](db)}
}

func (r *DonationStepRepository) GetDonationStepsByIdDonation(
	ctx context.Context,
	idDonation string,
) (*[]entities.DonationStep, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select("ds.*").
		From("donation_step", "ds").
		OrderBy(q.OrderBy{Column: "ds.created_at"}).
		PaginationPaged(1, entities.NUMBER_OF_DONATION_STEPS).
		WhereAnd(q.Where{Column: "ds.id_donation", Type: "=", Val: idDonation})

	return utils.ListQuery[entities.DonationStep](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(entities.NUMBER_OF_DONATION_STEPS),
		false,
	)
}
