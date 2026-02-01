package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniciuswilker/plataforma-cursos/internal/database"
	"github.com/viniciuswilker/plataforma-cursos/internal/models"
	"github.com/viniciuswilker/plataforma-cursos/internal/services"
)

type RequisicaoCurso struct {
	Titulo    string `json:"titulo" binding:"required"`
	Descricao string `json:"descricao"`
}

func CriarCurso(c *gin.Context) {

	idInterface, _ := c.Get("usuarioID")

	var instrutorID uint
	if val, ok := idInterface.(float64); ok {
		instrutorID = uint(val)
	} else if val, ok := idInterface.(uint); ok {
		instrutorID = val
	}

	var input struct {
		Titulo    string `json:"titulo" binding:"required"`
		Descricao string `json:"descricao"`
		CapaURL   string `json:"capa_url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Título é obrigatório"})
		return
	}

	curso := models.Curso{
		Titulo:      input.Titulo,
		Descricao:   input.Descricao,
		CapaURL:     input.CapaURL,
		InstrutorID: instrutorID,
	}

	if err := database.DB.Create(&curso).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar no banco"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Curso criado com sucesso",
		"id":      curso.ID,
	})
}

func ListarCursos(c *gin.Context) {
	cursos, err := services.BuscarTodosCursos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao listar cursos"})
		return
	}
	c.JSON(http.StatusOK, cursos)
}
