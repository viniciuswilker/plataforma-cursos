package repositorios

import (
	"github.com/viniciuswilker/plataforma-cursos/internal/database"
	"github.com/viniciuswilker/plataforma-cursos/internal/models"
)

type StatsProfessor struct {
	TotalAlunos int64
	TotalCursos int64
}

func CriarCurso(curso *models.Curso) error {
	return database.DB.Create(curso).Error
}

func ListarCursos() ([]models.Curso, error) {
	var cursos []models.Curso
	err := database.DB.Preload("Instrutor").Find(&cursos).Error
	return cursos, err
}

func ListarCursosPorInstrutor(instrutorID uint) ([]models.Curso, error) {
	var cursos []models.Curso

	err := database.DB.Table("cursos").
		Select("cursos.*, COUNT(matriculas.id) AS total_alunos").
		Joins("LEFT JOIN matriculas ON matriculas.curso_id = cursos.id").
		Where("cursos.instrutor_id = ? AND cursos.deleted_at IS NULL", instrutorID).
		Group("cursos.id").
		Scan(&cursos).Error

	return cursos, err
}

func CriarAula(aula *models.Aula) error {
	return database.DB.Create(aula).Error
}

func ObterStatsProfessor(instrutorID uint) (StatsProfessor, error) {
	var stats StatsProfessor

	err := database.DB.Table("matriculas").
		Joins("JOIN cursos ON cursos.id = matriculas.curso_id").
		Where("cursos.instrutor_id = ?", instrutorID).
		Count(&stats.TotalAlunos).Error

	if err != nil {
		return stats, err
	}

	err = database.DB.Model(&models.Curso{}).
		Where("instrutor_id = ?", instrutorID).
		Count(&stats.TotalCursos).Error

	return stats, err
}

func BuscarCursosMatriculados(usuarioID uint) ([]models.Curso, error) {
	var cursos []models.Curso

	err := database.DB.Table("cursos").
		Joins("JOIN matriculas ON matriculas.curso_id = cursos.id").
		Preload("Instrutor").
		Where("matriculas.usuario_id = ? AND matriculas.deleted_at IS NULL", usuarioID).
		Find(&cursos).Error

	return cursos, err
}
