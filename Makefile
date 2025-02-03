include .env

POSTGRES_URL='postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=${POSTGRES_SSL_MODE}'
POSTGRES_MIGRATIONS='./migrations/postgres'

all: compose_up migrate_up run

compose_up:
	@docker-compose up -d

compose_down:
	@docker-compose down

migrate_up:
	@migrate -path $(POSTGRES_MIGRATIONS) -database $(POSTGRES_URL) up $(n)

migrate_down:
	@migrate -path $(POSTGRES_MIGRATIONS) -database $(POSTGRES_URL) down $(n)

run:
	@go run ./cmd/app/main.go

mockgen:
	@mockgen -source $(SRC) -destination $(DEST) -package $(PKG)
