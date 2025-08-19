package constants

var (
	ErrInternalServer = "internal server error"
	ErrInvalidInput   = "invalid input"
	ErrUnauthorized   = "unauthorized"
	ErrForbidden      = "forbidden"
	ErrNotFound       = "resource not found"
	ErrConflict       = "conflict"

	// User-specific errors
	ErrEmailAlreadyExists = "Email already exists"
	ErrInvalidCredentials = "Invalid email or password"
	ErrTokenExpired       = "Token has expired"
	ErrInvalidToken       = "Invalid token"
	ErrUserNotFound       = "no user found with given identifier"

	// Otp-specific errors
	ErrOtpNotFound       = "Otp not found"
	ErrOtpExpired        = "Otp has expired"
	ErrOtpAlreadyUsed    = "Otp already used"
	ErrOtpInvalid        = "Otp is invalid"
	ErrOtpResendTooSoon  = "Otp resend not allowed yet, please wait"
	ErrOtpResendLimit    = "Otp resend limit reached"
	ErrOtpGenerationFail = "failed to generate otp"
	ErrOtpVerification   = "Otp verification failed"
)
