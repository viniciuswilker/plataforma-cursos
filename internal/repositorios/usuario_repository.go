package repositorios

import (
	"github.com/viniciuswilker/plataforma-cursos/internal/database"
	"github.com/viniciuswilker/plataforma-cursos/internal/models"
)

func CriarUsuario(usuario *models.Usuario) error {
	return database.DB.Create(usuario).Error
}

func BuscarPorEmail(email string) (*models.Usuario, error) {
	var usuario models.Usuario

	err := database.DB.Where("email = ?", email).First(&usuario).Error
	return &usuario, err
}

func BuscarPorID(id interface{}) (*models.Usuario, error) {
	var usuario models.Usuario
	err := database.DB.First(&usuario, id).Error
	return &usuario, err
}

func ObterResumoPerfil(usuarioID uint) (models.Usuario, int64, error) {
	var usuario models.Usuario
	var totalCursos int64

	if err := database.DB.First(&usuario, usuarioID).Error; err != nil {
		return usuario, 0, err
	}

	database.DB.Model(&models.Matricula{}).Where("usuario_id = ?", usuarioID).Count(&totalCursos)

	return usuario, totalCursos, nil
}
