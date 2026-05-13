package repositories

import (
	"context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
)

type UserBabyRepository struct {
	fluxgo.Repository[entities.UserBaby]
}

func UserBabyRepositoryStart(db *fluxgo.Database) *UserBabyRepository {
	return &UserBabyRepository{*fluxgo.NewRepository[entities.UserBaby](db)}
}

func (r *UserBabyRepository) GetUserBabyesByUserId(
	ctx context.Context,
	userId string,
) (*[]entities.UserBaby, int, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select("ub.*").
		From("user_baby", "ub").
		PaginationPaged(1, entities.MAX_BABY_QUANTITY_PER_USER).
		OrderBy(q.OrderBy{Column: "ub.created_at"}).
		WhereAnd(q.Where{Column: "ub.removed_at", Type: "IS NULL"}).
		WhereAnd(q.Where{Column: "ub.id_user", Type: "=", Val: userId})

	return utils.ListQuery[entities.UserBaby](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(entities.MAX_BABY_QUANTITY_PER_USER),
		false,
	)
}

func (r *UserBabyRepository) GetUserBabyById(
	ctx context.Context,
	userBabyId string,
) (*entities.UserBaby, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.UserBaby](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "user_baby" WHERE id_user_baby = $1 AND removed_at IS NULL`,
		userBabyId,
	)
}

type CreateUserBabyRepositoryReq struct {
	IdUserBaby string
	IdUser     string
	Name       *string
	BirthDate  time.Time
}

func (r *UserBabyRepository) CreateUserBaby(ctx context.Context, data *CreateUserBabyRepositoryReq) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		INSERT INTO user_baby (
			id_user_baby,
			id_user,
			name,
			birth_date,
			created_at
		) VALUES (
			:id_user_baby,
			:id_user,
			:name,
			:birth_date,
			now() 
		)
	`

	params := map[string]any{
		"id_user_baby": data.IdUserBaby,
		"id_user":      data.IdUser,
		"name":         data.Name,
		"birth_date":   data.BirthDate,
	}

	return utils.Insert(
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		query,
		params,
	)
}

func (r *UserBabyRepository) RemoveUserBaby(ctx context.Context, id, actionBy string) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	query := `
		UPDATE user_baby
		SET removed_at = now()
		WHERE id_user_baby = $1 AND id_user = $2 AND removed_at IS NULL 
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