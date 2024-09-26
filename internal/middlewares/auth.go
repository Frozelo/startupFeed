package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	httpwriter "github.com/Frozelo/startupFeed/pkg/http"
	jwter "github.com/Frozelo/startupFeed/pkg/jwt"
)

type ctxKey = string

const UserIdKey = ctxKey("userId")

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
		token, err := jwter.ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userId, ok := claims["userId"].(float64)
		fmt.Println(userId)
		if !ok {
			http.Error(w, "Invalid userId in token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIdKey, userId)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userIdCtx := ctx.Value(UserIdKey)
	if userIdFloat, ok := userIdCtx.(float64); ok {
		userId := int64(userIdFloat)
		return userId, true
	}
	return 0, false
}
