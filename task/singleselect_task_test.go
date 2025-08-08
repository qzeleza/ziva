// task/singleselect_task_test.go

package task

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// TestSingleSelectTaskCreation проверяет корректность создания задачи SingleSelectTask
func TestSingleSelectTaskCreation(t *testing.T) {
	// Создаем задачу SingleSelectTask
	title := "Выберите опцию"
	options := []string{"Опция 1", "Опция 2", "Опция 3"}
	
	// Без указания индекса по умолчанию
	selectTask := NewSingleSelectTask(title, options)
	
	// Проверяем, что задача создана корректно
	assert.NotNil(t, selectTask, "Задача не должна быть nil")
	assert.Equal(t, title, selectTask.Title(), "Заголовок задачи должен соответствовать переданному значению")
	assert.False(t, selectTask.IsDone(), "Новая задача не должна быть отмечена как завершенная")
	
	// Создаем еще одну задачу
	selectTaskWithDefault := NewSingleSelectTask(title, options)
	
	// Проверяем, что задача создана корректно
	assert.NotNil(t, selectTaskWithDefault, "Задача не должна быть nil")
	assert.Equal(t, title, selectTaskWithDefault.Title(), "Заголовок задачи должен соответствовать переданному значению")
	assert.False(t, selectTaskWithDefault.IsDone(), "Новая задача не должна быть отмечена как завершенная")
}

// TestSingleSelectTaskUpdate проверяет обработку различных клавиш в методе Update
func TestSingleSelectTaskUpdate(t *testing.T) {
	// Создаем задачу SingleSelectTask
	title := "Выберите опцию"
	options := []string{"Опция 1", "Опция 2", "Опция 3"}
	selectTask := NewSingleSelectTask(title, options)
	
	// Проверяем обработку клавиши 'down'
	updatedTask, _ := selectTask.Update(tea.KeyMsg{Type: tea.KeyDown})
	selectTaskAfterDown, ok := updatedTask.(*SingleSelectTask)
	assert.True(t, ok, "Обновленная задача должна быть типа *SingleSelectTask")
	assert.False(t, selectTaskAfterDown.IsDone(), "Задача не должна быть отмечена как завершенная после нажатия 'down'")
	
	// Проверяем обработку клавиши 'up'
	updatedTask, _ = selectTaskAfterDown.Update(tea.KeyMsg{Type: tea.KeyUp})
	selectTaskAfterUp, ok := updatedTask.(*SingleSelectTask)
	assert.True(t, ok, "Обновленная задача должна быть типа *SingleSelectTask")
	assert.False(t, selectTaskAfterUp.IsDone(), "Задача не должна быть отмечена как завершенная после нажатия 'up'")
	
	// Проверяем обработку клавиши 'enter'
	updatedTask, _ = selectTaskAfterUp.Update(tea.KeyMsg{Type: tea.KeyEnter})
	selectTaskAfterEnter, ok := updatedTask.(*SingleSelectTask)
	assert.True(t, ok, "Обновленная задача должна быть типа *SingleSelectTask")
	assert.True(t, selectTaskAfterEnter.IsDone(), "Задача должна быть отмечена как завершенная после нажатия 'enter'")
	
	// Проверяем, что выбрана правильная опция
	finalView := selectTaskAfterEnter.FinalView(80)
	assert.Contains(t, finalView, options[0], "Значение задачи должно содержать выбранную опцию")
}

// TestSingleSelectTaskView проверяет отображение задачи в активном состоянии
func TestSingleSelectTaskView(t *testing.T) {
	// Создаем задачу SingleSelectTask
	title := "Выберите опцию"
	options := []string{"Опция 1", "Опция 2", "Опция 3"}
	selectTask := NewSingleSelectTask(title, options)
	
	// Проверяем, что View содержит заголовок и опции
	view := selectTask.View(80)
	assert.Contains(t, view, title, "View должен содержать заголовок")
	for _, option := range options {
		assert.Contains(t, view, option, "View должен содержать опцию")
	}
	
	// Проверяем, что после завершения задачи View возвращает пустую строку
	updatedTask, _ := selectTask.Update(tea.KeyMsg{Type: tea.KeyEnter})
	selectTaskDone, _ := updatedTask.(*SingleSelectTask)
	assert.Equal(t, "", selectTaskDone.View(80), "View должен возвращать пустую строку для завершенной задачи")
}

// TestSingleSelectTaskWithDefaultIndex проверяет работу с выбором определенного индекса
func TestSingleSelectTaskWithDefaultIndex(t *testing.T) {
	// Создаем задачу SingleSelectTask
	title := "Выберите опцию"
	options := []string{"Опция 1", "Опция 2", "Опция 3"}
	defaultIndex := 1
	selectTask := NewSingleSelectTask(title, options)
	
	// Устанавливаем курсор на нужный индекс
	// Нажимаем 'down' один раз, чтобы перейти к опции с индексом 1
	updatedTask, _ := selectTask.Update(tea.KeyMsg{Type: tea.KeyDown})
	selectTask, _ = updatedTask.(*SingleSelectTask)
	
	// Нажимаем Enter для завершения задачи
	updatedTask, _ = selectTask.Update(tea.KeyMsg{Type: tea.KeyEnter})
	selectTaskDone, _ := updatedTask.(*SingleSelectTask)
	assert.True(t, selectTaskDone.IsDone(), "Задача должна быть отмечена как завершенная после нажатия 'enter'")
	
	// Проверяем, что выбрана правильная опция
	finalView := selectTaskDone.FinalView(80)
	assert.Contains(t, finalView, options[defaultIndex], "Значение задачи должно содержать опцию с выбранным индексом")
}

// TestSingleSelectTaskBoundaries проверяет обработку граничных случаев
func TestSingleSelectTaskBoundaries(t *testing.T) {
	// Создаем задачу SingleSelectTask
	title := "Выберите опцию"
	options := []string{"Опция 1", "Опция 2", "Опция 3"}
	selectTask := NewSingleSelectTask(title, options)
	
	// Проверяем, что курсор не выходит за нижнюю границу
	// Нажимаем 'down' несколько раз, чтобы достичь нижней границы
	for i := 0; i < len(options) + 2; i++ {
		updatedTask, _ := selectTask.Update(tea.KeyMsg{Type: tea.KeyDown})
		selectTask, _ = updatedTask.(*SingleSelectTask)
	}
	
	// Проверяем, что после достижения нижней границы и нажатия Enter
	// выбрана последняя опция
	updatedTask, _ := selectTask.Update(tea.KeyMsg{Type: tea.KeyEnter})
	selectTaskDone, _ := updatedTask.(*SingleSelectTask)
	finalView := selectTaskDone.FinalView(80)
	assert.Contains(t, finalView, options[len(options)-1], "Значение задачи должно содержать последнюю опцию")
	
	// Создаем новую задачу для проверки верхней границы
	selectTask = NewSingleSelectTask(title, options)
	
	// Нажимаем 'up' несколько раз, чтобы попытаться выйти за верхнюю границу
	for i := 0; i < 3; i++ {
		updatedTask, _ := selectTask.Update(tea.KeyMsg{Type: tea.KeyUp})
		selectTask, _ = updatedTask.(*SingleSelectTask)
	}
	
	// Проверяем, что после попытки выйти за верхнюю границу и нажатия Enter
	// выбрана первая опция
	updatedTask, _ = selectTask.Update(tea.KeyMsg{Type: tea.KeyEnter})
	selectTaskDone, _ = updatedTask.(*SingleSelectTask)
	finalView = selectTaskDone.FinalView(80)
	assert.Contains(t, finalView, options[0], "Значение задачи должно содержать первую опцию")
}
