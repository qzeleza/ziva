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
		selectAllText:   "Выбрать все",
		showHelpMessage: false,
		helpMessage:     "",
	}

	// Для embedded устройств: используем битсет для списков <= 32 элементов,
	// иначе fallback на карту
	if len(choices) > 32 {
		task.fallbackMap = make(map[int]struct{})
	}

	return task
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
	}
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
		// При нажатии клавиш НЕ сбрасываем таймер - пусть продолжает работать
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
		case "down", "j":
			t.stopTimeout()
			if t.hasSelectAll && t.cursor == -1 {
				// С опции "Выбрать все" переходим к первому элементу списка
				t.cursor = 0
			} else if t.cursor < len(t.choices)-1 {
				t.cursor++
			}
		case " ":
			t.stopTimeout()

			// В любом случае выполняем выбор/переключение
			if t.hasSelectAll && t.cursor == -1 {
				// Нажатие пробела на опции "Выбрать все"
				t.toggleSelectAll()
			} else if t.cursor >= 0 {
				// Обычная логика выбора для элементов списка
				t.toggleSelection(t.cursor)
			}
		case "q", "Q":
			// Отмена пользователем
			cancelErr := fmt.Errorf("отменено пользователем")
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
			t.finalValue = strings.Join(selectedChoices, DefaultSeparator)
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
			t.finalValue = strings.Join(selectedChoices, DefaultSeparator)
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

	// Отображаем опцию "Выбрать все" если она включена
	if t.hasSelectAll {
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

	for i, choice := range t.choices {
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

	// Добавляем подсказку о навигации и управлении с новым отступом
	helpIndent := performance.RepeatEfficient(" ", ui.MainLeftIndent)

	// Добавляем сообщение-подсказку если нужно
	var warning string
	if t.showHelpMessage && t.helpMessage != "" {
		warning = ui.GetErrorMessageStyle().Render(fmt.Sprintf("%s%s", helpIndent, t.helpMessage))
		warning += "\n"
	}

	helpText := "[↑/↓ для навигации, пробел для выбора, Enter для подтверждения]"
	if t.hasSelectAll {
		helpText = "[↑/↓ навигация, пробел выбор/переключение всех, Enter подтверждение]"
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
		parts := strings.Split(t.finalValue, DefaultSeparator)
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
