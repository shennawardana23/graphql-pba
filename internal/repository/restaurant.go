package repository

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/shennawardana23/graphql-pba/internal/entity"
	"github.com/shennawardana23/graphql-pba/internal/util/exception"
)

type RestaurantRepository struct {
	db *pg.DB
}

func NewRestaurantRepository(db *pg.DB) *RestaurantRepository {
	return &RestaurantRepository{db: db}
}

func (r *RestaurantRepository) WithContext(ctx context.Context) *pg.DB {
	return r.db.WithContext(ctx)
}

func (r *RestaurantRepository) FindAll(ctx context.Context) ([]entity.Restaurant, error) {
	var restaurants []entity.Restaurant
	err := r.WithContext(ctx).ModelContext(ctx, &restaurants).Select()
	return restaurants, exception.TranslatePostgresError(ctx, err)
}

func (r *RestaurantRepository) FindByID(ctx context.Context, id int64) (*entity.Restaurant, error) {
	restaurant := &entity.Restaurant{ID: id}
	err := r.WithContext(ctx).
		Model(restaurant).
		Relation("User").
		WherePK().
		Select()
	if err == pg.ErrNoRows {
		return nil, nil
	}
	return restaurant, exception.TranslatePostgresError(ctx, err)
}

func (r *RestaurantRepository) Create(ctx context.Context, restaurant *entity.Restaurant) error {
	restaurant.CreatedAt = time.Now()
	restaurant.UpdatedAt = time.Now()

	_, err := r.WithContext(ctx).ModelContext(ctx, restaurant).Insert()
	return exception.TranslatePostgresError(ctx, err)
}

func (r *RestaurantRepository) Update(ctx context.Context, restaurant *entity.Restaurant) error {
	restaurant.UpdatedAt = time.Now()

	_, err := r.WithContext(ctx).ModelContext(ctx, restaurant).WherePK().Update()
	if err == pg.ErrNoRows {
		return exception.ErrEmptyResult
	}
	return exception.TranslatePostgresError(ctx, err)
}

func (r *RestaurantRepository) Delete(ctx context.Context, id int64) error {
	restaurant := &entity.Restaurant{ID: id}
	_, err := r.WithContext(ctx).ModelContext(ctx, restaurant).WherePK().Delete()
	if err == pg.ErrNoRows {
		return exception.ErrEmptyResult
	}
	return exception.TranslatePostgresError(ctx, err)
}
