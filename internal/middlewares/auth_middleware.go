package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Autorizar(cargosPermitidos ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Token não fornecido"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte("sua_chave_secreta_super_segura"), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			cargoUsuario := claims["cargo"].(string)

			permitido := false
			for _, c := range cargosPermitidos {
				if c == cargoUsuario {
					permitido = true
					break
				}
			}

			if !permitido {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"erro": "Acesso negado"})
				return
			}

			c.Set("usuarioID", claims["sub"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"erro": "Token inválido"})
		}
	}
}
