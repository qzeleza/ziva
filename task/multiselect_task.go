package task

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/common"
)

// MultiSelectTask представляет задачу выбора нескольких элементов из списка
type MultiSelectTask struct {
	BaseTask
	options     []string
	selected    []bool
	cursorIdx   int
	confirmed   bool
}

// NewMultiSelectTask создает новую задачу множественного выбора
func NewMultiSelectTask(prompt string, options []string) *MultiSelectTask {
	return &MultiSelectTask{
		BaseTask:  NewBaseTask(prompt),
		options:   options,
		selected:  make([]bool, len(options)),
		cursorIdx: 0,
		confirmed: false,
	}
}

// View отображает текущее состояние задачи
func (t *MultiSelectTask) View(width int) string {
	if t.IsDone() {
		return t.FinalView(width)
	}
	
	result := t.Title() + ":\n"
	for i, option := range t.options {
		cursor := "  "
		if i == t.cursorIdx {
			cursor = "> "
		}
		
		checkbox := "[ ]"
		if t.selected[i] {
			checkbox = "[x]"
		}
		
		result += cursor + checkbox + " " + option + "\n"
	}
	
	result += "\n(Space to select, Enter to confirm)"
	return result
}

// Update обновляет состояние задачи
func (t *MultiSelectTask) Update(msg tea.Msg) (common.Task, tea.Cmd) {
	if t.IsDone() {
		return t, nil
	}
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if t.cursorIdx > 0 {
				t.cursorIdx--
			}
		case "down", "j":
			if t.cursorIdx < len(t.options)-1 {
				t.cursorIdx++
			}
		case " ":
			t.selected[t.cursorIdx] = !t.selected[t.cursorIdx]
		case "enter":
			t.confirmed = true
			t.done = true
			
			selectedOptions := []string{}
			for i, selected := range t.selected {
				if selected {
					selectedOptions = append(selectedOptions, t.options[i])
				}
			}
			
			if len(selectedOptions) == 0 {
				t.finalValue = "Ничего не выбрано"
			} else if len(selectedOptions) == 1 {
				t.finalValue = selectedOptions[0]
			} else {
				t.finalValue = fmt.Sprintf("%d элементов выбрано", len(selectedOptions))
			}
		}
	}
	
	return t, nil
}

// Run запускает выполнение задачи
func (t *MultiSelectTask) Run() tea.Cmd {
	return nil
}

// GetSelectedOptions возвращает выбранные элементы
func (t *MultiSelectTask) GetSelectedOptions() []string {
	if !t.confirmed {
		return nil
	}
	
	selectedOptions := []string{}
	for i, selected := range t.selected {
		if selected {
			selectedOptions = append(selectedOptions, t.options[i])
		}
	}
	return selectedOptions
}

// GetSelectedIndices возвращает индексы выбранных элементов
func (t *MultiSelectTask) GetSelectedIndices() []int {
	if !t.confirmed {
		return nil
	}
	
	selectedIndices := []int{}
	for i, selected := range t.selected {
		if selected {
			selectedIndices = append(selectedIndices, i)
		}
	}
	return selectedIndices
}