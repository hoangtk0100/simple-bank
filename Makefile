PSQL_IMAGE=postgres:14-alpine
DB_CONTAINER_NAME=postgres14
DB_NAME=simple_bank
DB_HOST=localhost
DB_PORT=5433
DB_SOURCE=postgresql://root:secret@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

network:
	docker network create bank-network

postgres:
	docker run --name $(DB_CONTAINER_NAME) --network bank-network -p $(DB_PORT):5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d $(PSQL_IMAGE)

mysql:
	docker run --name mysql8 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8

createdb:
	docker exec -it $(DB_CONTAINER_NAME) createdb --username=root --owner=root $(DB_NAME)

dropdb:
	docker exec -it $(DB_CONTAINER_NAME) dropdb $(DB_NAME)

db:
	docker exec -it $(DB_CONTAINER_NAME) psql -U root $(DB_NAME)

migrate:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_SOURCE)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/hoangtk0100/simple-bank/db/sqlc Store

.PHONY: network postgres createdb dropdb db migrate migrateup migratedown migrateup1 migratedown1 sqlc test server mock
