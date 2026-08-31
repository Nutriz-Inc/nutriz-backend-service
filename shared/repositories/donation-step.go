package repositories

import (
	c "context"
	dto "nutriz-backend-service/modules/donation/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"strings"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
	"github.com/jmoiron/sqlx"
)

type DonationStepRepository struct {
	fluxgo.Repository[entities.DonationStep]
}

func DonationStepRepositoryStart(db *fluxgo.Database) *DonationStepRepository {
	return &DonationStepRepository{*fluxgo.NewRepository[entities.DonationStep](db)}
}

func (r *DonationStepRepository) GetDonationStepsByIdDonation(
	ctx c.Context,
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
func (r *DonationStepRepository) GetDonationStepById(ctx c.Context, id string) (*entities.DonationStep, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.DonationStep](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "donation_step" WHERE id_donation_step = $1`,
		id,
	)
}

type CreateDonationStepRepositoryReq struct {
	IdDonationStep string
	IdDonation     string
	IdUser         string
	Name           entities.EnumDonationSteps
	Description    string
	Status         entities.EnumDonationStepStatus
	SetDate        *time.Time
	IdAddress      *string
}

func (r *DonationStepRepository) createDonationStep(
	ctx c.Context,
	exec sqlx.ExtContext,
	data *CreateDonationStepRepositoryReq,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		INSERT INTO donation_step (
			id_donation_step,
			id_donation,
			id_address,
			name,
			description,
			status,
			set_date,
			created_at,
			created_by
		) VALUES (
			:id_donation_step,
			:id_donation,
			:id_address,
			:name,
			:description,
			:status,
			:set_date,
			now(),
			:id_user
		)
	`

	params := map[string]any{
		"id_donation_step": data.IdDonationStep,
		"id_donation":      data.IdDonation,
		"id_user":          data.IdUser,
		"id_address":       data.IdAddress,
		"name":             data.Name,
		"description":      data.Description,
		"status":           data.Status,
		"set_date":         data.SetDate,
	}

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *DonationStepRepository) CreateDonationStepTx(
	ctx c.Context,
	tx *sqlx.Tx,
	data *CreateDonationStepRepositoryReq,
) error {
	return r.createDonationStep(
		ctx,
		tx,
		data,
	)
}

func (r *DonationStepRepository) CreateDonationStep(
	ctx c.Context,
	data *CreateDonationStepRepositoryReq,
) error {
	return r.createDonationStep(
		ctx,
		r.DB.WriteDB(),
		data,
	)
}

type UpdateDonationStepRepositoryReq struct {
	IdDonationStep string
	IdUser         string
	IdAddress      *string
	Description    string
	Status         *entities.EnumDonationStepStatus
	SetDate        *time.Time
	IsComplete     bool
}

func (r *DonationStepRepository) updateDonationStep(ctx c.Context, exec sqlx.ExtContext, data *UpdateDonationStepRepositoryReq) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	params := map[string]any{
		"id_donation_step": data.IdDonationStep,
		"updated_by":       data.IdUser,
		"description":      data.Description,
	}

	sets := []string{
		"description = :description",
	}

	if data.Status != nil {
		sets = append(sets, "status = :status")
		params["status"] = *data.Status
	}
	if data.SetDate != nil {
		sets = append(sets, "set_date = :set_date")
		params["set_date"] = *data.SetDate
	}
	if data.IsComplete {
		sets = append(sets, "completed_at = now()")
	}
	if data.IdAddress != nil {
		sets = append(sets, "id_address = :id_address")
		params["id_address"] = *data.IdAddress
	}

	query := `
		UPDATE donation_step
		SET ` + strings.Join(sets, ", ") + `,
		    updated_at = now(),
			updated_by = :updated_by
		WHERE id_donation_step = :id_donation_step
	`

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *DonationStepRepository) UpdateDonationStepTx(
	ctx c.Context,
	tx *sqlx.Tx,
	data *UpdateDonationStepRepositoryReq,
) error {
	return r.updateDonationStep(
		ctx,
		tx,
		data,
	)
}

func (r *DonationStepRepository) UpdateDonationStep(
	ctx c.Context,
	data *UpdateDonationStepRepositoryReq,
) error {
	return r.updateDonationStep(
		ctx,
		r.DB.WriteDB(),
		data,
	)
}

type DonationStepWithLocation struct {
	entities.DonationStep
	IsDonationActive bool     `db:"is_donation_active" json:"is_donation_active"`
	Latitude         *float64 `db:"latitude" json:"latitude"`
	Longitude        *float64 `db:"longitude" json:"longitude"`
	City             *string  `db:"city" json:"city"`
	Neighborhood     *string  `db:"neighborhood" json:"neighborhood"`
}

func (r *DonationStepRepository) GetDonationStepsWithLocationByIds(
	ctx c.Context,
	ids []string,
) (*[]DonationStepWithLocation, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query, args, err := sqlx.In(
		`SELECT
			ds.*,
			d.is_active AS is_donation_active,
			a.latitude AS latitude,
			a.longitude AS longitude,
			a.city AS city,
			a.neighborhood AS neighborhood
		 FROM donation_step ds
		 INNER JOIN donation d ON d.id_donation = ds.id_donation
		 LEFT JOIN address a ON a.id_address = ds.id_address AND a.removed_at IS NULL
		 WHERE ds.id_donation_step IN (?)`,
		ids,
	)
	if err != nil {
		return nil, err
	}

	query = r.DB.ReadOnlyDB().Rebind(query)

	return utils.List[DonationStepWithLocation](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		query,
		args...,
	)
}

type DonationStepWithAddress struct {
	entities.DonationStep
	AddrIdAddress       *string    `db:"addr_id_address"`
	AddrIdUser          *string    `db:"addr_id_user"`
	AddrIdDonationPoint *string    `db:"addr_id_donation_point"`
	AddrZipcode         *string    `db:"addr_zipcode"`
	AddrStreet          *string    `db:"addr_street"`
	AddrNumber          *string    `db:"addr_number"`
	AddrCity            *string    `db:"addr_city"`
	AddrState           *string    `db:"addr_state"`
	AddrNeighborhood    *string    `db:"addr_neighborhood"`
	AddrComplement      *string    `db:"addr_complement"`
	AddrLatitude        *float64   `db:"addr_latitude"`
	AddrLongitude       *float64   `db:"addr_longitude"`
	AddrCreatedAt       *time.Time `db:"addr_created_at"`
	AddrUpdatedAt       *time.Time `db:"addr_updated_at"`
}

func (s DonationStepWithAddress) Address() *entities.Address {
	if s.AddrIdAddress == nil {
		return nil
	}

	return &entities.Address{
		IdAddress:       *s.AddrIdAddress,
		IdUser:          s.AddrIdUser,
		IdDonationPoint: s.AddrIdDonationPoint,
		Zipcode:         utils.DerefString(s.AddrZipcode),
		Street:          utils.DerefString(s.AddrStreet),
		Number:          s.AddrNumber,
		City:            utils.DerefString(s.AddrCity),
		State:           utils.DerefString(s.AddrState),
		Neighborhood:    utils.DerefString(s.AddrNeighborhood),
		Complement:      s.AddrComplement,
		Latitude:        s.AddrLatitude,
		Longitude:       s.AddrLongitude,
		CreatedAt:       utils.DerefTime(s.AddrCreatedAt),
		UpdatedAt:       s.AddrUpdatedAt,
	}
}

func (r *DonationStepRepository) ListDonationStepsByFilters(
	ctx c.Context,
	filter *dto.ListDonationStepsReq,
) (*[]DonationStepWithAddress, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select(
			"ds.*",
			"a.id_address        AS addr_id_address",
			"a.id_user           AS addr_id_user",
			"a.id_donation_point AS addr_id_donation_point",
			"a.zipcode           AS addr_zipcode",
			"a.street            AS addr_street",
			"a.number            AS addr_number",
			"a.city              AS addr_city",
			"a.state             AS addr_state",
			"a.neighborhood      AS addr_neighborhood",
			"a.complement        AS addr_complement",
			"a.latitude          AS addr_latitude",
			"a.longitude         AS addr_longitude",
			"a.created_at        AS addr_created_at",
			"a.updated_at        AS addr_updated_at",
		).
		From("donation_step", "ds").
		Join(q.Join{
			Table: "address",
			As:    "a",
			On:    "a.id_address = ds.id_address AND a.removed_at IS NULL",
			Type:  q.LeftJoin,
		}).
		OrderBy(q.OrderBy{Column: "ds.created_at", Type: "DESC"}).
		PaginationPaged(filter.Page, filter.PageSize)

	if filter.Status != nil {
		qb.WhereAnd(q.Where{Column: "ds.status", Type: "=", Val: string(*filter.Status)})
	}
	if filter.IdDonation != nil {
		qb.WhereAnd(q.Where{Column: "ds.id_donation", Type: "=", Val: *filter.IdDonation})
	}
	if filter.Name != nil {
		qb.WhereAnd(q.Where{Column: "ds.name", Type: "=", Val: string(*filter.Name)})
	}
	if filter.SetDate != nil {
		qb.WhereAnd(q.Where{Column: "DATE(ds.set_date)", Type: "=", Val: *filter.SetDate})
	}
	if filter.Neighborhood != nil {
		qb.WhereAnd(q.Where{Column: "a.neighborhood", Type: "ILIKE", Val: "%" + *filter.Neighborhood + "%"})
	}
	if filter.City != nil {
		qb.WhereAnd(q.Where{Column: "a.city", Type: "ILIKE", Val: "%" + *filter.City + "%"})
	}

	return utils.ListQuery[DonationStepWithAddress](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(filter.PageSize),
		true,
	)
}
