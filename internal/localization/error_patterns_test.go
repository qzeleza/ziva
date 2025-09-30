package localization

import (
	"errors"
	"testing"

	"github.com/qzeleza/ziva/internal/defaults"
)

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		errMsg   string
		expected ErrorType
	}{
		// Русский язык
		{
			name:     "Russian cancel",
			lang:     "ru",
			errMsg:   "операция отменена пользователем",
			expected: ErrorTypeUserCancel,
		},
		{
			name:     "Russian validation",
			lang:     "ru",
			errMsg:   "поле должно содержать валидное значение",
			expected: ErrorTypeValidation,
		},
		{
			name:     "Russian timeout",
			lang:     "ru",
			errMsg:   "превышен таймаут операции",
			expected: ErrorTypeTimeout,
		},
		{
			name:     "Russian network",
			lang:     "ru",
			errMsg:   "ошибка сетевого соединения",
			expected: ErrorTypeNetwork,
		},
		{
			name:     "Russian filesystem",
			lang:     "ru",
			errMsg:   "файл не найден",
			expected: ErrorTypeFileSystem,
		},

		// Английский язык
		{
			name:     "English cancel",
			lang:     "en",
			errMsg:   "operation cancelled by user",
			expected: ErrorTypeUserCancel,
		},
		{
			name:     "English validation",
			lang:     "en",
			errMsg:   "field must contain valid value",
			expected: ErrorTypeValidation,
		},
		{
			name:     "English timeout",
			lang:     "en",
			errMsg:   "operation timed out",
			expected: ErrorTypeTimeout,
		},
		{
			name:     "English network",
			lang:     "en",
			errMsg:   "network connection failed",
			expected: ErrorTypeNetwork,
		},

		// Турецкий язык
		{
			name:     "Turkish cancel",
			lang:     "tr",
			errMsg:   "işlem iptal edildi",
			expected: ErrorTypeUserCancel,
		},
		{
			name:     "Turkish validation",
			lang:     "tr",
			errMsg:   "değer geçersiz",
			expected: ErrorTypeValidation,
		},

		// Украинский язык
		{
			name:     "Ukrainian cancel",
			lang:     "uk",
			errMsg:   "операція скасована",
			expected: ErrorTypeUserCancel,
		},
		{
			name:     "Ukrainian validation",
			lang:     "uk",
			errMsg:   "поле повинно містити значення",
			expected: ErrorTypeValidation,
		},

		// Белорусский язык
		{
			name:     "Belarusian cancel",
			lang:     "be",
			errMsg:   "аперацыя адменена",
			expected: ErrorTypeUserCancel,
		},

		// Неизвестная ошибка
		{
			name:     "Unknown error",
			lang:     "ru",
			errMsg:   "какая-то странная ошибка без ключевых слов",
			expected: ErrorTypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем язык
			defaults.SetLanguage(tt.lang)

			// Классифицируем ошибку
			err := errors.New(tt.errMsg)
			result := ClassifyError(err)

			if result != tt.expected {
				t.Errorf("ClassifyError() = %v, expected %v for message: %q", result, tt.expected, tt.errMsg)
			}
		})
	}
}

func TestGetKeywordsLazyLoading(t *testing.T) {
	registry := GetRegistry()

	// Очищаем реестр для теста
	registry.mu.Lock()
	registry.loaded = make(map[string]bool)
	registry.patterns = make(map[string]map[ErrorType]ErrorPattern)
	registry.mu.Unlock()

	// Проверяем, что паттерны загружаются только при первом обращении
	keywords := registry.GetKeywords("ru", ErrorTypeValidation)
	if len(keywords) == 0 {
		t.Error("Expected keywords for Russian validation, got empty slice")
	}

	// Проверяем, что язык отмечен как загруженный
	registry.mu.RLock()
	loaded := registry.loaded["ru"]
	registry.mu.RUnlock()

	if !loaded {
		t.Error("Expected Russian language to be marked as loaded")
	}
}

func TestContainsAny(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		keywords []string
		expected bool
	}{
		{
			name:     "Found keyword",
			text:     "операция отменена пользователем",
			keywords: []string{"отмен", "cancel"},
			expected: true,
		},
		{
			name:     "Not found",
			text:     "операция выполнена успешно",
			keywords: []string{"ошибк", "error"},
			expected: false,
		},
		{
			name:     "Empty keywords",
			text:     "some text",
			keywords: []string{},
			expected: false,
		},
		{
			name:     "Case insensitive",
			text:     "Operation CANCELLED by user",
			keywords: []string{"cancel"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsAny(tt.text, tt.keywords)
			if result != tt.expected {
				t.Errorf("ContainsAny() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestClassifyErrorNil(t *testing.T) {
	result := ClassifyError(nil)
	if result != ErrorTypeUnknown {
		t.Errorf("ClassifyError(nil) = %v, expected ErrorTypeUnknown", result)
	}
}