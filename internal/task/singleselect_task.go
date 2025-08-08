package task

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qzeleza/termos/internal/performance"
	"github.com/qzeleza/termos/internal/ui"
)

// SingleSelectTask - задача для выбора одного варианта из списка.
type SingleSelectTask struct {
	BaseTask
	choices     []string       // Список вариантов выбора
	cursor      int            // Текущая позиция курсора
	activeStyle lipgloss.Style // Стиль для активного элемента
}

// NewSingleSelectTask создает новую задачу выбора одного варианта из списка.
//
// @param title Заголовок задачи
// @param choices Список вариантов выбора
// @return Указатель на новую задачу выбора
func NewSingleSelectTask(title string, choices []string) *SingleSelectTask {
	return &SingleSelectTask{
		BaseTask:    NewBaseTask(title),
		choices:     choices,
		activeStyle: ui.ActiveStyle,
	}
}

// Update handles key presses for navigation and selection.
func (t *SingleSelectTask) Update(msg tea.Msg) (Task, tea.Cmd) {
	if t.done {
		return t, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if t.cursor > 0 {
				t.cursor--
			}
		case "down", "j":
			if t.cursor < len(t.choices)-1 {
				t.cursor++
			}
		case "q", "Q":
			// Отмена пользователем
			cancelErr := fmt.Errorf("отменено пользователем")
			t.done = true
			t.BaseTask.err = cancelErr
			t.icon = ui.IconCancelled
			t.finalValue = ui.GetErrorMessageStyle().Render(cancelErr.Error())
			t.SetStopOnError(true)
			return t, nil

		case "enter", " ":
			t.done = true
			t.icon = ui.IconDone
			t.finalValue = t.choices[t.cursor]
			return t, nil
		}
	}
	return t, nil
}

// View отрисовывает список вариантов выбора для пользователя с выделением активного элемента.
//
// @param width Ширина макета для отображения
// @return Строка с отформатированным представлением задачи
func (t *SingleSelectTask) View(width int) string {
	if t.done {
		return ""
	}
	var sb strings.Builder

	// Добавляем заголовок задачи с новым префиксом для текущей задачи
	titlePrefix := ui.GetCurrentTaskPrefix()
	sb.WriteString(fmt.Sprintf("%s%s\n", titlePrefix, ui.ActiveTitleStyle.Render(t.title)))

	for i, choice := range t.choices {
		checked := ui.IconRadioOff
		var itemPrefix string

		// Определяем тип элемента для получения правильного префикса
		if t.cursor == i {
			// Активный элемент
			itemPrefix = ui.GetSelectItemPrefix("active")
			checked = ui.IconRadioOn

			// Применяем стиль активного элемента
			choice = t.activeStyle.Render(choice)
			checked = t.activeStyle.Render(checked)
		} else if i < t.cursor {
			// Элемент выше активного
			itemPrefix = ui.GetSelectItemPrefix("above")
		} else {
			// Элемент ниже активного
			itemPrefix = ui.GetSelectItemPrefix("below")
		}

		// Формируем строку для отображения варианта выбора с новым префиксом
		if t.cursor == i {
			// Для активного элемента используем стилизованные скобки
			bracketsOpen := t.activeStyle.Render("(")
			bracketsClose := t.activeStyle.Render(")")
			sb.WriteString(fmt.Sprintf("%s%s%s%s %s\n", itemPrefix, bracketsOpen, checked, bracketsClose, choice))
		} else {
			// Для неактивных элементов
			sb.WriteString(fmt.Sprintf("%s(%s) %s\n", itemPrefix, checked, choice))
		}
	}

	// Добавляем подсказку о навигации с новым отступом
	helpIndent := performance.RepeatEfficient(" ", ui.MainLeftIndent)
	sb.WriteString("\n" + ui.DrawLine(width) +
		ui.SubtleStyle.Render(fmt.Sprintf("%s[↑/↓ для навигации, Enter для выбора]", helpIndent)))

	return sb.String()
}

func (t *SingleSelectTask) FinalView(width int) string {
	// Получаем базовое финальное представление
	result := t.BaseTask.FinalView(width)

	// Если задача завершилась успешно и есть дополнительные строки для вывода
	if t.icon == ui.IconDone && len(t.choices) > 0 {
		result += "\n" + ui.DrawSummaryLine(t.choices[t.cursor]) +
			performance.RepeatEfficient(" ", ui.MainLeftIndent) + ui.VerticalLineSymbol
	}

	return result
}

// GetSelected возвращает выбранное значение из списка
//
// @return string Выбранный пользователем вариант
func (t *SingleSelectTask) GetSelected() string {
	// Если задача не завершена, возвращаем текущий выбор
	if !t.done && t.cursor >= 0 && t.cursor < len(t.choices) {
		return t.choices[t.cursor]
	}
	// Если задача завершена, возвращаем сохраненное финальное значение
	if t.finalValue != "" {
		return t.finalValue
	}
	// Если нет финального значения, но есть выбор, возвращаем его
	if t.cursor >= 0 && t.cursor < len(t.choices) {
		return t.choices[t.cursor]
	}
	// Если ничего не выбрано, возвращаем пустую строку
	return ""
}

// GetSelectedIndex возвращает индекс выбранного элемента
func (t *SingleSelectTask) GetSelectedIndex() int {
	return t.cursor
}
