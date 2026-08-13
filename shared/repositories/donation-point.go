package repositories

import (
	c "context"
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
	ctx c.Context,
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
		a.removed_at         AS "address.removed_at"
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

	rows, total, err := utils.ListQuery[donationPointRow](
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

	donationPoints := make([]dto.DonationPointsRes, len(*rows))
	for i, row := range *rows {
		var address *entities.AddressOut
		if row.Address != nil {
			out := entities.NewAddressOut(*row.Address)
			address = &out
		}
		donationPoints[i] = dto.DonationPointsRes{
			DonationPointOut: entities.NewDonationPointOut(row.DonationPoint),
			Address:          address,
			DistanceFromYou:  row.DistanceFromYou,
		}
	}

	return &donationPoints, total, nil
}

// donationPointRow is a scan-only shape used to map SQL columns (including
// the optional joined address and the computed distance) via sqlx. It never
// leaves this file - callers only see dto.DonationPointsRes.
type donationPointRow struct {
	entities.DonationPoint
	Address         *entities.Address `db:"address"`
	DistanceFromYou *float64          `db:"distance_from_you"`
}
