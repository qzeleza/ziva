// task/base.go

package task

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/common"
	"github.com/qzeleza/termos/ui"
)

// Task - псевдоним интерфейса common.Task для обратной совместимости
type Task = common.Task

// BaseTask contains common fields for all tasks.
type BaseTask struct {
	title       string
	done        bool
	icon        string // Icon to show when done (e.g., check or cross)
	finalValue  string // The final value to display (e.g., "Yes", "Option 1")
	err         error  // Ошибка, если задача завершилась с ошибкой
	stopOnError bool   // Флаг, указывающий, нужно ли останавливать очередь при ошибке
}

func NewBaseTask(title string) BaseTask {
	return BaseTask{
		title:       title,
		stopOnError: true, // По умолчанию останавливаем очередь при ошибке
	}
}

func (t *BaseTask) Title() string                      { return t.title }
func (t *BaseTask) IsDone() bool                       { return t.done }
func (t *BaseTask) Run() tea.Cmd                       { return nil }
func (t *BaseTask) Update(msg tea.Msg) (Task, tea.Cmd) { return t, nil }

// HasError возвращает true, если при выполнении задачи произошла ошибка.
func (t *BaseTask) HasError() bool { return t.err != nil }

// Error возвращает ошибку, если она есть.
func (t *BaseTask) Error() error { return t.err }

// StopOnError возвращает true, если при возникновении ошибки в этой задаче
// нужно остановить выполнение всей очереди задач.
func (t *BaseTask) StopOnError() bool { return t.stopOnError }

// SetStopOnError устанавливает флаг остановки очереди при ошибке.
func (t *BaseTask) SetStopOnError(stop bool) { t.stopOnError = stop }

// SetError устанавливает ошибку для задачи
func (t *BaseTask) SetError(err error) { t.err = err }

// View provides a default implementation for active tasks.
func (t *BaseTask) View(_ int) string {
	// Most active tasks manage their own view, so this is a fallback.
	return t.title
}

// FinalView handles right-alignment for all tasks and formats error messages.
//
// @param width Ширина макета для выравнивания текста
// @return Отформатированное представление задачи с выравниванием
func (t *BaseTask) FinalView(width int) string {
	// Используем константы из пакета common для расчета оптимальной ширины
	// если переданная ширина меньше минимальной
	if width < common.DefaultWidth {
		width = common.DefaultWidth
	}

	// Определяем успешность выполнения задачи
	success := !t.HasError() && t.finalValue != "Отменено"

	// Определяем тип задачи для выбора правильного префикса
	isTextInputTask := IsTextInputTask(t)

	// Создаем префикс для завершенной задачи с новой системой отображения
	var prefix string
	if isTextInputTask {
		prefix = ui.GetCompletedInputTaskPrefix(success)
	} else {
		prefix = ui.GetCompletedTaskPrefix(success)
	}

	// Для простых значений Yes/No используем отдельные стили для "Да" и "Нет"
	if t.finalValue == "Да" || t.finalValue == "Нет" {
		left := fmt.Sprintf("%s %s", prefix, t.title)
		var right string
		if t.finalValue == "Да" {
			right = ui.SelectionStyle.Render(t.finalValue)
		} else {
			right = ui.SelectionNoStyle.Render(t.finalValue)
		}
		return ui.AlignTextToRight(left, right, width)
	}

	// Для ошибок выводим текст ошибки с отступом и слово "Ошибка" справа
	if t.icon == ui.IconError {
		// Создаем левую часть с заголовком и префиксом (prefix уже содержит ✕)
		left := fmt.Sprintf("%s %s", prefix, t.title)

		// Создаем правую часть со словом "Ошибка"
		right := ui.GetErrorStatusStyle().Render("ОШИБКА")

		// Создаем верхнюю строку с выравниванием
		result := ui.AlignTextToRight(left, right, width) + "\n"

		// Форматируем текст ошибки с отступом и переносами строк
		// Используем ширину на 4 символа меньше для отступа
		errText := ""
		// Получаем текст ошибки из finalValue, так как это уже отрендеренный текст
		if t.finalValue != "" {
			// Убираем стилизацию из текста ошибки
			errText = strings.ReplaceAll(t.finalValue, ui.IconError, "")
			errText = strings.TrimSpace(errText)
		}

		// Добавляем отформатированный текст ошибки
		result += ui.FormatErrorMessage(errText, common.CalculateLayoutWidth(width))

		return result
	}

	// Для обычных задач используем стандартное форматирование с новым префиксом
	if t.finalValue != "" && !strings.Contains(t.finalValue, t.title) {
		left := fmt.Sprintf("%s %s", prefix, t.title)
		// right := ui.SelectionStyle.Render(t.finalValue)
		ready := strings.ToUpper(DefaultSuccessLabel)
		right := ui.SelectionStyle.Render(ready)
		//
		return ui.AlignTextToRight(left, right, width)
	}

	// Если finalValue уже содержит полное форматирование, возвращаем как есть
	if t.finalValue != "" {
		return t.finalValue
	}

	// Запасной вариант - просто отображаем заголовок с префиксом
	return fmt.Sprintf("%s %s", prefix, t.title)
}

// IsTextInputTask определяет, является ли задача текстовой задачей ввода
// (не задачей выбора SingleSelect/MultiSelect)
func IsTextInputTask(task Task) bool {
	// Проверяем по названию типа через рефлексию
	switch task.(type) {
	case *SingleSelectTask, *MultiSelectTask:
		return false
	default:
		// Все остальные задачи (InputTaskNew, YesNoTask, FuncTask) являются текстовыми
		return true
	}
}