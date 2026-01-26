package repositorios

import (
	"github.com/viniciuswilker/plataforma-cursos/internal/database"
	"github.com/viniciuswilker/plataforma-cursos/internal/models"
)

func CriarCurso(curso *models.Curso) error {
	return database.DB.Create(curso).Error
}

func ListarCursos() ([]models.Curso, error) {
	var cursos []models.Curso
	err := database.DB.Preload("Instrutor").Find(&cursos).Error
	return cursos, err
}

func CriarAula(aula *models.Aula) error {
	return database.DB.Create(aula).Error
}
