DB_URL=postgres://postgres:password@localhost:5432/acacia?sslmode=disable

.PHONY: build run dev sqlc-generate migrate-up migrate-down test clean

build:
	go build -o bin/acacia ./main.go

run: build
	./bin/acacia

dev:
	go run main.go

sqlc-generate:
	sqlc generate

test:
	go test ./...

clean:
	rm -rf bin/