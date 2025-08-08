package task

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/common"
	"github.com/qzeleza/termos/ui"
)

// YesNoTask представляет задачу с выбором "да/нет"
type YesNoTask struct {
	BaseTask
	question    string
	description string
	choice      bool
	confirmed   bool
}

// NewYesNoTask создает новую задачу с выбором "да/нет"
func NewYesNoTask(question, description string) *YesNoTask {
	return &YesNoTask{
		BaseTask:    NewBaseTask(question),
		question:    question,
		description: description,
		choice:      false, // По умолчанию "нет"
		confirmed:   false,
	}
}

// View отображает текущее состояние задачи
func (t *YesNoTask) View(width int) string {
	if t.IsDone() {
		return t.FinalView(width)
	}

	// Используем ширину из common пакета если передана недостаточная ширина
	if width < common.DefaultWidth {
		width = common.DefaultWidth
	}

	prefix := ui.GetActiveTaskPrefix()
	title := ui.ActiveTitleStyle.Render(t.question)
	
	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s%s\n", prefix, title))
	
	if t.description != "" {
		desc := ui.SubtleStyle.Render(t.description)
		result.WriteString(fmt.Sprintf("   %s\n", desc))
	}
	
	// Показываем опции
	yesStyle := ui.SelectionStyle
	noStyle := ui.SelectionNoStyle
	cursor := ui.IconCursor
	
	if t.choice {
		result.WriteString(fmt.Sprintf("   %s %s\n", cursor, yesStyle.Render("Да")))
		result.WriteString(fmt.Sprintf("     %s\n", noStyle.Render("Нет")))
	} else {
		result.WriteString(fmt.Sprintf("     %s\n", yesStyle.Render("Да")))
		result.WriteString(fmt.Sprintf("   %s %s\n", cursor, noStyle.Render("Нет")))
	}
	
	return result.String()
}

// Update обновляет состояние задачи
func (t *YesNoTask) Update(msg tea.Msg) (common.Task, tea.Cmd) {
	if t.IsDone() {
		return t, nil
	}
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k", "left", "h":
			t.choice = true
		case "down", "j", "right", "l":
			t.choice = false
		case "enter", " ":
			t.confirmed = true
			t.done = true
			if t.choice {
				t.finalValue = "Да"
			} else {
				t.finalValue = "Нет"
			}
		case "q", "esc", "ctrl+c":
			t.done = true
			t.finalValue = "Отменено"
			t.SetError(fmt.Errorf("отменено пользователем"))
		}
	}
	
	return t, nil
}

// Run запускает выполнение задачи
func (t *YesNoTask) Run() tea.Cmd {
	return nil // YesNoTask не требует асинхронного выполнения
}

// GetChoice возвращает выбор пользователя
func (t *YesNoTask) GetChoice() bool {
	return t.choice && t.confirmed
}