package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/viniciuswilker/plataforma-cursos/internal/database"
	"github.com/viniciuswilker/plataforma-cursos/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Aviso: Arquivo .env não encontrado. Usando variáveis de ambiente do sistema.")
	}

	database.Conectar()

	servidor := gin.Default()

	servidor.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"mensagem": "Servidor rodando e banco conectado!"})
	})

	routes.ConfigurarRotas(servidor)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	servidor.Run(":" + port)

}
