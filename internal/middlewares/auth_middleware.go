package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/viniciuswilker/plataforma-cursos/internal/services"
)

func Autorizar(cargosPermitidos ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			tokenString = strings.Replace(authHeader, "Bearer ", "", 1)
		}

		if tokenString == "" {
			cookieToken, err := c.Cookie("token")
			if err == nil {
				tokenString = cookieToken
			}
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Token não fornecido"})
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return services.ChaveSecreta, nil
		})

		if err != nil || !token.Valid {
			fmt.Println("Erro no Parse do JWT:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Token inválido ou expirado"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			cargoUsuario := claims["cargo"].(string)

			permitido := false
			for _, cPermitido := range cargosPermitidos {
				if cPermitido == cargoUsuario {
					permitido = true
					break
				}
			}

			if !permitido {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"erro": "Você não tem permissão para acessar este recurso"})
				return
			}

			c.Set("usuarioID", claims["sub"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Falha ao processar permissões"})
		}
	}
}
