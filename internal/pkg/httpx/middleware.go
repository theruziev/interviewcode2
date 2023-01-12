package httpx

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/theruziev/oson_auth/internal/pkg/logging"
	"github.com/theruziev/oson_auth/internal/pkg/validatorx"
	"go.uber.org/zap"
)

// PopulateLogger populates the logger onto the context.
func PopulateLogger(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx = logging.WithLogger(ctx, logger)
			r = r.Clone(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func PopulateValidator(validate *validator.Validate) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			ctx = validatorx.WithValidator(ctx, validate)
			r = r.Clone(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func Recoverer(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
					logger.Errorf("panic: %s", rvr)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
