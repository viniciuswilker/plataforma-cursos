package services

import (
	"github.com/viniciuswilker/plataforma-cursos/internal/models"
	"github.com/viniciuswilker/plataforma-cursos/internal/repositorios"
)

func CriarNovoCurso(titulo, descricao string, instrutorID uint) error {
	curso := models.Curso{
		Titulo:      titulo,
		Descricao:   descricao,
		InstrutorID: instrutorID,
	}
	return repositorios.CriarCurso(&curso)

}

func BuscarTodosCursos() ([]models.Curso, error) {
	return repositorios.ListarCursos()
}

func AdicionarAula(cursoID uint, titulo, conteudo string, ordem int) error {
	aula := models.Aula{
		CursoID:  cursoID,
		Titulo:   titulo,
		Conteudo: conteudo,
		Ordem:    ordem,
	}
	return repositorios.CriarAula(&aula)
}
