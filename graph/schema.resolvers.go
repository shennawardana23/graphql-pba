package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/shennawardana23/graphql-pba/graph/generated"
	"github.com/shennawardana23/graphql-pba/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	user, err := r.UserService.CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:    int(user.ID),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUserInput) (*model.User, error) {
	user, err := r.UserService.UpdateUser(ctx, input)
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:    int(user.ID),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id int) (*model.User, error) {
	err := r.UserService.DeleteUser(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	return &model.User{}, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	users, err := r.UserService.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	var usersModel []*model.User
	for _, user := range users {
		usersModel = append(usersModel, &model.User{
			ID:    int(user.ID),
			Name:  user.Name,
			Email: user.Email,
		})
	}
	return usersModel, nil
}

func (r *queryResolver) User(ctx context.Context, id int) (*model.User, error) {
	user, err := r.UserService.GetUser(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	return &model.User{
		ID:    int(user.ID),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
