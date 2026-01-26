package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/viniciuswilker/plataforma-cursos/internal/handlers"
)

func ConfigurarRotas(router *gin.Engine) {

	// autenticacao
	auth := router.Group("/auth")
	{
		auth.POST("/cadastro", handlers.Cadastrar)
	}

	// cursos
	

}
