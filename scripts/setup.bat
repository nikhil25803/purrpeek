@echo off
setlocal

cd /d "%~dp0.."

where go >nul 2>&1
if errorlevel 1 (
    echo Error: Go is required but was not found in PATH. 1>&2
    exit /b 1
)

echo Downloading dependencies...
go mod download || exit /b 1

echo Running tests...
go test ./... || exit /b 1

echo Building purrpeek...
if not exist bin mkdir bin
go build -o bin\purrpeek.exe ./cmd/purrpeek || exit /b 1

echo Ready: bin\purrpeek.exe
