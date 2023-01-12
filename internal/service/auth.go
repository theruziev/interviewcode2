package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/theruziev/oson_auth/internal/model"
	"github.com/theruziev/oson_auth/internal/pkg/auth"
)

const (
	twoFARequiredExpireAt = 5 * time.Minute
)

func (s *UserService) Auth(ctx context.Context, username, password string) (*model.AuthToken, error) {
	user, err := s.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	if user.Status != model.UserStatusActivate {
		return nil, fmt.Errorf("user not active")
	}
	if !user.ValidatePassword(password) {
		return nil, fmt.Errorf("incorrect password or username")
	}
	expireAt := time.Now().Add(s.authOpt.JWTTtl)

	twoFARequired := s.authOpt.Otp.Enabled && user.OtpEnabled
	if twoFARequired {
		expireAt = time.Now().Add(twoFARequiredExpireAt)
	}

	claim := auth.Claim{
		PublicID: user.PublicID,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
		},
		Scopes: []auth.Scope{auth.UserScope},
	}

	if twoFARequired {
		claim.Scopes = []auth.Scope{auth.TwoFACheckScope}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(s.authOpt.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	return &model.AuthToken{
		AuthToken:     tokenString,
		TwoFARequired: twoFARequired,
		ExpireAt:      expireAt,
	}, nil
}

func (s *UserService) AuthTwoFA(ctx context.Context, claim *auth.Claim, code string) (*model.AuthToken, error) {
	if !s.authOpt.Otp.Enabled {
		return nil, nil
	}
	user, err := s.GetByUsername(ctx, claim.Email)
	if err != nil {
		return nil, err
	}
	if user.Status != model.UserStatusActivate {
		return nil, fmt.Errorf("user not active")
	}
	expireAt := time.Now().Add(s.authOpt.JWTTtl)

	isValid, err := s.otp.ValidateCode(ctx, user.OtpSecret, code)
	if err != nil {
		return nil, err
	}
	var foundInRecovery bool
	if !isValid {
		// try to use recovery code
		for _, recoveryCode := range user.OtpRecoveryCodes {
			if recoveryCode == code {
				foundInRecovery = true
				break
			}
		}
		if !foundInRecovery {
			return nil, fmt.Errorf("incorrect otp code")
		}
	}

	newClaim := auth.Claim{
		PublicID: user.PublicID,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
		},
		Scopes: []auth.Scope{auth.UserScope},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaim)
	tokenString, err := token.SignedString([]byte(s.authOpt.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	if foundInRecovery {
		err = s.userStore.RemoveCodeFromRecoveryCode(ctx, user.PublicID, code)
		if err != nil {
			return nil, err
		}
	}

	return &model.AuthToken{
		AuthToken: tokenString,
		ExpireAt:  expireAt,
	}, nil
}
