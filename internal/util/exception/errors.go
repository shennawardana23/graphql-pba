package exception

type CustomError struct {
	Code    string
	Message string
	Details string
}

func (e *CustomError) Error() string {
	return e.Message
}

var (
	ErrDuplicateEmail = &CustomError{
		Code:    "USER_EMAIL_EXISTS",
		Message: "Email address is already in use",
		Details: "Please use a different email address",
	}

	ErrInvalidInput = &CustomError{
		Code:    "INVALID_INPUT",
		Message: "Invalid input provided",
		Details: "Please check your input and try again",
	}

	ErrNotFound = &CustomError{
		Code:    "NOT_FOUND",
		Message: "Resource not found",
		Details: "The requested resource does not exist",
	}

	ErrInternalServer = &CustomError{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Internal server error",
		Details: "An unexpected error occurred",
	}
)

func NewCustomError(code, message, details string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewValidationError(details string) *CustomError {
	return &CustomError{
		Code:    "VALIDATION_ERROR",
		Message: "Validation failed",
		Details: details,
	}
}
