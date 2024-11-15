package graph

import (
	"github.com/go-pg/pg/v10"
	"github.com/shennawardana23/graphql-pba/internal/repository"
)

type Resolver struct {
	DB                   *pg.DB
	UserRepository       *repository.UserRepository
	RestaurantRepository *repository.RestaurantRepository
}

func NewResolver(db *pg.DB) *Resolver {
	return &Resolver{
		DB:                   db,
		UserRepository:       repository.NewUserRepository(db),
		RestaurantRepository: repository.NewRestaurantRepository(db),
	}
}
