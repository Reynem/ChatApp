package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtKey struct {
	SecretKey []byte
}

func (j *JwtKey) GenerateToken(username string, userId uint64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"user_id":  userId,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	return token.SignedString(j.SecretKey)
}

func (j *JwtKey) ValidateToken(tokenString string) (string, uint64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.SecretKey, nil
	})

	if err != nil {
		return "", 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", 0, errors.New("invalid claims")
		}
		userIdFloat, ok := claims["user_id"].(float64)
		if !ok {
			return "", 0, errors.New("invalid claims")
		}

		user_id := uint64(userIdFloat)

		return username, user_id, nil
	}

	return "", 0, errors.New("invalid token")
}
