network:
	docker network create campin-network

postgres:
	docker run --name db --network campin-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secure -d postgres:alpine

createdb:
	docker exec -it db createdb --username=root --owner=root campin_db

dropdb:
	docker exec -it db dropdb campin_db

migrateup:
	migrate -path db/migration -database "postgres://root:secure@localhost:5432/campin_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:secure@localhost:5432/campin_db?sslmode=disable" -verbose down

server:
	go run main.go

.PHONY: network, postgres, createdb, dropdb, migrateup, migratedown