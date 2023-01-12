package auth

import (
	"context"
	"time"
)

type contextKey string

const claimKey = contextKey("claim")

type AuthOption struct {
	JWTSecret string        `help:"listen string" env:"SECRET"`
	JWTTtl    time.Duration `help:"ttl" env:"TTL"`
	Otp       OtpConfig     `embed:"" prefix:"otp." envprefix:"OTP_" validate:"required,dive,required"`
}

func WithClaim(ctx context.Context, claim *Claim) context.Context {
	return context.WithValue(ctx, claimKey, claim)
}

func FromContext(ctx context.Context) *Claim {
	if claim, ok := ctx.Value(claimKey).(*Claim); ok {
		return claim
	}

	return nil
}
