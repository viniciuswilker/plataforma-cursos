package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniciuswilker/plataforma-cursos/internal/services"
)

type RequisicaoCadastro struct {
	Nome      string `json:"nome" binding:"required"`
	Sobrenome string `json:"sobrenome" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Senha     string `json:"senha" binding:"required,min=6"`
	Cargo     string `json:"cargo" binding:"required"`
}

type RequisicaoLogin struct {
	Email string `json:"email" binding:"required"`
	Senha string `json:"senha" binding:"required"`
}

func Cadastrar(c *gin.Context) {

	var req RequisicaoCadastro
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println("Erro de Bind:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	err := services.CadastrarUsuario(req.Nome, req.Sobrenome, req.Email, req.Senha, req.Cargo)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mensagem": "Usuário criado com sucesso!"})
}

func Login(c *gin.Context) {
	var req RequisicaoLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "Dados inválidos"})
		return
	}

	token, err := services.Login(req.Email, req.Senha)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": err.Error()})
		return
	}
	c.SetCookie("token", token, 3600*24, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"token": token})
}
