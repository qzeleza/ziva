// Package errors предоставляет типизированные ошибки и утилиты для их обработки в терминальном UI
package errors

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/qzeleza/ziva/internal/defaults"
	"github.com/qzeleza/ziva/internal/performance"
)

// TaskError представляет ошибку, возникшую в задаче
type TaskError struct {
	// Тип ошибки
	Type ErrorType

	// Название задачи где произошла ошибка
	TaskTitle string

	// Исходная ошибка
	Err error

	// Время возникновения ошибки
	Timestamp time.Time

	// Дополнительный контекст
	Context map[string]interface{}
}

// ErrorType определяет тип ошибки задачи
type ErrorType int

const (
	// ErrorTypeUnknown неизвестная ошибка
	ErrorTypeUnknown ErrorType = iota

	// ErrorTypeValidation ошибка валидации
	ErrorTypeValidation

	// ErrorTypeUserCancel пользователь отменил операцию
	ErrorTypeUserCancel

	// ErrorTypeTimeout таймаут операции
	ErrorTypeTimeout

	// ErrorTypeNetwork сетевая ошибка
	ErrorTypeNetwork

	// ErrorTypeFileSystem ошибка файловой системы
	ErrorTypeFileSystem

	// ErrorTypePermission ошибка прав доступа
	ErrorTypePermission

	// ErrorTypeConfiguration ошибка конфигурации
	ErrorTypeConfiguration
)

// Error реализует интерфейс error
func (te *TaskError) Error() string {
	builder := performance.GetBuffer()
	defer performance.PutBuffer(builder)

	// Добавляем тип ошибки если он не неизвестный
	if te.Type != ErrorTypeUnknown {
		fmt.Fprintf(builder, "[%s] ", te.Type.String())
	}

	// Добавляем контекст задачи
	if te.TaskTitle != "" {
		fmt.Fprintf(builder, defaults.ErrorMsgTaskPrefix, te.TaskTitle)
	}

	// Добавляем основное сообщение об ошибке
	if te.Err != nil {
		builder.WriteString(te.Err.Error())
	} else {
		builder.WriteString(defaults.ErrorMsgUnknown)
	}

	return builder.String()
}

// Unwrap возвращает исходную ошибку для работы с errors.Is и errors.As
func (te *TaskError) Unwrap() error {
	return te.Err
}

// String возвращает строковое представление типа ошибки
func (et ErrorType) String() string {
	switch et {
	case ErrorTypeValidation:
		return defaults.ErrorTypeValidation
	case ErrorTypeUserCancel:
		return defaults.ErrorTypeUserCancel
	case ErrorTypeTimeout:
		return defaults.ErrorTypeTimeout
	case ErrorTypeNetwork:
		return defaults.ErrorTypeNetwork
	case ErrorTypeFileSystem:
		return defaults.ErrorTypeFileSystem
	case ErrorTypePermission:
		return defaults.ErrorTypePermission
	case ErrorTypeConfiguration:
		return defaults.ErrorTypeConfig
	default:
		return defaults.ErrorTypeUnknown
	}
}

// NewTaskError создает новую ошибку задачи
func NewTaskError(taskTitle string, err error, errorType ErrorType) *TaskError {
	return &TaskError{
		Type:      errorType,
		TaskTitle: taskTitle,
		Err:       err,
		Timestamp: time.Now(),
		Context:   make(map[string]interface{}),
	}
}

// WithContext добавляет контекст к ошибке
func (te *TaskError) WithContext(key string, value interface{}) *TaskError {
	te.Context[key] = value
	return te
}

// GetContext возвращает значение из контекста
func (te *TaskError) GetContext(key string) (interface{}, bool) {
	value, exists := te.Context[key]
	return value, exists
}

// IsRetryable определяет, можно ли повторить операцию
func (te *TaskError) IsRetryable() bool {
	switch te.Type {
	case ErrorTypeNetwork, ErrorTypeTimeout:
		return true
	case ErrorTypeUserCancel, ErrorTypeValidation, ErrorTypePermission:
		return false
	default:
		return false
	}
}

// GetUserFriendlyMessage возвращает понятное пользователю сообщение
func (te *TaskError) GetUserFriendlyMessage() string {
	switch te.Type {
	case ErrorTypeValidation:
		return defaults.ErrorUserMsgValidation
	case ErrorTypeUserCancel:
		return defaults.ErrorUserMsgCancel
	case ErrorTypeTimeout:
		return defaults.ErrorUserMsgTimeout
	case ErrorTypeNetwork:
		return defaults.ErrorUserMsgNetwork
	case ErrorTypeFileSystem:
		return defaults.ErrorUserMsgFileSystem
	case ErrorTypePermission:
		return defaults.ErrorUserMsgPermission
	case ErrorTypeConfiguration:
		return defaults.ErrorUserMsgConfiguration
	default:
		if te.Err != nil {
			return te.Err.Error()
		}
		return defaults.ErrorUserMsgUnknown
	}
}

// Предопределенные конструкторы ошибок

// NewValidationError создает ошибку валидации
func NewValidationError(taskTitle string, err error) *TaskError {
	return NewTaskError(taskTitle, err, ErrorTypeValidation)
}

// NewCancelError создает ошибку отмены пользователем
func NewCancelError(taskTitle string) *TaskError {
	return NewTaskError(taskTitle, errors.New(defaults.ErrorMsgCanceled), ErrorTypeUserCancel)
}

// NewTimeoutError создает ошибку таймаута
func NewTimeoutError(taskTitle string, duration time.Duration) *TaskError {
	err := fmt.Errorf(defaults.ErrorMsgTimeout, duration)
	return NewTaskError(taskTitle, err, ErrorTypeTimeout).
		WithContext("duration", duration)
}

// NewNetworkError создает сетевую ошибку
func NewNetworkError(taskTitle string, err error) *TaskError {
	return NewTaskError(taskTitle, err, ErrorTypeNetwork)
}

// NewFileSystemError создает ошибку файловой системы
func NewFileSystemError(taskTitle string, err error, path string) *TaskError {
	return NewTaskError(taskTitle, err, ErrorTypeFileSystem).
		WithContext("path", path)
}

// NewPermissionError создает ошибку прав доступа
func NewPermissionError(taskTitle string, resource string) *TaskError {
	err := fmt.Errorf(defaults.ErrorMsgPermission, resource)
	return NewTaskError(taskTitle, err, ErrorTypePermission).
		WithContext("resource", resource)
}

// NewConfigurationError создает ошибку конфигурации
func NewConfigurationError(taskTitle string, err error, configKey string) *TaskError {
	return NewTaskError(taskTitle, err, ErrorTypeConfiguration).
		WithContext("config_key", configKey)
}

// ErrorHandler обрабатывает ошибки и предоставляет единообразное поведение
type ErrorHandler struct {
	// ShowStackTrace показывать ли стек трассировки (для отладки)
	ShowStackTrace bool

	// LogErrors логировать ли ошибки
	LogErrors bool

	// RetryAttempts количество попыток повтора для retryable ошибок
	RetryAttempts int
}

// NewErrorHandler создает новый обработчик ошибок
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		ShowStackTrace: false,
		LogErrors:      true,
		RetryAttempts:  3,
	}
}

// Handle обрабатывает ошибку и возвращает обогащенную TaskError
func (eh *ErrorHandler) Handle(taskTitle string, err error) *TaskError {
	// Если ошибка уже является TaskError, возвращаем её
	var taskErr *TaskError
	if errors.As(err, &taskErr) {
		return taskErr
	}

	// Определяем тип ошибки по содержанию
	errorType := eh.classifyError(err)

	// Создаем новую TaskError
	return NewTaskError(taskTitle, err, errorType)
}

// classifyError определяет тип ошибки по её содержанию
func (eh *ErrorHandler) classifyError(err error) ErrorType {
	if err == nil {
		return ErrorTypeUnknown
	}

	errStr := performance.ToLowerEfficient(err.Error())

	// Унифицированная проверка по наборам ключевых слов для всех поддерживаемых языков (ru, en, tr, be, uk)
	// Используем подстроки/стемы, чтобы охватить разные формы слов и склонения.
	containsAny := func(s string, subs []string) bool {
		for _, sub := range subs {
			if sub != "" && strings.Contains(s, sub) {
				return true
			}
		}
		return false
	}

	// Отмена пользователем
	cancelWords := []string{
		// ru
		"отмен", "отмена",
		// en
		"canceled", "cancelled", "cancel",
		// uk
		"скас",
		// be
		"адмен",
		// tr
		"iptal",
	}
	if containsAny(errStr, cancelWords) {
		return ErrorTypeUserCancel
	}

	// Ошибки валидации
	validationWords := []string{
		// ru
		"валидац", "некоррект", "неправиль", "должен", "обязательн",
		// en
		"validat", "invalid", "must contain", "must be", "required", "should",
		// uk
		"валідац", "некорект", "повинен", "має містити", "має бути",
		// be
		"валідац", "некарэкт", "павін", "змяшчаць",
		// tr
		"doğrula", "dogrula", "geçersiz", "gecersiz", "içermelidir", "icermelidir", "olmalıdır", "olmalidir",
	}
	if containsAny(errStr, validationWords) {
		return ErrorTypeValidation
	}

	// Таймаут
	timeoutWords := []string{
		// en
		"timeout", "deadline exceeded", "timed out",
		// ru
		"таймаут",
		// uk
		"тайм-аут",
		// be
		"таймаўт",
		// tr
		"zaman aş", "zaman asim",
	}
	if containsAny(errStr, timeoutWords) {
		return ErrorTypeTimeout
	}

	// Сетевые ошибки
	networkWords := []string{
		// en
		"network", "connection", "connect",
		// ru
		"сеть", "соединен", "подключен",
		// uk
		"мереж", "з'єднан", "зєднан",
		// be
		"сетк", "злучен",
		// tr
		"ağ", "ag", "bağlant", "baglant",
	}
	if containsAny(errStr, networkWords) {
		return ErrorTypeNetwork
	}

	// Ошибки файловой системы
	fileWords := []string{
		// en
		"file", "directory", "path", "not found", "no such file",
		// ru
		"файл", "каталог", "директори", "путь", "не найден", "нет такого файла",
		// uk
		"файл", "каталог", "шлях", "не знайден",
		// be
		"файл", "каталог", "дырэктор", "шлях", "не знойдзен",
		// tr
		"dosya", "dizin", "yol", "bulunamad",
	}
	if containsAny(errStr, fileWords) {
		return ErrorTypeFileSystem
	}

	// Ошибки прав доступа
	permissionWords := []string{
		// en
		"permission", "access", "forbidden", "unauthorized", "denied",
		// ru
		"права", "доступ", "запрещено", "отказано",
		// uk
		"доступ", "заборонено", "відмовлено", "прав",
		// be
		"даступ", "забаронена", "адмоўлена", "прав",
		// tr
		"izin", "erişim", "yasak", "yetkisiz", "reddedildi",
	}
	if containsAny(errStr, permissionWords) {
		return ErrorTypePermission
	}

	// Ошибки конфигурации
	configWords := []string{
		// en
		"config", "configuration",
		// ru
		"настройк", "конфиг",
		// uk
		"конфіг",
		// be
		"канфіг",
		// tr
		"yapılandır", "yapilandir",
	}
	if containsAny(errStr, configWords) {
		return ErrorTypeConfiguration
	}

	// По умолчанию - неизвестная ошибка
	return ErrorTypeUnknown
}

// FormatForUser возвращает отформатированное для пользователя сообщение об ошибке
func (eh *ErrorHandler) FormatForUser(err error) string {
	var taskErr *TaskError
	if errors.As(err, &taskErr) {
		return taskErr.GetUserFriendlyMessage()
	}

	// Для обычных ошибок возвращаем как есть
	return err.Error()
}

// ShouldRetry определяет, следует ли повторить операцию
func (eh *ErrorHandler) ShouldRetry(err error, attempt int) bool {
	if attempt >= eh.RetryAttempts {
		return false
	}

	var taskErr *TaskError
	if errors.As(err, &taskErr) {
		return taskErr.IsRetryable()
	}

	return false
}

// Глобальный обработчик ошибок
var DefaultErrorHandler = NewErrorHandler()
