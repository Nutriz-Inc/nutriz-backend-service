package repositories

import (
	"context"
	dto "nutriz-backend-service/modules/user/dtos"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
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
		`SELECT * FROM "user" WHERE id_user = $1`,
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
		`SELECT * FROM "user" WHERE email = $1`,
		email,
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

// func (r *UserRepository) CreateUser(ctx context.Context, data dto.CreateUserReq, idUser string) error {
// 	ctx, span := r.StartSpan(ctx)
// 	defer span.End()

// 	query := `
// 		INSERT INTO "user" (
// 			id_user,
// 			name,
// 			email,
// 			password,
// 			closing_date,
// 			created_at
// 		) VALUES (
// 	        :id_user,
// 			:name,
// 			:email,
// 			:password,
// 			:closing_date,
// 			now()
// 		)
// 	`
// 	params := map[string]any{
// 		"id_user":      idUser,
// 		"name":         data.Name,
// 		"email":        data.Email,
// 		"password":     data.Password,
// 		"closing_date": data.ClosingDate,
// 	}

// 	return utils.Insert(
// 		ctx,
// 		r.DB.ReadOnlyDB(),
// 		span,
// 		query,
// 		params,
// 	)
// }
