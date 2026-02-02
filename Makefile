APP_NAME=plataforma-cursos
DOCKER_COMPOSE=docker-compose

.PHONY: build run test docker-up docker-down clean

build:
	go build -o bin/$(APP_NAME) ./cmd/main.go

run:
	go run ./cmd/main.go

up:
	$(DOCKER_COMPOSE) up --build

down:
	$(DOCKER_COMPOSE) down

clean:
	rm -rf bin/
	rm -rf data/database.db