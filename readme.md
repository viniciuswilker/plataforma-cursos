🚀 EduStream - Plataforma de Cursos Online
O EduStream é uma aplicação Full Stack de streaming de vídeo educacional desenvolvida em Go (Golang). O projeto foi desenhado seguindo princípios de arquitetura limpa e conteinerização, garantindo que a aplicação seja escalável e fácil de implantar em qualquer ambiente de nuvem.

🛠️ Tecnologias Utilizadas
Backend: Go (Golang) com o framework Gin Gonic.

Banco de Dados: PostgreSQL (Produção/Render) e SQLite (Desenvolvimento).

ORM: GORM para abstração de banco de dados e migrações automáticas.

Frontend: HTML5, CSS3 e JavaScript.

Infraestrutura: Docker e Docker Compose para isolamento de ambiente.

Deploy: Render (CI/CD via GitHub).

Autenticação: JWT (JSON Web Tokens) com armazenamento seguro em Cookies HttpOnly.

✨ Funcionalidades Principais
Autenticação Robusta: Sistema de login e cadastro com diferenciação de cargos (Admin/Professor/Aluno).

Filtro de Conteúdo: Barra de busca dinâmica que realiza queries LIKE no banco de dados para encontrar cursos em tempo real.

Arquitetura Multi-Ambiente: Configuração inteligente que alterna entre SQLite e PostgreSQL dependendo do ambiente (local ou cloud).

Streaming de Vídeo: Upload e entrega de vídeos otimizada através de rotas estáticas do Gin.

Middleware de Autorização: Proteção de rotas que redireciona usuários não autenticados e bloqueia acessos indevidos conforme o cargo.

Como rodar o projeto (Docker)
Certifique-se de ter o Docker instalado e execute:

Bash

# Clone o repositório
git clone https://github.com/viniciuswilker/plataforma-cursos.git

# Entre na pasta
cd plataforma-curso

# Suba os containers
make up
A aplicação estará disponível em http://localhost:8080.

