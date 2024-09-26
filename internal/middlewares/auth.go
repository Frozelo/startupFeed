package middlewares

import (
	"errors"
	"net/http"
	"strings"

	httpwriter "github.com/Frozelo/startupFeed/pkg/http"
	"github.com/Frozelo/startupFeed/pkg/jwt"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().
			Add("Access-Control-Allow-Origin", "*")
		w.Header().
			Add("Access-Control-Allow-Credentials", "true")
		w.Header().
			Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().
			Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		if r.Method == "OPTIONS" {
			httpwriter.Error(w, http.StatusNoContent, nil, "No content", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

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
