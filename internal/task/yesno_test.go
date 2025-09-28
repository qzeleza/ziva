package task

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// TestYesNoTaskErrorHandling проверяет правильность обработки ошибок в YesNoTask
func TestYesNoTaskErrorHandling(t *testing.T) {
	// Создаем задачу
	task := NewYesNoTask("Тестовый вопрос", "Согласны ли вы?")

	// В начале задача не завершена и нет ошибок
	assert.False(t, task.IsDone(), "Задача не должна быть завершена в начале")
	assert.False(t, task.HasError(), "В начале не должно быть ошибок")

	// Симулируем навигацию к опции "Нет" (индекс 1)
	task.Update(tea.KeyMsg{Type: tea.KeyDown}) // Переходим к "Нет"

	// Симулируем выбор опции "Нет"
	task.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Проверяем, что задача завершена
	assert.True(t, task.IsDone(), "Задача должна быть завершена после выбора")

	// Проверяем, что выбрано "Нет"
	assert.True(t, task.IsNo(), "Должна быть выбрана опция 'Нет'")
	assert.False(t, task.IsYes(), "Не должна быть выбрана опция 'Да'")

	// Проверяем, что есть ошибка (для статистики)
	assert.True(t, task.HasError(), "При выборе 'Нет' должна быть установлена ошибка для статистики")
	assert.NotNil(t, task.Error(), "Ошибка не должна быть nil")
	assert.Contains(t, task.Error().Error(), "Нет", "Ошибка должна содержать информацию о выборе 'Нет'")

	// Проверяем, что очередь не останавливается
	assert.False(t, task.StopOnError(), "Выбор 'Нет' не должен останавливать очередь")
}

// TestYesNoTaskSuccessHandling проверяет правильность обработки успешного выбора "Да"
func TestYesNoTaskSuccessHandling(t *testing.T) {
	// Создаем задачу
	task := NewYesNoTask("Тестовый вопрос", "Согласны ли вы?")

	// Симулируем выбор опции "Да" (она выбрана по умолчанию, индекс 0)
	task.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Проверяем, что задача завершена
	assert.True(t, task.IsDone(), "Задача должна быть завершена после выбора")

	// Проверяем, что выбрано "Да"
	assert.True(t, task.IsYes(), "Должна быть выбрана опция 'Да'")
	assert.False(t, task.IsNo(), "Не должна быть выбрана опция 'Нет'")

	// Проверяем, что нет ошибки
	assert.False(t, task.HasError(), "При выборе 'Да' не должно быть ошибки")
	assert.Nil(t, task.Error(), "Ошибка должна быть nil")
}

// TestYesNoTaskOnYesHandler проверяет вызов обработчика для ответа "Да"
func TestYesNoTaskOnYesHandler(t *testing.T) {
	var called bool
	task := NewYesNoTask("Тестовый вопрос", "Согласны ли вы?")
	task.OnYes(func() error {
		called = true
		return nil
	})

	task.Update(tea.KeyMsg{Type: tea.KeyEnter})

	assert.True(t, called, "Обработчик 'Да' должен вызываться")
	assert.True(t, task.IsDone(), "Задача должна завершиться после выбора")
	assert.True(t, task.IsYes(), "Опция 'Да' должна быть выбрана")
	assert.False(t, task.HasError(), "После успешного обработчика не должно быть ошибки")
	assert.True(t, task.StopOnError(), "Флаг StopOnError должен оставаться по умолчанию")
}

// TestYesNoTaskOnNoHandlerError проверяет обработку ошибки из обработчика "Нет"
func TestYesNoTaskOnNoHandlerError(t *testing.T) {
	var called bool
	expectedErr := errors.New("custom no error")
	task := NewYesNoTask("Тестовый вопрос", "Согласны ли вы?")
	task.OnNo(func() error {
		called = true
		return expectedErr
	})

	task.Update(tea.KeyMsg{Type: tea.KeyDown})
	task.Update(tea.KeyMsg{Type: tea.KeyEnter})

	assert.True(t, called, "Обработчик 'Нет' должен вызываться")
	assert.True(t, task.IsDone(), "Задача должна завершиться после выбора")
	assert.True(t, task.IsNo(), "Опция 'Нет' должна быть выбрана")
	assert.True(t, task.HasError(), "Ошибка из обработчика должна сохраняться")
	assert.Equal(t, expectedErr, task.Error(), "Ошибка должна совпадать с возвращенной обработчиком")
	assert.True(t, task.StopOnError(), "Ошибка обработчика должна останавливать очередь")
	assert.Contains(t, task.FinalView(80), expectedErr.Error(), "Финальное представление должно содержать текст ошибки")
}

// TestYesNoTaskCancelHandling проверяет правильность обработки отмены
func TestYesNoTaskCancelHandling(t *testing.T) {
	// Создаем задачу
	task := NewYesNoTask("Тестовый вопрос", "Согласны ли вы?")

	// Симулируем отмену (нажатие 'q')
	task.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	// Проверяем, что задача завершена
	assert.True(t, task.IsDone(), "Задача должна быть завершена после отмены")

	// Проверяем, что есть ошибка отмены
	assert.True(t, task.HasError(), "При отмене должна быть ошибка")
	assert.NotNil(t, task.Error(), "Ошибка не должна быть nil")
	assert.Contains(t, task.Error().Error(), "отменено", "Ошибка должна содержать информацию об отмене")

	// Проверяем, что очередь останавливается при отмене
	assert.True(t, task.StopOnError(), "Отмена должна останавливать очередь")
}
