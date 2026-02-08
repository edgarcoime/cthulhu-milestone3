.PHONY: dev build lint clean test

dev: build
	./bin/service

build:
	go mod tidy
	go build -o bin/service ./cmd/service

lint:
	golangci-lint run

clean:
	rm -rf tmp/ bin/

test:
	go test ./... -v
