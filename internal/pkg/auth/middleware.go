package auth

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/theruziev/oson_auth/internal/pkg/httpx"
)

func Middleware(jwtSecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := getToken(r)
			if tokenString == "" {
				httpx.JSONError(w, http.StatusForbidden, "token is empty")
				return
			}
			token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})
			if err != nil {
				httpx.JSONError(w, http.StatusForbidden, "failed to decrypt jwt token")
				return
			}
			claim, ok := token.Claims.(*Claim)
			if !ok {
				httpx.JSONError(w, http.StatusForbidden, "failed to get claim from token")
				return
			}

			ctx := WithClaim(r.Context(), claim)
			r = r.Clone(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func CheckScope(scopes ...Scope) func(next http.Handler) http.Handler {
	scopeMap := make(map[Scope]struct{})
	for _, scp := range scopes {
		scopeMap[scp] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claim := FromContext(r.Context())
			if claim == nil {
				panic("incorrect usage of middleware")
			}
			for _, scp := range claim.Scopes {
				if _, ok := scopeMap[scp]; ok {
					next.ServeHTTP(w, r)
					return
				}
			}

			httpx.JSONError(w, http.StatusForbidden, "invalid scope")

		})
	}
}

func getToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	return strings.Replace(authHeader, "Bearer ", "", -1)
}
