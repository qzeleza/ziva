# Основные переменные для управления сборкой
BINARY_NAME:=ziva
BUILD_DIR:=bin
MAIN_PKG:=./cmd/ziva

.PHONY: help run build test test-coverage clean deps fmt lint cache

help:
	@echo "Доступные команды:"
	@echo "  make help            показать это сообщение"
	@echo "  make run             запустить приложение"
	@echo "  make build           собрать приложение"
	@echo "  make test            запустить тесты"
	@echo "  make test-coverage   запустить тесты с покрытием"
	@echo "  make clean           очистить сборку"
	@echo "  make deps            установить зависимости"
	@echo "  make fmt             форматировать код"
	@echo "  make lint            проверить код линтером"
	@echo "  make cache           очистить кэш"

run:
	go run $(MAIN_PKG)

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PKG)

test:
	go test ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	@echo "Отчёт сохранён в coverage.out"

clean:
	rm -rf $(BUILD_DIR) coverage.out

deps:
	go mod download

fmt:
	go fmt ./...

lint:
	go vet ./...

cache:
	go clean -cache -testcache
