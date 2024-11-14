package validator

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/go-playground/validator/v10"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func translateValidationError(ctx context.Context, err error) *gqlerror.Error {
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return &gqlerror.Error{
			Message: "Validation failed",
			Path:    graphql.GetPath(ctx),
		}
	}

	var messages []string
	for _, e := range validationErrors {
		messages = append(messages, translateFieldError(e))
	}

	return &gqlerror.Error{
		Message: "Validation failed",
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code":    "VALIDATION_ERROR",
			"details": strings.Join(messages, "; "),
		},
	}
}
