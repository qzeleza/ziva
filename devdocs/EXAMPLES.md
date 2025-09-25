# Примеры Жива

Практические примеры для встроенных задач. Все примеры соответствуют текущему API и привязкам клавиш.

## Запуск готовых примеров

В каталоге `cmd/` находятся готовые исполняемые примеры:

```bash
# Базовый пример использования задач
go run ./cmd/basic_usage

# Детектор embedded окружения
go run ./cmd/detector

# Демонстрация embedded конфигураций
go run ./cmd/embedded_config

# Полнофункциональный TUI с очередью задач
go run ./cmd/tui_full
```

- Конструкторы:
  - `task.NewYesNoTask(question, description string) *YesNoTask`
  - `task.NewSingleSelectTask(prompt string, options []string) *SingleSelectTask`
  - `task.NewMultiSelectTask(prompt string, options []string) *MultiSelectTask`
  - `task.NewInputTaskNew(prompt, placeholder string) *InputTaskNew`
  - `task.NewFuncTaskWithOptions(title string, action func() error, options ...task.FuncTaskOption) *FuncTask`

Примечание: Встроенный `Run()` сейчас возвращает `nil`. Интегрируйте задачи в модель Bubble Tea и управляйте ими через `Update`.

## Минимальная интеграция Bubble Tea для одной задачи

```go
package main

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/qzeleza/ziva/common"
    "github.com/qzeleza/ziva/task"
)

type model struct {
    t common.Task
    width int
}

func newModel() model {
    // Замените на любую задачу для демонстрации
    return model{t: task.NewYesNoTask("Продолжить?", "Подтвердите ваше действие"), width: 60}
}

func (m model) Init() tea.Cmd {
    return m.t.Run() // сейчас возвращает nil для встроенных задач
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    m.t, cmd = m.t.Update(msg)
    if m.t.IsDone() {
        return m, tea.Quit
    }
    return m, cmd
}

func (m model) View() string {
    if m.t.IsDone() {
        return m.t.FinalView(m.width)
    }
    return m.t.View(m.width)
}

func main() {
    if _, err := tea.NewProgram(newModel()).Run(); err != nil {
        panic(err)
    }
}
```

## Yes/No (только 2 опции: "Да" и "Нет")

```go
yn := task.NewYesNoTask(
    "Хотите продолжить?",
    "Подтвердите ваше действие",
)
// Клавиши: вверх/k/влево/h = Да, вниз/j/вправо/l = Нет
//          enter/пробел = подтвердить выбор + остановить таймер (если есть)
//          q/esc/ctrl+c = отмена
// После завершения:
choice := yn.GetValue() // true если выбран "Да", false для "Нет"
// Либо используйте:
// yn.IsYes()  // true если выбрано "Да"
// yn.IsNo()   // true если выбрано "Нет"
```

## Одиночный выбор

```go
options := []string{"Первый", "Второй", "Третий"}
ss := task.NewSingleSelectTask("Выберите один вариант", options)
// Клавиши: вверх/k, вниз/j, enter/пробел = подтвердить
// После завершения:
opt := ss.GetSelectedOption()
idx := ss.GetSelectedIndex()
_ = opt; _ = idx
```

## Множественный выбор

```go
options := []string{"A", "B", "C", "D"}
ms := task.NewMultiSelectTask("Выберите несколько вариантов", options)
// Клавиши: вверх/k, вниз/j, пробел = переключить, enter = подтвердить
// После завершения:
opts := ms.GetSelectedOptions()   // []string
idxs := ms.GetSelectedIndices()   // []int
_ = opts; _ = idxs
```

## Ввод текста (новый)

```go
// Базовый вариант
it := task.NewInputTaskNew("Введите имя", "")

// С валидатором через fluent API
it := task.NewInputTaskNew("Введите имя", "").
    WithValidator(func(s string) error {
        if len(s) == 0 { return fmt.Errorf("значение не может быть пустым") }
        return nil
    })

// Или через построитель
it := task.NewInputTaskBuilder("Введите имя", "").
    Required().
    Email().
    Build()

// Клавиши: enter для валидации+завершения, backspace для удаления, другие клавиши добавляют символы
// После завершения:
val := it.GetValue()
_ = val
```

## Выполнение функции (FuncTask)

```go
// Создание задачи, выполняющей произвольную функцию
fn := task.NewFuncTaskWithOptions(
    "Проверка соединения",
    func() error {
        // Здесь могла бы быть реальная проверка
        // Верните ошибку, если произошёл сбой
        return nil
    },
    // Вывод краткой сводки под заголовком после успеха
    task.WithSummaryFunction(func() []string {
        return []string{
            "Пинг: 12мс",
            "Потери пакетов: 0%",
        }
    }),
    // Не останавливать очередь при ошибке (для демонстрации)
    task.WithStopOnError(false),
    // Переопределение метки успешного завершения (по умолчанию "ГОТОВО")
    task.WithSuccessLabelOption("ЗАВЕРШЕНО"),
)

// Клавиши: q/esc/ctrl+c — отмена задачи пользователем

// Использование: добавьте в очередь задач или интегрируйте в вашу модель Bubble Tea
// queue.AddTasks([]common.Task{ fn })
```

## Очередь задач

```go
import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/qzeleza/ziva/common"
    "github.com/qzeleza/ziva/query"
    "github.com/qzeleza/ziva/task"
)

// Создаем очередь задач
queue := query.New("Заголовок очереди")

// Добавляем задачи
queue.AddTasks([]common.Task{
    task.NewYesNoTask("Хотите продолжить?", "Подтвердите действие"),
    task.NewInputTaskNew("Введите имя", ""),
    task.NewSingleSelectTask("Выберите вариант", []string{"A", "B", "C"}),
})

// Запускаем очередь через Bubble Tea
if _, err := tea.NewProgram(queue).Run(); err != nil {
    panic(err)
}
```

### Нумерация задач в очереди

`WithTasksNumbered(enable, keepFirstSymbol, format)` позволяет заменить стандартные маркеры `○/●` на числовую нумерацию.

```go
queue := ziva.NewQueue("CI Pipeline").
    WithAppName("Ziva").
    WithSummary(true).
    // Используем квадратные скобки и лидирующие нули: [01], [02], ...
    WithTasksNumbered(true, false, "[%02d]")

queue.AddTasks(buildTask, testTask, deployTask)

if err := queue.Run(); err != nil {
    log.Fatal(err)
}

// Оставляем первый маркер без номера и переходим на круглые скобки
queue.WithTasksNumbered(true, true, "(%d)")
```

- `enable` — включает нумерацию.
- `keepFirstSymbol` — сохранит `○/●` для первой задачи (остальные будут пронумерованы).
- `format` — любой шаблон `fmt.Sprintf`, например `"(%02d)"`, `"[0%d]"`.

### Локализация интерфейса

- Язык задаётся флагом `--lang` или переменной окружения `ZIVA_LANG` (поддерживаются `ru`, `en`, `tr`, `be`, `uk`).
- При отсутствии нужной локали (например, `ru_RU.UTF-8`) TUI автоматически переключится на английский и подскажет команды для установки.
- Для Entware/BusyBox пригодится установка пакетов `locale-full`, `glibc-binary-locales` и настройка шрифтов (`setterm -reset && setterm -store`).

## Embedded оптимизации

Включение оптимизаций для embedded теперь происходит автоматически при импортировании модуля благодаря внутреннему пакету `internal/autoconfig`. Ручные вызовы из примеров остаются опциональными и могут использоваться для демонстрации.

Полезные переменные окружения для управления поведением автодетекции:

- `ZIVA_EMBEDDED` — принудительно включает/выключает режим embedded (`true`/`1` или пусто/`0`).
- `ZIVA_MEMORY_LIMIT` — задаёт порог памяти для эвристики (например: `64MB`, `128KB`, `1GB`).
- `ZIVA_ASCII_ONLY` — форсирует ASCII-режим (значение `true`).

Примеры запуска с переменными окружения:

```bash
# Принудительно включить embedded-режим
ZIVA_EMBEDDED=1 go run ./cmd/tui_full

# Считать систему «малопамятной» при m.Sys < 64MB
ZIVA_MEMORY_LIMIT=64MB go run ./cmd/tui_full

# Отключить UTF-8 и использовать максимально совместимый вывод
ZIVA_ASCII_ONLY=true go run ./cmd/tui_full
