CGO_ENABLED=1
CC=x86_64-w64-mingw32-cc
OUTPUT_WINDOWS=./bin/Installer-Windows.exe
OUTPUT_LINUX=./bin/Installer-Linux
THEME=dark

windows:
	@echo "Сборка проекта для Windows..."
	@CGO_ENABLED=$(CGO_ENABLED) CC=$(CC) GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o $(OUTPUT_WINDOWS) ./src/*.go

linux:
	@echo "Сборка проекта для Linux..."
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(OUTPUT_LINUX) ./src/*.go

all:
	@echo "Сборка проекта для Windows..."
	@CGO_ENABLED=$(CGO_ENABLED) CC=$(CC) GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H=windowsgui" -o $(OUTPUT_WINDOWS) ./src/*.go
	@echo "Сборка проекта для Linux..."
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(OUTPUT_LINUX) ./src/*.go

clean:
	@echo "Очистка..."
	@rm -f $(OUTPUT_WINDOWS) $(OUTPUT_LINUX)

help:
	@echo "Используйте 'make windows' для сборки Windows версии."
	@echo "Используйте 'make linux' для сборки Linux версии."
	@echo "Используйте 'make clean' для очистки."

.PHONY: all windows linux clean help
