package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniciuswilker/plataforma-cursos/internal/services"
)

type RequisicaoCurso struct {
	Titulo    string `json:"titulo" binding:"required"`
	Descricao string `json:"descricao"`
}

func CriarCurso(c *gin.Context) {
	var req RequisicaoCurso
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	idUsuario := c.MustGet("usuarioID").(float64)

	err := services.CriarNovoCurso(req.Titulo, req.Descricao, uint(idUsuario))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao criar curso"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mensagem": "Curso criado com sucesso!"})
}

func ListarCursos(c *gin.Context) {
	cursos, err := services.BuscarTodosCursos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao listar cursos"})
		return
	}
	c.JSON(http.StatusOK, cursos)
}
