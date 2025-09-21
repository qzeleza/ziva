package task_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/internal/common"
	"github.com/qzeleza/termos/internal/defaults"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
	"github.com/stretchr/testify/assert"
)

// TestYesNoTaskInQueueStatistics проверяет интеграцию YesNoTask с очередью и подсчетом статистики
func TestYesNoTaskInQueueStatistics(t *testing.T) {
	// Создаем задачи: одну с "Да", одну с "Нет"
	task1 := task.NewYesNoTask("Первый вопрос", "Согласны с первым пунктом?")
	task2 := task.NewYesNoTask("Второй вопрос", "Согласны со вторым пунктом?")
	task3 := task.NewYesNoTask("Третий вопрос", "Согласны с третьим пунктом?")

	// Симулируем ответы пользователя
	// Задача 1: выбираем "Да" (по умолчанию выбрана опция 0)
	task1.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Задача 2: выбираем "Нет" (навигация к индексу 1)
	task2.Update(tea.KeyMsg{Type: tea.KeyDown})  // Переходим к "Нет"
	task2.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Выбираем "Нет"

	// Задача 3: выбираем "Да" (по умолчанию выбрана опция 0)
	task3.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Проверяем состояние задач
	assert.True(t, task1.IsDone(), "Задача 1 должна быть завершена")
	assert.True(t, task2.IsDone(), "Задача 2 должна быть завершена")
	assert.True(t, task3.IsDone(), "Задача 3 должна быть завершена")

	// Проверяем ошибки
	assert.False(t, task1.HasError(), "Задача 1 (Да) не должна иметь ошибку")
	assert.True(t, task2.HasError(), "Задача 2 (Нет) должна иметь ошибку для статистики")
	assert.False(t, task3.HasError(), "Задача 3 (Да) не должна иметь ошибку")

	// Создаем очередь с этими задачами
	model := query.New("Тест статистики YesNoTask").WithSummary(true)
	tasks := []common.Task{task1, task2, task3}
	model.AddTasks(tasks)

	// Симулируем завершение всех задач
	model.Update(tea.WindowSizeMsg{Width: 80, Height: 24}) // Устанавливаем размер окна

	// Имитируем завершение очереди
	for i := 0; i < len(tasks); i++ {
		// Обновляем модель для каждой завершенной задачи
		model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	}

	// Получаем финальное представление
	view := model.View()

	// Проверяем, что в статистике правильно подсчитываются ошибки
	// Должно быть: 1 успешных из 3 всего и статус ПРОБЛЕМА
	assert.Contains(t, view, "Успешно завершено 1 из 3 задач", "Статистика должна показывать 1 успешных из 3 всего")
	assert.Contains(t, view, defaults.StatusProblem, "Статус должен показывать проблему из-за task2")
}

// TestYesNoTaskInQueueAllSuccess проверяет очередь, где все YesNoTask выбирают "Да"
func TestYesNoTaskInQueueAllSuccess(t *testing.T) {
	// Создаем задачи
	task1 := task.NewYesNoTask("Первый вопрос", "Согласны?")
	task2 := task.NewYesNoTask("Второй вопрос", "Подтверждаете?")

	// Симулируем выбор "Да" для обеих задач
	task1.Update(tea.KeyMsg{Type: tea.KeyEnter})
	task2.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Проверяем, что нет ошибок
	assert.False(t, task1.HasError(), "Задача 1 не должна иметь ошибку")
	assert.False(t, task2.HasError(), "Задача 2 не должна иметь ошибку")

	// Создаем очередь
	model := query.New("Тест успешной очереди").WithSummary(true)
	tasks := []common.Task{task1, task2}
	model.AddTasks(tasks)

	// Имитируем завершение очереди
	model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	for i := 0; i < len(tasks); i++ {
		model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	}

	// Получаем финальное представление
	view := model.View()

	// Проверяем успешную статистику
	assert.Contains(t, view, "Успешно завершено 1 из 2 задач", "Статистика должна показывать успешных из всего")
}

// TestYesNoTaskInQueueAllNo проверяет очередь, где все YesNoTask выбирают "Нет"
func TestYesNoTaskInQueueAllNo(t *testing.T) {
	// Создаем задачи
	task1 := task.NewYesNoTask("Первый вопрос", "Согласны?")
	task2 := task.NewYesNoTask("Второй вопрос", "Подтверждаете?")

	// Симулируем выбор "Нет" для обеих задач
	task1.Update(tea.KeyMsg{Type: tea.KeyDown}) // К "Нет"
	task1.Update(tea.KeyMsg{Type: tea.KeyEnter})
	task2.Update(tea.KeyMsg{Type: tea.KeyDown}) // К "Нет"
	task2.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Проверяем, что есть ошибки
	assert.True(t, task1.HasError(), "Задача 1 должна иметь ошибку")
	assert.True(t, task2.HasError(), "Задача 2 должна иметь ошибку")

	// Создаем очередь
	model := query.New("Тест очереди с отказами").WithSummary(true)
	tasks := []common.Task{task1, task2}
	model.AddTasks(tasks)

	// Имитируем завершение очереди
	model.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	for i := 0; i < len(tasks); i++ {
		model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	}

	// Получаем финальное представление
	view := model.View()

	// Проверяем статистику с ошибками
	assert.Contains(t, view, "Успешно завершено 0 из 2 задач", "Статистика должна показывать 0 успешных из 2 всего")
	assert.Contains(t, view, defaults.StatusProblem, "Статус должен показывать проблему")
}
