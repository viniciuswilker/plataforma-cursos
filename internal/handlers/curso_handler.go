package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/viniciuswilker/plataforma-cursos/internal/database"
	"github.com/viniciuswilker/plataforma-cursos/internal/models"
	"github.com/viniciuswilker/plataforma-cursos/internal/repositorios"
	"github.com/viniciuswilker/plataforma-cursos/internal/services"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RequisicaoCurso struct {
	Titulo    string `json:"titulo" binding:"required"`
	Descricao string `json:"descricao"`
}

func CriarCurso(c *gin.Context) {
	idInterface, _ := c.Get("usuarioID")

	var instrutorID uint
	switch v := idInterface.(type) {
	case float64:
		instrutorID = uint(v)
	case uint:
		instrutorID = v
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

	cursoSlug := slug.Make(input.Titulo)

	curso := models.Curso{
		Titulo:      input.Titulo,
		Slug:        cursoSlug,
		Descricao:   input.Descricao,
		CapaURL:     input.CapaURL,
		InstrutorID: instrutorID,
	}

	if err := database.DB.Create(&curso).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar no banco ou slug já existente"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Curso criado com sucesso",
		"id":      curso.ID,
		"slug":    curso.Slug,
	})
}

func ExcluirCurso(c *gin.Context) {
	cursoID := c.Param("id")
	idRaw, _ := c.Get("usuarioID")

	var instrutorID uint
	switch v := idRaw.(type) {
	case float64:
		instrutorID = uint(v)
	case uint:
		instrutorID = v
	}

	var curso models.Curso
	result := database.DB.Where("id = ? AND instrutor_id = ?", cursoID, instrutorID).First(&curso)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Curso não encontrado"})
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado"})
		}
		return
	}

	if err := database.DB.Delete(&curso).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar exclusão"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Curso movido para a lixeira"})
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
	fmt.Println("--- Iniciando CriarAula ---")

	idRaw, _ := c.Get("usuarioID")
	instrutorID := uint(idRaw.(float64))
	fmt.Printf("[DEBUG] Instrutor ID: %d\n", instrutorID)

	moduloID, _ := strconv.ParseUint(c.PostForm("modulo_id"), 10, 32)
	titulo := c.PostForm("titulo")
	fmt.Printf("[DEBUG] Módulo ID: %d | Titulo: %s\n", moduloID, titulo)

	var modulo models.Modulo
	err := database.DB.Joins("JOIN cursos ON cursos.id = modulos.curso_id").
		Where("modulos.id = ? AND cursos.instrutor_id = ?", uint(moduloID), instrutorID).
		First(&modulo).Error

	if err != nil {
		fmt.Println("[ERRO] Falha na validação de permissão ou módulo inexistente")
		c.JSON(http.StatusForbidden, gin.H{"error": "Acesso negado ao módulo"})
		return
	}

	videoURL := ""
	file, err := c.FormFile("video")

	if err == nil {
		fmt.Printf("[DEBUG] Arquivo recebido: %s (%d bytes)\n", file.Filename, file.Size)

		cld, err := cloudinary.NewFromParams(
			os.Getenv("CLOUDINARY_CLOUD_NAME"),
			os.Getenv("CLOUDINARY_API_KEY"),
			os.Getenv("CLOUDINARY_API_SECRET"),
		)
		if err != nil {
			fmt.Println("[ERRO] Falha ao ler variáveis de ambiente do Cloudinary")
			c.JSON(500, gin.H{"error": "Erro na configuração do storage"})
			return
		}

		openedFile, _ := file.Open()
		if err != nil {
			fmt.Printf("[ERRO] Falha ao abrir arquivo: %v\n", err)
			return
		}
		defer openedFile.Close()

		fmt.Println("[STORAGE] Iniciando upload para o Cloudinary...")
		uploadRes, err := cld.Upload.Upload(c.Request.Context(), openedFile, uploader.UploadParams{
			Folder:       "aulas_plataforma",
			ResourceType: "video",
		})

		if uploadRes.SecureURL != "" {
			videoURL = uploadRes.SecureURL
		} else {
			videoURL = uploadRes.URL
		}
		fmt.Printf("[DEBUG] URL Retornada: %s\n", videoURL)

		if err != nil {
			fmt.Printf("[ERRO] Falha no Cloudinary: %v\n", err)
			c.JSON(500, gin.H{"error": "Falha ao enviar vídeo para nuvem"})
			return
		}

		videoURL = uploadRes.SecureURL
		fmt.Printf("[SUCESSO] Vídeo disponível em: %s\n", videoURL)
	}

	aula := models.Aula{
		ModuloID: uint(moduloID),
		Titulo:   titulo,
		Slug:     slug.Make(titulo),
		VideoURL: videoURL,
	}

	if err := database.DB.Create(&aula).Error; err != nil {
		fmt.Printf("[ERRO] Falha ao persistir aula no banco: %v\n", err)
		c.JSON(500, gin.H{"error": "Erro ao salvar aula"})
		return
	}

	fmt.Printf("[FINAL] Aula '%s' criada com sucesso no ID: %d\n", aula.Titulo, aula.ID)
	c.JSON(http.StatusCreated, aula)
}

func ExcluirAula(c *gin.Context) {

	id := c.Param("id")

	var aula models.Aula
	if err := database.DB.First(&aula, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "Aula não encontrada"})
		return
	}

	if aula.VideoURL != "" {
		if err := os.Remove(aula.VideoURL); err != nil {
			fmt.Printf("Aviso: não foi possivel deletar o arquivo %s: %v\n", aula.VideoURL, err)
		}
	}

	if err := database.DB.Unscoped().Delete(&aula).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Erro ao excluir aula"})
		return
	}

	fmt.Println("Aula excluida com sucesso")
	c.JSON(http.StatusOK, gin.H{"mensagem": "Aula e vídeo excluidos com sucesso"})

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
	slug := c.Param("slug")
	idRaw, _ := c.Get("usuarioID")

	var alunoID uint
	switch v := idRaw.(type) {
	case float64:
		alunoID = uint(v)
	case uint:
		alunoID = v
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Sessão inválida"})
		return
	}

	var curso models.Curso
	if err := database.DB.Where("slug = ?", slug).First(&curso).Error; err != nil {
		c.JSON(404, gin.H{"error": "Curso não encontrado"})
		return
	}

	novaMatricula := models.Matricula{
		CursoID:   curso.ID,
		UsuarioID: alunoID,
	}

	if err := database.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&novaMatricula).Error; err != nil {
		fmt.Println("Erro ao criar matrícula:", err)
		c.JSON(500, gin.H{"error": "Erro ao salvar matrícula"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/cursos/"+curso.Slug+"/assistir")
}

func AssistirCurso(c *gin.Context) {
	cursoSlug := c.Param("slug")
	idRaw, _ := c.Get("usuarioID")
	usuarioID := uint(idRaw.(float64))

	var curso models.Curso
	if err := database.DB.Preload("Modulos.Aulas.Materiais").Where("slug = ?", cursoSlug).First(&curso).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/feed")
		return
	}

	var matricula models.Matricula
	if err := database.DB.Where("curso_id = ? AND usuario_id = ?", curso.ID, usuarioID).First(&matricula).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/curso/preview/"+curso.Slug)
		return
	}

	aulaSlug := c.Query("aula")
	var aulaAtual models.Aula
	if aulaSlug != "" {
		database.DB.Preload("Materiais").Where("slug = ?", aulaSlug).First(&aulaAtual)
	} else if len(curso.Modulos) > 0 && len(curso.Modulos[0].Aulas) > 0 {
		aulaAtual = curso.Modulos[0].Aulas[0]
	}

	concluidasMap := make(map[uint]bool)
	var aulasConcluidasIDs []uint
	database.DB.Model(&models.ProgressoAula{}).Where("usuario_id = ?", usuarioID).Pluck("aula_id", &aulasConcluidasIDs)
	for _, id := range aulasConcluidasIDs {
		concluidasMap[id] = true
	}

	totalAulas := 0
	for _, m := range curso.Modulos {
		totalAulas += len(m.Aulas)
	}

	concluidasNesteCurso := 0
	for _, m := range curso.Modulos {
		for _, a := range m.Aulas {
			if concluidasMap[a.ID] {
				concluidasNesteCurso++
			}
		}
	}

	podeGerarCertificado := totalAulas > 0 && concluidasNesteCurso == totalAulas
	concluida := concluidasMap[aulaAtual.ID]

	usuario, _ := repositorios.BuscarPorID(usuarioID)

	c.HTML(http.StatusOK, "layout", gin.H{
		"title":                "Assistindo: " + curso.Titulo,
		"page":                 "ver_curso_aluno",
		"curso":                curso,
		"aulaAtual":            aulaAtual,
		"usuario":              usuario,
		"concluida":            concluida,
		"concluidasMap":        concluidasMap,
		"podeGerarCertificado": podeGerarCertificado,
		"totalAulas":           totalAulas,
	})
}

func ConcluirAula(c *gin.Context) {
	aulaIDStr := c.Param("id")

	usuarioID, existe := c.Get("usuarioID")
	if !existe {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": "Usuário não identificado"})
		return
	}

	aulaID, _ := strconv.ParseUint(aulaIDStr, 10, 32)

	progresso := models.ProgressoAula{
		UsuarioID: uint(usuarioID.(float64)),
		AulaID:    uint(aulaID),
		Data:      time.Now(),
	}

	err := database.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&progresso).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "Não foi possível salvar o progresso"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "Aula concluída com sucesso!"})
}

func parseID(id string) uint {
	val, _ := strconv.ParseUint(id, 10, 32)
	return uint(val)
}
