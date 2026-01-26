package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniciuswilker/plataforma-cursos/internal/handlers"
	"github.com/viniciuswilker/plataforma-cursos/internal/middlewares"
)

func ConfigurarRotas(router *gin.Engine) {

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	// ROTAS WEB
	web := router.Group("/")
	{
		web.GET("/cadastro", func(c *gin.Context) {
			c.HTML(http.StatusOK, "cadastro.html", nil)
		})

		web.GET("/hello", func(c *gin.Context) {
			c.HTML(http.StatusOK, "feed.html", nil)
		})

		web.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", nil)
		})
	}

	// ROTAS DE API

	// Autenticação
	auth := router.Group("/auth")
	{
		auth.POST("/cadastro", handlers.Cadastrar)
		auth.POST("/login", handlers.Login)
	}

	// Recursos Protegidos
	api := router.Group("/api")
	api.Use(middlewares.Autorizar("aluno", "instrutor", "admin"))
	{
		api.GET("/cursos", handlers.ListarCursos)

		// Apenas instrutor e admin
		admin := api.Group("/cursos")
		admin.Use(middlewares.Autorizar("instrutor", "admin"))
		{
			admin.POST("/", handlers.CriarCurso)
		}
	}
}
