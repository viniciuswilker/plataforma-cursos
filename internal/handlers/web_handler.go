package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/viniciuswilker/plataforma-cursos/internal/database"
	"github.com/viniciuswilker/plataforma-cursos/internal/models"
	"github.com/viniciuswilker/plataforma-cursos/internal/repositorios"
	"github.com/viniciuswilker/plataforma-cursos/internal/services"
)

// ExibirHome renderiza a página inicial
func ExibirHome(c *gin.Context) {
	c.HTML(http.StatusOK, "hello.html", nil)
}

// ExibirCadastro renderiza a página de registro
func ExibirCadastro(c *gin.Context) {
	tokenString, _ := c.Cookie("token")
	if services.TokenEhValido(tokenString) {
		c.Redirect(http.StatusSeeOther, "/feed/")
		return
	}

	c.HTML(http.StatusOK, "cadastro.html", nil)
}

// ExibirLogin renderiza a página de login
func ExibirLogin(c *gin.Context) {
	tokenString, _ := c.Cookie("token")
	if services.TokenEhValido(tokenString) {
		c.Redirect(http.StatusSeeOther, "/feed/")
		return
	}

	c.HTML(http.StatusOK, "login.html", nil)
}

// ExibirFeed renderiza o layout do feed com dados dinâmicos
func ExibirFeed(c *gin.Context) {

	id, existe := c.Get("usuarioID")
	if !existe {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	cursos, err := repositorios.ListarCursos()
	if err != nil {
		cursos = []models.Curso{}
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
		"cursos":  cursos,
	})
}

func PaginaProfessor(c *gin.Context) {

	idRaw, existe := c.Get("usuarioID")
	if !existe {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	var id uint
	if val, ok := idRaw.(float64); ok {
		id = uint(val)
	} else {
		id = idRaw.(uint)
	}

	usuario, err := repositorios.BuscarPorID(id)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	stats, _ := repositorios.ObterStatsProfessor(id)
	cursos, _ := repositorios.ListarCursosPorInstrutor(id)

	c.HTML(http.StatusOK, "layout", gin.H{
		"title":   "Pagina do professor",
		"page":    "professor",
		"usuario": usuario,
		"stats":   stats,
		"cursos":  cursos,
	})
}

func ExibirEdicaoCurso(c *gin.Context) {
	cursoID := c.Param("id")
	idRaw, _ := c.Get("usuarioID")

	var instrutorID uint
	if val, ok := idRaw.(float64); ok {
		instrutorID = uint(val)
	} else {
		instrutorID = idRaw.(uint)
	}

	usuario, err := repositorios.BuscarPorID(instrutorID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	var curso models.Curso
	err = database.DB.Preload("Modulos.Aulas").
		Where("id = ? AND instrutor_id = ?", cursoID, instrutorID).
		First(&curso).Error

	if err != nil {
		c.String(http.StatusNotFound, "Curso não encontrado")
		return
	}

	c.HTML(http.StatusOK, "layout", gin.H{
		"title":   "Editando: " + curso.Titulo,
		"page":    "editar_curso",
		"usuario": usuario,
		"curso":   curso,
	})
}

func ExibirDetalhesCurso(c *gin.Context) {
	id := c.Param("id")
	var curso models.Curso

	err := database.DB.Preload("Instrutor").
		Preload("Modulos.Aulas").
		First(&curso, id).Error

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/feed")
		return
	}

	idRaw, _ := c.Get("usuarioID")
	var usuarioID uint
	if val, ok := idRaw.(float64); ok {
		usuarioID = uint(val)
	} else {
		usuarioID = idRaw.(uint)
	}

	var matriculado bool
	if idRaw != nil {
		var count int64
		database.DB.Model(&models.Matricula{}).
			Where("curso_id = ? AND usuario_id = ?", id, idRaw).
			Count(&count)
		matriculado = count > 0
	}

	usuario, _ := repositorios.BuscarPorID(usuarioID)

	c.HTML(http.StatusOK, "layout", gin.H{
		"matriculado": matriculado,
		"title":       curso.Titulo,
		"page":        "conteudo_curso",
		"curso":       curso,
		"usuario":     usuario,
	})
}
