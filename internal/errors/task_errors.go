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

	// Проверяем на отмену пользователем
	if strings.Contains(errStr, "отменено") || strings.Contains(errStr, "canceled") ||
		strings.Contains(errStr, "cancelled") {
		return ErrorTypeUserCancel
	}

	// Проверяем на валидацию
	if strings.Contains(errStr, "валидация") || strings.Contains(errStr, "validation") ||
		strings.Contains(errStr, "некорректн") || strings.Contains(errStr, "неправильн") ||
		strings.Contains(errStr, "должен содержать") {
		return ErrorTypeValidation
	}

	// Проверяем на таймаут
	if strings.Contains(errStr, "timeout") || strings.Contains(errStr, "таймаут") ||
		strings.Contains(errStr, "deadline exceeded") {
		return ErrorTypeTimeout
	}

	// Проверяем на сетевые ошибки
	if strings.Contains(errStr, "network") || strings.Contains(errStr, "connection") ||
		strings.Contains(errStr, "сеть") || strings.Contains(errStr, "соединение") {
		return ErrorTypeNetwork
	}

	// Проверяем на файловые ошибки
	if strings.Contains(errStr, "file") || strings.Contains(errStr, "directory") ||
		strings.Contains(errStr, "path") || strings.Contains(errStr, "файл") {
		return ErrorTypeFileSystem
	}

	// Проверяем на ошибки прав доступа
	if strings.Contains(errStr, "permission") || strings.Contains(errStr, "access") ||
		strings.Contains(errStr, "forbidden") || strings.Contains(errStr, "unauthorized") ||
		strings.Contains(errStr, "права") || strings.Contains(errStr, "доступ") {
		return ErrorTypePermission
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
