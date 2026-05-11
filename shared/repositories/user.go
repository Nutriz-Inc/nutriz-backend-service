package repositories

import (
	"context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
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
