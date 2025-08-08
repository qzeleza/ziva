// task/multiselect_task_test.go

package task

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// TestMultiSelectTaskCreation проверяет корректность создания задачи MultiSelectTask
func TestMultiSelectTaskCreation(t *testing.T) {
	// Создаем задачу MultiSelectTask
	title := "Выберите опции"
	options := []string{"Опция 1", "Опция 2", "Опция 3"}
	
	// Создаем задачу
	multiSelectTask := NewMultiSelectTask(title, options)
	
	// Проверяем, что задача создана корректно
	assert.NotNil(t, multiSelectTask, "Задача не должна быть nil")
	assert.Equal(t, title, multiSelectTask.Title(), "Заголовок задачи должен соответствовать переданному значению")
	assert.False(t, multiSelectTask.IsDone(), "Новая задача не должна быть отмечена как завершенная")
	
	// Создаем еще одну задачу
	multiSelectTaskEmpty := NewMultiSelectTask(title, options)
	
	// Проверяем, что задача создана корректно
	assert.NotNil(t, multiSelectTaskEmpty, "Задача не должна быть nil")
	assert.Equal(t, title, multiSelectTaskEmpty.Title(), "Заголовок задачи должен соответствовать переданному значению")
	assert.False(t, multiSelectTaskEmpty.IsDone(), "Новая задача не должна быть отмечена как завершенная")
}

// TestMultiSelectTaskUpdate проверяет обработку различных клавиш в методе Update
func TestMultiSelectTaskUpdate(t *testing.T) {
	// Создаем задачу MultiSelectTask
	title := "Выберите опции"
	options := []string{"Опция 1", "Опция 2", "Опция 3"}
	multiSelectTask := NewMultiSelectTask(title, options)
	
	// Проверяем обработку клавиши 'down'
	updatedTask, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeyDown})
	multiSelectTaskAfterDown, ok := updatedTask.(*MultiSelectTask)
	assert.True(t, ok, "Обновленная задача должна быть типа *MultiSelectTask")
	assert.False(t, multiSelectTaskAfterDown.IsDone(), "Задача не должна быть отмечена как завершенная после нажатия 'down'")
	
	// Проверяем обработку клавиши 'up'
	updatedTask, _ = multiSelectTaskAfterDown.Update(tea.KeyMsg{Type: tea.KeyUp})
	multiSelectTaskAfterUp, ok := updatedTask.(*MultiSelectTask)
	assert.True(t, ok, "Обновленная задача должна быть типа *MultiSelectTask")
	assert.False(t, multiSelectTaskAfterUp.IsDone(), "Задача не должна быть отмечена как завершенная после нажатия 'up'")
	
	// Проверяем обработку клавиши 'space' для выбора опции
	updatedTask, _ = multiSelectTaskAfterUp.Update(tea.KeyMsg{Type: tea.KeySpace})
	multiSelectTaskAfterSpace, ok := updatedTask.(*MultiSelectTask)
	assert.True(t, ok, "Обновленная задача должна быть типа *MultiSelectTask")
	assert.False(t, multiSelectTaskAfterSpace.IsDone(), "Задача не должна быть отмечена как завершенная после нажатия 'space'")
	
	// Проверяем обработку клавиши 'enter' для завершения задачи
	updatedTask, _ = multiSelectTaskAfterSpace.Update(tea.KeyMsg{Type: tea.KeyEnter})
	multiSelectTaskAfterEnter, ok := updatedTask.(*MultiSelectTask)
	assert.True(t, ok, "Обновленная задача должна быть типа *MultiSelectTask")
	assert.True(t, multiSelectTaskAfterEnter.IsDone(), "Задача должна быть отмечена как завершенная после нажатия 'enter'")
}

// TestMultiSelectTaskView проверяет отображение задачи в активном состоянии
func TestMultiSelectTaskView(t *testing.T) {
	// Создаем задачу MultiSelectTask
	title := "Выберите опции"
	options := []string{"Опция 1", "Опция 2", "Опция 3"}
	multiSelectTask := NewMultiSelectTask(title, options)
	
	// Проверяем, что View содержит заголовок и опции
	view := multiSelectTask.View(80)
	assert.Contains(t, view, title, "View должен содержать заголовок")
	for _, option := range options {
		assert.Contains(t, view, option, "View должен содержать опцию")
	}
	
	// Сначала выбираем опцию с помощью пробела
	updatedTask1, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeySpace})
	multiSelectTaskWithSelection, _ := updatedTask1.(*MultiSelectTask)
	
	// Затем завершаем задачу с помощью Enter
	updatedTask2, _ := multiSelectTaskWithSelection.Update(tea.KeyMsg{Type: tea.KeyEnter})
	multiSelectTaskDone, _ := updatedTask2.(*MultiSelectTask)
	
	// Проверяем, что задача завершена
	assert.True(t, multiSelectTaskDone.IsDone(), "Задача должна быть завершена после нажатия Enter")
}

// TestMultiSelectTaskFinalValue проверяет корректность финального значения
func TestMultiSelectTaskFinalValue(t *testing.T) {
	// Создаем задачу MultiSelectTask
	title := "Выберите опции"
	options := []string{"Опция 1", "Опция 2", "Опция 3"}
	multiSelectTask := NewMultiSelectTask(title, options)
	
	// Выбираем первую опцию с помощью пробела
	updatedTask1, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeySpace})
	multiSelectTaskWithSelection, _ := updatedTask1.(*MultiSelectTask)
	
	// Завершаем задачу с помощью Enter
	updatedTask2, _ := multiSelectTaskWithSelection.Update(tea.KeyMsg{Type: tea.KeyEnter})
	multiSelectTaskDone, _ := updatedTask2.(*MultiSelectTask)
	
	// Проверяем, что задача завершена и финальное значение содержит выбранную опцию
	assert.True(t, multiSelectTaskDone.IsDone(), "Задача должна быть завершена после нажатия Enter")
	
	// Проверяем финальное представление
	finalView := multiSelectTaskDone.FinalView(80)
	assert.Contains(t, finalView, "Опция 1", "FinalView должен содержать выбранную опцию")
}

// TestMultiSelectTaskBoundaries проверяет обработку граничных случаев
func TestMultiSelectTaskBoundaries(t *testing.T) {
	// Создаем задачу MultiSelectTask
	title := "Выберите опции"
	options := []string{"Опция 1", "Опция 2", "Опция 3"}
	multiSelectTask := NewMultiSelectTask(title, options)
	
	// Проверяем, что курсор не выходит за нижнюю границу
	// Нажимаем 'down' несколько раз, чтобы достичь нижней границы
	for i := 0; i < len(options) + 2; i++ {
		updatedTask, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeyDown})
		multiSelectTask, _ = updatedTask.(*MultiSelectTask)
	}
	
	// Выбираем последнюю опцию
	updatedTask, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeySpace})
	multiSelectTask, _ = updatedTask.(*MultiSelectTask)
	
	// Завершаем задачу
	updatedTask, _ = multiSelectTask.Update(tea.KeyMsg{Type: tea.KeyEnter})
	multiSelectTaskDone, _ := updatedTask.(*MultiSelectTask)
	
	// Проверяем, что задача завершена
	assert.True(t, multiSelectTaskDone.IsDone(), "Задача должна быть завершена")
	// Проверяем финальное представление
	finalView := multiSelectTaskDone.FinalView(80)
	assert.Contains(t, finalView, "Опция 3", "FinalView должен содержать выбранную опцию")
	
	// Создаем новую задачу для проверки верхней границы
	multiSelectTask = NewMultiSelectTask(title, options)
	
	// Нажимаем 'up' несколько раз, чтобы попытаться выйти за верхнюю границу
	for i := 0; i < 3; i++ {
		updatedTask, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeyUp})
		multiSelectTask, _ = updatedTask.(*MultiSelectTask)
	}
	
	// Выбираем первую опцию
	updatedTask, _ = multiSelectTask.Update(tea.KeyMsg{Type: tea.KeySpace})
	multiSelectTask, _ = updatedTask.(*MultiSelectTask)
	
	// Завершаем задачу
	updatedTask, _ = multiSelectTask.Update(tea.KeyMsg{Type: tea.KeyEnter})
	multiSelectTaskDone, _ = updatedTask.(*MultiSelectTask)
	
	// Проверяем, что задача завершена
	assert.True(t, multiSelectTaskDone.IsDone(), "Задача должна быть завершена")
	// Проверяем финальное представление
	view := multiSelectTaskDone.FinalView(80)
	assert.Contains(t, view, "Опция 1", "FinalView должен содержать выбранную опцию")
}

// TestMultiSelectTaskWithSelectAll проверяет функциональность "Выбрать все"
func TestMultiSelectTaskWithSelectAll(t *testing.T) {
	// Создаем задачу MultiSelectTask с опцией "Выбрать все"
	title := "Выберите компоненты"
	options := []string{"API", "Frontend", "Database", "Worker"}
	multiSelectTask := NewMultiSelectTask(title, options).WithSelectAll()
	
	// Проверяем, что опция "Выбрать все" включена
	assert.True(t, multiSelectTask.hasSelectAll, "Опция 'Выбрать все' должна быть включена")
	assert.Equal(t, -1, multiSelectTask.cursor, "Курсор должен быть на опции 'Выбрать все'")
	assert.Equal(t, "Выбрать все", multiSelectTask.selectAllText, "Текст по умолчанию должен быть 'Выбрать все'")
	
	// Проверяем, что View содержит опцию "Выбрать все"
	view := multiSelectTask.View(80)
	assert.Contains(t, view, "Выбрать все", "View должен содержать опцию 'Выбрать все'")
	
	// Выбираем "Выбрать все" с помощью пробела
	updatedTask1, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeySpace})
	multiSelectTaskAfterSelectAll, _ := updatedTask1.(*MultiSelectTask)
	
	// Проверяем, что все элементы выбраны
	assert.True(t, multiSelectTaskAfterSelectAll.isAllSelected(), "Все элементы должны быть выбраны")
	assert.Equal(t, len(options), multiSelectTaskAfterSelectAll.selected.Count(), "Количество выбранных элементов должно равняться количеству опций")
	
	// Повторно выбираем "Выбрать все" чтобы снять выбор со всех
	updatedTask2, _ := multiSelectTaskAfterSelectAll.Update(tea.KeyMsg{Type: tea.KeySpace})
	multiSelectTaskAfterUnselectAll, _ := updatedTask2.(*MultiSelectTask)
	
	// Проверяем, что все элементы не выбраны
	assert.False(t, multiSelectTaskAfterUnselectAll.isAllSelected(), "Все элементы должны быть не выбраны")
	assert.Equal(t, 0, multiSelectTaskAfterUnselectAll.selected.Count(), "Количество выбранных элементов должно быть равно 0")
}

// TestMultiSelectTaskWithCustomSelectAllText проверяет кастомный текст для "Выбрать все"
func TestMultiSelectTaskWithCustomSelectAllText(t *testing.T) {
	// Создаем задачу MultiSelectTask с кастомным текстом для "Выбрать все"
	title := "Выберите модули"
	options := []string{"Модуль A", "Модуль B", "Модуль C"}
	customText := "Выделить всё"
	multiSelectTask := NewMultiSelectTask(title, options).WithSelectAll(customText)
	
	// Проверяем, что кастомный текст установлен
	assert.Equal(t, customText, multiSelectTask.selectAllText, "Кастомный текст должен быть установлен")
	
	// Проверяем, что View содержит кастомный текст
	view := multiSelectTask.View(80)
	assert.Contains(t, view, customText, "View должен содержать кастомный текст")
	assert.NotContains(t, view, "Выбрать все", "View не должен содержать текст по умолчанию")
}

// TestMultiSelectTaskNavigationWithSelectAll проверяет навигацию с опцией "Выбрать все"
func TestMultiSelectTaskNavigationWithSelectAll(t *testing.T) {
	// Создаем задачу MultiSelectTask с опцией "Выбрать все"
	title := "Выберите элементы"
	options := []string{"Элемент 1", "Элемент 2", "Элемент 3"}
	multiSelectTask := NewMultiSelectTask(title, options).WithSelectAll()
	
	// Изначально курсор должен быть на "Выбрать все"
	assert.Equal(t, -1, multiSelectTask.cursor, "Курсор должен быть на опции 'Выбрать все'")
	
	// Нажимаем "down" - переходим к первому элементу списка
	updatedTask1, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeyDown})
	multiSelectTask1, _ := updatedTask1.(*MultiSelectTask)
	assert.Equal(t, 0, multiSelectTask1.cursor, "Курсор должен быть на первом элементе")
	
	// Нажимаем "up" - возвращаемся к "Выбрать все"
	updatedTask2, _ := multiSelectTask1.Update(tea.KeyMsg{Type: tea.KeyUp})
	multiSelectTask2, _ := updatedTask2.(*MultiSelectTask)
	assert.Equal(t, -1, multiSelectTask2.cursor, "Курсор должен вернуться на опцию 'Выбрать все'")
	
	// Нажимаем "up" еще раз - должны остаться на "Выбрать все"
	updatedTask3, _ := multiSelectTask2.Update(tea.KeyMsg{Type: tea.KeyUp})
	multiSelectTask3, _ := updatedTask3.(*MultiSelectTask)
	assert.Equal(t, -1, multiSelectTask3.cursor, "Курсор должен остаться на опции 'Выбрать все'")
}

// TestMultiSelectTaskToggleSelectAllLogic проверяет логику переключения "Выбрать все"
func TestMultiSelectTaskToggleSelectAllLogic(t *testing.T) {
	// Создаем задачу MultiSelectTask с опцией "Выбрать все"
	title := "Выберите пункты"
	options := []string{"Пункт 1", "Пункт 2", "Пункт 3", "Пункт 4"}
	multiSelectTask := NewMultiSelectTask(title, options).WithSelectAll()
	
	// Выбираем некоторые элементы вручную
	// Переходим к первому элементу
	updatedTask1, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeyDown})
	multiSelectTask, _ = updatedTask1.(*MultiSelectTask)
	
	// Выбираем первый элемент
	updatedTask2, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeySpace})
	multiSelectTask, _ = updatedTask2.(*MultiSelectTask)
	
	// Переходим ко второму элементу
	updatedTask3, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeyDown})
	multiSelectTask, _ = updatedTask3.(*MultiSelectTask)
	
	// Выбираем второй элемент
	updatedTask4, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeySpace})
	multiSelectTask, _ = updatedTask4.(*MultiSelectTask)
	
	// Проверяем, что выбрано 2 элемента из 4
	assert.Equal(t, 2, multiSelectTask.selected.Count(), "Должно быть выбрано 2 элемента")
	assert.False(t, multiSelectTask.isAllSelected(), "Не все элементы должны быть выбраны")
	
	// Возвращаемся к опции "Выбрать все"
	for multiSelectTask.cursor != -1 {
		updatedTask, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeyUp})
		multiSelectTask, _ = updatedTask.(*MultiSelectTask)
	}
	
	// Нажимаем "Выбрать все" - должны выбраться все элементы
	updatedTask5, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeySpace})
	multiSelectTask, _ = updatedTask5.(*MultiSelectTask)
	
	// Проверяем, что все элементы выбраны
	assert.Equal(t, len(options), multiSelectTask.selected.Count(), "Все элементы должны быть выбраны")
	assert.True(t, multiSelectTask.isAllSelected(), "Все элементы должны быть выбраны")
	
	// Нажимаем "Выбрать все" еще раз - должны сняться все выборы
	updatedTask6, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeySpace})
	multiSelectTask, _ = updatedTask6.(*MultiSelectTask)
	
	// Проверяем, что все элементы не выбраны
	assert.Equal(t, 0, multiSelectTask.selected.Count(), "Все элементы должны быть не выбраны")
	assert.False(t, multiSelectTask.isAllSelected(), "Все элементы должны быть не выбраны")
}

// TestMultiSelectTaskEmptySelectionHandling проверяет обработку Enter без выбранных элементов
func TestMultiSelectTaskEmptySelectionHandling(t *testing.T) {
	// Создаем задачу MultiSelectTask
	title := "Выберите элементы"
	options := []string{"Элемент 1", "Элемент 2", "Элемент 3"}
	multiSelectTask := NewMultiSelectTask(title, options)
	
	// Нажимаем Enter без выбора элементов
	updatedTask1, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeyEnter})
	multiSelectTaskAfterEnter, _ := updatedTask1.(*MultiSelectTask)
	
	// Проверяем, что задача НЕ завершена и НЕТ ошибки
	assert.False(t, multiSelectTaskAfterEnter.IsDone(), "Задача не должна быть завершена при пустом выборе")
	assert.False(t, multiSelectTaskAfterEnter.HasError(), "Не должно быть ошибки при пустом выборе")
	assert.True(t, multiSelectTaskAfterEnter.showHelpMessage, "Должно показываться сообщение-подсказка")
	assert.NotEmpty(t, multiSelectTaskAfterEnter.helpMessage, "Сообщение-подсказка не должно быть пустым")
	
	// Проверяем, что сообщение-подсказка отображается в View
	view := multiSelectTaskAfterEnter.View(80)
	assert.Contains(t, view, "Необходимо выбрать хотя бы один элемент", "View должен содержать сообщение-подсказку")
	
	// Нажимаем пробел (любую другую клавишу) - сообщение должно исчезнуть
	updatedTask2, _ := multiSelectTaskAfterEnter.Update(tea.KeyMsg{Type: tea.KeySpace})
	multiSelectTaskAfterSpace, _ := updatedTask2.(*MultiSelectTask)
	
	assert.False(t, multiSelectTaskAfterSpace.showHelpMessage, "Сообщение-подсказка должно исчезнуть")
	assert.Empty(t, multiSelectTaskAfterSpace.helpMessage, "Текст сообщения-подсказки должен быть очищен")
	
	// Теперь выбираем элемент и нажимаем Enter - задача должна завершиться успешно
	updatedTask3, _ := multiSelectTaskAfterSpace.Update(tea.KeyMsg{Type: tea.KeyEnter})
	multiSelectTaskFinal, _ := updatedTask3.(*MultiSelectTask)
	
	assert.True(t, multiSelectTaskFinal.IsDone(), "Задача должна завершиться после выбора элемента")
	assert.False(t, multiSelectTaskFinal.HasError(), "Не должно быть ошибки при успешном завершении")
}
