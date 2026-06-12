package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shubham071122/collab/internal/config"
)

func GenerateJWT(userID string, isVerified bool) (string, error) {

	jwtSecret := config.LoadConfig().JWTSecret

	claims := jwt.MapClaims{
		"user_id":     userID,
		"is_verified": isVerified,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
