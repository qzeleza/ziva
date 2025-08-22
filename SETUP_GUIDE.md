# Инструкция по настройке Термос как отдельного модуля

## 1. Инициализация Git репозитория

Выполните следующие команды в директории `/Users/samovar/Documents/develop/go/termos`:

```bash
cd /Users/samovar/Documents/develop/go/termos

# Инициализация Git репозитория
git init

# Добавление всех файлов
git add .

# Первый коммит
git commit -m "Initial commit: Термос TUI framework

- Complete terminal user interface framework
- Support for Yes/No, Input, SingleSelect, MultiSelect tasks
- Performance optimizations for embedded devices
- Comprehensive styling and theming system
- Memory-efficient string utilities
- Full Bubble Tea integration"

# Создание основной ветки (если нужно)
git branch -M main
```

## 2. Связывание с GitHub репозиторием

```bash
# Добавление remote origin (замените на ваш GitHub репозиторий)
git remote add origin https://github.com/qzeleza/termos.git

# Первый push
git push -u origin main
```

## 3. Создание тегов для версионирования

```bash
# Создание первого релиза
git tag -a v1.0.0 -m "Release v1.0.0: Stable Термос TUI framework"
git push origin v1.0.0
```

## 4. Проверка модуля

```bash
# Проверка зависимостей
go mod tidy

# Проверка компиляции
go build ./...

# Запуск примера
go run examples/basic_usage.go
```

## 5. Использование в других проектах

После публикации на GitHub, другие проекты смогут использовать модуль:

```bash
go get github.com/qzeleza/termos
```

```go
import "github.com/qzeleza/termos/task"

// Использование YesNoTask (теперь только 2 опции: "Да" и "Нет")
yesNo := task.NewYesNoTask("Продолжить?", "Подтверждение")
```

## 6. Обновление KvasPro проекта

После создания GitHub репозитория обновите `kvaspro/backend/go.mod`:

```go
require (
    github.com/qzeleza/termos v1.0.0
    // ... остальные зависимости
)
```

И замените все импорты с:
```go
"github.com/qzeleza/termos/..."
```

на:
```go
"github.com/qzeleza/termos/..."
```

## Структура готового модуля

```
termos/
├── go.mod                    # Go модуль
├── go.sum                    # (создается автоматически)
├── README.md                 # Документация
├── LICENSE                   # MIT лицензия
├── SETUP_GUIDE.md           # Эта инструкция
├── common/                   # Общие интерфейсы
│   ├── layout.go
│   └── task.go
├── task/                     # Основные компоненты задач
│   ├── base.go
│   ├── defaults.go
│   ├── input_task_new.go
│   ├── multiselect_task.go
│   ├── singleselect_task.go
│   └── yesno_task.go
├── ui/                       # UI компоненты и стили
│   └── styles.go
├── performance/              # Оптимизации производительности
│   └── string_utils.go
└── examples/                 # Примеры использования
    └── basic_usage.go
```

Модуль готов к использованию! 🎉