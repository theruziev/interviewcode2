package http

import "time"

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}

type AuthTokenResponse struct {
	AuthToken     string    `json:"auth_token"`
	ExpireAt      time.Time `json:"expire_at"`
	TwoFARequired bool      `json:"twofa_required"`
}

type CredentialRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	PublicID  string    `json:"public_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserResetPasswordReqRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UserResetPasswordRequest struct {
	ResetCode string `json:"reset_code"`
	Password  string `json:"password"   validate:"required,min=6"`
}

type UserChangePasswordRequest struct {
	Password string `json:"password" validate:"required,min=6"`
}

type UserTwoFACodeRequest struct {
	Code string `json:"code"`
}

type UserOtpResponse struct {
	URL string `json:"url"`
}

type UserRecoveryCodeResponse struct {
	Codes []string `json:"codes"`
}

type ContentResponse struct {
	ID int64
}
