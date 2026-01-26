package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/viniciuswilker/plataforma-cursos/internal/handlers"
	"github.com/viniciuswilker/plataforma-cursos/internal/middlewares"
)

func ConfigurarRotas(router *gin.Engine) {

	// autenticacao
	auth := router.Group("/auth")
	{
		auth.POST("/cadastro", handlers.Cadastrar)
		auth.POST("/login", handlers.Login)
	}

	// protegidas
	api := router.Group("/api")

	api.Use(middlewares.Autorizar("aluno", "instrutor", "admin"))
	{
		api.GET("/cursos", handlers.ListarCursos)

		//  instrutor e admin
		admin := api.Group("/cursos")
		admin.Use(middlewares.Autorizar("instrutor", "admin"))
		{
			admin.POST("/", handlers.CriarCurso)
		}
	}

}
