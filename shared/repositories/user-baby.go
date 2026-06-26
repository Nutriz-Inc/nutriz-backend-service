package repositories

import (
	c "context"
	"nutriz-backend-service/shared/entities"
	"nutriz-backend-service/shared/utils"
	"strings"
	"time"

	fluxgo "github.com/MMortari/FluxGo"
	q "github.com/MMortari/go-query-builder"
	"github.com/jmoiron/sqlx"
)

type UserBabyRepository struct {
	fluxgo.Repository[entities.UserBaby]
}

func UserBabyRepositoryStart(db *fluxgo.Database) *UserBabyRepository {
	return &UserBabyRepository{*fluxgo.NewRepository[entities.UserBaby](db)}
}

func (r *UserBabyRepository) GetUserBabiesByUserId(
	ctx c.Context,
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
		true,
	)
}

func (r *UserBabyRepository) GetUserBabyById(
	ctx c.Context,
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

func (r *UserBabyRepository) createUserBaby(ctx c.Context, exec sqlx.ExtContext, data *CreateUserBabyRepositoryReq) error {
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

	_, err := sqlx.NamedExecContext(
		ctx,
		exec,
		query,
		params,
	)

	return err
}

func (r *UserBabyRepository) CreateUserBabyTx(
	ctx c.Context,
	tx *sqlx.Tx,
	data *CreateUserBabyRepositoryReq,
) error {
	return r.createUserBaby(
		ctx,
		tx,
		data,
	)
}

func (r *UserBabyRepository) CreateUserBaby(
	ctx c.Context,
	data *CreateUserBabyRepositoryReq,
) error {
	return r.createUserBaby(
		ctx,
		r.DB.WriteDB(),
		data,
	)
}

func (r *UserBabyRepository) RemoveUserBaby(ctx c.Context, id, actionBy string) error {
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

type UpdateUserBabyRepositoryReq struct {
	IdUserBaby string
	IdUser     string
	Name       *string
	BirthDate  *time.Time
}

func (r *UserBabyRepository) UpdateUserBaby(ctx c.Context, data UpdateUserBabyRepositoryReq) error {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	sets := []string{}
	params := map[string]any{
		"id_user_baby": data.IdUserBaby,
		"id_user":      data.IdUser,
	}

	if data.Name != nil {
		sets = append(sets, "name = :name")
		params["name"] = data.Name
	}

	if data.BirthDate != nil {
		sets = append(sets, "birth_date = :birth_date")
		params["birth_date"] = data.BirthDate
	}

	if len(sets) == 0 {
		return nil
	}

	query := `
        UPDATE user_baby
        SET ` + strings.Join(sets, ", ") + `,
            updated_at = now()
        WHERE id_user_baby = :id_user_baby
          AND removed_at IS NULL
          AND id_user = :id_user
	`

	return utils.Update(
		ctx,
		r.DB.ReadOnlyDB(),
		span,
		query,
		params,
	)
}
