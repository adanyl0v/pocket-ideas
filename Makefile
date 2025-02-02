include .env

POSTGRES_URL='postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=${POSTGRES_SSL_MODE}'
POSTGRES_MIGRATION_SOURCE='file:./migrations/postgres'

all: compose_up migrate_up run

compose_up:
	@docker-compose up -d

compose_down:
	@docker-compose down

NUP=1
migrate_up:
	@migrate -source $(POSTGRES_MIGRATION_SOURCE) -database $(POSTGRES_URL) up $(NUP)

NDOWN=1
migrate_down:
	@migrate -source $(POSTGRES_MIGRATION_SOURCE) -database $(POSTGRES_URL) down $(NDOWN)

run:
	@go run ./cmd/app/main.go

mockgen:
	@mockgen -source $(SRC) -destination $(DEST) -package $(PKG)
