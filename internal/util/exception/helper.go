package exception

import (
	"context"
	"database/sql"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/shennawardana23/graphql-pba/internal/util/logger"
)

func PanicOnError(err error) {
	if err != nil {
		logger.Log.Error(err)
		panic(err)
	}
}

func PanicOnErrorContext(ctx context.Context, err error) {
	if err != nil {
		logger.Error(ctx, err)
		panic(err)
	}
}

func TranslatePostgresError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	var customErr *CustomError

	switch {
	case err == sql.ErrNoRows || err == pg.ErrNoRows:
		customErr = ErrNotFound

	case isPgError(err, "23505"): // unique_violation
		if strings.Contains(err.Error(), "idx_users_email") {
			customErr = ErrDuplicateEmail
		} else {
			customErr = NewCustomError(
				"DUPLICATE_ENTRY",
				"Duplicate entry found",
				"A record with this value already exists",
			)
		}

	case isPgError(err, "23503"): // foreign_key_violation
		customErr = NewCustomError(
			"FOREIGN_KEY_VIOLATION",
			"Invalid reference",
			"The referenced record does not exist",
		)

	case isPgError(err, "23502"): // not_null_violation
		customErr = NewCustomError(
			"REQUIRED_FIELD",
			"Required field missing",
			"Please provide all required fields",
		)

	default:
		logger.Error(ctx, err)
		customErr = ErrInternalServer
	}

	return customErr
}

// isPgError checks if the error is a postgres error with the given code
func isPgError(err error, code string) bool {
	pgErr, ok := err.(pg.Error)
	return ok && pgErr.Field('C') == code
}

func CancelBackground(ctx context.Context, cancel context.CancelFunc, errorMessage string, successMessage string) {
	select {
	case <-ctx.Done():
		if len(errorMessage) > 0 {
			logger.Error(ctx, errorMessage)
		}
		cancel()
		return
	default:
		if len(successMessage) > 0 {
			logger.Info(ctx, successMessage)
		}
		return
	}
}
