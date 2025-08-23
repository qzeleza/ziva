// task/yesno_task_test.go

package task

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// TestYesNoTaskCreation проверяет корректность создания YesNoTask
func TestYesNoTaskCreation(t *testing.T) {
	// Создаем задачу YesNoTask
	title := "Подтверждение"
	question := "Вы согласны с условиями?"
	yesNoTask := NewYesNoTask(title, question)

	// Проверяем, что задача создана корректно
	assert.NotNil(t, yesNoTask, "Задача не должна быть nil")
	assert.Equal(t, title, yesNoTask.Title(), "Заголовок задачи должен соответствовать переданному значению")
	assert.False(t, yesNoTask.IsDone(), "Новая задача не должна быть отмечена как завершенная")

	// Проверяем значения по умолчанию
	assert.Equal(t, "Да", yesNoTask.yesLabel, "Метка 'Да' должна быть установлена по умолчанию")
	assert.Equal(t, "Нет", yesNoTask.noLabel, "Метка 'Нет' должна быть установлена по умолчанию")
	assert.Equal(t, YesOption, yesNoTask.selectedOption, "По умолчанию должна быть выбрана опция 'Да'")
}

// TestYesNoTaskWithCustomLabels проверяет метод WithCustomLabels
func TestYesNoTaskWithCustomLabels(t *testing.T) {
	// Создаем задачу YesNoTask
	yesNoTask := NewYesNoTask("Подтверждение", "Вы согласны с условиями?")

	// Устанавливаем пользовательские метки (только первые две)
	customYesLabel := "Согласен"
	customNoLabel := "Не согласен"
	yesNoTask = yesNoTask.WithCustomLabels(customYesLabel, customNoLabel)

	// Проверяем, что метки установлены корректно
	assert.Equal(t, customYesLabel, yesNoTask.yesLabel, "Метка 'Да' должна соответствовать пользовательской")
	assert.Equal(t, customNoLabel, yesNoTask.noLabel, "Метка 'Нет' должна соответствовать пользовательской")
}

// TestYesNoTaskWithCustomLabelsAll проверяет метод WithCustomLabelsAll
func TestYesNoTaskWithCustomLabelsAll(t *testing.T) {
	// Создаем задачу YesNoTask
	yesNoTask := NewYesNoTask("Лицензия", "Принять лицензионное соглашение?")

	// Устанавливаем все пользовательские метки
	customYesLabel := "Принять"
	customNoLabel := "Отклонить"
	yesNoTask = yesNoTask.WithCustomLabelsAll(customYesLabel, customNoLabel)

	// Проверяем, что все метки установлены корректно
	assert.Equal(t, customYesLabel, yesNoTask.yesLabel, "Метка 'Да' должна соответствовать пользовательской")
	assert.Equal(t, customNoLabel, yesNoTask.noLabel, "Метка 'Нет' должна соответствовать пользовательской")
}

// TestYesNoTaskSelectionAndNavigation проверяет навигацию и выбор
func TestYesNoTaskSelectionAndNavigation(t *testing.T) {
	// Создаем задачу YesNoTask
	yesNoTask := NewYesNoTask("Подтверждение", "Вы согласны с условиями?")

	// Проверяем начальное состояние (должно быть на первой опции - "Да")
	assert.Equal(t, 0, yesNoTask.GetSelectedIndex(), "Изначально должна быть выбрана первая опция")

	// Перемещаемся вниз к "Нет"
	updatedTask, _ := yesNoTask.Update(tea.KeyMsg{Type: tea.KeyDown})
	yesNoTask = updatedTask.(*YesNoTask)
	assert.Equal(t, 1, yesNoTask.GetSelectedIndex(), "После Down должна быть выбрана вторая опция")

	// Перемещаемся вверх обратно к "Нет"
	updatedTask, _ = yesNoTask.Update(tea.KeyMsg{Type: tea.KeyUp})
	yesNoTask = updatedTask.(*YesNoTask)
	assert.Equal(t, 1, yesNoTask.GetSelectedIndex(), "После Up должна быть выбрана вторая опция")
}

// TestYesNoTaskOptionSelection проверяет выбор различных опций
func TestYesNoTaskOptionSelection(t *testing.T) {
	// Тест выбора "Да"
	yesNoTask := NewYesNoTask("Подтверждение", "Вы согласны с условиями?")

	// Выбираем первую опцию (Да)
	updatedTask, _ := yesNoTask.Update(tea.KeyMsg{Type: tea.KeyEnter})
	yesNoTaskDone := updatedTask.(*YesNoTask)

	assert.True(t, yesNoTaskDone.IsDone(), "Задача должна быть завершена")
	assert.Equal(t, YesOption, yesNoTaskDone.GetSelectedOption(), "Должна быть выбрана опция 'Да'")
	assert.True(t, yesNoTaskDone.IsYes(), "IsYes() должен возвращать true")
	assert.False(t, yesNoTaskDone.IsNo(), "IsNo() должен возвращать false")
	assert.False(t, yesNoTaskDone.GetValue(), "GetValue() должен возвращать true для 'Да'")

	// Тест выбора "Нет"
	yesNoTask = NewYesNoTask("Подтверждение", "Вы согласны с условиями?")
	updatedTask, _ = yesNoTask.Update(tea.KeyMsg{Type: tea.KeyDown}) // переходим к "Нет"
	yesNoTask = updatedTask.(*YesNoTask)
	updatedTask, _ = yesNoTask.Update(tea.KeyMsg{Type: tea.KeyEnter}) // выбираем "Нет"
	yesNoTaskDone = updatedTask.(*YesNoTask)

	assert.True(t, yesNoTaskDone.IsDone(), "Задача должна быть завершена")
	assert.Equal(t, NoOption, yesNoTaskDone.GetSelectedOption(), "Должна быть выбрана опция 'Нет'")
	assert.False(t, yesNoTaskDone.IsYes(), "IsYes() должен возвращать false")
	assert.True(t, yesNoTaskDone.IsNo(), "IsNo() должен возвращать true")
	assert.False(t, yesNoTaskDone.GetValue(), "GetValue() должен возвращать false для 'Нет'")

	// Тест выбора "Выйти"
	yesNoTask = NewYesNoTask("Подтверждение", "Вы согласны с условиями?")
	updatedTask, _ = yesNoTask.Update(tea.KeyMsg{Type: tea.KeyDown}) // переходим к "Нет"
	yesNoTask = updatedTask.(*YesNoTask)
	updatedTask, _ = yesNoTask.Update(tea.KeyMsg{Type: tea.KeyDown}) // переходим к "Выйти"
	yesNoTask = updatedTask.(*YesNoTask)
	updatedTask, _ = yesNoTask.Update(tea.KeyMsg{Type: tea.KeyEnter}) // выбираем "Выйти"
	yesNoTaskDone = updatedTask.(*YesNoTask)

	assert.True(t, yesNoTaskDone.IsDone(), "Задача должна быть завершена")
	assert.Equal(t, NoOption, yesNoTaskDone.GetSelectedOption(), "Должна быть выбрана опция 'Нет'")
	assert.False(t, yesNoTaskDone.IsYes(), "IsYes() должен возвращать false")
	assert.False(t, yesNoTaskDone.IsNo(), "IsNo() должен возвращать false")

	// GetValue() должен паниковать для опции "Выйти"
	assert.Panics(t, func() {
		yesNoTaskDone.GetValue()
	}, "GetValue() должен паниковать для опции 'Выйти'")
}

// TestYesNoTaskView проверяет отображение задачи
func TestYesNoTaskView(t *testing.T) {
	// Создаем задачу YesNoTask
	title := "Подтверждение"
	question := "Вы согласны с условиями?"
	yesNoTask := NewYesNoTask(title, question)

	// Проверяем, что View содержит опции (делегируется SingleSelectTask)
	view := yesNoTask.View(80)
	assert.Contains(t, view, "Да", "View должен содержать опцию 'Да'")
	assert.Contains(t, view, "Нет", "View должен содержать опцию 'Нет'")
	assert.Contains(t, view, "Выйти", "View должен содержать опцию 'Выйти'")

	// Проверяем, что после завершения задачи View возвращает FinalView
	updatedTask, _ := yesNoTask.Update(tea.KeyMsg{Type: tea.KeyEnter})
	yesNoTaskDone := updatedTask.(*YesNoTask)
	finalView := yesNoTaskDone.View(80)
	assert.NotEmpty(t, finalView, "View должен возвращать FinalView для завершенной задачи")
	assert.Contains(t, finalView, "Да", "FinalView должен содержать выбранную опцию")
}

// TestYesNoTaskLegacyCompatibility проверяет совместимость со старым API
func TestYesNoTaskLegacyCompatibility(t *testing.T) {
	// Создаем задачу и эмулируем старое поведение
	yesNoTask := NewYesNoTask("Старый стиль", "Продолжить?")

	// Выбираем "Да" и завершаем
	updatedTask, _ := yesNoTask.Update(tea.KeyMsg{Type: tea.KeyEnter})
	yesNoTaskDone := updatedTask.(*YesNoTask)

	// Проверяем, что старый метод GetValue() все еще работает
	assert.True(t, yesNoTaskDone.GetValue(), "GetValue() должен работать для обратной совместимости")

	// Проверяем новые методы
	assert.True(t, yesNoTaskDone.IsYes(), "Новый API должен работать")
	assert.Equal(t, YesOption, yesNoTaskDone.GetSelectedOption(), "Новый API должен работать")
}
