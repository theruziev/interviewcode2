package model

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserStatus string

const (
	UserStatusActivate   UserStatus = "activated"
	UserStatusRegistered UserStatus = "registered"
)

type User struct {
	ID                uint64     `db:"id" json:"id"`
	PublicID          string     `db:"public_id" json:"public_id"`
	Email             string     `db:"email" json:"email"`
	FirstName         string     `db:"first_name" json:"first_name"`
	LastName          string     `db:"last_name" json:"last_name"`
	Password          string     `db:"password" json:"password"`
	Status            UserStatus `db:"status" json:"status"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
	ActivationCode    string     `db:"activation_code" json:"activation_code"`
	ResetPasswordCode string     `db:"reset_password_code" json:"reset_password_code"`

	OtpSecret        string   `json:"secret" db:"otp_secret"`
	OtpRecoveryCodes []string `json:"recovery_codes"  db:"otp_recovery_codes"`
	OtpEnabled       bool     `db:"otp_enabled"`
}

func (u *User) SetPassword(passwd string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	u.Password = string(hashed)
	return nil
}

func (u *User) Touch() {
	u.UpdatedAt = time.Now()
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type AuthToken struct {
	AuthToken     string    `json:"auth_token"`
	ExpireAt      time.Time `json:"expire_at"`
	TwoFARequired bool      `json:"two_fa_required"`
}

type OtpToken struct {
	OtpURL string
}

type OtpRecoveryCode struct {
	Codes []string
}
