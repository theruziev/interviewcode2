package v0

import (
	"time"
)

type UserStatus string

const (
	UserStatusActivate   UserStatus = "activated"
	UserStatusRegistered UserStatus = "registered"
)

type UserRegisteredEvent struct {
	PublicID       string `json:"public_id"`
	Email          string `json:"email"`
	ActivationCode string `json:"activation_code"`
}

type UserEvent struct {
	PublicID  string     `json:"public_id"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Status    UserStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type UserResetPasswordEvent struct {
	PublicID          string `json:"public_id"`
	Email             string `json:"email"`
	ResetPasswordCode string `json:"reset_password_code"`
}
