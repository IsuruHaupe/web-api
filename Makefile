postgres: 
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root web-api

dropdb:
	docker exec -it postgres12 dropdb web-api

migrateforce: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/web-api?sslmode=disable" force 1

migrateup: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/web-api?sslmode=disable" -verbose up

migratedown: 
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/web-api?sslmode=disable" -verbose down

sqlc: 
	sqlc generate

tests: 
	go test -v -cover ./...

server: 
	go run main.go

mock: 
	mockgen -package mockdb -destination db/mock/database.go github.com/IsuruHaupe/web-api/db/postgres Database

swagger:
	swag init

.PHONY: postgres createdb dropdb migrateforce migrateup migratedown sqlc tests server mock swagger