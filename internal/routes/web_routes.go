package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/viniciuswilker/plataforma-cursos/internal/handlers"
	"github.com/viniciuswilker/plataforma-cursos/internal/middlewares"
)

func RegistrarRotasWeb(router *gin.Engine) {
	web := router.Group("/")
	{
		web.GET("/cadastro", handlers.ExibirCadastro)
		web.GET("/login", handlers.ExibirLogin)
		web.GET("/logout", handlers.Logout)
		web.GET("", handlers.ExibirHome)

	}

	privado := router.Group("/")
	privado.Use(middlewares.Autorizar("aluno", "instrutor", "admin"))
	{
		privado.GET("/feed/", handlers.ExibirFeed)
		privado.GET("/professor/cursos/", handlers.PaginaProfessor)
	}
}
