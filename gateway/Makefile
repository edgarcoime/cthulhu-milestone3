.PHONY: dev lint clean test

dev:
	go mod tidy
	go run ./cmd/service

lint:
	golangci-lint run

clean:
	rm -rf tmp/

test:
	go test ./... -v
