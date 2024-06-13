package util

import (
	"backend/server/env"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type JWT struct{}

func (ref JWT) Generate(user_id string) (string, string, error) {
	secretKey := env.GetSecretKey()
	Generate := Generate{}
	jti := Generate.UUIDv4()
	claims := jwt.MapClaims{
		"user_id": user_id,
		"iat":     time.Now().Unix(),
		"jti":     jti, // Assign JTI to the JWT claims
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_str, err := token.SignedString([]byte(secretKey))
	return token_str, jti, err
}

func (ref JWT) Validate(tokenString string) (jwt.MapClaims, error) {
	secretKey := env.GetSecretKey()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid token")
}
