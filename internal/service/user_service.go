package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-playground/validator"
	"github.com/shennawardana23/graphql-pba/graph/model"
	"github.com/shennawardana23/graphql-pba/internal/app/database"
	"github.com/shennawardana23/graphql-pba/internal/entity"
	"github.com/shennawardana23/graphql-pba/internal/repository"
	"github.com/shennawardana23/graphql-pba/internal/util/logger"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	UserRepo repository.UserRepository
	DB       database.DB
	Validate *validator.Validate
}

func NewUserService(userRepo repository.UserRepository, db database.DB, validator *validator.Validate) *UserService {
	return &UserService{
		UserRepo: userRepo,
		DB:       db,
		Validate: validator,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req model.NewUser) (entity.User, error) {
	user, err := s.UserRepo.FindByEmail(ctx, s.DB, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.WithFields(logrus.Fields{
				"error": err,
				"file":  "service.go",
			}).Error("Error checking for existing user")
			return entity.User{}, fmt.Errorf("failed to check existing user: %w", err)
		}
	}
	if user.Email == req.Email {
		return entity.User{}, fmt.Errorf("user already exists")
	}

	createdUser, err := s.UserRepo.Create(ctx, s.DB, user.Name, user.Email)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

func (s *UserService) GetUser(ctx context.Context, id int64) (entity.User, error) {
	user, err := s.UserRepo.FindByID(ctx, s.DB, id)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *UserService) GetUsers(ctx context.Context) ([]entity.User, error) {
	users, err := s.UserRepo.FindAll(ctx, s.DB)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error":   err,
			"context": "GetUsers",
			"file":    "service.go",
		}).Error("Failed to retrieve users from repository")
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req model.UpdateUserInput) (entity.User, error) {
	// Check if user exists
	user, err := s.UserRepo.FindByID(ctx, s.DB, int64(req.ID))
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to check existing user: %w", err)
	}

	// Update user
	updatedUser, err := s.UserRepo.Update(ctx, s.DB, user.ID, *req.Name, *req.Email)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	err := s.UserRepo.Delete(ctx, s.DB, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
