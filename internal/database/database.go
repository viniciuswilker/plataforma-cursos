package database

import (
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/viniciuswilker/plataforma-cursos/internal/models"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Conectar() {
	var err error

	dbPath := os.Getenv("DATABASE_URL")
	if dbPath == "" {
		dbPath = "database.db"
	}

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

	if err != nil {
		log.Fatal("Falha ao conectar com o banco de dados: ", err)
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
