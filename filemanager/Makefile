.PHONY: dev lint clean test sqlc

dev:
	go mod tidy
	go run ./cmd/service

lint:
	golangci-lint run

clean:
	rm -rf tmp/

test:
	go test ./... -v

sqlc:
	cd internal/repository/sqlc && sqlc generate