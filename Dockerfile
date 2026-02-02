# Estágio de Build
FROM golang:1.24-bookworm AS builder

# Instalando dependências para SQLite e compilação
RUN apt-get update && apt-get install -y gcc musl-dev

WORKDIR /app

# 1. Copia apenas os módulos para aproveitar o cache do Render
COPY go.mod go.sum ./
RUN go mod download

# 2. Copia o código fonte (Certifique-se que .dockerignore NÃO ignora static e templates)
COPY . .

# 3. Debug: Lista os arquivos para o log do Render (ajuda muito a debugar se a pasta sumiu)
RUN ls -la /app

# 4. Build garantindo CGO para o SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/api/main.go

# --- Estágio Final (Runtime) ---
FROM debian:bookworm-slim

WORKDIR /app

# Instala certificados (Obrigatório para o Render acessar bancos externos/HTTPS)
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copia o binário do estágio builder
COPY --from=builder /app/main .

# 5. Copia os assets garantindo que o build não quebre se estiverem vazios
# Se o erro de 'not found' persistir aqui, o problema é o seu .gitignore ou .dockerignore
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/static /app/static

# Cria as pastas necessárias para persistência (especialmente importante para SQLite/Uploads)
RUN mkdir -p /app/data /app/uploads/videos

# Porta dinâmica que o Render fornece
EXPOSE 8080

# Executa o binário
CMD ["./main"]