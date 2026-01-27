package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ChaveSecreta = []byte("sua_chave_secreta_super_segura_123")

func GerarToken(usuarioID uint, cargo string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   usuarioID,
		"cargo": cargo,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(ChaveSecreta)
}
