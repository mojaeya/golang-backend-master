createdb:
	docker exec -it golang-backend-master-postgres-1 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it golang-backend-master-postgres-1 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createdb dropdb migrateup migratedown sqlc test server