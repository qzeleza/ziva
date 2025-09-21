package task

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qzeleza/termos/internal/defaults"
	"github.com/qzeleza/termos/internal/performance"
	"github.com/qzeleza/termos/internal/ui"
)

// SelectionBitset представляет битовый набор для оптимальной работы с embedded устройствами
// Использует uint32 для лучшей производительности на 32-битных embedded системах
// Поддерживает до 32 элементов выбора, что покрывает большинство практических случаев
type SelectionBitset uint32

// Set устанавливает бит в позиции index
func (s *SelectionBitset) Set(index int) {
	if index >= 0 && index < 32 {
		*s |= 1 << index
	}
}

// Clear очищает бит в позиции index
func (s *SelectionBitset) Clear(index int) {
	if index >= 0 && index < 32 {
		*s &^= 1 << index
	}
}

// IsSet проверяет, установлен ли бит в позиции index
func (s SelectionBitset) IsSet(index int) bool {
	if index < 0 || index >= 32 {
		return false
	}
	return s&(1<<index) != 0
}

// Toggle переключает бит в позиции index
func (s *SelectionBitset) Toggle(index int) {
	if index >= 0 && index < 32 {
		*s ^= 1 << index
	}
}

// Count возвращает количество установленных битов (оптимизировано для 32-бит)
func (s SelectionBitset) Count() int {
	// Используем встроенную функцию подсчета битов для оптимизации
	// На современных процессорах это компилируется в инструкцию POPCNT
	return popcount32(uint32(s))
}

// popcount32 подсчитывает установленные биты в 32-битном числе
func popcount32(x uint32) int {
	// Алгоритм Брайана Кернигана - эффективен для разреженных битов
	count := 0
	for x != 0 {
		x &= x - 1 // Убираем самый младший установленный бит
		count++
	}
	return count
}

// Clear очищает все биты
func (s *SelectionBitset) ClearAll() {
	*s = 0
}

// SetAll устанавливает биты для первых n позиций
func (s *SelectionBitset) SetAll(n int) {
	if n <= 0 {
		*s = 0
		return
	}
	if n >= 32 {
		*s = SelectionBitset(^uint32(0))
		return
	}
	*s = SelectionBitset((1 << n) - 1)
}

// MultiSelectTask позволяет выбрать несколько вариантов из списка.
type MultiSelectTask struct {
	BaseTask
	choices         []string         // Список вариантов выбора
	disabled        map[int]struct{} // Набор отключённых пунктов
	selected        SelectionBitset  // Битовый набор выбранных элементов (оптимизировано для embedded)
	fallbackMap     map[int]struct{} // Резервная карта для списков > 64 элементов
	cursor          int              // Текущая позиция курсора
	activeStyle     lipgloss.Style   // Стиль для активного элемента
	hasSelectAll    bool             // Включена ли опция "Выбрать все"
	selectAllText   string           // Текст опции "Выбрать все"
	showHelpMessage bool             // Показывать ли сообщение-подсказку
	helpMessage     string           // Текст сообщения-подсказки
	itemHelps       []string         // Справочные сообщения для элементов
	// Viewport (окно просмотра) для ограничения количества отображаемых элементов
	viewportSize  int // Размер viewport (количество видимых элементов), 0 = показать все
	viewportStart int // Начальная позиция viewport в списке элементов
	showCounters  bool
}

// NewMultiSelectTask создает новую задачу множественного выбора.
//
// @param title Заголовок задачи
// @param choices Список вариантов выбора
// @return Указатель на новую задачу множественного выбора
func NewMultiSelectTask(title string, choices []string) *MultiSelectTask {
	labels, helps := parseChoicesWithHelp(choices)

	task := &MultiSelectTask{
		BaseTask:        NewBaseTask(title),
		choices:         labels,
		itemHelps:       helps,
		disabled:        make(map[int]struct{}),
		cursor:          0, // Начинаем с первого элемента списка
		activeStyle:     ui.ActiveStyle,
		hasSelectAll:    false,
		selectAllText:   defaults.SelectAllDefaultText,
		showHelpMessage: false,
		helpMessage:     "",
		// Viewport по умолчанию отключен (показываем все элементы)
		viewportSize:  0,
		viewportStart: 0,
		showCounters:  true,
	}

	// Для embedded устройств: используем битсет для списков <= 32 элементов,
	// иначе fallback на карту
	if len(choices) > 32 {
		task.fallbackMap = make(map[int]struct{})
	}

	task.ensureCursorSelectable()
	return task
}

// isDisabled проверяет, помечен ли элемент как недоступный
func (t *MultiSelectTask) isDisabled(index int) bool {
	if index < 0 || index >= len(t.choices) {
		return true
	}
	_, exists := t.disabled[index]
	return exists
}

// ensureCursorSelectable пытается разместить курсор на ближайшем доступном элементе
func (t *MultiSelectTask) ensureCursorSelectable() bool {
	if len(t.choices) == 0 {
		if t.hasSelectAll {
			t.cursor = -1
		} else {
			t.cursor = -1
		}
		return false
	}

	if t.hasSelectAll && t.cursor == -1 {
		return true
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

	if t.hasSelectAll {
		t.cursor = -1
		return true
	}

	t.cursor = -1
	return false
}

// findEnabledForward возвращает индекс первого доступного элемента начиная с from (включительно)
func (t *MultiSelectTask) findEnabledForward(from int) (int, bool) {
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

// findEnabledBackward возвращает индекс предыдущего доступного элемента
func (t *MultiSelectTask) findEnabledBackward(from int) (int, bool) {
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
func (t *MultiSelectTask) moveCursorForward() bool {
	if t.cursor == -1 {
		if idx, ok := t.findEnabledForward(0); ok {
			t.cursor = idx
			return true
		}
		return false
	}

	start := t.cursor + 1
	if idx, ok := t.findEnabledForward(start); ok {
		t.cursor = idx
		return true
	}
	return false
}

// moveCursorBackward перемещает курсор на предыдущий доступный элемент
func (t *MultiSelectTask) moveCursorBackward() bool {
	if t.cursor == -1 {
		return false
	}
	if idx, ok := t.findEnabledBackward(t.cursor - 1); ok {
		t.cursor = idx
		return true
	}
	if t.hasSelectAll {
		t.cursor = -1
		return true
	}
	return false
}

// resolveDisabledIndices конвертирует произвольный ввод в индексы элементов списка
func (t *MultiSelectTask) resolveDisabledIndices(input interface{}) []int {
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
func (t *MultiSelectTask) choiceIndex(value string) int {
	normalized, _ := splitChoiceAndHelp(value)
	normalized = strings.TrimSpace(normalized)
	for i, choice := range t.choices {
		if choice == value || choice == normalized {
			return i
		}
	}
	return -1
}

// clearDisabledSelections снимает выбор с отключённых элементов
func (t *MultiSelectTask) clearDisabledSelections() {
	if len(t.choices) > 32 {
		if t.fallbackMap == nil {
			return
		}
		for idx := range t.disabled {
			delete(t.fallbackMap, idx)
		}
		return
	}

	for idx := range t.disabled {
		t.selected.Clear(idx)
	}
}

// clearAllSelections снимает выбор со всех элементов
func (t *MultiSelectTask) clearAllSelections() {
	if len(t.choices) > 32 {
		if t.fallbackMap == nil {
			return
		}
		for idx := range t.fallbackMap {
			delete(t.fallbackMap, idx)
		}
		return
	}

	t.selected.ClearAll()
}

// selectAllEnabled отмечает все доступные элементы выбранными
func (t *MultiSelectTask) selectAllEnabled() {
	if len(t.choices) > 32 {
		if t.fallbackMap == nil {
			t.fallbackMap = make(map[int]struct{})
		}
		for i := range t.choices {
			if t.isDisabled(i) {
				delete(t.fallbackMap, i)
				continue
			}
			t.fallbackMap[i] = struct{}{}
		}
		return
	}

	t.selected.ClearAll()
	for i := range t.choices {
		if t.isDisabled(i) {
			continue
		}
		t.selected.Set(i)
	}
}

// WithViewport устанавливает размер viewport (окна просмотра) для ограничения количества отображаемых элементов.
// Это полезно для длинных списков, когда нужно показывать только часть элементов.
//
// @param size Количество элементов для отображения одновременно (0 = показать все)
// @return Указатель на задачу для цепочки вызовов
func (t *MultiSelectTask) WithViewport(size int, showCounters ...bool) *MultiSelectTask {
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

// WithItemsDisabled помечает элементы меню как недоступные для выбора.
// Поддерживаются типы: int, []int, string, []string. Nil очищает список отключённых элементов.
func (t *MultiSelectTask) WithItemsDisabled(disabled interface{}) *MultiSelectTask {
	for idx := range t.disabled {
		delete(t.disabled, idx)
	}

	indices := t.resolveDisabledIndices(disabled)
	for _, idx := range indices {
		if idx >= 0 && idx < len(t.choices) {
			t.disabled[idx] = struct{}{}
		}
	}

	t.clearDisabledSelections()
	t.ensureCursorSelectable()
	t.updateViewport()
	return t
}

// updateViewport обновляет позицию viewport на основе текущего положения курсора
func (t *MultiSelectTask) updateViewport() {
	// Если viewport отключен, ничего не делаем
	if t.viewportSize <= 0 {
		return
	}

	// Получаем эффективную позицию курсора (с учетом опции "Выбрать все")
	effectiveCursor := t.cursor
	if t.hasSelectAll {
		effectiveCursor = t.cursor + 1 // +1 потому что опция "Выбрать все" занимает позицию -1
	}

	// Если курсор выше viewport, сдвигаем viewport вверх
	if effectiveCursor < t.viewportStart {
		t.viewportStart = effectiveCursor
	}

	// Если курсор ниже viewport, сдвигаем viewport вниз
	if effectiveCursor >= t.viewportStart+t.viewportSize {
		t.viewportStart = effectiveCursor - t.viewportSize + 1
	}

	// Убеждаемся, что viewport не выходит за границы списка
	if t.viewportStart < 0 {
		t.viewportStart = 0
	}

	maxStart := len(t.choices) - t.viewportSize
	if t.hasSelectAll {
		maxStart = len(t.choices) + 1 - t.viewportSize // +1 для опции "Выбрать все"
	}
	if maxStart < 0 {
		maxStart = 0
	}
	if t.viewportStart > maxStart {
		t.viewportStart = maxStart
	}
}

// getVisibleRange возвращает диапазон видимых элементов с учетом viewport
// Возвращает: startIdx, endIdx, showSelectAll
func (t *MultiSelectTask) getVisibleRange() (int, int, bool) {
	// Если viewport отключен, показываем все элементы
	if t.viewportSize <= 0 {
		return 0, len(t.choices), t.hasSelectAll
	}

	// Определяем, показывать ли опцию "Выбрать все"
	showSelectAll := t.hasSelectAll && t.viewportStart == 0

	// Вычисляем диапазон элементов списка
	startIdx := t.viewportStart
	if t.hasSelectAll && startIdx > 0 {
		startIdx-- // Компенсируем опцию "Выбрать все"
	}

	if startIdx < 0 {
		startIdx = 0
	}

	endIdx := startIdx + t.viewportSize
	if showSelectAll {
		endIdx-- // Уменьшаем на 1, так как одно место занимает опция "Выбрать все"
	}

	if endIdx > len(t.choices) {
		endIdx = len(t.choices)
	}

	return startIdx, endIdx, showSelectAll
}

// WithSelectAll добавляет опцию "Выбрать все" в начало списка.
// При выборе этой опции все остальные пункты автоматически помечаются/снимаются.
//
// @param text Текст для опции "Выбрать все" (по умолчанию "Выбрать все")
// @return Указатель на задачу для цепочки вызовов
func (t *MultiSelectTask) WithSelectAll(text ...string) *MultiSelectTask {
	t.hasSelectAll = true
	t.cursor = -1 // Начинаем с опции "Выбрать все"
	if len(text) > 0 && strings.TrimSpace(text[0]) != "" {
		t.selectAllText = text[0]
	} else {
		t.selectAllText = defaults.SelectAllDefaultText
	}
	t.ensureCursorSelectable()
	t.updateViewport()
	return t
}

// WithDefaultItems позволяет заранее отметить элементы списка выбранными при открытии задачи.
// Поддерживает выбор одного индекса/строки или списков значений ([]int, []string).
func (t *MultiSelectTask) WithDefaultItems(defauiltSelection interface{}) *MultiSelectTask {
	if defauiltSelection == nil || len(t.choices) == 0 {
		return t
	}

	// Сбрасываем текущий выбор
	if len(t.choices) > 32 {
		if t.fallbackMap == nil {
			t.fallbackMap = make(map[int]struct{})
		} else {
			for k := range t.fallbackMap {
				delete(t.fallbackMap, k)
			}
		}
		t.selected.ClearAll()
	} else {
		t.selected.ClearAll()
	}

	setSelected := func(index int) bool {
		if index < 0 || index >= len(t.choices) {
			return false
		}
		if t.isDisabled(index) {
			return false
		}
		if len(t.choices) > 32 && t.fallbackMap != nil {
			t.fallbackMap[index] = struct{}{}
		} else {
			t.selected.Set(index)
		}
		return true
	}

	anyApplied := false

	switch v := defauiltSelection.(type) {
	case int:
		anyApplied = setSelected(v) || anyApplied
	case string:
		if idx := t.choiceIndex(v); idx != -1 {
			if setSelected(idx) {
				anyApplied = true
			}
		}
	case []int:
		for _, idx := range v {
			if setSelected(idx) {
				anyApplied = true
			}
		}
	case []string:
		for _, val := range v {
			if idx := t.choiceIndex(val); idx != -1 {
				if setSelected(idx) {
					anyApplied = true
				}
			}
		}
	}

	if anyApplied {
		// Перемещаем курсор на первый выбранный элемент (если он вне диапазона)
		if t.cursor < 0 || t.cursor >= len(t.choices) {
			for i := range t.choices {
				if t.isSelected(i) {
					t.cursor = i
					break
				}
			}
		}
	}

	// После обновления выбора синхронизируем viewport
	t.updateViewport()
	return t
}

// toggleSelectAll переключает состояние всех элементов списка.
// Если все элементы выбраны - снимает выбор со всех.
// Если хотя бы один элемент не выбран - выбирает все.
func (t *MultiSelectTask) toggleSelectAll() {
	if len(t.choices) == 0 {
		return
	}

	if t.isAllSelected() {
		t.clearAllSelections()
		return
	}

	t.selectAllEnabled()
}

// isAllSelected проверяет, выбраны ли все элементы списка
func (t *MultiSelectTask) isAllSelected() bool {
	enabledCount := 0
	selectedCount := 0
	for i := range t.choices {
		if t.isDisabled(i) {
			continue
		}
		enabledCount++
		if t.isSelected(i) {
			selectedCount++
		}
	}

	if enabledCount == 0 {
		return false
	}

	return enabledCount == selectedCount
}

// toggleSelection переключает состояние выбора элемента по индексу (оптимизировано для embedded)
func (t *MultiSelectTask) toggleSelection(index int) {
	if t.isDisabled(index) {
		return
	}

	if len(t.choices) > 32 && t.fallbackMap != nil {
		// Используем fallback карту для больших списков
		if _, exists := t.fallbackMap[index]; exists {
			delete(t.fallbackMap, index)
		} else {
			t.fallbackMap[index] = struct{}{}
		}
	} else {
		// Используем битсет для оптимизации
		t.selected.Toggle(index)
	}
}

// isSelected проверяет, выбран ли элемент по индексу (оптимизировано для embedded)
func (t *MultiSelectTask) isSelected(index int) bool {
	if t.isDisabled(index) {
		return false
	}
	if len(t.choices) > 32 && t.fallbackMap != nil {
		_, exists := t.fallbackMap[index]
		return exists
	}
	return t.selected.IsSet(index)
}

// stopTimeout останавливает таймер
func (t *MultiSelectTask) stopTimeout() {
	// Если таймер активен, останавливаем его
	if t.timeoutEnabled && t.timeoutManager != nil && t.timeoutManager.IsActive() {
		t.timeoutManager.StopTimeout()
		t.showTimeout = false
	}
}

// Update handles key presses for navigation and selection.
func (t *MultiSelectTask) Update(msg tea.Msg) (Task, tea.Cmd) {
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
		// Сбрасываем сообщение-подсказку при любом нажатии клавиш (кроме Enter)
		if msg.String() != "enter" {
			t.showHelpMessage = false
			t.helpMessage = ""
		}

		switch msg.String() {
		case "up", "k":
			t.stopTimeout()
			if t.moveCursorBackward() {
				// Обновляем viewport после изменения позиции курсора
				t.updateViewport()
			}
			return t, nil
		case "down", "j":
			t.stopTimeout()
			if t.moveCursorForward() {
				// Обновляем viewport после изменения позиции курсора
				t.updateViewport()
			}
			return t, nil
		case " ", "right", "Right":
			// При выборе останавливаем таймер
			t.stopTimeout()

			// В любом случае выполняем выбор/переключение
			if t.hasSelectAll && t.cursor == -1 {
				// Нажатие пробела на опции "Выбрать все"
				t.toggleSelectAll()
			} else if t.cursor >= 0 && !t.isDisabled(t.cursor) {
				// Обычная логика выбора для элементов списка
				t.toggleSelection(t.cursor)
			}
		case "q", "Q", "esc", "Esc", "ctrl+c", "Ctrl+C", "left", "Left":
			// Отмена пользователем
			cancelErr := fmt.Errorf(defaults.ErrorMsgCanceled)
			t.done = true
			t.err = cancelErr
			t.icon = ui.IconCancelled
			t.finalValue = ui.CancelStyle.Render(cancelErr.Error())
			t.SetStopOnError(true)
			return t, nil

		case "enter":
			t.stopTimeout()
			// Собираем выбранные элементы
			var selectedChoices []string
			for i := range t.choices {
				if t.isSelected(i) {
					selectedChoices = append(selectedChoices, t.choices[i])
				}
			}
			if len(selectedChoices) == 0 {
				// Если ничего не выбрано, показываем сообщение-подсказку
				// но НЕ устанавливаем ошибку и НЕ завершаем задачу
				t.showHelpMessage = true
				t.helpMessage = defaults.NeedSelectAtLeastOne
				return t, nil
			}
			// Если есть выбранные элементы, завершаем задачу успешно
			t.done = true
			t.icon = ui.IconDone
			t.finalValue = strings.Join(selectedChoices, defaults.DefaultSeparator)
			// Убеждаемся, что ошибка очищена
			t.SetError(nil)
			t.showHelpMessage = false
			t.helpMessage = ""
			return t, nil
		}
		// После обработки клавиш возвращаем команду для продолжения тикера
		if t.timeoutEnabled && t.timeoutManager != nil && t.timeoutManager.IsActive() {
			return t, t.timeoutManager.StartTicker()
		}
	}

	// Запускаем таймер при первом обновлении, если он включен и еще не активен
	if t.timeoutEnabled && t.timeoutManager != nil && !t.timeoutManager.IsActive() {
		return t, t.timeoutManager.StartTickerAndTimeout()
	}

	return t, nil
}

// Run запускает задачу выбора
func (t *MultiSelectTask) Run() tea.Cmd {
	// Запускаем таймер и тикер, если они включены
	if t.timeoutEnabled && t.timeoutManager != nil {
		return t.timeoutManager.StartTickerAndTimeout()
	}
	return nil
}

// applyDefaultValue применяет значение по умолчанию при истечении таймера
func (t *MultiSelectTask) applyDefaultValue() {
	// Если есть значение по умолчанию
	if t.defaultValue != nil {
		switch val := t.defaultValue.(type) {
		case []int:
			// Если это список индексов для выбора
			for _, index := range val {
				// Выбираем только корректные индексы
				if index >= 0 && index < len(t.choices) && !t.isDisabled(index) {
					// Устанавливаем выбор для этого элемента
					if len(t.choices) > 32 && t.fallbackMap != nil {
						t.fallbackMap[index] = struct{}{}
					} else {
						t.selected.Set(index)
					}
				}
			}
		case []string:
			// Если это список строк для выбора
			for _, strVal := range val {
				if idx := t.choiceIndex(strVal); idx != -1 && !t.isDisabled(idx) {
					if len(t.choices) > 32 && t.fallbackMap != nil {
						t.fallbackMap[idx] = struct{}{}
					} else {
						t.selected.Set(idx)
					}
				}
			}
		}

		// Проверяем, есть ли хотя бы один выбранный элемент
		hasSelection := false
		if len(t.choices) > 32 && t.fallbackMap != nil {
			hasSelection = len(t.fallbackMap) > 0
		} else {
			hasSelection = t.selected.Count() > 0
		}

		// Если есть выбранные элементы, завершаем задачу
		if hasSelection {
			// Собираем выбранные элементы
			var selectedChoices []string
			for i := range t.choices {
				if t.isSelected(i) {
					selectedChoices = append(selectedChoices, t.choices[i])
				}
			}
			// Завершаем задачу
			t.done = true
			t.icon = ui.IconDone
			t.finalValue = strings.Join(selectedChoices, defaults.DefaultSeparator)
			t.SetError(nil)
		}
	}
}

// View отрисовывает список вариантов выбора для пользователя с выделением активного элемента.
//
// @param width Ширина макета для отображения
// @return Строка с отформатированным представлением задачи
func (t *MultiSelectTask) View(width int) string {
	// Если задача завершена, возвращаем FinalView
	if t.done {
		return t.FinalView(width)
	}

	var sb strings.Builder
	// Заголовок задачи с новым префиксом для текущей задачи
	titlePrefix := ui.GetCurrentTaskPrefix()

	// Формируем заголовок с префиксом
	title := ui.ActiveTaskStyle.Render(t.title)
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
	startIdx, endIdx, showSelectAll := t.getVisibleRange()

	// Отображаем опцию "Выбрать все" если она включена и видима в viewport
	if showSelectAll {
		checked := " "
		var itemPrefix string
		selectAllText := t.selectAllText

		// Проверяем, выбраны ли все элементы
		if t.isAllSelected() {
			checked = ui.IconSelected
		}

		// Определяем префикс для опции "Выбрать все"
		if t.cursor == -1 {
			// Опция "Выбрать все" активна
			itemPrefix = ui.GetSelectItemPrefix("active")
			selectAllText = t.activeStyle.Render(selectAllText)
		} else {
			// Опция "Выбрать все" не активна - всегда показываем префикс "above"
			// потому что курсор находится ниже неё (на элементах списка)
			itemPrefix = ui.GetSelectItemPrefix("above")
		}

		// Формируем строку для отображения опции "Выбрать все"
		sb.WriteString(fmt.Sprintf("%s[%s] %s\n", itemPrefix, checked, selectAllText))
	}

	// Добавляем индикатор прокрутки вверх, если есть скрытые элементы выше
	// При наличии пункта "Выбрать все" индикатор должен показываться даже если startIdx == 0
	if t.viewportSize > 0 && (startIdx > 0 || (t.hasSelectAll && t.viewportStart > 0)) {
		// Используем точно такой же префикс как у элементов "above"
		indentPrefix := ui.GetSelectItemPrefix("above")
		// Определяем количество элементов выше
		itemsAbove := startIdx
		if t.hasSelectAll && t.viewportStart > 0 {
			// Если есть пункт "Выбрать все" и он скрыт, добавляем +1 к счетчику
			itemsAbove = t.viewportStart
		}
		var indicator string
		if t.showCounters {
			arrow := ui.UpArrowSymbol + " "
			indicator = fmt.Sprintf(defaults.ScrollAboveFormat, indentPrefix, arrow, itemsAbove)
		} else {
			indicator = fmt.Sprintf("%s %s", indentPrefix, ui.UpArrowSymbol)
		}
		// Не добавляем перенос строки в самой строке, чтобы не нарушать форматирование
		sb.WriteString(ui.SubtleStyle.Render(indicator))
		// Добавляем перенос строки отдельно
		sb.WriteString("\n")
	}

	activeHelp := ""

	// Отображаем только видимые элементы списка
	for i := startIdx; i < endIdx; i++ {
		if i >= len(t.choices) {
			break
		}

		choice := t.choices[i]
		checked := " "
		var itemPrefix string
		itemDisabled := t.isDisabled(i)
		helpsAvailable := i < len(t.itemHelps)

		// Проверяем, выбран ли этот элемент
		if t.isSelected(i) {
			checked = ui.IconSelected
		}

		if itemDisabled {
			choice = ui.DisabledStyle.Render(choice)
			checked = ui.DisabledStyle.Render(checked)
		}

		// Определяем тип элемента для получения правильного префикса
		if t.cursor == i {
			// Активный элемент
			itemPrefix = ui.GetSelectItemPrefix("active")
			// Применяем стиль активного элемента
			choice = t.activeStyle.Render(choice)
		} else if t.hasSelectAll && t.cursor == -1 {
			// Если активна опция "Выбрать все", все элементы списка должны быть "below"
			itemPrefix = ui.GetSelectItemPrefix("below")
		} else if i < t.cursor {
			// Элемент выше текущего активного элемента
			itemPrefix = ui.GetSelectItemPrefix("above")
		} else {
			// Элемент ниже текущего активного элемента
			itemPrefix = ui.GetSelectItemPrefix("below")
		}

		// Формируем строку для отображения варианта выбора с новым префиксом
		openBracket := "["
		closeBracket := "]"
		if itemDisabled {
			openBracket = ui.DisabledStyle.Render(openBracket)
			closeBracket = ui.DisabledStyle.Render(closeBracket)
		}
		sb.WriteString(fmt.Sprintf("%s%s%s%s %s\n", itemPrefix, openBracket, checked, closeBracket, choice))

		if t.cursor == i && helpsAvailable {
			help := strings.TrimSpace(t.itemHelps[i])
			if help != "" {
				activeHelp = help
			}
		}
	}

	// Добавляем индикатор прокрутки вниз, если есть скрытые элементы ниже
	if t.viewportSize > 0 && endIdx < len(t.choices) {
		// Используем точно такой же префикс как у элементов "below"
		indentPrefix := ui.GetSelectItemPrefix("below")
		// Не добавляем перенос строки в конце, чтобы не нарушать форматирование
		remaining := len(t.choices) - endIdx
		var indicator string
		if t.showCounters {
			arrow := ui.DownArrowSymbol + " "
			indicator = fmt.Sprintf(defaults.ScrollBelowFormat, indentPrefix, arrow, remaining)
		} else {
			indicator = fmt.Sprintf("%s %s", indentPrefix, ui.DownArrowSymbol)
		}
		sb.WriteString(ui.SubtleStyle.Render(indicator))
		// Добавляем перенос строки отдельно
		sb.WriteString("\n")
	}

	// Формируем отступ для подсказки
	helpIndent := performance.RepeatEfficient(" ", ui.MainLeftIndent)

	// Новая строка для подсказки
	var helpLine string
	if activeHelp != "" {
		helpLine = "\n"
	}

	// Добавляем сообщение-подсказку если нужно
	var warning string
	if t.showHelpMessage && t.helpMessage != "" {
		activeHelp = ""
		helpLine = ""
		warning = ui.GetErrorMessageStyle().Render(fmt.Sprintf("%s%s", helpIndent, t.helpMessage))
	}

	// Если есть опция "Выбрать все", добавляем её в подсказку
	helpText := defaults.MultiSelectHelp
	if t.hasSelectAll {
		helpText = defaults.MultiSelectHelpSelectAll
	}
	// Добавляем разделительную линию
	sb.WriteString("\n" + ui.DrawLine(width))
	// Добавляем сообщение-подсказку если нужно
	if warning != "" {
		sb.WriteString(warning + "\n")
	}
	// Если есть активный элемент, добавляем его подсказку
	if activeHelp != "" {
		sb.WriteString(ui.HelpTextStyle.Render(fmt.Sprintf("%s%s", helpIndent, activeHelp)))
	}
	// Добавляем подсказку
	sb.WriteString(ui.SubtleStyle.Render(fmt.Sprintf("%s%s%s", helpLine, helpIndent, helpText)))

	return sb.String()
}

func (t *MultiSelectTask) FinalView(width int) string {
	// Получаем базовое финальное представление
	result := t.BaseTask.FinalView(width)

	// Если задача завершилась успешно и есть дополнительные строки для вывода
	if t.icon == ui.IconDone && len(t.choices) > 0 {
		selected := t.GetSelected()
		if len(selected) > 0 {
			result += "\n"
			for _, value := range selected {
				result += ui.DrawSummaryLine(value)
			}
		}
		result += performance.RepeatEfficient(" ", ui.MainLeftIndent) + ui.VerticalLineSymbol
	}

	return result
}

// WithDefaultOptions устанавливает варианты по умолчанию, которые будут выбраны при тайм-ауте.
// @param defauiltOptions Варианты по умолчанию (список индексов или строк)
// @param timeout Длительность тайм-аута
// @return Указатель на задачу для цепочки вызовов
func (t *MultiSelectTask) WithDefaultOptions(defauiltOptions interface{}, timeout time.Duration) *MultiSelectTask {
	t.WithTimeout(timeout, defauiltOptions)
	return t
}

// WithTimeout устанавливает тайм-аут для задачи множественного выбора
// @param duration Длительность тайм-аута
// @param defaultValue Значение по умолчанию (список индексов или строк)
// @return Указатель на задачу для цепочки вызовов
func (t *MultiSelectTask) WithTimeout(duration time.Duration, defaultValue interface{}) *MultiSelectTask {
	t.BaseTask.WithTimeout(duration, defaultValue)
	return t
}

// GetSelected возвращает список выбранных элементов
//
// @return []string Список выбранных пользователем вариантов
func (t *MultiSelectTask) GetSelected() []string {
	// Если задача завершена и есть финальное значение, разбираем его на список
	if t.done && t.finalValue != "" {
		// Разбиваем строку по запятой и пробелу
		parts := strings.Split(t.finalValue, defaults.DefaultSeparator)
		// Удаляем пустые элементы
		var result []string
		for _, part := range parts {
			if strings.TrimSpace(part) != "" {
				result = append(result, strings.TrimSpace(part))
			}
		}
		return result
	}

	// Если задача не завершена или нет финального значения, собираем выбранные элементы
	var selectedChoices []string
	for i := range t.choices {
		if t.isSelected(i) {
			selectedChoices = append(selectedChoices, t.choices[i])
		}
	}
	return selectedChoices
}
