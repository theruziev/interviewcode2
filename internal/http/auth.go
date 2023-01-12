package http

import (
	"net/http"

	"github.com/theruziev/oson_auth/internal/pkg/auth"
	"github.com/theruziev/oson_auth/internal/pkg/httpx"
	"github.com/theruziev/oson_auth/internal/pkg/logging"
	"github.com/theruziev/oson_auth/internal/pkg/validatorx"
)

func (s *UserHandler) Auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	validate := validatorx.FromContext(ctx)

	req, err := httpx.ParseJSON[CredentialRequest](r)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = validate.Struct(req); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := s.userService.Auth(ctx, req.Email, req.Password)
	if err != nil {
		logger.Warnf("failed to auth: %s", err)
		httpx.JSONError(w, http.StatusForbidden, "incorrect user and password")
		return
	}

	tokenResponse := AuthTokenResponse{
		AuthToken:     token.AuthToken,
		ExpireAt:      token.ExpireAt,
		TwoFARequired: token.TwoFARequired,
	}
	httpx.JSONResponse(w, http.StatusOK, tokenResponse)
}

func (s *UserHandler) AuthTwoFA(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	validate := validatorx.FromContext(ctx)
	claim := auth.FromContext(ctx)

	if claim == nil {
		httpx.JSONError(w, http.StatusForbidden, "not permitted")
		return
	}
	if !claim.CheckScope(auth.TwoFACheckScope) {
		httpx.JSONError(w, http.StatusForbidden, "not permitted")
		return
	}

	req, err := httpx.ParseJSON[UserTwoFACodeRequest](r)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = validate.Struct(req); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := s.userService.AuthTwoFA(ctx, claim, req.Code)
	if err != nil {
		logger.Warnf("failed to 2fa: %s", err)
		httpx.JSONError(w, http.StatusForbidden, "incorrect otp code")
		return
	}

	tokenResponse := AuthTokenResponse{
		AuthToken: token.AuthToken,
		ExpireAt:  token.ExpireAt,
	}
	httpx.JSONResponse(w, http.StatusOK, tokenResponse)
}
