package models

import (
	"time"

	"gorm.io/gorm"
)

type Usuario struct {
	gorm.Model
	Nome      string `gorm:"size:100;not null"`
	Sobrenome string `gorm:"size:100;not null"`
	Email     string `gorm:"size:100;unique;not null"`
	Senha     string `gorm:"not null"`
	Cargo     string `gorm:"default:'aluno'"`
}

type Curso struct {
	gorm.Model
	Titulo      string `gorm:"size:200;not null"`
	Descricao   string `gorm:"type:text"`
	CapaURL     string
	InstrutorID uint
	Instrutor   Usuario  `gorm:"foreignKey:InstrutorID"`
	Modulos     []Modulo `gorm:"constraint:OnDelete:CASCADE;"`
	TotalAlunos int64    `gorm:"column:total_alunos"`
}

type Modulo struct {
	gorm.Model
	CursoID uint
	Titulo  string `gorm:"size:200;not null"`
	Ordem   int
	Aulas   []Aula `gorm:"constraint:OnDelete:CASCADE;"`
}

type Aula struct {
	gorm.Model
	ModuloID  uint
	Titulo    string `gorm:"size:200;not null"`
	VideoURL  string
	Conteudo  string `gorm:"type:text"`
	Ordem     int
	Materiais []Material `gorm:"constraint:OnDelete:CASCADE;"`
}

type Material struct {
	gorm.Model
	AulaID     uint
	Nome       string
	ArquivoURL string
}

type Matricula struct {
	gorm.Model
	UsuarioID uint
	CursoID   uint
	Usuario   Usuario
	Curso     Curso
	Concluido bool `gorm:"default:false"`
}

type ProgressoAula struct {
	UsuarioID uint `gorm:"primaryKey"`
	AulaID    uint `gorm:"primaryKey"`
	Data      time.Time
}
