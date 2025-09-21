// Package validation предоставляет интерфейсы и реализации для валидации различных типов данных.
// Этот пакет следует принципу Single Responsibility и позволяет легко добавлять новые типы валидаторов.
package validation

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/qzeleza/termos/internal/defauilt"
	"github.com/qzeleza/termos/internal/performance"
)

// Validator определяет интерфейс для валидации входных данных
type Validator interface {
	// Validate проверяет входную строку и возвращает ошибку если валидация не прошла
	Validate(input string) error

	// Description возвращает описание требований для валидации
	Description() string
}

// ValidatorFunc тип функции для создания простых валидаторов
type ValidatorFunc func(string) error

// Validate реализует интерфейс Validator для ValidatorFunc
func (vf ValidatorFunc) Validate(input string) error {
	return vf(input)
}

// Description возвращает базовое описание для функции-валидатора
func (vf ValidatorFunc) Description() string {
	return defauilt.ValidatorCustomValidation
}

// PasswordValidator валидатор для паролей
type PasswordValidator struct {
	MinLength int
}

// NewPasswordValidator создает новый валидатор паролей
func NewPasswordValidator(minLength int) *PasswordValidator {
	if minLength < 1 {
		minLength = 8 // значение по умолчанию
	}
	return &PasswordValidator{MinLength: minLength}
}

// Validate проверяет надежность пароля
func (pv *PasswordValidator) Validate(password string) error {
	if len(password) < pv.MinLength {
		return fmt.Errorf(defauilt.ValidatorPasswordMinLength, pv.MinLength)
	}

	// Проверка на наличие кириллических символов
	hasCyrillic := false
	for _, char := range password {
		// Диапазон кириллических символов в Unicode
		if (char >= '\u0400' && char <= '\u04FF') || (char >= '\u0500' && char <= '\u052F') {
			hasCyrillic = true
			break
		}
	}

	if hasCyrillic {
		return errors.New(defauilt.ValidatorPasswordCyrillic)
	}

	hasDigit := false
	hasSpecial := false
	hasUpper := false
	hasLower := false

	for _, char := range password {
		switch {
		case '0' <= char && char <= '9':
			hasDigit = true
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case strings.ContainsRune("!@#$%^&*()-_=+[]{}|;:'\",.<>/?\\~`", char):
			hasSpecial = true
		}
	}

	var missing []string
	if !hasDigit {
		missing = append(missing, defauilt.ValidatorPasswordRequirementDigits)
	}
	if !hasSpecial {
		missing = append(missing, defauilt.ValidatorPasswordRequirementSpecial)
	}
	if !hasUpper {
		missing = append(missing, defauilt.ValidatorPasswordRequirementUpper)
	}
	if !hasLower {
		missing = append(missing, defauilt.ValidatorPasswordRequirementLower)
	}

	if len(missing) > 0 {
		return fmt.Errorf(defauilt.ValidatorPasswordMissingRequirements, performance.JoinEfficient(missing, defauilt.ValidatorListSeparator))
	}

	return nil
}

// Description возвращает описание требований к паролю
func (pv *PasswordValidator) Description() string {
	return fmt.Sprintf(defauilt.ValidatorPasswordDescription, pv.MinLength)
}

// EmailValidator валидатор для email адресов
type EmailValidator struct {
	pattern *regexp.Regexp
}

// NewEmailValidator создает новый валидатор email
func NewEmailValidator() *EmailValidator {
	// RFC 5322 совместимая регулярка (упрощенная версия)
	pattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return &EmailValidator{pattern: pattern}
}

// Validate проверяет корректность email адреса
func (ev *EmailValidator) Validate(email string) error {
	if !ev.pattern.MatchString(email) {
		return errors.New(defauilt.ValidatorEmailInvalid)
	}
	return nil
}

// Description возвращает описание требований к email
func (ev *EmailValidator) Description() string {
	return defauilt.ValidatorEmailDescription
}

// NumberValidator валидатор для чисел в диапазоне
type NumberValidator struct {
	Min int
	Max int
}

// NewNumberValidator создает новый валидатор чисел
func NewNumberValidator(min, max int) *NumberValidator {
	return &NumberValidator{Min: min, Max: max}
}

// Validate проверяет, что строка содержит число в заданном диапазоне
func (nv *NumberValidator) Validate(s string) error {
	num, err := strconv.Atoi(s)
	if err != nil {
		return errors.New(defauilt.ValidatorNumberInvalid)
	}
	if num < nv.Min || num > nv.Max {
		return fmt.Errorf(defauilt.ValidatorNumberRange, nv.Min, nv.Max)
	}
	return nil
}

// Description возвращает описание требований к числу
func (nv *NumberValidator) Description() string {
	return fmt.Sprintf(defauilt.ValidatorNumberDescription, nv.Min, nv.Max)
}

// IPValidator валидатор для IP адресов
type IPValidator struct {
	allowIPv4 bool
	allowIPv6 bool
}

// NewIPValidator создает новый валидатор IP адресов
func NewIPValidator(allowIPv4, allowIPv6 bool) *IPValidator {
	// По умолчанию разрешаем IPv4 если ничего не указано
	if !allowIPv4 && !allowIPv6 {
		allowIPv4 = true
	}
	return &IPValidator{allowIPv4: allowIPv4, allowIPv6: allowIPv6}
}

// NewIPv4Validator создает валидатор только для IPv4
func NewIPv4Validator() *IPValidator {
	return NewIPValidator(true, false)
}

// NewIPv6Validator создает валидатор только для IPv6
func NewIPv6Validator() *IPValidator {
	return NewIPValidator(false, true)
}

// Validate проверяет корректность IP-адреса
func (iv *IPValidator) Validate(ip string) error {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return errors.New(defauilt.ValidatorIPInvalid)
	}

	// Проверяем тип IP адреса
	isIPv4 := parsedIP.To4() != nil
	isIPv6 := !isIPv4

	if isIPv4 && !iv.allowIPv4 {
		return errors.New(defauilt.ValidatorIPv4NotAllowed)
	}
	if isIPv6 && !iv.allowIPv6 {
		return errors.New(defauilt.ValidatorIPv6NotAllowed)
	}

	return nil
}

// Description возвращает описание требований к IP адресу
func (iv *IPValidator) Description() string {
	if iv.allowIPv4 && iv.allowIPv6 {
		return defauilt.ValidatorIPBothDescription
	} else if iv.allowIPv4 {
		return defauilt.ValidatorIPv4Description
	} else if iv.allowIPv6 {
		return defauilt.ValidatorIPv6Description
	}
	return defauilt.ValidatorIPGenericDescription
}

// DomainValidator валидатор для доменных имен
type DomainValidator struct {
	pattern *regexp.Regexp
}

// NewDomainValidator создает новый валидатор доменных имен
func NewDomainValidator() *DomainValidator {
	// RFC 1035 совместимая регулярка для доменных имен
	pattern := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	return &DomainValidator{pattern: pattern}
}

// Validate проверяет корректность доменного имени
func (dv *DomainValidator) Validate(domain string) error {
	if !dv.pattern.MatchString(domain) {
		return errors.New(defauilt.ValidatorDomainInvalid)
	}
	return nil
}

// Description возвращает описание требований к домену
func (dv *DomainValidator) Description() string {
	return defauilt.ValidatorDomainDescription
}

// TextValidator базовый валидатор для текста
type TextValidator struct {
	MinLength int
	MaxLength int
	Pattern   *regexp.Regexp
}

// NewTextValidator создает новый валидатор текста
func NewTextValidator(minLen, maxLen int) *TextValidator {
	return &TextValidator{MinLength: minLen, MaxLength: maxLen}
}

// WithPattern добавляет регулярное выражение к валидатору текста
func (tv *TextValidator) WithPattern(pattern string) *TextValidator {
	tv.Pattern = regexp.MustCompile(pattern)
	return tv
}

// Validate проверяет текст по заданным критериям
func (tv *TextValidator) Validate(text string) error {
	if tv.MinLength > 0 && len(text) < tv.MinLength {
		return fmt.Errorf(defauilt.ValidatorTextMin, tv.MinLength)
	}

	if tv.MaxLength > 0 && len(text) > tv.MaxLength {
		return fmt.Errorf(defauilt.ValidatorTextMax, tv.MaxLength)
	}

	if tv.Pattern != nil && !tv.Pattern.MatchString(text) {
		return errors.New(defauilt.ValidatorTextPattern)
	}

	return nil
}

// Description возвращает описание требований к тексту
func (tv *TextValidator) Description() string {
	desc := defauilt.ValidatorTextBase
	if tv.MinLength > 0 || tv.MaxLength > 0 {
		if tv.MinLength > 0 && tv.MaxLength > 0 {
			desc += fmt.Sprintf(defauilt.ValidatorTextRange, tv.MinLength, tv.MaxLength)
		} else if tv.MinLength > 0 {
			desc += fmt.Sprintf(defauilt.ValidatorTextMinOnly, tv.MinLength)
		} else {
			desc += fmt.Sprintf(defauilt.ValidatorTextMaxOnly, tv.MaxLength)
		}
	}
	return desc
}

// CompositeValidator объединяет несколько валидаторов
type CompositeValidator struct {
	validators []Validator
	mode       CompositeMode
}

// CompositeMode определяет режим работы композитного валидатора
type CompositeMode int

const (
	// AllMustPass - все валидаторы должны пройти проверку
	AllMustPass CompositeMode = iota
	// AnyCanPass - достаточно одного прошедшего валидатора
	AnyCanPass
)

// NewCompositeValidator создает новый композитный валидатор
func NewCompositeValidator(mode CompositeMode, validators ...Validator) *CompositeValidator {
	return &CompositeValidator{
		validators: validators,
		mode:       mode,
	}
}

// Validate выполняет валидацию согласно выбранному режиму
func (cv *CompositeValidator) Validate(input string) error {
	if len(cv.validators) == 0 {
		return nil
	}

	var errors []string

	switch cv.mode {
	case AllMustPass:
		for _, validator := range cv.validators {
			if err := validator.Validate(input); err != nil {
				errors = append(errors, err.Error())
			}
		}
		if len(errors) > 0 {
			return fmt.Errorf(defauilt.ValidatorCompositeAllErrors, performance.JoinEfficient(errors, defauilt.ValidatorCompositeAllSeparator))
		}
		return nil

	case AnyCanPass:
		for _, validator := range cv.validators {
			if err := validator.Validate(input); err == nil {
				return nil // Один из валидаторов прошел успешно
			} else {
				errors = append(errors, err.Error())
			}
		}
		return fmt.Errorf(defauilt.ValidatorCompositeNonePassed, performance.JoinEfficient(errors, defauilt.ValidatorCompositeAllSeparator))

	default:
		return fmt.Errorf(defauilt.ValidatorCompositeUnknownMode)
	}
}

// Description возвращает описание композитного валидатора
func (cv *CompositeValidator) Description() string {
	if len(cv.validators) == 0 {
		return defauilt.ValidatorCompositeNoValidation
	}

	var descriptions []string
	for _, validator := range cv.validators {
		descriptions = append(descriptions, validator.Description())
	}

	switch cv.mode {
	case AllMustPass:
		return fmt.Sprintf(defauilt.ValidatorCompositeAllDescription, performance.JoinEfficient(descriptions, defauilt.ValidatorCompositeAllSeparator))
	case AnyCanPass:
		return fmt.Sprintf(defauilt.ValidatorCompositeAnyDescription, performance.JoinEfficient(descriptions, defauilt.ValidatorCompositeAnySeparator))
	default:
		return defauilt.ValidatorCompositeDescription
	}
}
