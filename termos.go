// Package termos предоставляет Terminal Multitask Orchestrator System - библиотеку для создания интерактивных TUI приложений.
// Эта библиотека позволяет создавать разнообразные задачи пользовательского интерфейса такие как:
// - Задачи выбора Да/Нет
// - Задачи одиночного выбора из списка
// - Задачи множественного выбора
// - Задачи ввода текста с валидацией
// - Задачи выполнения функций
// - Очереди задач с автоматическим выполнением
//
// Пример использования:
//
//	queue := termos.NewQueue("Пример использования Termos")
//
//	// Создаем задачу выбора Да/Нет
//	confirm := termos.NewYesNoTask("Подтверждение", "Продолжить выполнение?")
//
//	// Создаем задачу выбора из списка
//	options := []string{"development", "staging", "production"}
//	env := termos.NewSingleSelectTask("Выбор среды", options)
//
//	// Добавляем задачи в очередь
//	queue.AddTasks(confirm, env)
//
//	// Запускаем очередь
//	err := queue.Run()
//	if err != nil {
//		log.Fatal(err)
//	}
package termos

import (
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/qzeleza/termos/internal/common"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
	"github.com/qzeleza/termos/internal/ui"
	"github.com/qzeleza/termos/internal/validation"
)

// Task представляет собой интерфейс для выполнения задач в очереди.
// Этот интерфейс используется как в пакете task, так и в пакете query.
type Task = common.Task

// YesNoOption представляет варианты выбора для YesNoTask
type YesNoOption = task.YesNoOption

const (
	// YesOption - опция "Да"
	YesOption = task.YesOption
	// NoOption - опция "Нет"
	NoOption = task.NoOption
)

// Queue представляет очередь задач для выполнения
type Queue struct {
	model *query.Model
}

// ----------------------------------------------------------------------------
// Queue
// ----------------------------------------------------------------------------

// NewQueue создает новую очередь задач с заданным заголовком
//
// @param title Заголовок очереди
// @return Указатель на новую очередь задач
func NewQueue(title string) *Queue {
	return &Queue{
		model: query.New(title),
	}
}

// AddTasks добавляет задачи в очередь
//
// @param tasks Задачи, которые будут добавлены в очередь
// @return Указатель на очередь задач
func (q *Queue) AddTasks(tasks ...Task) *Queue {
	q.model.AddTasks(tasks)
	return q
}

// WithAppName устанавливает название приложения в заголовке
//
// @param appName Название приложения
// @return Указатель на очередь задач
func (q *Queue) WithAppName(appName string) *Queue {
	q.model.WithAppName(appName)
	return q
}

// WithSummary включает/выключает отображение сводки по завершению
//
// @param show Флаг, указывающий, нужно ли отображать сводку
// @param lineEnabled Флаг, указывающий, нужно ли отображать линию сводки, под которой выводится сводка
// @return Указатель на очередь задач
func (q *Queue) WithOutSummary(lineEnabled bool) *Queue {
	q.model.WithSummary(false)

	return q
}

// WithResultFormatting включает форматирование результатов задач с разделительными линиями.
// Если enabled=true, то перед каждым результатом задачи будет добавляться разделительная линия
// из префикса и указанного количества символов "─".
// При enabled=false поведение остается как и раньше - результаты выводятся сразу после строки с задачей.
// @param enabled - включить/выключить форматирование
func (q *Queue) WithOutResultLine() *Queue {
	q.model.WithResultFormatting(false)
	return q
}

// WithTasksNumbered включает отображение номеров для задач и задаёт формат (например "[%02d]" или "(%d)").
//
// @param keepFirstSymbol Флаг, указывающий, нужно ли сохранять первый символ формата
// @param numberFormat Формат номеров задач
// @return Указатель на очередь задач
func (q *Queue) WithTasksNumbered(keepFirstSymbol bool, numberFormat string) *Queue {
	q.model.WithTasksNumbered(true, keepFirstSymbol, numberFormat)
	return q
}

// WithAppNameColor устанавливает цвет текста и стиль названия приложения.
//
// @param textColor Цвет текста
// @param bold Флаг, указывающий, нужно ли сделать текст жирным
// @return Указатель на очередь задач
func (q *Queue) WithAppNameColor(textColor lipgloss.TerminalColor, bold bool) *Queue {
	q.model.WithAppNameColor(textColor, bold)
	return q
}

// WithTitleColor устанавливает цвет заголовка.
//
// @param titleColor Цвет заголовка
// @param bold Флаг, указывающий, нужно ли сделать заголовок жирным
// @return Указатель на очередь задач
func (q *Queue) WithTitleColor(titleColor lipgloss.TerminalColor, bold bool) *Queue {
	q.model.WithTitleColor(titleColor, bold)
	return q
}

// WithClearScreen включает/выключает очистку экрана перед запуском очереди задач
//
// @param clear Флаг, указывающий, нужно ли очищать экран перед запуском очереди задач
// @return Указатель на очередь задач
func (q *Queue) WithClearScreen(clear bool) *Queue {
	q.model.WithClearScreen(clear)
	return q
}

// SetErrorColor устанавливает цвет для отображения ошибок в очереди.
//
// @param color Цвет для отображения ошибок
// @return Указатель на очередь задач
func (q *Queue) SetErrorColor(color query.ErrorColor) *Queue {
	q.model.SetErrorColor(color)
	return q
}

// Run запускает выполнение очереди задач
//
// @return Ошибка, если она возникла
func (q *Queue) Run() error {
	return q.model.Run()
}

// ----------------------------------------------------------------------------
// YesNoTask
// ----------------------------------------------------------------------------

// NewYesNoTask создает новую задачу выбора Да/Нет
//
// @param title Заголовок задачи
// @param question Вопрос, который будет задан пользователю
// @return Указатель на новую задачу выбора Да/Нет
func NewYesNoTask(title, question string) *YesNoTask {
	return &YesNoTask{task.NewYesNoTask(title, question)}
}

// YesNoTask представляет задачу выбора из двух опций: Да, Нет
type YesNoTask struct {
	*task.YesNoTask
}

// WithTimeout устанавливает тайм-аут для задачи с значением по умолчанию
//
// @param duration Тайм-аут
// @param defaultValue Значение по умолчанию
// @return Указатель на задачу для цепочки вызовов
func (t *YesNoTask) WithTimeout(duration time.Duration, defaultValue interface{}) *YesNoTask {
	t.WithDefaultOption(defaultValue, duration)
	return t
}

// WithDefaultItem задает опцию, которая будет подсвечена при открытии задачи.
//
// @param option Опция, которая будет подсвечена при открытии задачи
// @return Указатель на задачу для цепочки вызовов
func (t *YesNoTask) WithDefaultItem(option interface{}) *YesNoTask {
	t.YesNoTask.WithDefaultItem(option)
	return t
}

// WithCustomLabels позволяет изменить текст опций
//
// @param yesLabel Текст для опции "Да"
// @param noLabel Текст для опции "Нет"
// @return Указатель на задачу для цепочки вызовов
func (t *YesNoTask) WithCustomLabels(yesLabel, noLabel string) *YesNoTask {
	t.YesNoTask.WithCustomLabels(yesLabel, noLabel)
	return t
}

// SetChoiceHelpDelimiter задаёт разделитель, используемый для встроенных подсказок в строках выбора.
// По умолчанию используется "::".
func SetChoiceHelpDelimiter(delim string) {
	task.SetChoiceHelpDelimiter(delim)
}

// GetSelectedOption возвращает выбранную опцию
//
// @return Выбранная опция
func (t *YesNoTask) GetSelectedOption() YesNoOption {
	return t.YesNoTask.GetSelectedOption()
}

// IsYes возвращает true если выбрано "Да"
//
// @return true если выбрано "Да"
func (t *YesNoTask) IsYes() bool {
	return t.YesNoTask.IsYes()
}

// IsNo возвращает true если выбрано "Нет"
//
// @return true если выбрано "Нет"
func (t *YesNoTask) IsNo() bool {
	return t.YesNoTask.IsNo()
}

// ----------------------------------------------------------------------------
// SingleSelectTask
// ----------------------------------------------------------------------------

// NewSingleSelectTask создает новую задачу выбора одного варианта из списка
//
// @param title Заголовок задачи
// @param choices Список вариантов выбора
// @return Указатель на новую задачу выбора
func NewSingleSelectTask(title string, choices []string) *SingleSelectTask {
	return &SingleSelectTask{task.NewSingleSelectTask(title, choices)}
}

// SingleSelectTask представляет задачу для выбора одного варианта из списка
type SingleSelectTask struct {
	*task.SingleSelectTask
}

// WithViewport устанавливает размер viewport (окна просмотра) для ограничения количества отображаемых элементов.
// Это полезно для длинных списков, когда нужно показывать только часть элементов.
//
// @param size Количество элементов для отображения одновременно (0 = показать все)
// @param showCounters Флаг, указывающий, нужно ли отображать счетчики элементов
// @return Указатель на задачу для цепочки вызовов
func (t *SingleSelectTask) WithViewport(size int, showCounters ...bool) *SingleSelectTask {
	t.SingleSelectTask.WithViewport(size, showCounters...)
	return t
}

// WithTimeout устанавливает тайм-аут для задачи с значением по умолчанию
//
// @param duration Тайм-аут в миллисекундах
// @param defaultValue Значение по умолчанию
// @return Указатель на задачу для цепочки вызовов
func (t *SingleSelectTask) WithTimeout(duration time.Duration, defaultValue interface{}) *SingleSelectTask {
	t.SingleSelectTask.WithTimeout(duration, defaultValue)
	return t
}

// WithItemsDisabled помечает элементы меню как недоступные для выбора.
// Поддерживаются типы: int, []int, string, []string. Nil очищает список отключённых элементов.
//
// @param disabled Список индексов или значений, которые нужно отключить
// @return Указатель на задачу для цепочки вызовов
func (t *SingleSelectTask) WithItemsDisabled(disabled interface{}) *SingleSelectTask {
	t.SingleSelectTask.WithItemsDisabled(disabled)
	return t
}

// WithDefaultItem задает элемент, который будет подсвечен при открытии задачи.
//
// @param selection Индекс или значение элемента, который будет подсвечен
// @return Указатель на задачу для цепочки вызовов
func (t *SingleSelectTask) WithDefaultItem(selection interface{}) *SingleSelectTask {
	t.SingleSelectTask.WithDefaultItem(selection)
	return t
}

// GetSelected возвращает выбранное значение
//
// @return Выбранное значение
func (t *SingleSelectTask) GetSelected() string {
	return t.SingleSelectTask.GetSelected()
}

// GetSelectedIndex возвращает индекс выбранного элемента
//
// @return Индекс выбранного элемента
func (t *SingleSelectTask) GetSelectedIndex() int {
	return t.SingleSelectTask.GetSelectedIndex()
}

// ----------------------------------------------------------------------------
// MultiSelectTask
// ----------------------------------------------------------------------------

// NewMultiSelectTask создает новую задачу множественного выбора
//
// @param title Заголовок задачи
// @param choices Список вариантов выбора
// @return Указатель на новую задачу множественного выбора
func NewMultiSelectTask(title string, choices []string) *MultiSelectTask {
	return &MultiSelectTask{task.NewMultiSelectTask(title, choices)}
}

// MultiSelectTask представляет задачу для выбора нескольких вариантов из списка
type MultiSelectTask struct {
	*task.MultiSelectTask
}

// WithViewport устанавливает размер viewport (окна просмотра) для ограничения количества отображаемых элементов.
// Это полезно для длинных списков, когда нужно показывать только часть элементов.
//
// @param size Количество элементов для отображения одновременно (0 = показать все)
// @param showCounters Флаг, указывающий, нужно ли отображать счетчики элементов
// @return Указатель на задачу для цепочки вызовов
func (t *MultiSelectTask) WithViewport(size int, showCounters ...bool) *MultiSelectTask {
	t.MultiSelectTask.WithViewport(size, showCounters...)
	return t
}

// WithTimeout устанавливает тайм-аут для задачи с значениями по умолчанию
//
// @param duration Тайм-аут в миллисекундах
// @param defaultValues Значения по умолчанию
// @return Указатель на задачу для цепочки вызовов
func (t *MultiSelectTask) WithTimeout(duration time.Duration, defaultValues interface{}) *MultiSelectTask {
	t.MultiSelectTask.WithTimeout(duration, defaultValues)
	return t
}

// WithItemsDisabled помечает элементы меню как недоступные для выбора.
// Поддерживаются типы: int, []int, string, []string. Nil очищает список отключённых элементов.
//
// @param disabled Список индексов или значений, которые нужно отключить
// @return Указатель на задачу для цепочки вызовов
func (t *MultiSelectTask) WithItemsDisabled(disabled interface{}) *MultiSelectTask {
	t.MultiSelectTask.WithItemsDisabled(disabled)
	return t
}

// WithSelectAll добавляет опцию "Выбрать все" в начало списка
//
// @param text Текст опции "Выбрать все" (по умолчанию "Выбрать все")
// @return Указатель на задачу для цепочки вызовов
func (t *MultiSelectTask) WithSelectAll(text ...string) *MultiSelectTask {
	t.MultiSelectTask.WithSelectAll(text...)
	return t
}

// WithDefaultItems задает элементы, которые будут отмечены при открытии задачи.
//
// @param defauiltSelection Индекс или значение элемента, который будет подсвечен
// @return Указатель на задачу для цепочки вызовов
func (t *MultiSelectTask) WithDefaultItems(defauiltSelection interface{}) *MultiSelectTask {
	t.MultiSelectTask.WithDefaultItems(defauiltSelection)
	return t
}

// GetSelected возвращает список выбранных элементов
//
// @return Список выбранных элементов
func (t *MultiSelectTask) GetSelected() []string {
	return t.MultiSelectTask.GetSelected()
}

// ----------------------------------------------------------------------------
// InputTask
// ----------------------------------------------------------------------------

// NewInputTask создает новую задачу ввода текста
//
// @param title Заголовок задачи
// @param prompt Подсказка для ввода
// @return Указатель на новую задачу ввода текста
func NewInputTask(title, prompt string) *InputTask {
	return &InputTask{task.NewInputTaskNew(title, prompt)}
}

// InputTask представляет задачу ввода текста
type InputTask struct {
	*task.InputTaskNew
}

// WithTimeout устанавливает тайм-аут для задачи с значением по умолчанию
//
// @param duration Тайм-аут в миллисекундах
// @param defaultValue Значение по умолчанию
// @return Указатель на задачу для цепочки вызовов
func (t *InputTask) WithTimeout(duration time.Duration, defaultValue string) *InputTask {
	t.InputTaskNew.WithTimeout(duration, defaultValue)
	return t
}

// WithValidator устанавливает валидатор для ввода
//
// @param validator Валидатор для ввода
// @return Указатель на задачу для цепочки вызовов
func (t *InputTask) WithValidator(validator validation.Validator) *InputTask {
	t.InputTaskNew.WithValidator(validator)
	return t
}

// WithInputType устанавливает тип ввода (пароль, email и т.д.)
//
// @param inputType Тип ввода
// @return Указатель на задачу для цепочки вызовов
func (t *InputTask) WithInputType(inputType task.InputType) *InputTask {
	t.InputTaskNew.WithInputType(inputType)
	return t
}

// GetValue возвращает введенное значение
//
// @return Введенное значение
func (t *InputTask) GetValue() string {
	return t.InputTaskNew.GetValue()
}

// InputType представляет тип поля ввода
type InputType = task.InputType

const (
	InputTypeText     = task.InputTypeText
	InputTypePassword = task.InputTypePassword
	InputTypeEmail    = task.InputTypeEmail
	InputTypeNumber   = task.InputTypeNumber
	InputTypeIP       = task.InputTypeIP
	InputTypeDomain   = task.InputTypeDomain
)

// ----------------------------------------------------------------------------
// FuncTask
// ----------------------------------------------------------------------------

// NewFuncTask создает новую задачу выполнения функции
//
// @param title Заголовок задачи
// @param fn Функция, которая будет выполнена
// @param opts Опции для конфигурации задачи
// @return Указатель на новую задачу выполнения функции
func NewFuncTask(title string, fn func() error, opts ...task.FuncTaskOption) *FuncTask {
	return &FuncTask{task.NewFuncTask(title, fn, opts...)}
}

// FuncTask представляет задачу выполнения функции
type FuncTask struct {
	*task.FuncTask
}

// WithStopOnError устанавливает флаг остановки очереди при ошибке
//
// @param stop Флаг остановки очереди при ошибке
// @return Указатель на задачу для цепочки вызовов
func (t *FuncTask) WithStopOnError(stop bool) *FuncTask {
	t.SetStopOnError(stop)
	return t
}

// Валидаторы - экспорт фабрики валидаторов
var DefaultValidators = validation.DefaultFactory

// FuncTaskOption представляет опцию для конфигурации FuncTask
type FuncTaskOption = task.FuncTaskOption

// Опции для FuncTask
var (
	WithSummaryFunction = task.WithSummaryFunction
	WithStopOnError     = task.WithStopOnError
)

// Стили для текста
var (
	ErrorStatusStyle   = ui.ErrorStatusStyle
	ErrorMessageStyle  = ui.ErrorMessageStyle
	CancelStyle        = ui.CancelStyle
	SubtleStyle        = ui.SubtleStyle
	SelectionStyle     = ui.SelectionStyle
	SelectionNoStyle   = ui.SelectionNoStyle
	ActiveStyle        = ui.ActiveStyle
	InputStyle         = ui.InputStyle
	SpinnerStyle       = ui.SpinnerStyle
	ActiveTitleStyle   = ui.ActiveTitleStyle
	ActiveTaskStyle    = ui.ActiveTaskStyle
	SuccessLabelStyle  = ui.SuccessLabelStyle
	FinishedLabelStyle = ui.FinishedLabelStyle
)

// Коды цветов
var (
	GreenBright   = ui.ColorBrightGreen
	RedBright     = ui.ColorBrightRed
	RedDark       = ui.ColorDarkRed
	YellowBright  = ui.ColorBrightYellow
	YellowDark    = ui.ColorDarkYellow
	OrangeBright  = ui.ColorBrightOrange
	OrangeDark    = ui.ColorDarkOrange
	BlueBright    = ui.ColorBrightBlue
	BlueDark      = ui.ColorDarkBlue
	CyanBright    = ui.ColorBrightCyan
	CyanDark      = ui.ColorDarkCyan
	MagentaBright = ui.ColorBrightMagenta
	WhiteBright   = ui.ColorBrightWhite
	GrayBright    = ui.ColorBrightGray
	GrayDark      = ui.ColorDarkGray
	LightBlue     = ui.ColorLightBlue
	Black         = ui.ColorBlack
	DarkGreen     = ui.ColorDarkGreen
)
