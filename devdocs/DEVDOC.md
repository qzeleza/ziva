# Руководство разработчика Термос

Это руководство объясняет архитектуру Термос, жизненный цикл `Task`, возможности стилизации и то, как расширять фреймворк. Все ссылки и примеры кода соответствуют текущей кодовой базе.

- Корневая директория репозитория: `github.com/qzeleza/termos`
- Минимальная версия Go: 1.22+
- Платформы: Linux, macOS, Windows

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

Канонический интерфейс `Task` определен в `common/task.go`:

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

`task/base.go` предоставляет общее состояние и вспомогательные функции через `BaseTask`:
- Обработка ошибок: `SetError(err)`, `HasError()`, `Error()`
- Завершение: `done` flag, `IsDone()`
- Финальный текст и выравнивание для финального представления с помощью `ui.AlignTextToRight(left, right, width)`
- Флаг stop-on-error: `StopOnError()` / `SetStopOnError(bool)`
- Форматирование ошибок: `WithNewLinesInErrors(preserve bool)` - управляет сохранением переносов строк в сообщениях об ошибках

Как правило, встраивайте `BaseTask` в вашу структуру задачи и делегируйте ему базовые вызовы жизненного цикла.

## Встроенные задачи

Все конструкторы возвращают конкретные типы задач, встраивающие `BaseTask` и реализующие `common.Task`.

- Yes/No – `task.NewYesNoTask(question, description string) *YesNoTask` (только 2 опции: "Да" и "Нет")
  - Клавиши: up/k/left/h для выбора "Да", down/j/right/l для выбора "Нет", enter/space для подтверждения + остановки таймера, q/esc/ctrl+c для отмены (устанавливает ошибку)
  - Результаты: `GetValue() bool` (true для "Да", false для "Нет"), `IsYes()`, `IsNo()`

- Single Select – `task.NewSingleSelectTask(prompt string, options []string) *SingleSelectTask`
  - Клавиши: up/k, down/j, enter/space для подтверждения + остановки таймера
  - Результаты: `GetSelectedOption() string`, `GetSelectedIndex() int`

- Multi Select – `task.NewMultiSelectTask(prompt string, options []string) *MultiSelectTask`
  - Клавиши: up/k, down/j, space для переключения текущего + остановки таймера, enter для подтверждения
  - Результаты: `GetSelectedOptions() []string`, `GetSelectedIndices() []int`

- Input (new) – `task.NewInputTaskNew(prompt, placeholder string, validator func(string) error) *InputTaskNew`
  - Клавиши: enter для подтверждения (validator runs if provided), backspace для удаления, любая другая клавиша добавляет символ в значение
  - Результат: `GetValue() string`

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
// ... позже чтобы вернуть значения по умолчанию
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

1) Определите структуру, встраивающую `BaseTask` и ваши поля.
2) Реализуйте `Title()`, `Run() tea.Cmd`, `Update(msg tea.Msg)`, `View(width int)`, `IsDone() bool`, `FinalView(width int)` и методы обработки ошибок/остановки (или делегируйте их `BaseTask`).
3) Предоставьте конструктор `NewYourTask(...)` возвращающий ваш тип задачи.
4) В `Update` обрабатывайте релевантные клавиши/события и устанавливайте `done`, `finalValue`, ошибки по необходимости.

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

`InputTaskNew` может принимать необязательный параметр `validator func(string) error`. При нажатии Enter, если валидатор возвращает ошибку, задача устанавливает ошибку через `SetError(err)` и остается активной.

## Очереди и оркестрация (планируется)

Пакет `query/` зарезервирован для очередей задач (последовательные пайплайны, распространение ошибок через `StopOnError`, статистика). Интеграции стоит строить поверх `common.Task` и цикла обновления Bubble Tea.

---

Если вам нужны дополнительные примеры или вы хотите предложить улучшения — создайте issue или PR.
