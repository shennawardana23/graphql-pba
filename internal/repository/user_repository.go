package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shennawardana23/graphql-pba/internal/app/database"
	"github.com/shennawardana23/graphql-pba/internal/entity"
	"github.com/shennawardana23/graphql-pba/internal/util/exception"
)

type UserRepository interface {
	FindAll(ctx context.Context, q database.Q) ([]entity.User, error)
	FindByID(ctx context.Context, q database.Q, id int64) (entity.User, error)
	FindByEmail(ctx context.Context, q database.Q, email string) (entity.User, error)
	Create(ctx context.Context, q database.Q, name, email string) (entity.User, error)
	Update(ctx context.Context, q database.Q, id int64, name, email string) (entity.User, error)
	Delete(ctx context.Context, q database.Q, id int64) error
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) FindAll(ctx context.Context, q database.Q) ([]entity.User, error) {
	query := `SELECT id, name, email, created_at, updated_at FROM users`

	rows, err := q.QueryContext(ctx, query)
	if err != nil {
		return nil, exception.TranslatePostgresError(ctx, err)
	}
	defer func() {
		err = rows.Close()
		exception.PanicOnErrorContext(ctx, err)
	}()

	var users []entity.User
	for rows.Next() {
		user := entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, exception.TranslatePostgresError(ctx, err)
		}
		users = append(users, user)
	}

	fmt.Println("In Repository ================")
	fmt.Println(users, "================")

	return users, nil
}

func (r *userRepository) FindByID(ctx context.Context, q database.Q, id int64) (entity.User, error) {
	user := entity.User{}
	query := `SELECT id, name, email FROM users WHERE id = $1`

	err := q.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)

	return user, exception.TranslatePostgresError(ctx, err)
}

func (r *userRepository) FindByEmail(ctx context.Context, q database.Q, email string) (entity.User, error) {
	query := `
		SELECT id, name, email
		FROM users
		WHERE email = $1
	`

	user := entity.User{}
	err := q.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email)

	return user, exception.TranslatePostgresError(ctx, err)
}

func (r *userRepository) Create(ctx context.Context, q database.Q, name, email string) (entity.User, error) {
	query := `
		INSERT INTO users (name, email, created_at, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, name, email
	`

	var user entity.User
	err := q.QueryRowContext(ctx, query, name, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	return user, exception.TranslatePostgresError(ctx, err)
}

func (r *userRepository) Update(ctx context.Context, q database.Q, id int64, name, email string) (entity.User, error) {
	query := `
		UPDATE users
		SET name = $1, email = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, name, email`

	var filterValues []interface{}

	filterValues = append(filterValues, name, email, id)

	result := entity.User{}
	err := q.QueryRowContext(ctx, query, filterValues...).Scan(
		&result.ID,
		&result.Name,
		&result.Email,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return result, exception.TranslatePostgresError(ctx, err)
}

func (r *userRepository) Delete(ctx context.Context, q database.Q, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := q.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %w", sql.ErrNoRows)
	}

	return nil
}
