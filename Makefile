.PHONY: run-purrpeek build-purrpeek build-windows test nix-check nix-build

run-purrpeek:
	@echo "Running purrpeek..."
	go run cmd/purrpeek/main.go

build-purrpeek:
	@echo "Building purrpeek..."
	go build -o bin/purrpeek cmd/purrpeek/main.go

build-windows:
	@echo "Building purrpeek for Windows..."
	GOOS=windows GOARCH=amd64 go build -o bin/purrpeek.exe ./cmd/purrpeek

test:
	go test ./...

nix-check:
	nix flake check --all-systems

nix-build:
	nix build