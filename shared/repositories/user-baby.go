package repositories

import (
	"context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"

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
		WhereAnd(q.Where{Column: "ub.removed_at", Type: "IS NULL"}).
		WhereAnd(q.Where{Column: "ub.id_user_baby", Type: "=", Val: userBabyId})

	return utils.ListQuery[entities.UserBaby](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		qb,
		utils.IntPtr(entities.MAX_BABY_QUANTITY_PER_USER),
		false,
	)
}
