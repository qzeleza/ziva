package task

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/common"
)

// SingleSelectTask представляет задачу выбора одного элемента из списка
type SingleSelectTask struct {
	BaseTask
	options      []string
	selectedIdx  int
	confirmed    bool
}

// NewSingleSelectTask создает новую задачу выбора одного элемента
func NewSingleSelectTask(prompt string, options []string) *SingleSelectTask {
	return &SingleSelectTask{
		BaseTask:    NewBaseTask(prompt),
		options:     options,
		selectedIdx: 0,
		confirmed:   false,
	}
}

// View отображает текущее состояние задачи
func (t *SingleSelectTask) View(width int) string {
	if t.IsDone() {
		return t.FinalView(width)
	}
	
	result := t.Title() + ":\n"
	for i, option := range t.options {
		if i == t.selectedIdx {
			result += "  > " + option + "\n"
		} else {
			result += "    " + option + "\n"
		}
	}
	return result
}

// Update обновляет состояние задачи
func (t *SingleSelectTask) Update(msg tea.Msg) (common.Task, tea.Cmd) {
	if t.IsDone() {
		return t, nil
	}
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if t.selectedIdx > 0 {
				t.selectedIdx--
			}
		case "down", "j":
			if t.selectedIdx < len(t.options)-1 {
				t.selectedIdx++
			}
		case "enter", " ":
			t.confirmed = true
			t.done = true
			t.finalValue = t.options[t.selectedIdx]
		}
	}
	
	return t, nil
}

// Run запускает выполнение задачи
func (t *SingleSelectTask) Run() tea.Cmd {
	return nil
}

// GetSelectedOption возвращает выбранный элемент
func (t *SingleSelectTask) GetSelectedOption() string {
	if t.confirmed && t.selectedIdx < len(t.options) {
		return t.options[t.selectedIdx]
	}
	return ""
}

// GetSelectedIndex возвращает индекс выбранного элемента
func (t *SingleSelectTask) GetSelectedIndex() int {
	if t.confirmed {
		return t.selectedIdx
	}
	return -1
}