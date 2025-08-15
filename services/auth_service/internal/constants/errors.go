package constants

var (
	ErrInternalServer = "internal server error"
	ErrInvalidInput   = "invalid input"
	ErrUnauthorized   = "unauthorized"
	ErrForbidden      = "forbidden"
	ErrNotFound       = "resource not found"
	ErrConflict       = "conflict"

	// Auth-specific errors
	ErrEmailAlreadyExists = "Email already exists"
	ErrInvalidCredentials = "Invalid email or password"
	ErrTokenExpired       = "Token has expired"
	ErrInvalidToken       = "Invalid token"
)
