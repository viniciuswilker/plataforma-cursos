package database

import (
	"log"
	"os"
	"strings"

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

	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		dialector = postgres.Open(dsn)
	} else {

		if dsn == "" {
			dsn = "data/database.db"
		}
		dialector = sqlite.Open(dsn)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}

	log.Println("Banco de dados conectado com sucesso!")
	DB.AutoMigrate(
		&models.Usuario{},
		&models.Curso{},
		&models.Aula{},
		&models.ProgressoAula{},
		&models.Material{},
		&models.Matricula{},
	)
}
