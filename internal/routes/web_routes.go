package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegistrarRotasWeb(router *gin.Engine) {
	web := router.Group("/")
	{
		web.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "hello.html", nil)
		})

		web.GET("/cadastro", func(c *gin.Context) {
			c.HTML(http.StatusOK, "cadastro.html", nil)
		})

		web.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", nil)
		})

		web.GET("/feed/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "layout", gin.H{
				"title": "Feed de cursos",
				"page":  "feed",
			})
		})
	}
}
