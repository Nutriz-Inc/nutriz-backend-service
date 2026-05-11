package repositories

import (
	"context"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"strconv"

	q "github.com/MMortari/go-query-builder"

	fluxgo "github.com/MMortari/FluxGo"
)

type DonationPointRepository struct {
	fluxgo.Repository[entities.DonationPoint]
}

func DonationPointRepositoryStart(db *fluxgo.Database) *DonationPointRepository {
	return &DonationPointRepository{*fluxgo.NewRepository[entities.DonationPoint](db)}
}

func (r *DonationPointRepository) ListDonationPointsByFilters(
	ctx context.Context,
	filter *dto.ListDonationPointsReq,
) (*[]dto.DonationPointsRes, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select("dp.*").
		From("donation_point", "dp").
		PaginationPaged(filter.Page, filter.PageSize).
		WhereAnd(q.Where{
			Column: "dp.removed_at",
			Type:   "IS NULL",
		})

	if filter.Name != nil {
		qb.WhereAnd(q.Where{
			Column: "dp.name",
			Type:   "ILIKE",
			Val:    "%" + *filter.Name + "%",
		})
	}
	if filter.HasHome != nil {
		qb.WhereAnd(q.Where{
			Column: "dp.has_home",
			Type:   "=",
			Val:    *filter.HasHome,
		})
	}

	if filter.ShowAddress {
		qb.Join(q.Join{
			Table: "address",
			As:    "a",
			On:    "a.id_donation_point = dp.id_donation_point AND a.removed_at IS NULL",
			Type:  q.LeftJoin,
		})

		qb.Select(`
		a.id_address         AS "address.id_address",
		a.id_user            AS "address.id_user",
		a.id_donation_point  AS "address.id_donation_point",
		a.zipcode            AS "address.zipcode",
		a.street             AS "address.street",
		a.number             AS "address.number",
		a.city               AS "address.city",
		a.state              AS "address.state",
		a.neighborhood       AS "address.neighborhood",
		a.complement         AS "address.complement",
		a.latitude           AS "address.latitude",
		a.longitude          AS "address.longitude",
		a.created_at         AS "address.created_at",
		a.updated_at         AS "address.updated_at",
		a.updated_by         AS "address.updated_by",
		a.removed_at         AS "address.removed_at",
		a.removed_by         AS "address.removed_by"
	`)
	}

	if filter.Longitude != nil && filter.Latitude != nil {
		lat := strconv.FormatFloat(*filter.Latitude, 'f', 6, 64)
		lng := strconv.FormatFloat(*filter.Longitude, 'f', 6, 64)

		if !filter.ShowAddress {
			qb.Join(q.Join{
				Table: "address",
				As:    "a",
				On:    "a.id_donation_point = dp.id_donation_point AND a.removed_at IS NULL",
				Type:  q.LeftJoin,
			})
		}

		qb.Select(`
			(
				6371 * acos(
					cos(radians(` + lat + `)) *
					cos(radians(a.latitude)) *
					cos(radians(a.longitude) - radians(` + lng + `)) +
					sin(radians(` + lat + `)) *
					sin(radians(a.latitude))
				)
			) AS distance_from_you
		`)

		qb.OrderBy(q.OrderBy{
			Column: "distance_from_you",
			Type:   "ASC",
		})
	}

	return utils.ListQuery[dto.DonationPointsRes](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(filter.PageSize),
		true,
	)
}
