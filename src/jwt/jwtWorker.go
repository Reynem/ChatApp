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
			"exp":      time.Now().Add(time.Hour * 1).Unix(),
		})

	return token.SignedString(j.SecretKey)
}

func (j *JwtKey) ValidateToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.SecretKey, nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", "", errors.New("invalid claims")
		}
		user_id, ok := claims["user_id"].(string)
		if !ok {
			return "", "", errors.New("invalid claims")
		}
		return username, user_id, nil
	}

	return "", "", errors.New("invalid token")
}
