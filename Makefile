.PHONY: run-purrpeek build-purrpeek


run-purrpeek:
	@echo "Running purrpeek..."
	go run cmd/purrpeek/main.go


build-purrpeek:
	@echo "Building purrpeek..."
	go build -o bin/purrpeek cmd/purrpeek/main.go