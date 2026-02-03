package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/viniciuswilker/plataforma-cursos/internal/handlers"
	"github.com/viniciuswilker/plataforma-cursos/internal/middlewares"
)

func RegistrarRotasAPI(router *gin.Engine) {
	auth := router.Group("/auth")
	{
		auth.POST("/cadastro", handlers.Cadastrar)
		auth.POST("/login", handlers.Login)
	}

	api := router.Group("/api")
	api.Use(middlewares.Autorizar("aluno", "instrutor", "admin"))
	{
		api.POST("/cursos/:id/matricular", handlers.MatricularAluno)

		api.GET("/cursos", handlers.ListarCursos)

		professor := api.Group("/professor")
		professor.Use(middlewares.Autorizar("instrutor", "admin"))
		{
			professor.POST("/cursos/criar", handlers.CriarCurso)

			professor.POST("/aulas/criar", handlers.CriarAula)
			professor.DELETE("/aulas/:id", handlers.ExcluirAula)

			professor.POST("/modulos/criar", handlers.CriarModulo)
			professor.DELETE("/modulos/:id", handlers.ExcluirModulo)

		}
	}
}
