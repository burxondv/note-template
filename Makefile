POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=12345
POSTGRES_DATABASE=note_db

-include .env

DB_URL="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable"

print:
	echo "$(DB_URL)"

swag-init:
	swag init -g api/api.go -o api/docs

start:
	go run main.go

migrateup:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path migrations -database "$(DB_URL)" -verbose down

.PHONY: start migrateup migratedown