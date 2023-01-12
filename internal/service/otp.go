package service

import (
	"context"
	"fmt"

	"github.com/theruziev/oson_auth/internal/model"
)

func (s *UserService) RequestEnableOTPStep1(ctx context.Context, publicID string) (*model.OtpToken, error) {
	user, err := s.userStore.Get(ctx, publicID)
	if err != nil {
		return nil, err
	}
	otpRes, err := s.otp.Generate(ctx, user.PublicID)
	if err != nil {
		return nil, err
	}

	err = s.userStore.SetOtpSecret(ctx, user.PublicID, otpRes.Secret)
	if err != nil {
		return nil, err
	}

	return &model.OtpToken{
		OtpURL: otpRes.URL,
	}, nil
}

func (s *UserService) RequestEnableOTPStep2(ctx context.Context, publicID, code string) (*model.OtpRecoveryCode, error) {
	user, err := s.userStore.Get(ctx, publicID)
	if err != nil {
		return nil, err
	}
	isValid, err := s.otp.ValidateCode(ctx, user.OtpSecret, code)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, fmt.Errorf("invalid otp code")
	}
	codes, err := s.otp.GenerateRecoveryCodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recovery code: %w", err)
	}

	if err := s.userStore.SetOtpRecoveryCodes(ctx, user.PublicID, codes); err != nil {
		return nil, err
	}
	if err := s.userStore.SetOtpEnabled(ctx, user.PublicID, true); err != nil {
		return nil, err
	}
	return &model.OtpRecoveryCode{
		Codes: codes,
	}, nil
}

func (s *UserService) DisableOTP(ctx context.Context, publicID string) error {
	user, err := s.userStore.Get(ctx, publicID)
	if err != nil {
		return err
	}

	if err := s.userStore.SetOtpEnabled(ctx, user.PublicID, false); err != nil {
		return err
	}
	return nil
}
