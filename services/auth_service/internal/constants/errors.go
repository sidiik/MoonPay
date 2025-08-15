package constants

var (
	ErrInternalServer = "internal server error"
	ErrInvalidInput   = "invalid input"
	ErrUnauthorized   = "unauthorized"
	ErrForbidden      = "forbidden"
	ErrNotFound       = "resource not found"
	ErrConflict       = "conflict"

	// Auth-specific errors
	ErrEmailAlreadyExists = "email already exists"
	ErrInvalidCredentials = "invalid email or password"
	ErrTokenExpired       = "token has expired"
	ErrInvalidToken       = "invalid token"
)
