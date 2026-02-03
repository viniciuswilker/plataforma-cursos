🚀 EduStream - Plataforma de Cursos Online
O EduStream é uma aplicação Full Stack de streaming de vídeo educacional desenvolvida em Go (Golang). O projeto foi desenhado seguindo princípios de arquitetura limpa e conteinerização, garantindo que a aplicação seja escalável e fácil de implantar em qualquer ambiente de nuvem.

🛠️ Tecnologias Utilizadas
Backend: Go (Golang) com framework Gin Gonic.

Banco de Dados: PostgreSQL (Produção no Render) e SQLite (Desenvolvimento).

ORM: GORM para abstração de dados e migrações automáticas.

Storage de Mídia: Cloudinary (Armazenamento de vídeos na nuvem).

Frontend: HTML5, CSS3, JavaScript e Google Material Symbols.

Autenticação: JWT (JSON Web Tokens) com Cookies HttpOnly.

Infraestrutura: Docker & Docker Compose.

Autenticação: JWT (JSON Web Tokens) com armazenamento seguro em Cookies HttpOnly.

✨ Funcionalidades Principais
Gestão de Cursos (Slugs): URLs amigáveis para melhor SEO e navegação.

Streaming na Nuvem: Integração com Cloudinary para entrega de vídeos escalável e persistente.

Sistema de Progresso: Acompanhamento de aulas concluídas e visualização de barra de progresso em tempo real.

Emissão de Certificados: Geração dinâmica de certificados de conclusão após atingir 100% do curso.

Autenticação Hierárquica: Níveis de acesso distintos para Alunos, Professores e Administradores.



# Clone o repositório
git clone https://github.com/viniciuswilker/plataforma-cursos.git

# Entre na pasta
cd plataforma-curso

# Suba os containers
make up
A aplicação estará disponível em http://localhost:8080.

