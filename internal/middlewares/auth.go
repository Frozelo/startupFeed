package middlewares

import (
	"errors"
	"net/http"
	"strings"

	httpwriter "github.com/Frozelo/startupFeed/pkg/http"
	"github.com/Frozelo/startupFeed/pkg/jwt"
)

func JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			err := errors.New("missing auth header")
			httpwriter.Error(w, http.StatusUnauthorized, err, err.Error(), nil)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if err := jwt.VerifyToken(tokenString); err != nil {
			httpwriter.Error(
				w,
				http.StatusUnauthorized,
				err,
				"Invalid token",
				nil,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
