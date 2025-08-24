package defauilt

import (
	"time"
)

// Константы для времени выполнения задач
const (
	// DefaultCompletionDelay время ожидания перед завершением задачи для плавности анимации
	DefaultCompletionDelay = 300 * time.Millisecond

	// SpinnerTickInterval интервал обновления спиннера
	SpinnerTickInterval = 100 * time.Millisecond

	// InputValidationDelay задержка валидации ввода для предотвращения слишком частых проверок
	InputValidationDelay = 200 * time.Millisecond
)

// Константы для валидации
const (
	// DefaultPasswordMinLength минимальная длина пароля по умолчанию
	DefaultPasswordMinLength = 8

	// StrongPasswordMinLength минимальная длина для сильного пароля
	StrongPasswordMinLength = 12

	// MaxInputLength максимальная длина ввода по умолчанию
	MaxInputLength = 255

	// DefaultNumberMin минимальное значение числа по умолчанию
	DefaultNumberMin = 0

	// DefaultNumberMax максимальное значение числа по умолчанию
	DefaultNumberMax = 1000000
)

// Константы для текстового ввода
const (
	// DefaultInputWidth ширина поля ввода по умолчанию
	DefaultInputWidth = 40

	// MaxInputWidth максимальная ширина поля ввода
	MaxInputWidth = 80

	// MinInputWidth минимальная ширина поля ввода
	MinInputWidth = 10
)

// Константы для отображения
const (
	// DefaultLayoutWidth ширина макета по умолчанию
	DefaultLayoutWidth = 80

	// ErrorDisplayMaxLines максимальное количество строк для отображения ошибки
	ErrorDisplayMaxLines = 3

	// HelpTextMaxWidth максимальная ширина текста справки
	HelpTextMaxWidth = 60
)

// Константы для русскоязычных строк
const (
	// Строки статусов
	StatusSuccess    = "УСПЕШНО"
	StatusProblem    = "ПРОБЛЕМА"
	StatusInProgress = "В ПРОЦЕССЕ"

	// Строки для сводки
	SummaryCompleted = "Успешно завершено"
	SummaryOf        = "из"
	SummaryTasks     = "задач"
)

const (
	DefaultNo  = "Нет"
	DefaultYes = "Да"
)

// Константы успешного завершения задач
const (
	// DefaultSuccessLabel метка успешного завершения по умолчанию
	DefaultSuccessLabel = "Готово"

	DefaultFromSummaryLabel  = "из"
	DefaultTasksSummaryLabel = "задач"

	// DefaultErrorLabel метка ошибки по умолчанию
	DefaultErrorLabel = "Ошибка"

	// DefaultCancelLabel метка отмены по умолчанию
	DefaultCancelLabel = "Отменено пользователем"

	DefaultSelectedLabel = "пользователь выбрал"

	// DefaultYesNoLabel метка успешного завершения по умолчанию
	DefaultYesLabel = StatusSuccess
	DefaultNoLabel  = "ОТКАЗ"

	// TaskCancelledByUser сообщение об отмене задачи пользователем
	TaskCancelledByUser = "[отменено пользователем]"

	// TaskExitHint подсказка о выходе из задачи
	TaskExitHint = "[Для выхода из задачи нажмите Ctrl+C]"
)

// Константы для задач ввода
const (
	// DefaultPrompt подсказка по умолчанию для ввода
	DefaultPrompt = "Введите значение"

	// PasswordMask символ маскировки для паролей
	PasswordMask = '*'

	// DefaultPlaceholder текст-заполнитель по умолчанию
	DefaultPlaceholder = "..."

	// DefaultSeparator разделитель для списка
	DefaultSeparator = "♀ "
)

// Константы для сообщений об ошибках
const (
	// Типы ошибок
	ErrorTypeValidation = "ВАЛИДАЦИЯ"
	ErrorTypeUserCancel = "ОТМЕНА"
	ErrorTypeTimeout    = "ТАЙМАУТ"
	ErrorTypeNetwork    = "СЕТЬ"
	ErrorTypeFileSystem = "ФАЙЛ"
	ErrorTypePermission = "ДОСТУП"
	ErrorTypeConfig     = "КОНФИГ"
	ErrorTypeUnknown    = "ОШИБКА"

	// Сообщения об ошибках
	ErrorMsgUnknown       = "неизвестная ошибка"
	ErrorMsgTaskPrefix    = "задача '%s': "
	ErrorMsgCanceled      = "отменено пользователем"
	ErrorMsgTimeout       = "операция не завершилась за %v"
	ErrorMsgPermission    = "недостаточно прав для доступа к %s"
	
	// Пользовательские сообщения об ошибках
	ErrorUserMsgValidation    = "Проверьте правильность введенных данных"
	ErrorUserMsgCancel        = "Операция отменена"
	ErrorUserMsgTimeout       = "Операция заняла слишком много времени"
	ErrorUserMsgNetwork       = "Проблема с сетевым соединением"
	ErrorUserMsgFileSystem    = "Проблема доступа к файлу"
	ErrorUserMsgPermission    = "Недостаточно прав для выполнения операции"
	ErrorUserMsgConfiguration = "Ошибка в настройках"
	ErrorUserMsgUnknown       = "Произошла неизвестная ошибка"

	// Константы для статусов задач
	TaskStatusError = "ОШИБКА"
	TaskStatusCancelled = "Отменено"
)

const ClearScreen = "\033[H\033[2J"
