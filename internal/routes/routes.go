package routes

import "github.com/gin-gonic/gin"

func ConfigurarRotas(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	router.Static("/uploads", "./uploads")

	RegistrarRotasWeb(router)
	RegistrarRotasAPI(router)
}
