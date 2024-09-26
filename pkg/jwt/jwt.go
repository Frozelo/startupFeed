package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Frozelo/startupFeed/internal/models"
)

var (
	secretKey      = []byte("very-secret-key")
	expirationTime = time.Now().Add(24 * time.Hour)
)

type claims struct {
	UserId int64  `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func CreateToken(user *models.User) (string, error) {
	numericDate := jwt.NewNumericDate(expirationTime)
	claimsData := &claims{
		UserId: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: numericDate,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsData, nil)
	tokenString, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrInvalidKey
			}
			return []byte(secretKey), nil
		},
	)
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return token, err
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	)
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("token is not valid")
	}
	return nil
}
