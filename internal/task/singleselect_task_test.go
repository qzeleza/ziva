// task/singleselect_task_test.go

package task

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/qzeleza/termos/internal/defauilt"
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

func TestSingleSelectTaskWithDefaultItemByIndex(t *testing.T) {
	title := "Выберите опцию"
	options := []string{"Опция 1", "Опция 2", "Опция 3"}

	task := NewSingleSelectTask(title, options).WithDefaultItem(2)

	assert.Equal(t, 2, task.cursor, "Курсор должен указывать на элемент с индексом 2")
	assert.Equal(t, "Опция 3", task.GetSelected(), "Выбранным значением по умолчанию должна быть 'Опция 3'")
}

func TestSingleSelectTaskWithDefaultItemByValue(t *testing.T) {
	title := "Выберите опцию"
	options := []string{"Опция 1", "Опция 2", "Опция 3"}

	task := NewSingleSelectTask(title, options).WithDefaultItem("Опция 2")

	assert.Equal(t, 1, task.cursor, "Курсор должен указывать на элемент с индексом 1")
	assert.Equal(t, "Опция 2", task.GetSelected(), "Выбранным значением по умолчанию должна быть 'Опция 2'")
}

func TestSingleSelectTaskLeftCancels(t *testing.T) {
	title := "Выберите опцию"
	options := []string{"Опция 1", "Опция 2"}

	task := NewSingleSelectTask(title, options)

	updated, _ := task.Update(tea.KeyMsg{Type: tea.KeyLeft})
	canceledTask, ok := updated.(*SingleSelectTask)
	assert.True(t, ok, "Обновленная задача должна быть типа *SingleSelectTask")
	assert.True(t, canceledTask.IsDone(), "Задача должна завершиться после нажатия ←")
	if err := canceledTask.Error(); assert.NotNil(t, err, "Ошибка должна быть установлена") {
		assert.Equal(t, defauilt.ErrorMsgCanceled, err.Error())
	}
}

func TestSingleSelectTaskRightSelects(t *testing.T) {
	title := "Выберите опцию"
	options := []string{"Опция 1", "Опция 2"}

	task := NewSingleSelectTask(title, options)

	updated, _ := task.Update(tea.KeyMsg{Type: tea.KeyRight})
	selectedTask, ok := updated.(*SingleSelectTask)
	assert.True(t, ok, "Обновленная задача должна быть типа *SingleSelectTask")
	assert.True(t, selectedTask.IsDone(), "Задача должна завершиться после нажатия →")
	assert.Equal(t, options[0], selectedTask.GetSelected(), "Должна быть выбрана текущая опция")
}

// TestSingleSelectTaskBoundaries проверяет обработку граничных случаев
func TestSingleSelectTaskBoundaries(t *testing.T) {
	// Создаем задачу SingleSelectTask
	title := "Выберите опцию"
	options := []string{"Опция 1", "Опция 2", "Опция 3"}
	selectTask := NewSingleSelectTask(title, options)

	// Проверяем, что курсор не выходит за нижнюю границу
	// Нажимаем 'down' несколько раз, чтобы достичь нижней границы
	for i := 0; i < len(options)+2; i++ {
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

func TestSingleSelectTaskDisabledItems(t *testing.T) {
	title := "Выберите опцию"
	options := []string{"Опция 1", "Опция 2", "Опция 3"}

	task := NewSingleSelectTask(title, options)
	task = task.WithItemsDisabled([]int{1})

	assert.Equal(t, 0, task.GetSelectedIndex(), "Курсор должен оставаться на первом доступном элементе")

	updated, _ := task.Update(tea.KeyMsg{Type: tea.KeyDown})
	task, _ = updated.(*SingleSelectTask)
	assert.Equal(t, 2, task.GetSelectedIndex(), "Курсор должен перепрыгивать через отключённый элемент")

	task.WithDefaultItem(1)
	assert.Equal(t, 2, task.GetSelectedIndex(), "Значение по умолчанию не должно указывать на выключенный элемент")

	task = task.WithItemsDisabled(nil)
	updated, _ = task.Update(tea.KeyMsg{Type: tea.KeyUp})
	task, _ = updated.(*SingleSelectTask)
	assert.Equal(t, 1, task.GetSelectedIndex(), "После включения элемента курсор должен уметь на него переходить")
}

func TestSingleSelectTaskViewportIndicators(t *testing.T) {
	title := "Выберите опцию"
	options := []string{"Опция 1", "Опция 2", "Опция 3", "Опция 4"}

	task := NewSingleSelectTask(title, options).WithViewport(2)
	for i := 0; i < 3; i++ {
		updated, _ := task.Update(tea.KeyMsg{Type: tea.KeyDown})
		task, _ = updated.(*SingleSelectTask)
	}
	viewWithCounters := task.View(80)
	assert.True(t, strings.Contains(viewWithCounters, "▲  1"), "Индикатор должен содержать двойной пробел и количество")

	task = NewSingleSelectTask(title, options).WithViewport(2, false)
	for i := 0; i < 3; i++ {
		updated, _ := task.Update(tea.KeyMsg{Type: tea.KeyDown})
		task, _ = updated.(*SingleSelectTask)
	}
	viewWithoutCounters := task.View(80)
	assert.Contains(t, viewWithoutCounters, "▲", "Индикатор должен содержать символ стрелки")
	assert.NotContains(t, viewWithoutCounters, "above", "При отключении счётчиков текст не должен отображаться")
	assert.NotContains(t, viewWithoutCounters, "выше", "При отключении счётчиков текст не должен отображаться")
}
