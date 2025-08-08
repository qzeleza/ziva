# Руководство разработчика Термос

Это руководство объясняет архитектуру Термос, жизненный цикл `Task`, возможности стилизации и то, как расширять фреймворк. Все ссылки и примеры кода соответствуют текущей кодовой базе.

- Repo root: `github.com/qzeleza/termos`
- Minimum Go: 1.22+
- Platforms: Linux, macOS, Windows

## Обзор архитектуры

```
termos/
├── common/        # Shared interfaces and helpers
├── errors/        # Error related utilities
├── examples/      # Usage examples
├── performance/   # Performance/ANSI utilities
├── query/         # Queues (planned/placeholder)
├── task/          # Built-in tasks (Yes/No, Single, Multi, Input)
├── ui/            # Styles, icons, layout helpers
├── validation/    # Validation (planned/placeholder)
└── docs/          # Project documentation
```

Ключевые файлы:
- `common/task.go` – определяет интерфейс `Task`, используемый в Термос
- `task/base.go` – общее поведение для задач (`BaseTask`)
- `task/yesno_task.go`, `task/singleselect_task.go`, `task/multiselect_task.go`, `task/input_task_new.go` – встроенные задачи
- `ui/styles.go` – экспортируемые стили lipgloss, иконки и утилиты

## Интерфейс Task и жизненный цикл

The canonical `Task` interface is defined in `common/task.go`:

```go
// Task представляет собой интерфейс для выполнения задач в очереди.
type Task interface {
    // Title возвращает заголовок задачи.
    Title() string
    
    // Run запускает выполнение задачи и возвращает команду bubbletea.
    Run() tea.Cmd
    
    // Update обновляет состояние задачи на основе полученного сообщения.
    Update(msg tea.Msg) (Task, tea.Cmd)
    
    // View отображает текущее состояние задачи с учетом указанной ширины.
    View(width int) string
    
    // IsDone возвращает true, если задача завершена.
    IsDone() bool
    
    // FinalView отображает финальное состояние задачи с учетом указанной ширины.
    FinalView(width int) string
    
    // HasError возвращает true, если при выполнении задачи произошла ошибка.
    HasError() bool
    
    // Error возвращает ошибку, если она есть.
    Error() error
    
    // StopOnError возвращает true, если при возникновении ошибки в этой задаче
    // нужно остановить выполнение всей очереди задач.
    StopOnError() bool
    
    // SetStopOnError устанавливает флаг остановки очереди при ошибке.
    SetStopOnError(stop bool)
}
```

Типичный жизненный цикл:
1) Создайте задачу (например, `task.NewYesNoTask(...)`).
2) При необходимости вызовите `Run()` для асинхронных команд (во встроенных задачах сейчас возвращает `nil`).
3) Передавайте сообщения Bubble Tea через `Update(msg)` до тех пор, пока `IsDone()` не станет `true`.
4) Во время взаимодействия рендерьте `View(width)`; после завершения используйте `FinalView(width)`.
5) Получайте результаты через геттеры конкретных задач (например, `GetChoice()`, `GetSelectedOption()`).

## BaseTask

`task/base.go` provides shared state and helpers via `BaseTask`:
- Error handling: `SetError(err)`, `HasError()`, `Error()`
- Completion: `done` flag, `IsDone()`
- Final label and alignment for final view using `ui.AlignTextToRight(left, right, width)`
- Stop-on-error flag: `StopOnError()` / `SetStopOnError(bool)`

Как правило, встраивайте `BaseTask` в вашу структуру задачи и делегируйте ему базовые вызовы жизненного цикла.

## Встроенные задачи

All constructors return concrete task types that embed `BaseTask` and implement `common.Task`.

- Yes/No – `task.NewYesNoTask(question, description string) *YesNoTask`
  - Keys: up/k/left/h to choose “Да”, down/j/right/l to choose “Нет”, enter/space confirm, q/esc/ctrl+c cancel (sets error)
  - Result: `GetChoice() bool` returns true only if “Да” was chosen and confirmed

- Single Select – `task.NewSingleSelectTask(prompt string, options []string) *SingleSelectTask`
  - Keys: up/k, down/j, enter/space confirm
  - Results: `GetSelectedOption() string`, `GetSelectedIndex() int`

- Multi Select – `task.NewMultiSelectTask(prompt string, options []string) *MultiSelectTask`
  - Keys: up/k, down/j, space toggle current, enter confirm
  - Results: `GetSelectedOptions() []string`, `GetSelectedIndices() []int`

- Input (new) – `task.NewInputTaskNew(prompt, placeholder string, validator func(string) error) *InputTaskNew`
  - Keys: enter to accept (validator runs if provided), backspace to delete, any other key appends to value
  - Result: `GetValue() string`

Примечания:
- Все реализации `Run()` во встроенных задачах сейчас возвращают `nil` (нет фоновой работы). Интегрируйте задачи в модель Bubble Tea и управляйте ими через `Update`.
- `FinalView(width)` показывает финальное состояние, включая выравнивание заголовка и результата; при ошибке применяются стили ошибок.

## Стилизация и UI-утилиты

See `ui/styles.go`.

Экспортируемые стили и цвета (неполный список):
- Заголовки и статусы: `ui.TitleStyle`, `ui.FinishedLabelStyle`, `ui.SuccessLabelStyle`
- Ошибки: `ui.ErrorMessageStyle`, `ui.ErrorStatusStyle`, `ui.CancelStyle`
- Выделение/активные: `ui.SelectionStyle`, `ui.SelectionNoStyle`, `ui.ActiveStyle`, `ui.InputStyle`, `ui.SpinnerStyle`
- Палитра цветов: например, `ui.ColorBrightGreen`, `ui.ColorBrightRed`, ...

Кастомизация стилей (переопределяйте методами lipgloss):

```go
ui.TitleStyle = ui.TitleStyle.Foreground(ui.ColorBrightGreen).Bold(true)
ui.SelectionNoStyle = ui.SelectionNoStyle.Foreground(ui.ColorBrightRed).Bold(true)
```

Хелперы для цветов ошибок:

```go
ui.SetErrorColor(ui.ColorDarkYellow, ui.ColorBrightYellow)
// ... later to reset defaults
ui.ResetErrorColors()
```

Иконки и префиксы для дерева задач:
- Иконки: `ui.IconDone`, `ui.IconError`, `ui.IconCancelled`, `ui.IconCursor` и др.
- Префиксы: `ui.GetActiveTaskPrefix()`, `ui.GetTaskBelowPrefix()`, `ui.GetCompletedTaskPrefix(success bool)`, `ui.GetCompletedInputTaskPrefix(success bool)`

Утилиты разметки/форматирования:
- `ui.AlignTextToRight(left, right, width int) string`
- `ui.FormatErrorMessage(text string, width int) string`
- ANSI-утилиты: `ui.GetPlainTextLength(text)`, `ui.StripANSI(text)`, `ui.WrapText(text, width)`, `ui.GetRuneWidth(r rune)` (внутри использует `performance.StripANSILength`)

## Сборка и запуск примеров

- Сборка библиотеки:
  ```bash
  go build ./...
  ```

- Запуск примера:
  ```bash
  go run ./examples/basic_usage.go
  ```

- Использование в вашем приложении: выполните `go get github.com/qzeleza/termos` и импортируйте нужные пакеты.

## Создание новой задачи

1) Define a struct embedding `BaseTask` and your fields.
2) Implement `Title()`, `Run() tea.Cmd`, `Update(msg tea.Msg)`, `View(width int)`, `IsDone() bool`, `FinalView(width int)` and error/stop methods (or delegate to `BaseTask`).
3) Provide a constructor `NewYourTask(...)` returning your task type.
4) In `Update`, handle relevant keys/events and set `done`, `finalValue`, errors as needed.

Пример «скелета»:

```go
package task

type MyTask struct {
    BaseTask
    // your fields
}

func NewMyTask(title string) *MyTask {
    return &MyTask{BaseTask: NewBaseTask(title)}
}

func (t *MyTask) View(width int) string { /* ... */ return "" }
func (t *MyTask) Update(msg tea.Msg) (common.Task, tea.Cmd) { /* ... */ return t, nil }
func (t *MyTask) Run() tea.Cmd { return nil }
```

## Валидация

`InputTaskNew` optionally accepts `validator func(string) error`. On Enter, if validator returns an error, the task sets the error via `SetError(err)` and remains active.

## Очереди и оркестрация (планируется)

Пакет `query/` зарезервирован для очередей задач (последовательные пайплайны, распространение ошибок через `StopOnError`, статистика). Интеграции стоит строить поверх `common.Task` и цикла обновления Bubble Tea.

---

Если вам нужны дополнительные примеры или вы хотите предложить улучшения — создайте issue или PR.
