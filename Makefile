include app.env
export

up:
	docker compose up --build --detach

down:
	docker compose down

network:
	docker network create $(NETWORK_NAME)

postgres:
	docker run --name $(DB_CONTAINER_NAME) --network $(NETWORK_NAME) -p $(DB_PORT):5432 -e POSTGRES_USER=$(DB_USERNAME) -e POSTGRES_PASSWORD=$(DB_PASSWORD) -d $(PSQL_IMAGE)

mysql:
	docker run --name mysql8 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -d mysql:8

createdb:
	docker exec -it $(DB_CONTAINER_NAME) createdb --username=$(DB_USERNAME) --owner=$(DB_USERNAME) $(DB_NAME)

dropdb:
	docker exec -it $(DB_CONTAINER_NAME) dropdb $(DB_NAME)

db:
	docker exec -it $(DB_CONTAINER_NAME) psql -U $(DB_USERNAME) $(DB_NAME)

new_migration:
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
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/hoangtk0100/simple-bank/db/sqlc Store
	mockgen -package mockwk -destination worker/mock/distributor.go github.com/hoangtk0100/simple-bank/worker TaskDistributor

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
	docker run --name redis -e REDIS_PASSWORD=$(REDIS_PASSWORD) -p 6379:6379 -d redis:7-alpine --requirepass $(REDIS_PASSWORD)

.PHONY: up down network postgres createdb dropdb db new_migration migrateup migratedown migrateup1 migratedown1 db_docs db_schema  sqlc test server mock proto evans redis
