package task

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/internal/performance"
	"github.com/qzeleza/termos/internal/ui"
)

// YesNoOption представляет варианты выбора для YesNoTask
type YesNoOption int

const (
	YesOption YesNoOption = iota
	NoOption
)

// YesNoTask представляет задачу выбора из двух опций: Да, Нет
// Это обертка над SingleSelectTask для консистентности UI
type YesNoTask struct {
	*SingleSelectTask
	question       string
	yesLabel       string
	noLabel        string
	selectedOption YesNoOption
}

// NewYesNoTask создает новую задачу выбора с двумя опциями
func NewYesNoTask(title, question string) *YesNoTask {
	// Создаем варианты выбора
	options := []string{"Да", "Нет"}

	// Создаем базовую задачу выбора
	selectTask := NewSingleSelectTask(title, options)

	return &YesNoTask{
		SingleSelectTask: selectTask,
		question:         question,
		yesLabel:         "Да",
		noLabel:          "Нет",
		selectedOption:   YesOption,
	}
}

// Update обновляет состояние задачи, делегируя логику базовому SingleSelectTask
func (t *YesNoTask) Update(msg tea.Msg) (Task, tea.Cmd) {
	if t.IsDone() {
		return t, nil
	}

	// Рассматриваем специальные случаи для YesNoTask
	switch msg := msg.(type) {
	case TimeoutMsg:
		// Когда истекает тайм-аут, применяем значение по умолчанию
		t.applyDefaultValue()
		return t, nil
	case tea.KeyMsg:
		// Обрабатываем нажатие пробела для отключения таймера
		if msg.String() == " " && t.timeoutEnabled && t.timeoutManager != nil && t.timeoutManager.IsActive() {
			t.DisableTimeout()
			return t, nil
		}
	}

	// Делегируем обработку базовому SingleSelectTask
	updatedTask, cmd := t.SingleSelectTask.Update(msg)
	t.SingleSelectTask = updatedTask.(*SingleSelectTask)

	// Если задача завершена, определяем выбранную опцию
	if t.IsDone() {
		selectedIndex := t.SingleSelectTask.GetSelectedIndex()
		switch selectedIndex {
		case 0:
			t.selectedOption = YesOption
		case 1:
			t.selectedOption = NoOption
			// Выбор "Нет" считается ошибкой для статистики, но не останавливает очередь
			t.SetError(fmt.Errorf("пользователь выбрал \"Нет\""))
			t.SetStopOnError(false) // Не останавливаем очередь при выборе "Нет"
		}
	}

	return t, cmd
}

// Run запускает задачу выбора
func (t *YesNoTask) Run() tea.Cmd {
	// Запускаем таймер и тикер, если они включены
	if t.timeoutEnabled && t.timeoutManager != nil {
		return t.timeoutManager.StartTickerAndTimeout()
	}
	return nil
}

// View отображает текущее состояние задачи, делегируя логику базовому SingleSelectTask
// applyDefaultValue применяет значение по умолчанию при истечении таймера
func (t *YesNoTask) applyDefaultValue() {
	// Если есть значение по умолчанию
	if t.defaultValue != nil {
		switch val := t.defaultValue.(type) {
		case int:
			// Если это индекс (0 - Да, 1 - Нет)
			if val >= 0 && val < len(t.choices) {
				// Устанавливаем курсор на выбранный индекс
				t.cursor = val
				// Выбираем соответствующий вариант
				t.done = true
				t.icon = ui.IconDone
				t.finalValue = t.choices[t.cursor]
				
				// Устанавливаем выбранную опцию
				switch val {
				case 0: // Да
					t.selectedOption = YesOption
				case 1: // Нет
					t.selectedOption = NoOption
					// Выбор "Нет" считается ошибкой для статистики, но не останавливает очередь
					t.SetError(fmt.Errorf("пользователь выбрал \"Нет\""))
					t.SetStopOnError(false) // Не останавливаем очередь при выборе "Нет"
				}
			}
		case string:
			// Если это строка (Да, Нет)
			var selectedIndex int
			for i, choice := range t.choices {
				if strings.EqualFold(choice, val) { // Сравниваем без учета регистра
					// Устанавливаем курсор на выбранный вариант
					t.cursor = i
					selectedIndex = i
					// Выбираем вариант
					t.done = true
					t.icon = ui.IconDone
					t.finalValue = choice
					break
				}
			}
			
			// Устанавливаем выбранную опцию
			switch selectedIndex {
			case 0: // Да
				t.selectedOption = YesOption
			case 1: // Нет
				t.selectedOption = NoOption
				// Выбор "Нет" считается ошибкой для статистики, но не останавливает очередь
				t.SetError(fmt.Errorf("пользователь выбрал \"Нет\""))
				t.SetStopOnError(false) // Не останавливаем очередь при выборе "Нет"
			}
		}
	}
}

func (t *YesNoTask) View(width int) string {
	if t.IsDone() {
		return t.FinalView(width)
	}

	// Делегируем отображение базовому SingleSelectTask
	return t.SingleSelectTask.View(width)
}

// FinalView переопределяет отображение завершённой задачи.
// Если пользователь выбрал "Нет", с правой стороны выводится слово "ОТКАЗ"
// ярко-жёлтым цветом. Для ответа "Да" используется стандартное зелёное
// оформление выбранной опции.
func (t *YesNoTask) FinalView(width int) string {
	// Используем новый префикс завершённой задачи
	// Успешной считается только выбор "Да" и отсутствие других ошибок
	success := t.selectedOption == YesOption && !t.HasError()
	prefix := ui.GetCompletedTaskPrefix(success)

	left := fmt.Sprintf("%s %s", prefix, t.title)

	var right string
	if t.selectedOption == YesOption {
		right = ui.SelectionStyle.Render(DefaultYesLabel)
	} else if t.selectedOption == NoOption {
		// Для "Нет" выводим слово ОТКАЗ стилем ошибки
		right = ui.GetErrorStatusStyle().Render(DefaultNoLabel)
	}

	var result string
	// Если задача завершилась успешно и есть дополнительные строки для вывода
	if t.icon == ui.IconDone && len(t.choices) > 0 {
		result = "\n" + ui.DrawSummaryLine(t.choices[t.cursor]) +
			performance.RepeatEfficient(" ", ui.MainLeftIndent) + ui.VerticalLineSymbol
	}

	// Выравниваем по ширине макета
	return ui.AlignTextToRight(left, right, width) + result
}

// WithCustomLabels позволяет изменить текст опций (перегрузка для 2 параметров)
// Возвращает *YesNoTask для возможности цепочки вызовов
func (t *YesNoTask) WithCustomLabels(yesLabel, noLabel string) *YesNoTask {
	if strings.TrimSpace(yesLabel) != "" {
		t.yesLabel = yesLabel
		t.SingleSelectTask.choices[0] = yesLabel
	}
	if strings.TrimSpace(noLabel) != "" {
		t.noLabel = noLabel
		t.SingleSelectTask.choices[1] = noLabel
	}
	return t
}

// WithCustomLabelsAll позволяет изменить все три опции
// Возвращает *YesNoTask для возможности цепочки вызовов
func (t *YesNoTask) WithCustomLabelsAll(yesLabel, noLabel string) *YesNoTask {
	if strings.TrimSpace(yesLabel) != "" {
		t.yesLabel = yesLabel
		t.SingleSelectTask.choices[0] = yesLabel
	}
	if strings.TrimSpace(noLabel) != "" {
		t.noLabel = noLabel
		t.SingleSelectTask.choices[1] = noLabel
	}
	return t
}

// GetValue возвращает ответ пользователя (true для "Да", false для "Нет")
// @deprecated Рекомендуется использовать GetSelectedOption() для более ясной семантики
func (t *YesNoTask) GetValue() bool {
	switch t.selectedOption {
	case YesOption:
		return true
	case NoOption:
		return false
	default:
		return false
	}
}

// GetSelectedOption возвращает выбранную опцию
func (t *YesNoTask) GetSelectedOption() YesNoOption {
	return t.selectedOption
}

// IsYes возвращает true если выбрано "Да"
func (t *YesNoTask) IsYes() bool {
	return t.selectedOption == YesOption
}

// IsNo возвращает true если выбрано "Нет"
func (t *YesNoTask) IsNo() bool {
	return t.selectedOption == NoOption
}

// SetError устанавливает ошибку для задачи
func (t *YesNoTask) SetError(err error) {
	t.SingleSelectTask.SetError(err)
}

// HasError возвращает true, если при выполнении задачи произошла ошибка
func (t *YesNoTask) HasError() bool {
	return t.SingleSelectTask.HasError()
}

// Error возвращает ошибку, если она есть
func (t *YesNoTask) Error() error {
	return t.SingleSelectTask.Error()
}

// WithDefaultOption устанавливает вариант по умолчанию и тайм-аут для задачи
// @param defaultOption Может быть индексом (0 - Да, 1 - Нет) или строкой ("Да", "Нет")
// @param timeout Время ожидания в секундах до автовыбора
// @return Указатель на задачу для цепочки вызовов
func (t *YesNoTask) WithDefaultOption(defaultOption interface{}, timeout time.Duration) *YesNoTask {
	// Используем метод базового класса для установки тайм-аута
	t.SingleSelectTask.WithTimeout(timeout, defaultOption)
	return t
}
