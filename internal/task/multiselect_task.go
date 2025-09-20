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
	selected        SelectionBitset  // Битовый набор выбранных элементов (оптимизировано для embedded)
	fallbackMap     map[int]struct{} // Резервная карта для списков > 64 элементов
	cursor          int              // Текущая позиция курсора
	activeStyle     lipgloss.Style   // Стиль для активного элемента
	hasSelectAll    bool             // Включена ли опция "Выбрать все"
	selectAllText   string           // Текст опции "Выбрать все"
	showHelpMessage bool             // Показывать ли сообщение-подсказку
	helpMessage     string           // Текст сообщения-подсказки
	// Viewport (окно просмотра) для ограничения количества отображаемых элементов
	viewportSize  int // Размер viewport (количество видимых элементов), 0 = показать все
	viewportStart int // Начальная позиция viewport в списке элементов
}

// NewMultiSelectTask создает новую задачу множественного выбора.
//
// @param title Заголовок задачи
// @param choices Список вариантов выбора
// @return Указатель на новую задачу множественного выбора
func NewMultiSelectTask(title string, choices []string) *MultiSelectTask {
	task := &MultiSelectTask{
		BaseTask:        NewBaseTask(title),
		choices:         choices,
		cursor:          0, // Начинаем с первого элемента списка
		activeStyle:     ui.ActiveStyle,
		hasSelectAll:    false,
		selectAllText:   defauilt.SelectAllDefaultText,
		showHelpMessage: false,
		helpMessage:     "",
		// Viewport по умолчанию отключен (показываем все элементы)
		viewportSize:  0,
		viewportStart: 0,
	}

	// Для embedded устройств: используем битсет для списков <= 32 элементов,
	// иначе fallback на карту
	if len(choices) > 32 {
		task.fallbackMap = make(map[int]struct{})
	}

	return task
}

// WithViewport устанавливает размер viewport (окна просмотра) для ограничения количества отображаемых элементов.
// Это полезно для длинных списков, когда нужно показывать только часть элементов.
//
// @param size Количество элементов для отображения одновременно (0 = показать все)
// @return Указатель на задачу для цепочки вызовов
func (t *MultiSelectTask) WithViewport(size int) *MultiSelectTask {
	if size < 0 {
		size = 0
	}
	t.viewportSize = size
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
		t.selectAllText = defauilt.SelectAllDefaultText
	}
	return t
}

// WithDefaultItems позволяет заранее отметить элементы списка выбранными при открытии задачи.
// Поддерживает выбор одного индекса/строки или списков значений ([]int, []string).
func (t *MultiSelectTask) WithDefaultItems(defaultSelection interface{}) *MultiSelectTask {
	if defaultSelection == nil || len(t.choices) == 0 {
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
		if len(t.choices) > 32 && t.fallbackMap != nil {
			t.fallbackMap[index] = struct{}{}
		} else {
			t.selected.Set(index)
		}
		return true
	}

	anyApplied := false

	switch v := defaultSelection.(type) {
	case int:
		anyApplied = setSelected(v) || anyApplied
	case string:
		for i, choice := range t.choices {
			if choice == v {
				if setSelected(i) {
					anyApplied = true
				}
				break
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
			for i, choice := range t.choices {
				if choice == val {
					if setSelected(i) {
						anyApplied = true
					}
					break
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
	if len(t.choices) > 32 && t.fallbackMap != nil {
		// Используем fallback карту для больших списков
		allSelected := len(t.fallbackMap) == len(t.choices)
		if allSelected {
			t.fallbackMap = make(map[int]struct{})
		} else {
			for i := range t.choices {
				t.fallbackMap[i] = struct{}{}
			}
		}
	} else {
		// Используем битсет для оптимизации
		allSelected := t.selected.Count() == len(t.choices)
		if allSelected {
			t.selected.ClearAll()
		} else {
			t.selected.SetAll(len(t.choices))
		}
	}
}

// isAllSelected проверяет, выбраны ли все элементы списка
func (t *MultiSelectTask) isAllSelected() bool {
	if len(t.choices) == 0 {
		return false
	}

	if len(t.choices) > 32 && t.fallbackMap != nil {
		return len(t.fallbackMap) == len(t.choices)
	}

	return t.selected.Count() == len(t.choices)
}

// toggleSelection переключает состояние выбора элемента по индексу (оптимизировано для embedded)
func (t *MultiSelectTask) toggleSelection(index int) {
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
			if t.cursor > 0 {
				t.cursor--
			} else if t.hasSelectAll {
				// Если включена опция "Выбрать все" и курсор на первом элементе,
				// переходим к опции "Выбрать все" (позиция -1)
				t.cursor = -1
			}
			// Обновляем viewport после изменения позиции курсора
			t.updateViewport()
			return t, nil
		case "down", "j":
			t.stopTimeout()
			if t.hasSelectAll && t.cursor == -1 {
				// С опции "Выбрать все" переходим к первому элементу списка
				t.cursor = 0
			} else if t.cursor < len(t.choices)-1 {
				t.cursor++
			}
			// Обновляем viewport после изменения позиции курсора
			t.updateViewport()
			return t, nil
		case " ", "right", "Right":
			// При выборе останавливаем таймер
			t.stopTimeout()

			// В любом случае выполняем выбор/переключение
			if t.hasSelectAll && t.cursor == -1 {
				// Нажатие пробела на опции "Выбрать все"
				t.toggleSelectAll()
			} else if t.cursor >= 0 {
				// Обычная логика выбора для элементов списка
				t.toggleSelection(t.cursor)
			}
		case "q", "Q", "esc", "Esc", "ctrl+c", "Ctrl+C", "left", "Left":
			// Отмена пользователем
			cancelErr := fmt.Errorf(defauilt.ErrorMsgCanceled)
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
				t.helpMessage = "! Необходимо выбрать хотя бы один элемент"
				return t, nil
			}
			// Если есть выбранные элементы, завершаем задачу успешно
			t.done = true
			t.icon = ui.IconDone
			t.finalValue = strings.Join(selectedChoices, defauilt.DefaultSeparator)
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
				if index >= 0 && index < len(t.choices) {
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
				// Ищем строку в списке вариантов
				for i, choice := range t.choices {
					if choice == strVal {
						// Устанавливаем выбор для этого элемента
						if len(t.choices) > 32 && t.fallbackMap != nil {
							t.fallbackMap[i] = struct{}{}
						} else {
							t.selected.Set(i)
						}
						break
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
			t.finalValue = strings.Join(selectedChoices, defauilt.DefaultSeparator)
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
		// Не добавляем перенос строки в конце, чтобы не нарушать форматирование
		sb.WriteString(ui.SubtleStyle.Render(fmt.Sprintf(defauilt.ScrollAboveFormat, indentPrefix, ui.UpArrowSymbol, itemsAbove)))
		// Добавляем перенос строки отдельно
		sb.WriteString("\n")
	}

	// Отображаем только видимые элементы списка
	for i := startIdx; i < endIdx; i++ {
		if i >= len(t.choices) {
			break
		}

		choice := t.choices[i]
		checked := " "
		var itemPrefix string

		// Проверяем, выбран ли этот элемент
		if t.isSelected(i) {
			checked = ui.IconSelected
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
		sb.WriteString(fmt.Sprintf("%s[%s] %s\n", itemPrefix, checked, choice))
	}

	// Добавляем индикатор прокрутки вниз, если есть скрытые элементы ниже
	if t.viewportSize > 0 && endIdx < len(t.choices) {
		// Используем точно такой же префикс как у элементов "below"
		indentPrefix := ui.GetSelectItemPrefix("below")
		// Не добавляем перенос строки в конце, чтобы не нарушать форматирование
		sb.WriteString(ui.SubtleStyle.Render(fmt.Sprintf(defauilt.ScrollBelowFormat, indentPrefix, ui.DownArrowSymbol, len(t.choices)-endIdx)))
		// Добавляем перенос строки отдельно
		sb.WriteString("\n")
	}

	// Добавляем подсказку о навигации и управлении с новым отступом
	helpIndent := performance.RepeatEfficient(" ", ui.MainLeftIndent)

	// Добавляем сообщение-подсказку если нужно
	var warning string
	if t.showHelpMessage && t.helpMessage != "" {
		warning = ui.GetErrorMessageStyle().Render(fmt.Sprintf("%s%s", helpIndent, t.helpMessage))
		warning += "\n"
	}

	helpText := defauilt.MultiSelectHelp
	if t.hasSelectAll {
		helpText = defauilt.MultiSelectHelpSelectAll
	}
	sb.WriteString("\n" + ui.DrawLine(width) +
		ui.SubtleStyle.Render(fmt.Sprintf("%s%s%s", warning, helpIndent, helpText)))

	return sb.String()
}

func (t *MultiSelectTask) FinalView(width int) string {
	// Получаем базовое финальное представление
	result := t.BaseTask.FinalView(width) + "\n"

	// Если задача завершилась успешно и есть дополнительные строки для вывода
	if t.icon == ui.IconDone && len(t.choices) > 0 {
		for _, selected := range t.GetSelected() {
			result += ui.DrawSummaryLine(selected)
		}
		result += performance.RepeatEfficient(" ", ui.MainLeftIndent) + ui.VerticalLineSymbol
	}

	return result
}

// WithDefaultOptions устанавливает варианты по умолчанию, которые будут выбраны при тайм-ауте.
// @param defaultOptions Варианты по умолчанию (список индексов или строк)
// @param timeout Длительность тайм-аута
// @return Указатель на задачу для цепочки вызовов
func (t *MultiSelectTask) WithDefaultOptions(defaultOptions interface{}, timeout time.Duration) *MultiSelectTask {
	t.WithTimeout(timeout, defaultOptions)
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
		parts := strings.Split(t.finalValue, defauilt.DefaultSeparator)
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
