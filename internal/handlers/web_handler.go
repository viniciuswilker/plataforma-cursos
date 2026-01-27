package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniciuswilker/plataforma-cursos/internal/repositorios"
)

// ExibirHome renderiza a página inicial
func ExibirHome(c *gin.Context) {
	c.HTML(http.StatusOK, "hello.html", nil)
}

// ExibirCadastro renderiza a página de registro
func ExibirCadastro(c *gin.Context) {
	c.HTML(http.StatusOK, "cadastro.html", nil)
}

// ExibirLogin renderiza a página de login
func ExibirLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// ExibirFeed renderiza o layout do feed com dados dinâmicos
func ExibirFeed(c *gin.Context) {

	id, existe := c.Get("usuarioID")
	if !existe {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	usuario, err := repositorios.BuscarPorID(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	c.HTML(http.StatusOK, "layout", gin.H{
		"title":   "Feed de cursos",
		"page":    "feed",
		"usuario": usuario,
	})
}

func PaginaProfessor(c *gin.Context) {

	id, existe := c.Get("usuarioID")
	if !existe {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	usuario, err := repositorios.BuscarPorID(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	c.HTML(http.StatusOK, "layout", gin.H{
		"title":   "Pagina do professor",
		"page":    "professor",
		"usuario": usuario,
	})
}
