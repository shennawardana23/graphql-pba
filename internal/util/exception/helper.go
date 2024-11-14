package exception

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
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
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrEmptyResult
		} else if myErr, isMyErr := err.(*pq.Error); isMyErr {
			if myErr.Code.Name() == "23503" {
				err = ErrForeignKeyViolation
			} else if myErr.Code.Name() == "23505" {
				err = ErrUniqueViolation
			} else if myErr.Code.Name() == "42710" {
				err = ErrDupeKey
			} else {
				logger.Error(ctx, err)
			}
		} else {
			logger.Error(ctx, err)
		}
	}
	return err
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
