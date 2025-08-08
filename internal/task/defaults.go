package task

import "time"

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

// Константы успешного завершения задач
const (
	// DefaultSuccessLabel метка успешного завершения по умолчанию
	DefaultSuccessLabel = "Готово"

	// DefaultErrorLabel метка ошибки по умолчанию
	DefaultErrorLabel = "Ошибка"

	// DefaultCancelLabel метка отмены по умолчанию
	DefaultCancelLabel = "Отменено пользователем"

	// DefaultYesNoLabel метка успешного завершения по умолчанию
	DefaultYesLabel = "КОНЕЧНО"
	DefaultNoLabel  = "ОТКАЗ"
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
