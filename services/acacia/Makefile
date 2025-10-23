DB_URL=postgres://postgres:password@localhost:5432/acacia?sslmode=disable

.PHONY: build run dev sqlc-generate test clean

# Build API server
build:
	go build -o bin/acacia ./main.go

# Run API server
run: build
	./bin/acacia

# Development mode for API server
dev:
	go run main.go

sqlc-generate:
	sqlc generate

test:
	go test ./...

clean:
	rm -rf bin/