package task

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qzeleza/termos/internal/defauilt"
	"github.com/qzeleza/termos/internal/performance"
	"github.com/qzeleza/termos/internal/ui"
)

// SingleSelectTask - задача для выбора одного варианта из списка.
type SingleSelectTask struct {
	BaseTask
	choices     []string         // Список вариантов выбора
	disabled    map[int]struct{} // Набор отключённых пунктов меню
	cursor      int              // Текущая позиция курсора
	activeStyle lipgloss.Style   // Стиль для активного элемента
	// Viewport (окно просмотра) для ограничения количества отображаемых элементов
	viewportSize  int // Размер viewport (количество видимых элементов), 0 = показать все
	viewportStart int // Начальная позиция viewport в списке элементов
}

// NewSingleSelectTask создает новую задачу выбора одного варианта из списка.
//
// @param title Заголовок задачи
// @param choices Список вариантов выбора
// @return Указатель на новую задачу выбора
func NewSingleSelectTask(title string, choices []string) *SingleSelectTask {
	task := &SingleSelectTask{
		BaseTask:    NewBaseTask(title),
		choices:     choices,
		disabled:    make(map[int]struct{}),
		activeStyle: ui.ActiveStyle,
		// Viewport по умолчанию отключен (показываем все элементы)
		viewportSize:  0,
		viewportStart: 0,
	}

	task.ensureCursorSelectable()
	return task
}

// isDisabled проверяет, помечен ли элемент как недоступный
func (t *SingleSelectTask) isDisabled(index int) bool {
	if index < 0 || index >= len(t.choices) {
		return true
	}
	_, exists := t.disabled[index]
	return exists
}

// ensureCursorSelectable подбирает ближайший доступный элемент для курсора.
// Возвращает true, если удалось найти активный элемент.
func (t *SingleSelectTask) ensureCursorSelectable() bool {
	if len(t.choices) == 0 {
		t.cursor = -1
		return false
	}

	if t.cursor >= 0 && t.cursor < len(t.choices) && !t.isDisabled(t.cursor) {
		return true
	}

	start := t.cursor
	if start < 0 {
		start = 0
	}

	if idx, ok := t.findEnabledForward(start); ok {
		t.cursor = idx
		return true
	}

	if idx, ok := t.findEnabledBackward(start - 1); ok {
		t.cursor = idx
		return true
	}

	t.cursor = -1
	return false
}

// findEnabledForward возвращает индекс первого доступного элемента, начиная с from (включительно)
func (t *SingleSelectTask) findEnabledForward(from int) (int, bool) {
	if from < 0 {
		from = 0
	}
	for i := from; i < len(t.choices); i++ {
		if !t.isDisabled(i) {
			return i, true
		}
	}
	return -1, false
}

// findEnabledBackward возвращает индекс ближайшего доступного элемента, двигаясь в обратном порядке
func (t *SingleSelectTask) findEnabledBackward(from int) (int, bool) {
	if from >= len(t.choices) {
		from = len(t.choices) - 1
	}
	for i := from; i >= 0; i-- {
		if !t.isDisabled(i) {
			return i, true
		}
	}
	return -1, false
}

// moveCursorForward перемещает курсор на следующий доступный элемент
func (t *SingleSelectTask) moveCursorForward() bool {
	start := t.cursor + 1
	if t.cursor < 0 {
		start = 0
	}
	if idx, ok := t.findEnabledForward(start); ok {
		t.cursor = idx
		return true
	}
	return false
}

// moveCursorBackward перемещает курсор на предыдущий доступный элемент
func (t *SingleSelectTask) moveCursorBackward() bool {
	start := t.cursor - 1
	if idx, ok := t.findEnabledBackward(start); ok {
		t.cursor = idx
		return true
	}
	return false
}

// WithItemsDisabled помечает элементы меню как недоступные для выбора.
// Поддерживаются типы: int, []int, string, []string. Nil очищает список отключённых элементов.
func (t *SingleSelectTask) WithItemsDisabled(disabled interface{}) *SingleSelectTask {
	for idx := range t.disabled {
		delete(t.disabled, idx)
	}

	indices := t.resolveDisabledIndices(disabled)
	for _, idx := range indices {
		if idx >= 0 && idx < len(t.choices) {
			t.disabled[idx] = struct{}{}
		}
	}

	t.ensureCursorSelectable()
	t.updateViewport()
	return t
}

// resolveDisabledIndices конвертирует произвольный ввод в индексы элементов списка
func (t *SingleSelectTask) resolveDisabledIndices(input interface{}) []int {
	var result []int
	if input == nil {
		return result
	}

	addIndex := func(idx int) {
		if idx < 0 || idx >= len(t.choices) {
			return
		}
		for _, existing := range result {
			if existing == idx {
				return
			}
		}
		result = append(result, idx)
	}

	switch v := input.(type) {
	case int:
		addIndex(v)
	case []int:
		for _, idx := range v {
			addIndex(idx)
		}
	case string:
		if idx := t.choiceIndex(v); idx != -1 {
			addIndex(idx)
		}
	case []string:
		for _, val := range v {
			if idx := t.choiceIndex(val); idx != -1 {
				addIndex(idx)
			}
		}
	case []interface{}:
		for _, item := range v {
			for _, idx := range t.resolveDisabledIndices(item) {
				addIndex(idx)
			}
		}
	}

	return result
}

// choiceIndex возвращает индекс элемента по значению или -1, если элемент не найден
func (t *SingleSelectTask) choiceIndex(value string) int {
	for i, choice := range t.choices {
		if choice == value {
			return i
		}
	}
	return -1
}

// WithDefaultItem устанавливает элемент, который будет подсвечен курсором при открытии списка.
// Поддерживает выбор по индексу (int) или строковому значению (string). Некорректные значения игнорируются.
func (t *SingleSelectTask) WithDefaultItem(selection interface{}) *SingleSelectTask {
	if selection == nil || len(t.choices) == 0 {
		return t
	}

	setCursor := func(index int) {
		if index < 0 || index >= len(t.choices) {
			return
		}
		t.cursor = index
	}

	switch v := selection.(type) {
	case int:
		setCursor(v)
	case string:
		for i, choice := range t.choices {
			if choice == v {
				setCursor(i)
				break
			}
		}
	}

	t.ensureCursorSelectable()
	// После обновления курсора синхронизируем viewport
	t.updateViewport()
	return t
}

// WithViewport устанавливает размер viewport (окна просмотра) для ограничения количества отображаемых элементов.
// Это полезно для длинных списков, когда нужно показывать только часть элементов.
//
// @param size Количество элементов для отображения одновременно (0 = показать все)
// @return Указатель на задачу для цепочки вызовов
func (t *SingleSelectTask) WithViewport(size int) *SingleSelectTask {
	if size < 0 {
		size = 0
	}
	t.viewportSize = size
	return t
}

// updateViewport обновляет позицию viewport на основе текущего положения курсора
func (t *SingleSelectTask) updateViewport() {
	// Если viewport отключен, ничего не делаем
	if t.viewportSize <= 0 {
		return
	}

	// Если курсор выше viewport, сдвигаем viewport вверх
	if t.cursor < t.viewportStart {
		t.viewportStart = t.cursor
	}

	// Если курсор ниже viewport, сдвигаем viewport вниз
	if t.cursor >= t.viewportStart+t.viewportSize {
		t.viewportStart = t.cursor - t.viewportSize + 1
	}

	// Убеждаемся, что viewport не выходит за границы списка
	if t.viewportStart < 0 {
		t.viewportStart = 0
	}

	maxStart := len(t.choices) - t.viewportSize
	if maxStart < 0 {
		maxStart = 0
	}
	if t.viewportStart > maxStart {
		t.viewportStart = maxStart
	}
}

// getVisibleRange возвращает диапазон видимых элементов с учетом viewport
// Возвращает: startIdx, endIdx
func (t *SingleSelectTask) getVisibleRange() (int, int) {
	// Если viewport отключен, показываем все элементы
	if t.viewportSize <= 0 {
		return 0, len(t.choices)
	}

	startIdx := t.viewportStart
	if startIdx < 0 {
		startIdx = 0
	}

	endIdx := startIdx + t.viewportSize
	if endIdx > len(t.choices) {
		endIdx = len(t.choices)
	}

	return startIdx, endIdx
}

// stopTimeout останавливает таймер
func (t *SingleSelectTask) stopTimeout() {
	// Если таймер активен, останавливаем его
	if t.timeoutEnabled && t.timeoutManager != nil && t.timeoutManager.IsActive() {
		t.timeoutManager.StopTimeout()
		t.showTimeout = false
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
		// При нажатии клавиш сбрасываем таймер
		switch msg.String() {
		case "up", "k":
			// Если таймер активен, останавливаем его
			t.stopTimeout()
			if t.moveCursorBackward() {
				// Обновляем viewport после изменения позиции курсора
				t.updateViewport()
			}
			return t, nil
		case "down", "j":
			// Если таймер активен, останавливаем его
			t.stopTimeout()
			if t.moveCursorForward() {
				// Обновляем viewport после изменения позиции курсора
				t.updateViewport()
			}
			return t, nil
		case "q", "Q", "esc", "Esc", "ctrl+c", "Ctrl+C", "left", "Left":
			// Отмена пользователем
			cancelErr := fmt.Errorf(defauilt.ErrorMsgCanceled)
			t.done = true
			t.err = cancelErr
			t.icon = ui.IconCancelled
			t.finalValue = ui.GetErrorMessageStyle().Render(cancelErr.Error())
			t.SetStopOnError(true)
			return t, nil

		case "enter", "right", "Right":
			// Если таймер активен, останавливаем его
			t.stopTimeout()
			if t.cursor < 0 || t.cursor >= len(t.choices) || t.isDisabled(t.cursor) {
				return t, nil
			}
			t.done = true
			t.icon = ui.IconDone
			t.finalValue = t.choices[t.cursor]
			return t, nil
		case " ":
			// Если таймер активен, останавливаем его
			t.stopTimeout()
			// В любом случае выбираем текущий элемент
			if t.cursor < 0 || t.cursor >= len(t.choices) || t.isDisabled(t.cursor) {
				return t, nil
			}
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
		targetIndex := -1
		switch val := t.defaultValue.(type) {
		case int:
			// Если индекс в допустимом диапазоне
			if val >= 0 && val < len(t.choices) {
				targetIndex = val
			}
		case string:
			// Ищем строку в списке вариантов
			for i, choice := range t.choices {
				if choice == val {
					targetIndex = i
					break
				}
			}
		}

		if targetIndex >= 0 {
			t.cursor = targetIndex
			if t.ensureCursorSelectable() && t.cursor >= 0 {
				// Устанавливаем задачу как завершенную
				t.done = true
				t.icon = ui.IconDone
				t.finalValue = t.choices[t.cursor]
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

	// Получаем диапазон видимых элементов с учетом viewport
	startIdx, endIdx := t.getVisibleRange()

	// Добавляем индикатор прокрутки вверх, если есть скрытые элементы выше
	if t.viewportSize > 0 && startIdx > 0 {
		indentPrefix := ui.GetSelectItemPrefix("above")
		sb.WriteString(ui.SubtleStyle.Render(fmt.Sprintf(defauilt.ScrollAboveFormat, indentPrefix, ui.UpArrowSymbol, startIdx)))
		sb.WriteString("\n")
	}

	// Отображаем только видимые элементы списка
	for i := startIdx; i < endIdx; i++ {
		if i >= len(t.choices) {
			break
		}

		choice := t.choices[i]
		checked := ui.IconRadioOff
		var itemPrefix string
		isDisabled := t.isDisabled(i)
		if isDisabled {
			choice = ui.DisabledStyle.Render(choice)
			checked = ui.DisabledStyle.Render(checked)
		}

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
			openBracket := "("
			closeBracket := ")"
			if isDisabled {
				openBracket = ui.DisabledStyle.Render(openBracket)
				closeBracket = ui.DisabledStyle.Render(closeBracket)
			}
			sb.WriteString(fmt.Sprintf("%s%s%s%s %s\n", itemPrefix, openBracket, checked, closeBracket, choice))
		}
	}

	// Добавляем индикатор прокрутки вниз, если есть скрытые элементы ниже
	if t.viewportSize > 0 && endIdx < len(t.choices) {
		indentPrefix := ui.GetSelectItemPrefix("below")
		sb.WriteString(ui.SubtleStyle.Render(fmt.Sprintf(defauilt.ScrollBelowFormat, indentPrefix, ui.DownArrowSymbol, len(t.choices)-endIdx)))
		sb.WriteString("\n")
	}

	// Добавляем подсказку о навигации с новым отступом
	helpIndent := performance.RepeatEfficient(" ", ui.MainLeftIndent)
	sb.WriteString("\n" + ui.DrawLine(width) +
		ui.SubtleStyle.Render(fmt.Sprintf("%s%s", helpIndent, defauilt.SingleSelectHelp)))

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
