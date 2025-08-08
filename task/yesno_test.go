package task

import (
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