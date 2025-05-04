package utils

import (
	"banking_ledger/config"
	"banking_ledger/models"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(config.JWT_SECRET)

func GenerateJWTForUser(userId int, role string, name string) (string, error) {
	claims := &models.CustomClaims{
		Role: role,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userId),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
