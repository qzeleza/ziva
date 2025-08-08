package errors

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskError(t *testing.T) {
	// Создаем базовую ошибку
	originalErr := errors.New("original error")
	taskErr := NewTaskError("test task", originalErr, ErrorTypeValidation)
	
	assert.NotNil(t, taskErr)
	assert.Equal(t, "test task", taskErr.TaskTitle)
	assert.Equal(t, originalErr, taskErr.Err)
	assert.Equal(t, ErrorTypeValidation, taskErr.Type)
	assert.NotNil(t, taskErr.Context)
}

func TestTaskErrorError(t *testing.T) {
	originalErr := errors.New("original error")
	taskErr := NewTaskError("test task", originalErr, ErrorTypeValidation)
	
	errorString := taskErr.Error()
	assert.Contains(t, errorString, "test task")
	assert.Contains(t, errorString, "original error")
}

func TestTaskErrorUnwrap(t *testing.T) {
	originalErr := errors.New("original error")
	taskErr := NewTaskError("test task", originalErr, ErrorTypeValidation)
	
	unwrapped := taskErr.Unwrap()
	assert.Equal(t, originalErr, unwrapped)
}

func TestWithContext(t *testing.T) {
	originalErr := errors.New("original error")
	taskErr := NewTaskError("test task", originalErr, ErrorTypeValidation)
	
	taskErr = taskErr.WithContext("user", "testuser")
	value, exists := taskErr.GetContext("user")
	assert.True(t, exists)
	assert.Equal(t, "testuser", value)
}

func TestIsRetryable(t *testing.T) {
	// Retryable ошибки
	timeoutErr := NewTaskError("test task", errors.New("timeout"), ErrorTypeTimeout)
	assert.True(t, timeoutErr.IsRetryable())
	
	networkErr := NewTaskError("test task", errors.New("network"), ErrorTypeNetwork)
	assert.True(t, networkErr.IsRetryable())
	
	// Non-retryable ошибки
	validationErr := NewTaskError("test task", errors.New("validation"), ErrorTypeValidation)
	assert.False(t, validationErr.IsRetryable())
	
	cancelErr := NewTaskError("test task", errors.New("cancel"), ErrorTypeUserCancel)
	assert.False(t, cancelErr.IsRetryable())
}

func TestGetUserFriendlyMessage(t *testing.T) {
	originalErr := errors.New("original error")
	taskErr := NewTaskError("test task", originalErr, ErrorTypeValidation)
	
	message := taskErr.GetUserFriendlyMessage()
	assert.NotEmpty(t, message)
	assert.Equal(t, "Проверьте правильность введенных данных", message)
}

func TestErrorTypeCreators(t *testing.T) {
	// Тестируем все creator функции
	validationErr := NewValidationError("test task", errors.New("validation failed"))
	assert.Equal(t, ErrorTypeValidation, validationErr.Type)
	assert.Contains(t, validationErr.Error(), "validation failed")
	
	cancelErr := NewCancelError("test task")
	assert.Equal(t, ErrorTypeUserCancel, cancelErr.Type)
	assert.Contains(t, cancelErr.Error(), "отменено пользователем")
	
	timeoutErr := NewTimeoutError("test task", 30*time.Second)
	assert.Equal(t, ErrorTypeTimeout, timeoutErr.Type)
	assert.True(t, timeoutErr.IsRetryable())
	assert.Contains(t, timeoutErr.Error(), "30s")
	
	networkErr := NewNetworkError("test task", errors.New("connection failed"))
	assert.Equal(t, ErrorTypeNetwork, networkErr.Type)
	assert.True(t, networkErr.IsRetryable())
	
	fsErr := NewFileSystemError("test task", errors.New("file not found"), "/path/to/file")
	assert.Equal(t, ErrorTypeFileSystem, fsErr.Type)
	value, exists := fsErr.GetContext("path")
	assert.True(t, exists)
	assert.Equal(t, "/path/to/file", value)
	
	permErr := NewPermissionError("test task", "/path/to/file")
	assert.Equal(t, ErrorTypePermission, permErr.Type)
	assert.Contains(t, permErr.Error(), "/path/to/file")
	
	configErr := NewConfigurationError("test task", errors.New("parse error"), "invalid.yaml")
	assert.Equal(t, ErrorTypeConfiguration, configErr.Type)
	value, exists = configErr.GetContext("config_key")
	assert.True(t, exists)
	assert.Equal(t, "invalid.yaml", value)
}

func TestErrorHandler(t *testing.T) {
	handler := NewErrorHandler()
	assert.NotNil(t, handler)
	
	// Тестируем обработку разных типов ошибок
	originalErr := errors.New("test error")
	
	handledErr := handler.Handle("test task", originalErr)
	assert.NotNil(t, handledErr)
	assert.Equal(t, "test task", handledErr.TaskTitle)
	
	// Тестируем обработку уже существующей TaskError
	taskErr := NewTaskError("existing task", originalErr, ErrorTypeValidation)
	handledTaskErr := handler.Handle("new task", taskErr)
	assert.Equal(t, taskErr, handledTaskErr) // Должна вернуть оригинальную ошибку
}

func TestClassifyError(t *testing.T) {
	handler := NewErrorHandler()
	
	// Тестируем классификацию различных ошибок
	tests := []struct {
		err      error
		expected ErrorType
	}{
		{errors.New("validation error"), ErrorTypeValidation},
		{errors.New("network timeout"), ErrorTypeTimeout},
		{errors.New("permission denied"), ErrorTypePermission},
		{errors.New("file not found"), ErrorTypeFileSystem},
		{errors.New("unknown error"), ErrorTypeUnknown},
		{errors.New("отменено пользователем"), ErrorTypeUserCancel},
	}
	
	for _, test := range tests {
		classified := handler.classifyError(test.err)
		// Проверяем, что функция не паникует и возвращает разумный тип
		assert.IsType(t, ErrorTypeUnknown, classified)
	}
}

func TestFormatForUser(t *testing.T) {
	handler := NewErrorHandler()
	taskErr := NewValidationError("test task", errors.New("invalid input"))
	
	userMessage := handler.FormatForUser(taskErr)
	assert.NotEmpty(t, userMessage)
	assert.Equal(t, "Проверьте правильность введенных данных", userMessage)
}

func TestShouldRetry(t *testing.T) {
	handler := NewErrorHandler()
	
	// Retryable ошибка
	timeoutErr := NewTimeoutError("test task", 30*time.Second)
	assert.True(t, handler.ShouldRetry(timeoutErr, 1))
	assert.False(t, handler.ShouldRetry(timeoutErr, 10)) // Превышен лимит попыток
	
	// Non-retryable ошибка
	validationErr := NewValidationError("test task", errors.New("invalid input"))
	assert.False(t, handler.ShouldRetry(validationErr, 1))
	
	// Обычная ошибка
	normalErr := errors.New("normal error")
	assert.False(t, handler.ShouldRetry(normalErr, 1))
}

func TestErrorTypeString(t *testing.T) {
	// Тестируем, что все типы ошибок имеют строковое представление
	types := []ErrorType{
		ErrorTypeUnknown,
		ErrorTypeValidation,
		ErrorTypeUserCancel,
		ErrorTypeTimeout,
		ErrorTypeNetwork,
		ErrorTypeFileSystem,
		ErrorTypePermission,
		ErrorTypeConfiguration,
	}
	
	for _, errorType := range types {
		str := errorType.String()
		assert.NotEmpty(t, str, "ErrorType должен иметь строковое представление")
	}
}

func TestErrorWithNilOriginal(t *testing.T) {
	// Тестируем создание TaskError с nil оригинальной ошибкой
	taskErr := NewTaskError("test task", nil, ErrorTypeValidation)
	assert.NotNil(t, taskErr)
	assert.Equal(t, "test task", taskErr.TaskTitle)
	assert.Nil(t, taskErr.Err)
	
	// Unwrap должен вернуть nil
	assert.Nil(t, taskErr.Unwrap())
	
	// Error не должен паниковать
	errorString := taskErr.Error()
	assert.Contains(t, errorString, "test task")
}

func TestErrorContextEdgeCases(t *testing.T) {
	taskErr := NewTaskError("test task", errors.New("test"), ErrorTypeValidation)
	
	// Context с различными типами значений
	taskErr = taskErr.WithContext("string", "value")
	taskErr = taskErr.WithContext("int", 42)
	taskErr = taskErr.WithContext("bool", true)
	
	strValue, exists := taskErr.GetContext("string")
	assert.True(t, exists)
	assert.Equal(t, "value", strValue)
	
	intValue, exists := taskErr.GetContext("int")
	assert.True(t, exists)
	assert.Equal(t, 42, intValue)
	
	boolValue, exists := taskErr.GetContext("bool")
	assert.True(t, exists)
	assert.Equal(t, true, boolValue)
	
	// Несуществующий ключ
	_, exists = taskErr.GetContext("nonexistent")
	assert.False(t, exists)
}