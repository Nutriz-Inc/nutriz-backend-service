package repositories

import (
	"context"
	"database/sql"
	"nutriz-backend-service/shared/entities"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
	"go.opentelemetry.io/otel/codes"
)

type UserBabyRepository struct {
	fluxgo.Repository[entities.UserBaby]
}

func UserBabyRepositoryStart(db *fluxgo.Database) *UserBabyRepository {
	return &UserBabyRepository{*fluxgo.NewRepository[entities.UserBaby](db)}
}

func (r *UserBabyRepository) GetUserBabybyUserBabyId(
	ctx context.Context,
	userBabyId string,
) (*entities.UserBaby, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select("ub.*").
		From("user_baby", "ub").
		WhereAnd(q.Where{Column: "ub.removed_at", Type: "IS NULL"}).
		WhereAnd(q.Where{Column: "ub.id_user_baby", Type: "=", Val: userBabyId})

	query, args := qb.ToSelectSql()

	var resp entities.UserBaby
	err := r.DB.ReadOnlyDB().GetContext(ctx, &resp, query, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	
	return &resp, nil
}

func (r *UserBabyRepository) GetUserBabiesByUserId(
	ctx context.Context,
	userId string,
) (*[]entities.UserBaby, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	qb := q.NewQueryBuilder(q.SetOtelSpan(span)).
		Select("ub.*").
		From("user_baby", "ub").
		PaginationPaged(1, entities.MAX_BABY_QUANTITY_PER_USER).
		OrderBy(q.OrderBy{Column: "ub.created_at"}).
		WhereAnd(q.Where{Column: "ub.removed_at", Type: "IS NULL"}).
		WhereAnd(q.Where{Column: "ub.id_user", Type: "=", Val: userId})

	query, args := qb.ToSelectSql()
	resp := make([]entities.UserBaby, 0)

	err := r.DB.ReadOnlyDB().SelectContext(ctx, &resp, query, args...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	queryTotal, argsTotal := qb.ToSelectTotalSql()

	var total int
	err = r.DB.ReadOnlyDB().GetContext(ctx, &total, queryTotal, argsTotal...)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	return &resp, nil
}

type CreateUserBabyRepositoryReq struct {
	IdUserBaby		string
	IdUser			string
	Name 			*string
	BirthDate		time.Time
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
		"id_user_baby":	data.IdUserBaby,
		"id_user":		data.IdUser,
		"name":			data.Name,
		"birth_date":	data.BirthDate,
	}

	_, err := r.DB.ReadOnlyDB().NamedExecContext(ctx, query, params)
	if err != nil {
		span.SetError(err)
		return nil
	}

	return nil
}