package repositories

import (
	"context"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	fluxgo.Repository[entities.User]
}

func UserRepositoryStart(db *fluxgo.Database) *UserRepository {
	return &UserRepository{*fluxgo.NewRepository[entities.User](db)}
}

func (r *UserRepository) GetUserById(ctx context.Context, id string) (*entities.User, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.User](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "user" WHERE id_user = $1 AND removed_at IS NULL`,
		id,
	)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.User](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "user" WHERE email = $1 AND removed_at IS NULL`,
		email,
	)
}

func (r *UserRepository) GetUserByCpf(ctx context.Context, cpf string) (*entities.User, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.User](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "user" WHERE cpf = $1 AND removed_at IS NULL`,
		cpf,
	)
}

func (r *UserRepository) ListUsersByFilters(
	ctx context.Context,
	filter *dto.ListUsersReq,
) (*[]entities.User, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select("u.*").
		From("user", "u").
		OrderBy(q.OrderBy{Column: "u.created_at"}).
		PaginationPaged(filter.Page, filter.PageSize).
		WhereAnd(q.Where{Column: "u.removed_at", Type: "IS NULL"})

	if filter.Name != nil {
		qb.WhereAnd(q.Where{
			Column: "u.name",
			Type:   "ILIKE",
			Val:    "%" + *filter.Name + "%",
		})
	}

	if filter.Type != nil {
		qb.WhereAnd(q.Where{
			Column: "u.type",
			Type:   "=",
			Val:    *filter.Type,
		})
	}

	if filter.InternalIdentifier != nil {
		qb.WhereAnd(q.Where{
			Column: "u.internal_identifier",
			Type:   "ILIKE",
			Val:    "%" + *filter.InternalIdentifier + "%",
		})
	}

	return utils.ListQuery[entities.User](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(filter.PageSize),
		true,
	)
}

type CreateUserRepositoryReq struct {
	IdUser             string
	InternalIdentifier *string
	Type               entities.EnumUserType
	Name               string
	Cpf                string
	BirthDate          time.Time
	PhoneNumber        string
	Email              string
	Password           string
	ActionBy           string
}

func (r *UserRepository) createUser(
	ctx context.Context,
	exec sqlx.ExtContext,
	data *CreateUserRepositoryReq,
) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		INSERT INTO "user" (
			id_user,
			internal_identifier,
			type,
			name,
			cpf,
			birth_date,
			phone_number,
			email,
			password,
			created_by,
			created_at
		) VALUES (
	        :id_user,
		 	:internal_identifier,
		 	:type,
			:name,
			:cpf,
			:birth_date,
			:phone_number,
			:email,
			:password,
			:action_by,
			now()
		)
	`

	params := map[string]any{
		"id_user":             data.IdUser,
		"internal_identifier": data.InternalIdentifier,
		"type":                data.Type,
		"name":                data.Name,
		"cpf":                 data.Cpf,
		"birth_date":          data.BirthDate,
		"phone_number":        data.PhoneNumber,
		"email":               data.Email,
		"password":            data.Password,
		"action_by":           data.ActionBy,
	}

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *UserRepository) CreateUserTx(
	ctx context.Context,
	tx *sqlx.Tx,
	data *CreateUserRepositoryReq,
) error {
	return r.createUser(
		ctx,
		tx,
		data,
	)
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	data *CreateUserRepositoryReq,
) error {
	return r.createUser(
		ctx,
		r.DB.WriteDB(),
		data,
	)
}

func (r *UserRepository) removeUser(ctx context.Context, exec sqlx.ExtContext, id, actionBy string) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		UPDATE "user"
		SET removed_at = now(),
		    removed_by = :action_by
		WHERE id_user = :id_user AND removed_at IS NULL
	`

	params := map[string]any{
		"id_user":   id,
		"action_by": actionBy,
	}

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *UserRepository) RemoveUserTx(
	ctx context.Context,
	tx *sqlx.Tx,
	id, actionBy string,
) error {
	return r.removeUser(
		ctx,
		tx,
		id,
		actionBy,
	)
}

func (r *UserRepository) RemoveUser(
	ctx context.Context,
	id, actionBy string,
) error {
	return r.removeUser(
		ctx,
		r.DB.WriteDB(),
		id,
		actionBy,
	)
}

type UpdateUserRepositoryReq struct {
	IdUser					string
	ActionBy				string
	InternalIdentifier		*string
	Type					*entities.EnumUserType
	Name					*string
	PhoneNumber				*string
	Email					*string
	Password				*string
}

func (r *UserRepository) UpdateUser(ctx context.Context, data *UpdateUserRepositoryReq) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		UPDATE "user"
		SET
			internal_identifier = COALESCE(:internal_identifier, internal_identifier),
			type                = COALESCE(:type, type),
			name                = COALESCE(:name, name),
			phone_number        = COALESCE(:phone_number, phone_number),
			email               = COALESCE(:email, email),
			password            = COALESCE(:password, password),
			updated_at          = now(),
			updated_by          = :action_by
		WHERE id_user = :id_user AND removed_at IS NULL
	`

	params := map[string]any{
		"id_user":             data.IdUser,
		"action_by":           data.ActionBy,
		"internal_identifier": data.InternalIdentifier,
		"type":                data.Type,
		"name":                data.Name,
		"phone_number":        data.PhoneNumber,
		"email":               data.Email,
		"password":            data.Password,
	}

	_, err := sqlx.NamedExecContext(ctx, r.DB.WriteDB(), query, params)
	if err != nil {
		span.SetError(err)
		return err
	}

	return nil
}