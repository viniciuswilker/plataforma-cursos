package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniciuswilker/plataforma-cursos/internal/services"
)

type RequisicaoCadastro struct {
	Nome      string `json:"nome" binding:"required"`
	Sobrenome string `json:"sobrenome" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Senha     string `json:"senha" binding:"required,min=6"`
}

func Cadastrar(c *gin.Context) {

	var req RequisicaoCadastro
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	err := services.CadastrarUsuario(req.Nome, req.Sobrenome, req.Email, req.Senha)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mensagem": "Usuário criado com sucesso!"})
}
