package services

import (
	"errors"

	"github.com/viniciuswilker/plataforma-cursos/internal/models"
	"github.com/viniciuswilker/plataforma-cursos/internal/repositorios"
	"golang.org/x/crypto/bcrypt"
)

func CadastrarUsuario(nome, sobrenome, email, senha string) error {

	_, err := repositorios.BuscarPorEmail(email)
	if err != nil {
		return errors.New("este e-mail já esta cadastrado")
	}

	senhaHash, _ := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	novoUsuario := models.Usuario{
		Nome:      nome,
		Sobrenome: sobrenome,
		Email:     email,
		Senha:     string(senhaHash),
		Cargo:     "aluno",
	}

	return repositorios.CriarUsuario(&novoUsuario)

}

func Login(email, senha string) (string, error) {
	usuario, err := repositorios.BuscarPorEmail(email)

	if err != nil {
		return "", errors.New("credenciais inválidas")
	}

	err = bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(senha))
	if err != nil {
		return "", errors.New("credenciais inválidas")
	}

	return GerarToken(usuario.ID, usuario.Cargo)
}
