package dto

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"fullname" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type SigninRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type RequestResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Code        string `json:"code" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}
