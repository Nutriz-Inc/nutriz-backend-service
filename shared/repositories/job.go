package repositories

import (
	"context"
	"database/sql"
	"nutriz-backend-service/shared/entities"

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

	var job entities.Job

	err := r.DB.ReadOnlyDB().GetContext(ctx, &job, `SELECT * FROM "job" WHERE id_job = $1`, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		span.SetError(err)
		return nil, err
	}

	return &job, nil
}
