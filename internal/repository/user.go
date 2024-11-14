package repository

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/shennawardana23/graphql-pba/internal/entity"
	"github.com/shennawardana23/graphql-pba/internal/util/exception"
)

type UserRepository struct {
	db *pg.DB
}

func NewUserRepository(db *pg.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) WithContext(ctx context.Context) *pg.DB {
	return r.db.WithContext(ctx)
}

func (r *UserRepository) FindAll(ctx context.Context) ([]entity.User, error) {
	var users []entity.User
	err := r.WithContext(ctx).ModelContext(ctx, &users).Select()
	return users, exception.TranslatePostgresError(ctx, err)
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	user := &entity.User{ID: id}
	err := r.WithContext(ctx).Model(user).WherePK().Select()
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return user, exception.TranslatePostgresError(ctx, err)
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.WithContext(ctx).ModelContext(ctx, user).Insert()
	return exception.TranslatePostgresError(ctx, err)
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()

	_, err := r.WithContext(ctx).ModelContext(ctx, user).WherePK().Update()
	if err == pg.ErrNoRows {
		return exception.ErrEmptyResult
	}
	return exception.TranslatePostgresError(ctx, err)
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	user := &entity.User{ID: id}
	_, err := r.WithContext(ctx).ModelContext(ctx, user).WherePK().Delete()
	if err == pg.ErrNoRows {
		return exception.ErrEmptyResult
	}
	return exception.TranslatePostgresError(ctx, err)
}

// Additional helper methods for specific error cases
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := new(entity.User)
	err := r.WithContext(ctx).
		Model(user).
		Where("email = ?", email).
		Select()

	if err == pg.ErrNoRows {
		return nil, nil
	}
	return user, exception.TranslatePostgresError(ctx, err)
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.WithContext(ctx).
		Model((*entity.User)(nil)).
		Where("email = ?", email).
		Exists()

	return exists, exception.TranslatePostgresError(ctx, err)
}

// Transaction support
func (r *UserRepository) WithTransaction(ctx context.Context, fn func(*pg.Tx) error) error {
	return r.db.WithContext(ctx).RunInTransaction(ctx, func(tx *pg.Tx) error {
		err := fn(tx)
		return exception.TranslatePostgresError(ctx, err)
	})
}

// Batch operations
func (r *UserRepository) CreateBatch(ctx context.Context, users []*entity.User) error {
	if len(users) == 0 {
		return nil
	}

	now := time.Now()
	for _, user := range users {
		user.CreatedAt = now
		user.UpdatedAt = now
	}

	_, err := r.WithContext(ctx).
		ModelContext(ctx, &users).
		Insert()

	return exception.TranslatePostgresError(ctx, err)
}
