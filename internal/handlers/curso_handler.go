package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/viniciuswilker/plataforma-cursos/internal/database"
	"github.com/viniciuswilker/plataforma-cursos/internal/models"
	"github.com/viniciuswilker/plataforma-cursos/internal/repositorios"
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

func CriarModulo(c *gin.Context) {
	idRaw, _ := c.Get("usuarioID")
	var instrutorID uint
	if val, ok := idRaw.(float64); ok {
		instrutorID = uint(val)
	} else {
		instrutorID = idRaw.(uint)
	}

	var input struct {
		CursoID uint   `json:"curso_id" binding:"required"`
		Titulo  string `json:"titulo" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	var curso models.Curso
	if err := database.DB.Where("id = ? AND instrutor_id = ?", input.CursoID, instrutorID).First(&curso).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Você não tem permissão para alterar este curso"})
		return
	}

	novoModulo := models.Modulo{
		CursoID: input.CursoID,
		Titulo:  input.Titulo,
		Ordem:   0,
	}

	if err := database.DB.Create(&novoModulo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar módulo"})
		return
	}

	c.JSON(http.StatusCreated, novoModulo)
}

func ExcluirModulo(c *gin.Context) {
	moduloID := c.Param("id")
	idRaw, _ := c.Get("usuarioID")
	instrutorID := uint(idRaw.(float64))

	var modulo models.Modulo
	err := database.DB.Joins("JOIN cursos ON cursos.id = modulos.curso_id").
		Where("modulos.id = ? AND cursos.instrutor_id = ?", moduloID, instrutorID).
		First(&modulo).Error

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Você não tem permissão para excluir este módulo"})
		return
	}

	if err := database.DB.Delete(&modulo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao excluir módulo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Módulo e aulas removidos com sucesso"})
}

func CriarAula(c *gin.Context) {

	idRaw, _ := c.Get("usuarioID")
	instrutorID := uint(idRaw.(float64))

	moduloID, _ := strconv.ParseUint(c.PostForm("modulo_id"), 10, 32)
	titulo := c.PostForm("titulo")

	var modulo models.Modulo
	err := database.DB.Joins("JOIN cursos ON cursos.id = modulos.curso_id").
		Where("modulos.id = ? AND cursos.instrutor_id = ?", uint(moduloID), instrutorID).
		First(&modulo).Error

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado ao módulo"})
		return
	}

	videoPath := ""
	file, err := c.FormFile("video")
	if err == nil {
		videoPath = "uploads/videos/" + strconv.FormatInt(time.Now().Unix(), 10) + "_" + file.Filename
		if err := c.SaveUploadedFile(file, videoPath); err != nil {
			c.JSON(500, gin.H{"error": "Falha ao salvar vídeo"})
			return
		}
	}

	aula := models.Aula{
		ModuloID: uint(moduloID),
		Titulo:   titulo,
		VideoURL: videoPath,
	}

	database.DB.Create(&aula)
	c.JSON(http.StatusCreated, aula)
}

func AssistirCursoAdm(c *gin.Context) {
	cursoID := c.Param("id")
	aulaIDQuery := c.Query("aula")

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
		c.Redirect(http.StatusSeeOther, "/professor/cursos/")
		return
	}

	var aulaAtual models.Aula
	if aulaIDQuery != "" {
		database.DB.Preload("Materiais").First(&aulaAtual, aulaIDQuery)
	} else if len(curso.Modulos) > 0 && len(curso.Modulos[0].Aulas) > 0 {
		aulaAtual = curso.Modulos[0].Aulas[0]
		database.DB.Model(&aulaAtual).Association("Materiais").Find(&aulaAtual.Materiais)
	}

	c.HTML(http.StatusOK, "layout", gin.H{
		"title":     "Assistindo: " + curso.Titulo,
		"page":      "ver_curso_adm",
		"curso":     curso,
		"aulaAtual": aulaAtual,
		"usuario":   usuario,
	})
}

func MatricularAluno(c *gin.Context) {
	cursoID := c.Param("id")
	idRaw, _ := c.Get("usuarioID")

	var alunoID uint
	if val, ok := idRaw.(float64); ok {
		alunoID = uint(val)
	} else {
		alunoID = idRaw.(uint)
	}

	var matriculaExistente models.Matricula
	err := database.DB.Where("curso_id = ? AND usuario_id = ?", cursoID, alunoID).First(&matriculaExistente).Error

	if err == nil {
		c.Redirect(http.StatusSeeOther, "/cursos/"+cursoID+"/assistir")
		return
	}

	novaMatricula := models.Matricula{
		CursoID:   uint(parseID(cursoID)),
		UsuarioID: alunoID,
	}

	if err := database.DB.Create(&novaMatricula).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao realizar matrícula"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/cursos/"+cursoID+"/assistir")
}

func AssistirCurso(c *gin.Context) {
	cursoID := c.Param("id")
	idRaw, _ := c.Get("usuarioID")
	usuarioID := uint(idRaw.(float64))

	var matricula models.Matricula
	if err := database.DB.Where("curso_id = ? AND usuario_id = ?", cursoID, usuarioID).First(&matricula).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/cursos/"+cursoID+"/")
		return
	}

	usuario, _ := repositorios.BuscarPorID(usuarioID)

	var curso models.Curso
	database.DB.Preload("Modulos.Aulas.Materiais").First(&curso, cursoID)

	aulaID := c.Query("aula")
	var aulaAtual models.Aula
	if aulaID != "" {
		database.DB.Preload("Materiais").First(&aulaAtual, aulaID)
	} else if len(curso.Modulos) > 0 && len(curso.Modulos[0].Aulas) > 0 {
		aulaAtual = curso.Modulos[0].Aulas[0]
	}

	c.HTML(http.StatusOK, "layout", gin.H{
		"title":     "Assistindo: " + curso.Titulo,
		"page":      "ver_curso_aluno",
		"curso":     curso,
		"aulaAtual": aulaAtual,
		"usuario":   usuario,
	})
}

func parseID(id string) uint {
	val, _ := strconv.ParseUint(id, 10, 32)
	return uint(val)
}
