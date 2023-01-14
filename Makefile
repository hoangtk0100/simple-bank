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

db_docs:
	dbdocs build docs/db.dbml

db_schema:
	dbml2sql --posgres -o docs/schema.sql docs/db.dbml

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/hoangtk0100/simple-bank/db/sqlc Store

proto:
	rm -rf pb/*.go
	rm -rf docs/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto
	statik -src=./docs/swagger -dest=./docs

evans:
	evans --host localhost --port 9099 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: network postgres createdb dropdb db migrate migrateup migratedown migrateup1 migratedown1 db_docs db_schema  sqlc test server mock proto evans redis
