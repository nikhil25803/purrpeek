.PHONY: run-purrpeek build-purrpeek build-windows test nix-check nix-build aur-check aur-build

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

aur-check:
	docker run --rm -it \
		--platform linux/amd64 \
		-v "$$(pwd)/purrpeek-bin:/pkg" \
		archlinux:latest \
		bash -c "pacman --disable-sandbox-syscalls -Sy --noconfirm base-devel && useradd -m builder && chown -R builder:builder /pkg && su builder -c 'cd /pkg && makepkg --syncdeps --noconfirm -f'"


aur-build:
	docker run --rm \
		--platform linux/amd64 \
		-v "$$(pwd)/purrpeek-bin:/pkg" \
		archlinux:latest \
		bash -c "pacman --disable-sandbox-syscalls -Sy --noconfirm base-devel && useradd -m builder && chown -R builder:builder /pkg && su builder -c 'cd /pkg && makepkg --syncdeps --noconfirm -f'"