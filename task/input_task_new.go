package task

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/common"
)

// InputTaskNew представляет задачу ввода текста
type InputTaskNew struct {
	BaseTask
	prompt      string
	value       string
	placeholder string
	validator   func(string) error
}

// NewInputTaskNew создает новую задачу ввода текста
func NewInputTaskNew(prompt, placeholder string, validator func(string) error) *InputTaskNew {
	return &InputTaskNew{
		BaseTask:    NewBaseTask(prompt),
		prompt:      prompt,
		placeholder: placeholder,
		validator:   validator,
	}
}

// View отображает текущее состояние задачи
func (t *InputTaskNew) View(width int) string {
	if t.IsDone() {
		return t.FinalView(width)
	}
	return t.prompt + ": " + t.value
}

// Update обновляет состояние задачи
func (t *InputTaskNew) Update(msg tea.Msg) (common.Task, tea.Cmd) {
	if t.IsDone() {
		return t, nil
	}
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if t.validator != nil {
				if err := t.validator(t.value); err != nil {
					t.SetError(err)
					return t, nil
				}
			}
			t.done = true
			t.finalValue = t.value
		case "backspace":
			if len(t.value) > 0 {
				t.value = t.value[:len(t.value)-1]
			}
		default:
			t.value += msg.String()
		}
	}
	
	return t, nil
}

// Run запускает выполнение задачи
func (t *InputTaskNew) Run() tea.Cmd {
	return nil
}

// GetValue возвращает введенное значение
func (t *InputTaskNew) GetValue() string {
	return t.value
}