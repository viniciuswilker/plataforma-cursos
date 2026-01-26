package models

import "gorm.io/gorm"

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
	Instrutor   Usuario `gorm:"foreignKey:InstrutorID"`
	Aulas       []Aula  `gorm:"constraint:OnDelete:CASCADE;"`
}

type Aula struct {
	gorm.Model
	CursoID  uint
	Titulo   string `gorm:"size:200;not null"`
	Conteudo string `gorm:"type:text"`
	Ordem    int
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
	gorm.Model
	UsuarioID uint
	AulaID    uint
	Concluido bool `gorm:"default:true"`
}
