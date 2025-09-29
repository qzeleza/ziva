package task

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qzeleza/ziva/internal/defaults"
	"github.com/qzeleza/ziva/internal/performance"
	"github.com/qzeleza/ziva/internal/ui"
)

// SingleSelectTask - задача для выбора одного варианта из списка.
type SingleSelectTask struct {
	BaseTask
	items       []choice         // Список вариантов выбора
	disabled    map[int]struct{} // Набор отключённых пунктов меню
	cursor      int              // Текущая позиция курсора
	activeStyle lipgloss.Style   // Стиль для активного элемента
	selectedKey string           // Сохранённый ключ выбранного элемента
	// Viewport (окно просмотра) для ограничения количества отображаемых элементов
	viewportSize  int // Размер viewport (количество видимых элементов), 0 = показать все
	viewportStart int // Начальная позиция viewport в списке элементов
	showCounters  bool
}

// NewSingleSelectTask создает новую задачу выбора одного варианта из списка.
//
// @param title Заголовок задачи
// @param items Список элементов выбора
// @return Указатель на новую задачу выбора
func NewSingleSelectTask(title string, items []Item) *SingleSelectTask {
	normalized := normalizeItems(items)

	task := &SingleSelectTask{
		BaseTask:    NewBaseTask(title),
		items:       normalized,
		disabled:    make(map[int]struct{}),
		activeStyle: ui.ActiveStyle,
		// Viewport по умолчанию отключен (показываем все элементы)
		viewportSize:  0,
		viewportStart: 0,
		showCounters:  true,
	}

	task.ensureCursorSelectable()
	return task
}

func (t *SingleSelectTask) captureSelection(index int) {
	if index < 0 || index >= len(t.items) {
		return
	}
	selection := t.items[index]
	t.selectedKey = selection.key
	t.finalValue = selection.name
}

// isDisabled проверяет, помечен ли элемент как недоступный
func (t *SingleSelectTask) isDisabled(index int) bool {
	if index < 0 || index >= len(t.items) {
		return true
	}
	_, exists := t.disabled[index]
	return exists
}

// ensureCursorSelectable подбирает ближайший доступный элемент для курсора.
// Возвращает true, если удалось найти активный элемент.
func (t *SingleSelectTask) ensureCursorSelectable() bool {
	if len(t.items) == 0 {
		t.cursor = -1
		return false
	}

	if t.cursor >= 0 && t.cursor < len(t.items) && !t.isDisabled(t.cursor) {
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
	for i := from; i < len(t.items); i++ {
		if !t.isDisabled(i) {
			return i, true
		}
	}
	return -1, false
}

// findEnabledBackward возвращает индекс ближайшего доступного элемента, двигаясь в обратном порядке
func (t *SingleSelectTask) findEnabledBackward(from int) (int, bool) {
	if from >= len(t.items) {
		from = len(t.items) - 1
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
		if idx >= 0 && idx < len(t.items) {
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
		if idx < 0 || idx >= len(t.items) {
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
	normalized := strings.TrimSpace(value)
	if normalized == "" {
		return -1
	}
	for i, item := range t.items {
		if strings.EqualFold(item.key, normalized) || strings.EqualFold(item.name, normalized) {
			return i
		}
	}
	return -1
}

// WithDefaultItem устанавливает элемент, который будет подсвечен курсором при открытии списка.
// Поддерживает выбор по индексу (int) или строковому значению (string). Некорректные значения игнорируются.
func (t *SingleSelectTask) WithDefaultItem(selection interface{}) *SingleSelectTask {
	if selection == nil || len(t.items) == 0 {
		return t
	}

	setCursor := func(index int) {
		if index < 0 || index >= len(t.items) {
			return
		}
		t.cursor = index
	}

	switch v := selection.(type) {
	case int:
		setCursor(v)
	case string:
		if idx := t.choiceIndex(v); idx != -1 {
			setCursor(idx)
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
func (t *SingleSelectTask) WithViewport(size int, showCounters ...bool) *SingleSelectTask {
	if size < 0 {
		size = 0
	}
	t.viewportSize = size
	t.showCounters = true
	if len(showCounters) > 0 {
		t.showCounters = showCounters[0]
	}
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

	maxStart := len(t.items) - t.viewportSize
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
		return 0, len(t.items)
	}

	startIdx := t.viewportStart
	if startIdx < 0 {
		startIdx = 0
	}

	endIdx := startIdx + t.viewportSize
	if endIdx > len(t.items) {
		endIdx = len(t.items)
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
			cancelErr := fmt.Errorf(defaults.ErrorMsgCanceled)
			t.done = true
			t.err = cancelErr
			t.icon = ui.IconCancelled
			t.finalValue = ui.GetErrorMessageStyle().Render(cancelErr.Error())
			t.SetStopOnError(true)
			return t, nil

		case "enter", "right", "Right":
			// Если таймер активен, останавливаем его
			t.stopTimeout()
			if t.cursor < 0 || t.cursor >= len(t.items) || t.isDisabled(t.cursor) {
				return t, nil
			}
			t.done = true
			t.icon = ui.IconDone
			t.captureSelection(t.cursor)
			return t, nil
		case " ":
			// Если таймер активен, останавливаем его
			t.stopTimeout()
			// В любом случае выбираем текущий элемент
			if t.cursor < 0 || t.cursor >= len(t.items) || t.isDisabled(t.cursor) {
				return t, nil
			}
			t.done = true
			t.icon = ui.IconDone
			t.captureSelection(t.cursor)
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
			if val >= 0 && val < len(t.items) {
				targetIndex = val
			}
		case string:
			if idx := t.choiceIndex(val); idx != -1 {
				targetIndex = idx
			}
		}

		if targetIndex >= 0 {
			t.cursor = targetIndex
			if t.ensureCursorSelectable() && t.cursor >= 0 {
				// Устанавливаем задачу как завершенную
				t.done = true
				t.icon = ui.IconDone
				t.captureSelection(t.cursor)
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

	// Добавляем заголовок задачи с префиксом для активной задачи (учитываем нумерацию)
	titlePrefix := t.InProgressPrefix()

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

	sb.WriteString(renderSelectionSeparator(width, t.showSelectionSeparator, titlePrefix))

	// Получаем диапазон видимых элементов с учетом viewport
	startIdx, endIdx := t.getVisibleRange()

	// Добавляем индикатор прокрутки вверх, если есть скрытые элементы выше
	if t.viewportSize > 0 && startIdx > 0 {
		indentPrefix := ui.GetSelectItemPrefix("above")
		var indicator string
		if t.showCounters {
			arrow := ui.UpArrowSymbol + " "
			indicator = fmt.Sprintf(defaults.ScrollAboveFormat, indentPrefix, arrow, startIdx)
		} else {
			indicator = fmt.Sprintf("%s %s", indentPrefix, ui.UpArrowSymbol)
		}
		appendIndicatorWithPlainPipe(&sb, indicator)
		sb.WriteString("\n")
	}

	activeHelp := ""

	// Отображаем только видимые элементы списка
	for i := startIdx; i < endIdx; i++ {
		if i >= len(t.items) {
			break
		}

		item := t.items[i]
		label := item.displayName()
		description := item.helpText()
		checked := ui.IconRadioOff
		var itemPrefix string
		isDisabled := t.isDisabled(i)
		if isDisabled {
			label = ui.DisabledStyle.Render(label)
			checked = ui.DisabledStyle.Render(checked)
		}

		if t.cursor == i {
			itemPrefix = ui.GetSelectItemPrefix("active")
			checked = ui.IconRadioOn
			label = t.activeStyle.Render(label)
			checked = t.activeStyle.Render(checked)
		} else if i < t.cursor {
			itemPrefix = ui.GetSelectItemPrefix("above")
		} else {
			itemPrefix = ui.GetSelectItemPrefix("below")
		}

		if t.cursor == i {
			bracketsOpen := t.activeStyle.Render("(")
			bracketsClose := t.activeStyle.Render(")")
			sb.WriteString(fmt.Sprintf("%s%s%s%s %s\n", itemPrefix, bracketsOpen, checked, bracketsClose, label))
		} else {
			openBracket := "("
			closeBracket := ")"
			if isDisabled {
				openBracket = ui.DisabledStyle.Render(openBracket)
				closeBracket = ui.DisabledStyle.Render(closeBracket)
			}
			sb.WriteString(fmt.Sprintf("%s%s%s%s %s\n", itemPrefix, openBracket, checked, closeBracket, label))
		}

		if t.cursor == i && strings.TrimSpace(description) != "" {
			activeHelp = description
		}
	}

	// Добавляем индикатор прокрутки вниз, если есть скрытые элементы ниже
	if t.viewportSize > 0 && endIdx < len(t.items) {
		indentPrefix := ui.GetSelectItemPrefix("below")
		var indicator string
		remaining := len(t.items) - endIdx
		if t.showCounters {
			arrow := ui.DownArrowSymbol + " "
			indicator = fmt.Sprintf(defaults.ScrollBelowFormat, indentPrefix, arrow, remaining)
		} else {
			indicator = fmt.Sprintf("%s %s", indentPrefix, ui.DownArrowSymbol)
		}
		appendIndicatorWithPlainPipe(&sb, indicator)
		sb.WriteString("\n")
	}

	// Добавляем подсказку о навигации с новым отступом
	helpIndent := performance.RepeatEfficient(" ", ui.MainLeftIndent)

	sb.WriteString("\n" + ui.DrawLine(width))
	if activeHelp != "" {
		sb.WriteString(ui.HelpTextStyle.Render(fmt.Sprintf("%s%s", helpIndent, activeHelp)))
		sb.WriteString("\n")
	}
	sb.WriteString(ui.SubtleStyle.Render(fmt.Sprintf("%s%s", helpIndent, defaults.SingleSelectHelp)))

	return sb.String()
}

func (t *SingleSelectTask) FinalView(width int) string {
	// Получаем базовое финальное представление
	result := t.BaseTask.FinalView(width)

	// Если задача завершилась успешно и есть дополнительные строки для вывода
	if t.icon == ui.IconDone && len(t.items) > 0 && t.cursor >= 0 && t.cursor < len(t.items) {
		result += "\n" + ui.DrawSummaryLine(t.items[t.cursor].displayName())
	}

	return result
}

// GetSelected возвращает выбранное значение из списка
//
// @return string Выбранный пользователем вариант
func (t *SingleSelectTask) GetSelected() string {
	if !t.done && t.cursor >= 0 && t.cursor < len(t.items) {
		return t.items[t.cursor].valueKey()
	}
	if t.selectedKey != "" {
		return t.selectedKey
	}
	if t.cursor >= 0 && t.cursor < len(t.items) {
		return t.items[t.cursor].valueKey()
	}
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
