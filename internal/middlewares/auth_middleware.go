package middlewares

import (
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
			abortarComErro(c, http.StatusUnauthorized, "Token não fornecido")
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return services.ChaveSecreta, nil
		})

		if err != nil || !token.Valid {
			c.SetCookie("token", "", -1, "/", "", false, true)
			abortarComErro(c, http.StatusUnauthorized, "Sessão expirada. Faça login novamente")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			cargoUsuario, okCargo := claims["cargo"].(string)
			if !okCargo {
				abortarComErro(c, http.StatusUnauthorized, "Falha ao processar permissões")
				return
			}

			permitido := false
			for _, cPermitido := range cargosPermitidos {
				if cPermitido == cargoUsuario {
					permitido = true
					break
				}
			}

			if !permitido {
				abortarComErro(c, http.StatusForbidden, "Você não tem permissão para acessar este recurso")
				return
			}

			c.Set("usuarioID", claims["sub"])
			c.Next()
		} else {
			abortarComErro(c, http.StatusUnauthorized, "Falha ao processar claims")
		}
	}
}

func abortarComErro(c *gin.Context, status int, mensagem string) {
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		c.AbortWithStatusJSON(status, gin.H{"erro": mensagem})
		return
	}

	c.Abort()
	c.Redirect(http.StatusSeeOther, "/login")
}
