package database

import (
	"log"

	"github.com/glebarez/sqlite"
	"github.com/viniciuswilker/plataforma-cursos/internal/models"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Conectar() {
	var err error

	DB, err = gorm.Open(sqlite.Open("lms.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Falha ao conectar com o banco de dados: ", err)

	}

	DB.AutoMigrate(
		&models.Usuario{},
		&models.Curso{},
		&models.Aula{},
		&models.ProgressoAula{},
	)

}
