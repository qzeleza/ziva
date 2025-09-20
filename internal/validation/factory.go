package validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qzeleza/termos/internal/defauilt"
	"github.com/qzeleza/termos/internal/performance"
)

// ValidatorFactory предоставляет удобные методы для создания часто используемых валидаторов
type ValidatorFactory struct{}

// NewFactory создает новую фабрику валидаторов
func NewFactory() *ValidatorFactory {
	return &ValidatorFactory{}
}

// StrongPassword создает валидатор для сильного пароля (минимум 12 символов)
func (f *ValidatorFactory) StrongPassword() Validator {
	return NewPasswordValidator(12)
}

// StandardPassword создает валидатор для стандартного пароля (минимум 8 символов)
func (f *ValidatorFactory) StandardPassword() Validator {
	return NewPasswordValidator(8)
}

// Email создает валидатор для email адресов
func (f *ValidatorFactory) Email() Validator {
	return NewEmailValidator()
}

// Port создает валидатор для номеров портов (1-65535)
func (f *ValidatorFactory) Port() Validator {
	return NewNumberValidator(1, 65535)
}

// HTTPPort создает валидатор для HTTP портов (1-65535, обычно 80, 443, 8080 и т.д.)
func (f *ValidatorFactory) HTTPPort() Validator {
	return NewNumberValidator(1, 65535)
}

// IPv4 создает валидатор только для IPv4 адресов
func (f *ValidatorFactory) IPv4() Validator {
	return NewIPv4Validator()
}

// IPv6 создает валидатор только для IPv6 адресов
func (f *ValidatorFactory) IPv6() Validator {
	return NewIPv6Validator()
}

// IP создает валидатор для любых IP адресов (IPv4 и IPv6)
func (f *ValidatorFactory) IP() Validator {
	return NewIPValidator(true, true)
}

// Domain создает валидатор для доменных имен
func (f *ValidatorFactory) Domain() Validator {
	return NewDomainValidator()
}

// Username создает валидатор для имен пользователей (3-32 символа, буквы, цифры, подчеркивания)
func (f *ValidatorFactory) Username() Validator {
	return NewTextValidator(3, 32).WithPattern(`^[a-zA-Z0-9_]+$`)
}

// Required создает валидатор для обязательных полей (не пустых)
func (f *ValidatorFactory) Required() Validator {
	return ValidatorFunc(func(input string) error {
		if len(performance.TrimSpaceEfficient(input)) == 0 {
			return errors.New(defauilt.ErrFieldRequired)
		}
		return nil
	})
}

// OptionalEmail создает валидатор для опционального email (может быть пустым, но если заполнен - должен быть валидным)
func (f *ValidatorFactory) OptionalEmail() Validator {
	return ValidatorFunc(func(input string) error {
		trimmed := performance.TrimSpaceEfficient(input)
		if trimmed == "" {
			return nil // Пустое значение допустимо
		}
		return NewEmailValidator().Validate(trimmed)
	})
}

// Path создает валидатор для путей файловой системы
func (f *ValidatorFactory) Path() Validator {
	return ValidatorFunc(func(input string) error {
		if len(performance.TrimSpaceEfficient(input)) == 0 {
			return errors.New(defauilt.ErrPathEmpty)
		}
		// Проверяем на недопустимые символы для большинства ОС
		invalidChars := `<>:"|?*`
		if performance.ContainsAnyEfficient(input, invalidChars) {
			for _, char := range input {
				if strings.ContainsRune(invalidChars, char) {
					return fmt.Errorf(defauilt.ErrPathInvalidChar, char)
				}
			}
		}
		return nil
	})
}

// URL создает валидатор для URL адресов
func (f *ValidatorFactory) URL() Validator {
	return ValidatorFunc(func(input string) error {
		trimmed := performance.TrimSpaceEfficient(input)
		if trimmed == "" {
			return errors.New(defauilt.ErrURLEmpty)
		}

		// Простая проверка URL
		lowerTrimmed := performance.ToLowerEfficient(trimmed)
		if !strings.HasPrefix(lowerTrimmed, "http://") &&
			!strings.HasPrefix(lowerTrimmed, "https://") {
			return errors.New(defauilt.ErrURLScheme)
		}

		return nil
	})
}

// Range создает валидатор для чисел в заданном диапазоне
func (f *ValidatorFactory) Range(min, max int) Validator {
	return NewNumberValidator(min, max)
}

// MinLength создает валидатор для минимальной длины текста
func (f *ValidatorFactory) MinLength(minLen int) Validator {
	return NewTextValidator(minLen, 0)
}

// MaxLength создает валидатор для максимальной длины текста
func (f *ValidatorFactory) MaxLength(maxLen int) Validator {
	return NewTextValidator(0, maxLen)
}

// Length создает валидатор для точной длины текста
func (f *ValidatorFactory) Length(exactLen int) Validator {
	return NewTextValidator(exactLen, exactLen)
}

// AlphaNumeric создает валидатор для алфавитно-цифровых строк
func (f *ValidatorFactory) AlphaNumeric() Validator {
	return ValidatorFunc(func(input string) error {
		if len(performance.TrimSpaceEfficient(input)) == 0 {
			return errors.New(defauilt.ErrValueEmpty)
		}

		for _, char := range input {
			if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < '0' || char > '9') {
				return errors.New(defauilt.ErrValueAlphaNumeric)
			}
		}
		return nil
	})
}

// Предустановленные валидаторы
var (
	// DefaultFactory глобальный экземпляр фабрики
	DefaultFactory = NewFactory()
)

// Удобные функции для быстрого доступа к часто используемым валидаторам
var (
	StrongPassword   = DefaultFactory.StrongPassword
	StandardPassword = DefaultFactory.StandardPassword
	Email            = DefaultFactory.Email
	Port             = DefaultFactory.Port
	IPv4             = DefaultFactory.IPv4
	IPv6             = DefaultFactory.IPv6
	IP               = DefaultFactory.IP
	Domain           = DefaultFactory.Domain
	Username         = DefaultFactory.Username
	Required         = DefaultFactory.Required
	OptionalEmail    = DefaultFactory.OptionalEmail
	Path             = DefaultFactory.Path
	URL              = DefaultFactory.URL
)
