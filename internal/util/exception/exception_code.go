package exception

import "errors"

var (
	ErrEmptyResult        = errors.New("empty result")
	ErrTokenInvalid       = errors.New("token invalid")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidRequest     = errors.New("invalid request")
	ErrFailedReadEnum     = errors.New("failed to read enum")
	ErrUnableToLock       = errors.New("unable to lock")
	ErrMethodNotAllowed   = errors.New("method not allowed")
	ErrUnauthorizedMethod = errors.New("unauthorized method")

	// spesific mysql error
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrDupeKey             = errors.New("duplicate key value")
	ErrUniqueViolation     = errors.New("unique violation")
)

const (
	CodeDataNotFound         = "DATA_NOT_FOUND"
	CodeAccountLocked        = "ACCOUNT_LOCKED"
	CodeOtpFailed            = "OTP_FAILED"
	CodeOtpInvalid           = "OTP_INVALID"
	CodeRequestTooFast       = "REQUEST_TOO_FAST"
	CodeInvalidCredential    = "INVALID_CREDENTIAL"
	CodeInvalidRequest       = "INVALID_REQUEST"
	CodeEmailExist           = "EMAIL_EXIST"
	CodeInvalidValidation    = "INVALID_VALIDATION"
	CodeInvalidAge           = "INVALID_AGE"
	CodeForbidden            = "FORBIDDEN"
	CodeDataLocked           = "DATA_LOCKED"
	CodeBadRequest           = "BAD_REQUEST"
	CodeDateExpired          = "DATE_EXPIRED"
	CodeMissingRequiredData  = "MISSING_REQUIRED_DATA"
	CodeAccountBlocked       = "ACCOUNT_BLOCKED"
	CodeAccountDeleted       = "ACCOUNT_DELETED"
	CodeEmailNotVerified     = "EMAIL_NOT_VERIFIED"
	CodeInvalidData          = "INVALID_DATA"
	CodeConflict             = "CONFLICT"
	CodeQuotaLimitReached    = "QUOTA_LIMIT_REACHED"
	CodeServiceUnavailable   = "SERVICE_UNAVAILABLE"
	CodeRequestFailed        = "REQUEST_FAILED"
	CodeDataAlreadyExist     = "DATA_ALREADY_EXIST"
	CodeInternalServerError  = "INTERNAL_SERVER_ERROR"
	CodeFileNotProcessedYet  = "FILE_NOT_PROCESSED_YET"
	CodeFileTypeNotSupported = "FILE_TYPE_NOT_SUPPORTED"
	CodeTokenExpired         = "TOKEN_EXPIRED"
	CodeTokenInvalid         = "TOKEN_INVALID"
	CodeInvalidFormat        = "INVALID_FORMAT"
	CodeUnauthorized         = "UNAUTHORIZED"
	CodeMethodNotAllowed     = "METHOD_NOT_ALLOWED"
)

type (
	errorType struct {
		ErrorMessage string
	}
	ErrorNotFound            errorType
	ErrorOTPFailed           errorType
	ErrorForeignKeyViolation errorType
)

type NotFoundError struct {
	Error string
}

func NewNotFoundError(error string) NotFoundError {
	return NotFoundError{Error: error}
}
