package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

// TestSetErrorColor проверяет функцию SetErrorColor
func TestSetErrorColor(t *testing.T) {
	// Сохраняем исходные стили
	originalMessageStyle := ErrorMessageStyle
	originalStatusStyle := ErrorStatusStyle

	// Устанавливаем новый цвет
	testColor := lipgloss.Color("#ff0000") // Красный
	SetErrorColor(testColor, testColor)

	// Проверяем, что цвета изменились
	currentMessageStyle := GetErrorMessageStyle()
	currentStatusStyle := GetErrorStatusStyle()

	// Проверяем цвет (извлекаем из стилей)
	assert.Equal(t, testColor, currentMessageStyle.GetForeground())
	assert.Equal(t, testColor, currentStatusStyle.GetForeground())

	// Проверяем, что StatusStyle сохранил свой Bold
	assert.True(t, currentStatusStyle.GetBold())

	// Восстанавливаем исходные стили
	ErrorMessageStyle = originalMessageStyle
	ErrorStatusStyle = originalStatusStyle
}

// TestResetErrorColors проверяет функцию ResetErrorColors
func TestResetErrorColors(t *testing.T) {
	// Устанавливаем кастомный цвет
	testColor := lipgloss.Color("#00ff00") // Зеленый
	SetErrorColor(testColor, testColor)

	// Проверяем, что цвет изменился
	assert.Equal(t, testColor, GetErrorMessageStyle().GetForeground())
	assert.Equal(t, testColor, GetErrorStatusStyle().GetForeground())

	// Сбрасываем к значениям по умолчанию
	ResetErrorColors()

	// Проверяем, что цвета сбросились к значениям по умолчанию
	assert.Equal(t, ColorDarkYellow, GetErrorMessageStyle().GetForeground())
	assert.Equal(t, ColorBrightYellow, GetErrorStatusStyle().GetForeground())
	assert.True(t, GetErrorStatusStyle().GetBold())
}

// TestGetErrorStyles проверяет функции получения стилей
func TestGetErrorStyles(t *testing.T) {
	// Устанавливаем тестовый цвет
	testColor := lipgloss.Color("#0000ff") // Синий
	SetErrorColor(testColor, testColor)

	// Проверяем, что функции возвращают правильные стили
	messageStyle := GetErrorMessageStyle()
	statusStyle := GetErrorStatusStyle()

	assert.Equal(t, testColor, messageStyle.GetForeground())
	assert.Equal(t, testColor, statusStyle.GetForeground())
	assert.True(t, statusStyle.GetBold())

	// Проверяем, что стили не одинаковые (у StatusStyle есть Bold)
	assert.NotEqual(t, messageStyle.GetBold(), statusStyle.GetBold())

	// Сбрасываем стили
	ResetErrorColors()
}

// TestErrorColorPersistence проверяет, что изменения цвета сохраняются в глобальных переменных
func TestErrorColorPersistence(t *testing.T) {
	// Устанавливаем цвет
	testColor := lipgloss.Color("#ff00ff") // Фиолетовый
	SetErrorColor(testColor, testColor)

	// Проверяем через глобальные переменные
	assert.Equal(t, testColor, ErrorMessageStyle.GetForeground())
	assert.Equal(t, testColor, ErrorStatusStyle.GetForeground())

	// Проверяем через геттеры
	assert.Equal(t, testColor, GetErrorMessageStyle().GetForeground())
	assert.Equal(t, testColor, GetErrorStatusStyle().GetForeground())

	// Сбрасываем
	ResetErrorColors()
	
	// Проверяем сброс
	assert.Equal(t, ColorDarkYellow, ErrorMessageStyle.GetForeground())
	assert.Equal(t, ColorBrightYellow, ErrorStatusStyle.GetForeground())
}