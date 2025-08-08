# Примеры Termos

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

Примечание: Встроенный `Run()` сейчас возвращает `nil`. Интегрируйте задачи в модель Bubble Tea и управляйте ими через `Update`.

## Минимальная интеграция Bubble Tea для одной задачи

```go
package main

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/qzeleza/termos/common"
    "github.com/qzeleza/termos/task"
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

## Yes/No

```go
yn := task.NewYesNoTask(
    "Хотите продолжить?",
    "Подтвердите ваше действие",
)
// Клавиши: вверх/k/влево/h = Да, вниз/j/вправо/l = Нет, enter/пробел = подтвердить, q/esc/ctrl+c = отмена
// После завершения:
choice := yn.GetChoice() // true если подтверждено "Да"
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

## Очередь задач

```go
import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/qzeleza/termos/common"
    "github.com/qzeleza/termos/query"
    "github.com/qzeleza/termos/task"
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

## Embedded оптимизации

```go
import (
    "github.com/qzeleza/termos/examples"
    "github.com/qzeleza/termos/ui"
)

// Автоматическое определение embedded окружения
if examples.IsEmbeddedEnvironment() {
    // Получаем оптимизированную конфигурацию
    config := examples.OptimizedEmbeddedConfig()
    
    // Применяем оптимизации
    examples.ApplyEmbeddedConfig(config)
    ui.EnableEmbeddedMode()
}
