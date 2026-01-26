package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("123456789")

func GerarToken(usuarioID uint, cargo string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   usuarioID,
		"cargo": cargo,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(secretKey)
}
