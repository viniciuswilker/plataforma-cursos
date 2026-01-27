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
