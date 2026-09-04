.PHONY: run-purrpeek build-purrpeek test


run-purrpeek:
	@echo "Running purrpeek..."
	go run cmd/purrpeek/main.go


build-purrpeek:
	@echo "Building purrpeek..."
	go build -o bin/purrpeek cmd/purrpeek/main.go

test:
	go test ./...
