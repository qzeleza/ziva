package task

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/internal/performance"
	"github.com/qzeleza/termos/internal/ui"
)

// YesNoExitOption представляет варианты выбора для YesNoTask
type YesNoExitOption int

const (
	YesOption YesNoExitOption = iota
	NoOption
	ExitOption
)

// YesNoTask представляет задачу выбора из трех опций: Да, Нет, Выйти
// Теперь это обертка над SingleSelectTask для консистентности UI
type YesNoTask struct {
	*SingleSelectTask
	question       string
	yesLabel       string
	noLabel        string
	selectedOption YesNoExitOption
}

// NewYesNoTask создает новую задачу выбора с тремя опциями
func NewYesNoTask(title, question string) *YesNoTask {
	// Создаем варианты выбора
	options := []string{"Да", "Нет", "Выйти"}

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
		case 2:
			t.selectedOption = ExitOption
			// Выбор "Выйти" останавливает выполнение очереди
			t.SetError(fmt.Errorf("пользователь выбрал \"Выйти\""))
			t.SetStopOnError(true) // Останавливаем очередь при выборе "Выйти"
		}
	}

	return t, cmd
}

// View отображает текущее состояние задачи, делегируя логику базовому SingleSelectTask
func (t *YesNoTask) View(width int) string {
	if t.IsDone() {
		return t.FinalView(width)
	}

	// return ""
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
	} else {
		// Для "Выйти" выводим слово ВЫХОД стилем ошибки
		right = ui.GetErrorStatusStyle().Render("ВЫХОД")
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

// GetValue возвращает ответ пользователя (true для "Да", false для "Нет", panic для "Выйти")
// @deprecated Рекомендуется использовать GetSelectedOption() для более ясной семантики
func (t *YesNoTask) GetValue() bool {
	switch t.selectedOption {
	case YesOption:
		return true
	case NoOption:
		return false
	case ExitOption:
		panic("GetValue() вызван для опции 'Выйти'")
	default:
		return false
	}
}

// GetSelectedOption возвращает выбранную опцию
func (t *YesNoTask) GetSelectedOption() YesNoExitOption {
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

// IsExit возвращает true если выбрано "Выйти"
func (t *YesNoTask) IsExit() bool {
	return t.selectedOption == ExitOption
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
