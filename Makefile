PSQL_IMAGE=postgres:14-alpine
DB_CONTAINER_NAME=postgres14
DB_NAME=simple_bank
DB_PORT=5433
DB_URL=postgresql://root:secret@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable

postgres:
	lima nerdctl run --name $(DB_CONTAINER_NAME) -p $(DB_PORT):5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d $(PSQL_IMAGE)

mysql:
	lima nerdctl run --name mysql8 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8

createdb:
	lima nerdctl exec -it $(DB_CONTAINER_NAME) createdb --username=root --owner=root $(DB_NAME)

dropdb:
	lima nerdctl exec -it $(DB_CONTAINER_NAME) dropdb $(DB_NAME)

db:
	lima nerdctl exec -it $(DB_CONTAINER_NAME) psql -U root $(DB_NAME)

migrate:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/hoangtk0100/simple-bank/db/sqlc Store

.PHONY: postgres createdb dropdb db migrate migrateup migratedown migrateup1 migratedown1 sqlc test server mock
