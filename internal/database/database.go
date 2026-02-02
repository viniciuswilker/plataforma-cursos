package database

import (
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/viniciuswilker/plataforma-cursos/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Conectar() {
	var err error
	var dialector gorm.Dialector

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		dialector = sqlite.Open("data/database.db")
	} else {
		dialector = postgres.Open(dsn)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}

	DB.AutoMigrate(
		&models.Usuario{},
		&models.Curso{},
		&models.Aula{},
		&models.ProgressoAula{},
		&models.Material{},
		&models.Matricula{},
	)
}
