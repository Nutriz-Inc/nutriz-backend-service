package repositories

import (
	"context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"

	fluxgo "github.com/MMortari/FluxGo"
)

type JobRepository struct {
	fluxgo.Repository[entities.Job]
}

func JobRepositoryStart(db *fluxgo.Database) *JobRepository {
	return &JobRepository{*fluxgo.NewRepository[entities.Job](db)}
}

func (r *JobRepository) GetJobById(ctx context.Context, id string) (*entities.Job, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	return utils.Get[entities.Job](
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		`SELECT * FROM "job" WHERE id_job = $1 AND removed_at IS NULL`,
		id,
	)
}
