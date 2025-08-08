package task

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

// TestMultiSelectTaskPrefixRendering проверяет правильность формирования префиксов
func TestMultiSelectTaskPrefixRendering(t *testing.T) {
	// Создаем задачу MultiSelectTask с опцией "Выбрать все"
	title := "Тест префиксов"
	options := []string{"Элемент 1", "Элемент 2", "Элемент 3"}
	multiSelectTask := NewMultiSelectTask(title, options).WithSelectAll("Выбрать все")

	// Тест 1: Курсор на "Выбрать все" (позиция -1)
	// Ожидаем: "Выбрать все" с префиксом "active", все элементы с префиксом "below" (только отступы)
	view1 := multiSelectTask.View(80)
	
	// Проверяем, что "Выбрать все" отображается как активная
	assert.Contains(t, view1, "└─>", "Должен быть префикс активного элемента для 'Выбрать все'")
	
	// Проверяем, что элементы списка имеют префикс "below" (только отступы без вертикальной линии)
	lines := strings.Split(view1, "\n")
	elementLines := make([]string, 0)
	for _, line := range lines {
		if strings.Contains(line, "Элемент") {
			elementLines = append(elementLines, line)
		}
	}
	
	// Все элементы должны начинаться только с отступов (без вертикальной линии)
	for i, line := range elementLines {
		assert.True(t, strings.HasPrefix(line, "      "), 
			"Элемент %d должен иметь префикс 'below' (только отступы): %s", i+1, line)
		assert.False(t, strings.Contains(line[:6], "│"), 
			"Элемент %d не должен содержать вертикальную линию в префиксе: %s", i+1, line)
	}

	// Тест 2: Переходим к первому элементу списка (позиция 0)
	updatedTask1, _ := multiSelectTask.Update(tea.KeyMsg{Type: tea.KeyDown})
	multiSelectTask1, _ := updatedTask1.(*MultiSelectTask)
	view2 := multiSelectTask1.View(80)
	
	// Проверяем, что "Выбрать все" теперь имеет префикс "above" (с вертикальной линией)
	lines2 := strings.Split(view2, "\n")
	selectAllLine := ""
	for _, line := range lines2 {
		if strings.Contains(line, "Выбрать все") {
			selectAllLine = line
			break
		}
	}
	
	assert.True(t, strings.Contains(selectAllLine[:6], "│"), 
		"'Выбрать все' должно иметь префикс 'above' с вертикальной линией: %s", selectAllLine)

	// Проверяем, что первый элемент активен (с префиксом "active")
	assert.Contains(t, view2, "└─>", "Первый элемент должен иметь префикс активного элемента")
	
	// Тест 3: Переходим ко второму элементу (позиция 1)
	updatedTask2, _ := multiSelectTask1.Update(tea.KeyMsg{Type: tea.KeyDown})
	multiSelectTask2, _ := updatedTask2.(*MultiSelectTask)
	view3 := multiSelectTask2.View(80)
	
	lines3 := strings.Split(view3, "\n")
	elementLines3 := make([]string, 0)
	for _, line := range lines3 {
		if strings.Contains(line, "Элемент") {
			elementLines3 = append(elementLines3, line)
		}
	}
	
	// Первый элемент должен иметь префикс "above" (с вертикальной линией)
	assert.True(t, strings.Contains(elementLines3[0][:6], "│"), 
		"Элемент 1 должен иметь префикс 'above': %s", elementLines3[0])
	
	// Второй элемент должен быть активным (с префиксом "active")
	assert.Contains(t, elementLines3[1], "└─>", "Элемент 2 должен иметь префикс активного элемента")
	
	// Третий элемент должен иметь префикс "below" (только отступы)
	assert.True(t, strings.HasPrefix(elementLines3[2], "      "), 
		"Элемент 3 должен иметь префикс 'below': %s", elementLines3[2])
	assert.False(t, strings.Contains(elementLines3[2][:6], "│"), 
		"Элемент 3 не должен содержать вертикальную линию: %s", elementLines3[2])
}

// TestMultiSelectTaskPrefixWithoutSelectAll проверяет префиксы без опции "Выбрать все"
func TestMultiSelectTaskPrefixWithoutSelectAll(t *testing.T) {
	// Создаем обычную задачу без "Выбрать все"
	title := "Тест без 'Выбрать все'"
	options := []string{"Пункт 1", "Пункт 2", "Пункт 3"}
	multiSelectTask := NewMultiSelectTask(title, options)

	// Курсор должен быть на первом элементе (позиция 0)
	view := multiSelectTask.View(80)
	
	// Первый элемент должен быть активным
	assert.Contains(t, view, "└─>", "Первый элемент должен иметь префикс активного элемента")
	
	lines := strings.Split(view, "\n")
	elementLines := make([]string, 0)
	for _, line := range lines {
		if strings.Contains(line, "Пункт") {
			elementLines = append(elementLines, line)
		}
	}
	
	// Второй и третий элементы должны иметь префикс "below"
	for i := 1; i < len(elementLines); i++ {
		assert.True(t, strings.HasPrefix(elementLines[i], "      "), 
			"Пункт %d должен иметь префикс 'below': %s", i+1, elementLines[i])
	}
}