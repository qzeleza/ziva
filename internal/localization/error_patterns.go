// Package localization предоставляет централизованное управление паттернами ошибок
// для классификации ошибок на основе языка интерфейса.
// Оптимизирован для встраиваемых систем с минимальным потреблением памяти.
package localization

import (
	"strings"
	"sync"

	"github.com/qzeleza/ziva/internal/defaults"
	"github.com/qzeleza/ziva/internal/performance"
)

// ErrorType определяет тип ошибки (дублирует errors.ErrorType для избежания циклических зависимостей)
type ErrorType int

const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeValidation
	ErrorTypeUserCancel
	ErrorTypeTimeout
	ErrorTypeNetwork
	ErrorTypeFileSystem
	ErrorTypePermission
	ErrorTypeConfig
)

// ErrorPattern содержит ключевые слова для определения типа ошибки
type ErrorPattern struct {
	Keywords []string
}

// ErrorPatternRegistry хранит паттерны для всех типов ошибок по языкам
type ErrorPatternRegistry struct {
	patterns map[string]map[ErrorType]ErrorPattern
	mu       sync.RWMutex
	loaded   map[string]bool // Отслеживание загруженных языков
}

var (
	registry     *ErrorPatternRegistry
	registryOnce sync.Once
)

// GetRegistry возвращает глобальный реестр паттернов (singleton)
func GetRegistry() *ErrorPatternRegistry {
	registryOnce.Do(func() {
		registry = &ErrorPatternRegistry{
			patterns: make(map[string]map[ErrorType]ErrorPattern),
			loaded:   make(map[string]bool),
		}
	})
	return registry
}

// GetKeywords возвращает ключевые слова для конкретного типа ошибки и языка
// Использует ленивую загрузку - загружает паттерны только при первом обращении
func (r *ErrorPatternRegistry) GetKeywords(lang string, errType ErrorType) []string {
	r.mu.RLock()
	if r.loaded[lang] {
		if langPatterns, ok := r.patterns[lang]; ok {
			if pattern, ok := langPatterns[errType]; ok {
				r.mu.RUnlock()
				return pattern.Keywords
			}
		}
		r.mu.RUnlock()
		return nil
	}
	r.mu.RUnlock()

	// Загружаем паттерны для языка
	r.loadLanguage(lang)

	// Повторная попытка чтения
	r.mu.RLock()
	defer r.mu.RUnlock()
	if langPatterns, ok := r.patterns[lang]; ok {
		if pattern, ok := langPatterns[errType]; ok {
			return pattern.Keywords
		}
	}
	return nil
}

// loadLanguage загружает паттерны для указанного языка
func (r *ErrorPatternRegistry) loadLanguage(lang string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверяем, не загружен ли уже
	if r.loaded[lang] {
		return
	}

	// Загружаем паттерны в зависимости от языка
	switch lang {
	case "ru":
		r.patterns[lang] = russianPatterns()
	case "en":
		r.patterns[lang] = englishPatterns()
	case "tr":
		r.patterns[lang] = turkishPatterns()
	case "be":
		r.patterns[lang] = belarusianPatterns()
	case "uk":
		r.patterns[lang] = ukrainianPatterns()
	default:
		// Fallback на английский
		r.patterns[lang] = englishPatterns()
	}

	r.loaded[lang] = true
}

// ContainsAny проверяет, содержит ли строка хотя бы одно из ключевых слов
func ContainsAny(s string, keywords []string) bool {
	lowerStr := performance.ToLowerEfficient(s)
	for _, keyword := range keywords {
		if keyword != "" && strings.Contains(lowerStr, keyword) {
			return true
		}
	}
	return false
}

// ClassifyError классифицирует ошибку на основе текущего языка
func ClassifyError(err error) ErrorType {
	if err == nil {
		return ErrorTypeUnknown
	}

	lang := defaults.CurrentLanguage()
	registry := GetRegistry()

	// Проверяем каждый тип ошибки в порядке приоритета
	errorTypes := []ErrorType{
		ErrorTypeUserCancel,
		ErrorTypeValidation,
		ErrorTypeTimeout,
		ErrorTypeNetwork,
		ErrorTypeFileSystem,
		ErrorTypePermission,
		ErrorTypeConfig,
	}

	for _, errType := range errorTypes {
		keywords := registry.GetKeywords(lang, errType)
		if ContainsAny(err.Error(), keywords) {
			return errType
		}
	}

	return ErrorTypeUnknown
}

// russianPatterns возвращает паттерны ошибок для русского языка
func russianPatterns() map[ErrorType]ErrorPattern {
	return map[ErrorType]ErrorPattern{
		ErrorTypeUserCancel: {
			Keywords: []string{"отмен", "отмена"},
		},
		ErrorTypeValidation: {
			Keywords: []string{
				"валидац", "некоррект", "неправиль",
				"должен", "обязательн", "недопустим",
				"содержать", "валидн",
			},
		},
		ErrorTypeTimeout: {
			Keywords: []string{"таймаут", "тайм-аут", "превышен"},
		},
		ErrorTypeNetwork: {
			Keywords: []string{"сеть", "соединен", "подключен"},
		},
		ErrorTypeFileSystem: {
			Keywords: []string{
				"файл", "директор", "папк", "путь",
				"не найден", "не существует",
			},
		},
		ErrorTypePermission: {
			Keywords: []string{"доступ", "запрещен", "прав"},
		},
		ErrorTypeConfig: {
			Keywords: []string{"конфигурац", "настройк", "параметр"},
		},
	}
}

// englishPatterns возвращает паттерны ошибок для английского языка
func englishPatterns() map[ErrorType]ErrorPattern {
	return map[ErrorType]ErrorPattern{
		ErrorTypeUserCancel: {
			Keywords: []string{"canceled", "cancelled", "cancel"},
		},
		ErrorTypeValidation: {
			Keywords: []string{
				"validat", "invalid", "must contain",
				"must be", "required", "should",
			},
		},
		ErrorTypeTimeout: {
			Keywords: []string{"timeout", "deadline exceeded", "timed out"},
		},
		ErrorTypeNetwork: {
			Keywords: []string{"network", "connection", "connect"},
		},
		ErrorTypeFileSystem: {
			Keywords: []string{
				"file", "director", "folder", "path",
				"not found", "does not exist", "no such",
			},
		},
		ErrorTypePermission: {
			Keywords: []string{"permission", "denied", "access denied", "forbidden"},
		},
		ErrorTypeConfig: {
			Keywords: []string{"config", "configuration", "setting"},
		},
	}
}

// turkishPatterns возвращает паттерны ошибок для турецкого языка
func turkishPatterns() map[ErrorType]ErrorPattern {
	return map[ErrorType]ErrorPattern{
		ErrorTypeUserCancel: {
			Keywords: []string{"iptal"},
		},
		ErrorTypeValidation: {
			Keywords: []string{
				"doğrula", "dogrula", "geçersiz", "gecersiz",
				"içermelidir", "icermelidir", "olmalıdır", "olmalidir",
			},
		},
		ErrorTypeTimeout: {
			Keywords: []string{"zaman aş", "zaman asim"},
		},
		ErrorTypeNetwork: {
			Keywords: []string{"ağ", "ag", "bağlant", "baglant"},
		},
		ErrorTypeFileSystem: {
			Keywords: []string{
				"dosya", "dizin", "yol",
				"bulunamadı", "bulunamadi", "mevcut değil",
			},
		},
		ErrorTypePermission: {
			Keywords: []string{"izin", "reddedildi", "yasak"},
		},
		ErrorTypeConfig: {
			Keywords: []string{"yapılandır", "yapilandir", "ayar"},
		},
	}
}

// belarusianPatterns возвращает паттерны ошибок для белорусского языка
func belarusianPatterns() map[ErrorType]ErrorPattern {
	return map[ErrorType]ErrorPattern{
		ErrorTypeUserCancel: {
			Keywords: []string{"адмен"},
		},
		ErrorTypeValidation: {
			Keywords: []string{
				"валідац", "некарэкт", "няправільн",
				"павін", "змяшчаць", "абавязков",
			},
		},
		ErrorTypeTimeout: {
			Keywords: []string{"таймаўт", "перавышан"},
		},
		ErrorTypeNetwork: {
			Keywords: []string{"сетк", "злучен", "падключ"},
		},
		ErrorTypeFileSystem: {
			Keywords: []string{
				"файл", "дырэктор", "шлях",
				"не знойдз", "не існу",
			},
		},
		ErrorTypePermission: {
			Keywords: []string{"доступ", "забаронен", "прав"},
		},
		ErrorTypeConfig: {
			Keywords: []string{"канфігурац", "налад", "параметр"},
		},
	}
}

// ukrainianPatterns возвращает паттерны ошибок для украинского языка
func ukrainianPatterns() map[ErrorType]ErrorPattern {
	return map[ErrorType]ErrorPattern{
		ErrorTypeUserCancel: {
			Keywords: []string{"скас", "скасуван"},
		},
		ErrorTypeValidation: {
			Keywords: []string{
				"валідац", "некорект", "неправильн",
				"повинен", "має містити", "має бути", "обов'язков",
				"містити", "значення",
			},
		},
		ErrorTypeTimeout: {
			Keywords: []string{"тайм-аут", "таймаут", "перевищен"},
		},
		ErrorTypeNetwork: {
			Keywords: []string{"мереж", "з'єднан", "зєднан", "підключ"},
		},
		ErrorTypeFileSystem: {
			Keywords: []string{
				"файл", "директор", "тек", "шлях",
				"не знайден", "не існує",
			},
		},
		ErrorTypePermission: {
			Keywords: []string{"доступ", "заборонен", "прав"},
		},
		ErrorTypeConfig: {
			Keywords: []string{"конфігурац", "налаштув", "параметр"},
		},
	}
}