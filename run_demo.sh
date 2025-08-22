#!/bin/bash

echo "🔨 Компиляция демонстрации таймаутов..."

# Очищаем кэш и пересобираем
go clean -cache
go clean -modcache
go mod tidy
go build -a -o demo cmd/main.go

if [ $? -eq 0 ]; then
    echo "✅ Компиляция успешна! Запускаем демонстрацию..."
    echo ""
    ./demo
else
    echo "❌ Ошибка компиляции!"
    exit 1
fi
