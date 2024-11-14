package graph

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/shennawardana23/graphql-pba/internal/util/exception"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func ErrorPresenter(ctx context.Context, err error) *gqlerror.Error {
	customErr, ok := err.(*exception.CustomError)
	if !ok {
		return &gqlerror.Error{
			Message: "Validation failed",
			Path:    graphql.GetPath(ctx),
			Extensions: map[string]interface{}{
				"code":    "VALIDATION_ERROR",
				"details": err.Error(),
			},
		}
	}

	return &gqlerror.Error{
		Message: customErr.Message,
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code":    customErr.Code,
			"details": customErr.Details,
		},
	}
}
