package task

import (
	"fmt"
	"strings"
	"time"

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

// Update обрабатывает нажатия клавиш для навигации и выбора.
func (t *SingleSelectTask) Update(msg tea.Msg) (Task, tea.Cmd) {
	if t.done {
		return t, nil
	}

	switch msg := msg.(type) {
	// Обработка сообщения о тайм-ауте
	case TimeoutMsg:
		// Когда истекает тайм-аут, применяем значение по умолчанию (если есть)
		t.applyDefaultValue()
		return t, nil
	// Обработка периодического обновления для счетчика времени
	case TickMsg:
		// Если таймер активен, продолжаем обновления
		if t.timeoutEnabled && t.timeoutManager != nil && t.timeoutManager.IsActive() {
			return t, t.timeoutManager.StartTicker()
		}
		return t, nil
	case tea.KeyMsg:
		// При нажатии клавиш НЕ сбрасываем таймер - пусть продолжает работать
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
			t.err = cancelErr
			t.icon = ui.IconCancelled
			t.finalValue = ui.GetErrorMessageStyle().Render(cancelErr.Error())
			t.SetStopOnError(true)
			return t, nil

		case "enter":
			t.done = true
			t.icon = ui.IconDone
			t.finalValue = t.choices[t.cursor]
			return t, nil
		case " ":
			// Если таймер активен, останавливаем его
			if t.timeoutEnabled && t.timeoutManager != nil && t.timeoutManager.IsActive() {
				t.timeoutManager.StopTimeout()
				t.showTimeout = false
			}
			// В любом случае выбираем текущий элемент
			t.done = true
			t.icon = ui.IconDone
			t.finalValue = t.choices[t.cursor]
			return t, nil
		}
		// После обработки клавиш возвращаем команду для продолжения тикера
		if t.timeoutEnabled && t.timeoutManager != nil && t.timeoutManager.IsActive() {
			return t, t.timeoutManager.StartTicker()
		}
	}
	return t, nil
}

// Run запускает задачу выбора
func (t *SingleSelectTask) Run() tea.Cmd {
	// Запускаем таймер и тикер, если они включены
	if t.timeoutEnabled && t.timeoutManager != nil {
		return t.timeoutManager.StartTickerAndTimeout()
	}
	return nil
}

// applyDefaultValue применяет значение по умолчанию при истечении таймера
func (t *SingleSelectTask) applyDefaultValue() {
	// Если есть значение по умолчанию и это число (индекс)
	if t.defaultValue != nil {
		switch val := t.defaultValue.(type) {
		case int:
			// Если индекс в допустимом диапазоне
			if val >= 0 && val < len(t.choices) {
				t.cursor = val
				// Устанавливаем задачу как завершенную
				t.done = true
				t.icon = ui.IconDone
				t.finalValue = t.choices[t.cursor]
			}
		case string:
			// Ищем строку в списке вариантов
			for i, choice := range t.choices {
				if choice == val {
					t.cursor = i
					// Устанавливаем задачу как завершенную
					t.done = true
					t.icon = ui.IconDone
					t.finalValue = t.choices[t.cursor]
					break
				}
			}
		}
	}
}

// View отрисовывает список вариантов выбора для пользователя с выделением активного элемента.
//
// @param width Ширина макета для отображения
// @return Строка с отформатированным представлением задачи
func (t *SingleSelectTask) View(width int) string {
	// Если задача завершена, возвращаем FinalView
	if t.done {
		return t.FinalView(width)
	}
	var sb strings.Builder

	// Добавляем заголовок задачи с префиксом для активной задачи
	titlePrefix := ui.GetCurrentTaskPrefix()

	// Формируем заголовок с префиксом
	title := ui.ActiveTitleStyle.Render(t.title)
	titleWithPrefix := fmt.Sprintf("%s%s", titlePrefix, title)

	// Получаем отформатированный таймер (если он активен)
	timerStr := t.RenderTimer()

	// Если есть таймер, выравниваем заголовок и таймер по правому краю
	if timerStr != "" {
		titleLine := ui.AlignTextToRight(titleWithPrefix, timerStr, width)
		sb.WriteString(titleLine + "\n")
	} else {
		sb.WriteString(titleWithPrefix + "\n")
	}

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

// WithTimeout устанавливает тайм-аут для задачи выбора
// @param duration Длительность тайм-аута
// @param defaultValue Значение по умолчанию (индекс или строка)
// @return Указатель на задачу для цепочки вызовов
func (t *SingleSelectTask) WithTimeout(duration time.Duration, defaultValue interface{}) *SingleSelectTask {
	t.BaseTask.WithTimeout(duration, defaultValue)
	return t
}
