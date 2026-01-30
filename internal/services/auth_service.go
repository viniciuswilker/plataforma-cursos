package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ChaveSecreta = []byte("sua_chave_secreta_super_segura_123")

// GerarToken cria um token novo
func GerarToken(usuarioID uint, cargo string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   usuarioID,
		"cargo": cargo,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(ChaveSecreta)
}

// TokenEhValido verifica se o token foi criado/se o usuario está logado
func TokenEhValido(tokenString string) bool {
	if tokenString == "" {
		return false
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado")
		}
		return ChaveSecreta, nil
	})

	if err != nil || token == nil {
		return false
	}

	return token.Valid
}
