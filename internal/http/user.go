package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/theruziev/oson_auth/internal/model"
	"github.com/theruziev/oson_auth/internal/pkg/auth"
	"github.com/theruziev/oson_auth/internal/pkg/errz"
	"github.com/theruziev/oson_auth/internal/pkg/httpx"
	"github.com/theruziev/oson_auth/internal/pkg/logging"
	"github.com/theruziev/oson_auth/internal/pkg/validatorx"
	"github.com/theruziev/oson_auth/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (s *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	validate := validatorx.FromContext(ctx)

	req, err := httpx.ParseJSON[RegisterRequest](r)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = validate.Struct(req); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := s.userService.Register(ctx, &model.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		if errz.ConflictErr.Is(err) {
			httpx.JSONError(w, http.StatusConflict, "user already exist")
			return
		}
		httpx.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Debugf("user %d succesfull registered", user.ID)
	httpx.JSONOKResponse(w)
}

func (s *UserHandler) Activate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	activationCode := chi.URLParam(r, "aid")
	err := s.userService.Activate(ctx, activationCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			httpx.JSONError(w, http.StatusNotFound, err.Error())
			return
		}
		httpx.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.JSONOKResponse(w)
}

func (s *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claim := auth.FromContext(ctx)

	user, err := s.userService.GetByUsername(ctx, claim.Email)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userResponse := UserResponse{
		PublicID:  user.PublicID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	httpx.JSONResponse(w, http.StatusOK, userResponse)
}

func (s *UserHandler) RequestEnableOTPStep1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claim := auth.FromContext(ctx)
	otpToken, err := s.userService.RequestEnableOTPStep1(ctx, claim.PublicID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.JSONResponse(w, http.StatusOK, UserOtpResponse{
		URL: otpToken.OtpURL,
	})
}

func (s *UserHandler) RequestEnableOTPStep2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claim := auth.FromContext(ctx)

	req, err := httpx.ParseJSON[UserTwoFACodeRequest](r)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	otpRes, err := s.userService.RequestEnableOTPStep2(ctx, claim.PublicID, req.Code)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.JSONResponse(w, http.StatusOK, UserRecoveryCodeResponse{
		Codes: otpRes.Codes,
	})
}

func (s *UserHandler) DisableOTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claim := auth.FromContext(ctx)

	err := s.userService.DisableOTP(ctx, claim.PublicID)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.JSONOKResponse(w)
}
