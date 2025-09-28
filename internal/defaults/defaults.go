package defaults

import (
	"sync/atomic"
	"time"
)

var completionDelayEnabled atomic.Bool

func init() {
	completionDelayEnabled.Store(true)
}

// SetCompletionDelayEnabled глобально включает или выключает задержку завершения задачи.
func SetCompletionDelayEnabled(enabled bool) {
	completionDelayEnabled.Store(enabled)
}

// IsCompletionDelayEnabled сообщает, включена ли задержка завершения задачи.
func IsCompletionDelayEnabled() bool {
	return completionDelayEnabled.Load()
}

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

// Переменные для локализуемых строк
var (
	// Строки статусов
	StatusSuccess    = "УСПЕШНО"
	StatusProblem    = "ПРОБЛЕМА"
	StatusInProgress = "В ПРОЦЕССЕ"

	// Строки для сводки
	SummaryCompleted = "Успешно завершено"
	SummaryOf        = "из"
	SummaryTasks     = "задач"
)

var (
	DefaultNo  = "Нет"
	DefaultYes = "Да"
)

// Константы успешного завершения задач
var (
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
var (
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
var (
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
	ErrorMsgUnknown    = "неизвестная ошибка"
	ErrorMsgTaskPrefix = "задача '%s': "
	ErrorMsgCanceled   = "отменено пользователем"
	ErrorMsgTimeout    = "операция не завершилась за %v"
	ErrorMsgPermission = "недостаточно прав для доступа к %s"

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
	TaskStatusError     = "ОШИБКА"
	TaskStatusCancelled = "Отменено"
)

// Переменные для сообщений валидации и подсказок
var (
	ErrFieldRequired         = "поле обязательно для заполнения"
	ErrPathEmpty             = "путь не может быть пустым"
	ErrPathInvalidChar       = "путь содержит недопустимый символ: %c"
	ErrURLEmpty              = "URL не может быть пустым"
	ErrURLScheme             = "URL должен начинаться с http:// или https://"
	ErrValueEmpty            = "значение не может быть пустым"
	ErrValueAlphaNumeric     = "значение должно содержать только буквы и цифры"
	ErrDefaultValueInvalid   = "значение по умолчанию невалидно"
	ErrDefaultValueEmpty     = "значение по умолчанию пусто"
	CancelShort              = "Отменено"
	NeedSelectAtLeastOne     = "! Необходимо выбрать хотя бы один элемент"
	ScrollAboveFormat        = "%s %s %d выше"
	ScrollBelowFormat        = "%s %s %d ниже"
	SingleSelectHelp         = "[↑/↓ навигация, Enter - выбор, Q/Esc - Выход]"
	MultiSelectHelp          = "[↑/↓ навигация, пробел выбор, Enter подтверждение, Q/Esc - Выход]"
	MultiSelectHelpSelectAll = "[↑/↓ навигация, пробел выбор/переключение всех, Enter подтверждение, Q/Esc - Выход]"
	SelectAllDefaultText     = "Выбрать все"
	InputConfirmHint         = "[Enter - подтвердить, Ctrl+C - отменить]"
	InputFormatLabel         = "Формат:"
	InputHintPassword        = "Используйте надежный пароль"
	InputHintEmail           = "Пример: user@example.com"
	InputHintNumber          = "Введите число"
	InputHintIP              = "Пример: 192.168.1.1"
	InputHintDomain          = "Пример: example.com"
)

// Переменные для сообщений валидации
var (
	ValidatorCustomValidation            = "Пользовательская валидация"
	ValidatorPasswordMinLength           = "пароль должен содержать не менее %d символов"
	ValidatorPasswordCyrillic            = "пароль содержит кириллические символы.\n  пожалуйста, переключитесь на английскую раскладку клавиатуры"
	ValidatorPasswordRequirementDigits   = "цифры"
	ValidatorPasswordRequirementSpecial  = "специальные символы"
	ValidatorPasswordRequirementUpper    = "заглавные буквы"
	ValidatorPasswordRequirementLower    = "строчные буквы"
	ValidatorPasswordMissingRequirements = "пароль должен содержать %s"
	ValidatorPasswordDescription         = "Пароль должен содержать не менее %d символов,\n  включая цифры, специальные символы, заглавные и строчные буквы"
	ValidatorEmailInvalid                = "некорректный email адрес"
	ValidatorEmailDescription            = "Email адрес в формате user@domain.com"
	ValidatorNumberInvalid               = "введите корректное число"
	ValidatorNumberRange                 = "число должно быть в диапазоне от %d до %d"
	ValidatorNumberDescription           = "Число в диапазоне от %d до %d"
	ValidatorIPInvalid                   = "некорректный IP-адрес"
	ValidatorIPv4NotAllowed              = "IPv4 адреса не разрешены"
	ValidatorIPv6NotAllowed              = "IPv6 адреса не разрешены"
	ValidatorIPBothDescription           = "IPv4 или IPv6 адрес"
	ValidatorIPv4Description             = "IPv4 адрес (например, 192.168.1.1)"
	ValidatorIPv6Description             = "IPv6 адрес (например, 2001:db8::1)"
	ValidatorIPGenericDescription        = "IP адрес"
	ValidatorDomainInvalid               = "некорректное доменное имя"
	ValidatorDomainDescription           = "Доменное имя (например, example.com)"
	ValidatorTextMin                     = "текст должен содержать не менее %d символов"
	ValidatorTextMax                     = "текст должен содержать не более %d символов"
	ValidatorTextPattern                 = "текст не соответствует требуемому формату.\n  попробуйте переключить раскладку клавиатуры"
	ValidatorTextBase                    = "Текст"
	ValidatorTextRange                   = " длиной от %d до %d символов"
	ValidatorTextMinOnly                 = " не менее %d символов"
	ValidatorTextMaxOnly                 = " не более %d символов"
	ValidatorCompositeAllErrors          = "ошибки валидации: %s"
	ValidatorCompositeNonePassed         = "ни один валидатор не прошел проверку: %s"
	ValidatorCompositeUnknownMode        = "неизвестный режим композитного валидатора"
	ValidatorCompositeNoValidation       = "Без валидации"
	ValidatorCompositeAllDescription     = "Все требования: %s"
	ValidatorCompositeAnyDescription     = "Любое из требований: %s"
	ValidatorCompositeDescription        = "Композитная валидация"
	ValidatorListSeparator               = ", "
	ValidatorCompositeAllSeparator       = "; "
	ValidatorCompositeAnySeparator       = " ИЛИ "
)

const ClearScreen = "\033[H\033[2J"
