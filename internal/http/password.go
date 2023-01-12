package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/theruziev/oson_auth/internal/pkg/auth"
	"github.com/theruziev/oson_auth/internal/pkg/errz"
	"github.com/theruziev/oson_auth/internal/pkg/httpx"
	"github.com/theruziev/oson_auth/internal/pkg/validatorx"
)

func (s *UserHandler) ResetPasswordRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	validate := validatorx.FromContext(ctx)

	req, err := httpx.ParseJSON[UserResetPasswordReqRequest](r)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = validate.Struct(req); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = s.userService.ResetPasswordRequest(ctx, req.Email)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpx.JSONOKResponse(w)

}

func (s *UserHandler) GetByResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resetCode := chi.URLParam(r, "rcode")

	_, err := s.userService.GetByResetPassword(ctx, resetCode)
	if err != nil {
		if errz.NotFoundErr.Is(err) {
			httpx.JSONError(w, http.StatusNotFound, err.Error())
			return
		}
		httpx.JSONError(w, http.StatusNotFound, err.Error())
		return
	}

	httpx.JSONOKResponse(w)
}

func (s *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	validate := validatorx.FromContext(ctx)

	req, err := httpx.ParseJSON[UserResetPasswordRequest](r)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = validate.Struct(req); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = s.userService.ResetPassword(ctx, req.ResetCode, req.Password)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.JSONOKResponse(w)
}

func (s *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	validate := validatorx.FromContext(ctx)
	claim := auth.FromContext(ctx)

	req, err := httpx.ParseJSON[UserResetPasswordRequest](r)
	if err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = validate.Struct(req); err != nil {
		httpx.JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = s.userService.ChangePassword(ctx, claim.PublicID, req.Password)
	if err != nil {
		httpx.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpx.JSONOKResponse(w)
}
