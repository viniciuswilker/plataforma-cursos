package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/viniciuswilker/plataforma-cursos/internal/database"
	"github.com/viniciuswilker/plataforma-cursos/internal/models"
	"github.com/viniciuswilker/plataforma-cursos/internal/repositorios"
	"github.com/viniciuswilker/plataforma-cursos/internal/services"
)

// ExibirHome renderiza a página inicial
func ExibirHome(c *gin.Context) {
	busca := c.Query("q")
	var cursos []models.Curso

	if busca != "" {
		database.DB.Where("titulo LIKE ? OR descricao LIKE ?", "%"+busca+"%", "%"+busca+"%").Find(&cursos)
	} else {
		// Busca todos
		database.DB.Find(&cursos)
	}

	c.HTML(http.StatusOK, "hello.html", gin.H{
		"cursos": cursos,
		"busca":  busca,
	})
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

func PublicoExibirDetalhesCurso(c *gin.Context) {
	id := c.Param("id")
	var curso models.Curso

	err := database.DB.Preload("Instrutor").
		Preload("Modulos.Aulas").
		First(&curso, id).Error

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	totalAulas := 0
	for _, m := range curso.Modulos {
		totalAulas += len(m.Aulas)
	}

	c.HTML(http.StatusOK, "preview_curso_publico.html", gin.H{
		"title":      curso.Titulo,
		"curso":      curso,
		"totalAulas": totalAulas,
	})
}
func ExibirDetalhesCurso(c *gin.Context) {
	slugParam := c.Param("slug")
	var curso models.Curso

	err := database.DB.Preload("Instrutor").
		Preload("Modulos").
		Preload("Modulos.Aulas").
		Where("slug = ?", slugParam).
		First(&curso).Error

	if err != nil {
		c.Redirect(http.StatusSeeOther, "/feed")
		return
	}

	idRaw, logado := c.Get("usuarioID")
	var usuarioID uint
	var matriculado bool
	var porcentagem int

	if logado && idRaw != nil {
		switch v := idRaw.(type) {
		case float64:
			usuarioID = uint(v)
		case uint:
			usuarioID = v
		}

		var countMatricula int64
		database.DB.Model(&models.Matricula{}).
			Where("curso_id = ? AND usuario_id = ?", curso.ID, usuarioID).
			Count(&countMatricula)
		matriculado = countMatricula > 0

		if matriculado {
			totalAulas := 0
			for _, m := range curso.Modulos {
				totalAulas += len(m.Aulas)
			}

			if totalAulas > 0 {
				var concluidas int64
				database.DB.Model(&models.ProgressoAula{}).
					Joins("JOIN aulas ON aulas.id = progresso_aulas.aula_id").
					Joins("JOIN modulos ON modulos.id = aulas.modulo_id").
					Where("modulos.curso_id = ? AND progresso_aulas.usuario_id = ?", curso.ID, usuarioID).
					Count(&concluidas)

				porcentagem = int((float64(concluidas) / float64(totalAulas)) * 100)
			}
		}
	}

	var usuario *models.Usuario
	if usuarioID > 0 {
		usuario, _ = repositorios.BuscarPorID(usuarioID)
	}

	c.HTML(http.StatusOK, "layout", gin.H{
		"title":       curso.Titulo,
		"page":        "conteudo_curso",
		"curso":       curso,
		"usuario":     usuario,
		"matriculado": matriculado,
		"progresso":   porcentagem,
	})
}
func ExibirMeusCursos(c *gin.Context) {
	idRaw, _ := c.Get("usuarioID")
	var usuarioID uint
	if val, ok := idRaw.(float64); ok {
		usuarioID = uint(val)
	} else {
		usuarioID = idRaw.(uint)
	}

	usuario, _ := repositorios.BuscarPorID(usuarioID)
	cursos, err := repositorios.BuscarCursosMatriculados(usuarioID)

	if err != nil {
		cursos = []models.Curso{}
	}

	c.HTML(http.StatusOK, "layout", gin.H{
		"title":   "Meus Cursos",
		"page":    "meus_cursos",
		"usuario": usuario,
		"cursos":  cursos,
	})
}

func ExibirPerfil(c *gin.Context) {
	idRaw, _ := c.Get("usuarioID")
	var usuarioID uint
	if val, ok := idRaw.(float64); ok {
		usuarioID = uint(val)
	} else {
		usuarioID = idRaw.(uint)
	}

	usuario, total, err := repositorios.ObterResumoPerfil(usuarioID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	c.HTML(http.StatusOK, "layout", gin.H{
		"title":       "Meu Perfil",
		"page":        "perfil",
		"usuario":     usuario,
		"totalCursos": total,
	})
}

func GerarCertificado(c *gin.Context) {
	cursoSlug := c.Param("slug")
	idRaw, _ := c.Get("usuarioID")
	usuarioID := uint(idRaw.(float64))

	var curso models.Curso
	if err := database.DB.Preload("Modulos.Aulas").Where("slug = ?", cursoSlug).First(&curso).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/feed/")
		return
	}

	var matricula models.Matricula
	if err := database.DB.Where("curso_id = ? AND usuario_id = ? ", curso.ID, usuarioID).First(&matricula).Error; err != nil {
		c.String(http.StatusForbidden, "Você não está matriculado nesse curso.")
		return
	}

	var aulasConcluidasCount int64
	database.DB.Model(&models.ProgressoAula{}).
		Joins("JOIN aulas ON aulas.id = progresso_aulas.aula_id").
		Joins("JOIN modulos ON modulos.id = aulas.modulo_id").
		Where("modulos.curso_id = ? AND progresso_aulas.usuario_id = ?", curso.ID, usuarioID).
		Count(&aulasConcluidasCount)

	totalAulas := 0
	for _, m := range curso.Modulos {
		totalAulas += len(m.Aulas)
	}

	if totalAulas == 0 || int(aulasConcluidasCount) < totalAulas {
		c.String(http.StatusForbidden, "Conclua todas as aulas antes de gerar o certificado.")
		return
	}

	usuario, _ := repositorios.BuscarPorID(usuarioID)

	c.HTML(http.StatusOK, "certificado.html", gin.H{
		"usuario": usuario,
		"curso":   curso,
		"data":    time.Now().Format("02/01/2006"),
	})

}
