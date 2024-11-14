package validation_model

import (
	"github.com/shennawardana23/graphql-pba/graph/model"
)

type ValidationNewUser struct {
	Name  string `validate:"required,min=2,max=100"`
	Email string `validate:"required,email,max=255"`
}

func NewUserFromGQL(input model.NewUser) ValidationNewUser {
	return ValidationNewUser{
		Name:  input.Name,
		Email: input.Email,
	}
}
