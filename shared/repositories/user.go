package repositories

import (
	"context"
	"database/sql"
	"nutriz-backend-service/shared/entities"

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

	var user entities.User

	err := r.DB.ReadOnlyDB().GetContext(
		ctx,
		&user,
		`SELECT * FROM "user" WHERE id_user = $1`,
		id,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		span.SetError(err)
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	ctx, span := r.StartSpan(ctx)
	defer span.End()

	var user entities.User

	err := r.DB.ReadOnlyDB().GetContext(
		ctx,
		&user,
		`SELECT * FROM "user" WHERE email = $1`,
		email,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		span.SetError(err)
		return nil, err
	}

	return &user, nil
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

// 	_, err := r.DB.ReadOnlyDB().NamedExecContext(ctx, query, params)
// 	if err != nil {
// 		span.SetError(err)
// 		return err
// 	}

// 	return nil
// }
