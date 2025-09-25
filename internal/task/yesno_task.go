package task

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/ziva/internal/defaults"
	"github.com/qzeleza/ziva/internal/performance"
	"github.com/qzeleza/ziva/internal/ui"
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
	question        string
	yesLabel        string
	noLabel         string
	selectedOption  YesNoOption
	showResultLine  bool
	noCountsAsError bool
}

// NewYesNoTask создает новую задачу выбора с двумя опциями
func NewYesNoTask(title, question string) *YesNoTask {
	// Создаем варианты выбора
	options := []string{defaults.DefaultYes, defaults.DefaultNo}

	// Создаем базовую задачу выбора
	selectTask := NewSingleSelectTask(title, options)

	return &YesNoTask{
		SingleSelectTask: selectTask,
		question:         question,
		yesLabel:         defaults.DefaultYes,
		noLabel:          defaults.DefaultNo,
		selectedOption:   YesOption,
		showResultLine:   true,
		noCountsAsError:  true,
	}
}

// stopTimeout останавливает таймер
func (t *YesNoTask) stopTimeout() {
	// Если таймер активен, останавливаем его
	if t.timeoutEnabled && t.timeoutManager != nil && t.timeoutManager.IsActive() {
		t.timeoutManager.StopTimeout()
		t.showTimeout = false
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
		// Обрабатываем локально для надежности (как в SingleSelectTask)
		// Значение по умолчанию уже нормализовано в WithTimeout
		t.applyDefaultValue()
		return t, nil
	case TickMsg:
		// Если таймер активен, продолжаем обновления
		if t.timeoutEnabled && t.timeoutManager != nil && t.timeoutManager.IsActive() {
			return t, t.timeoutManager.StartTicker()
		}
		return t, nil
	case tea.KeyMsg:
		// Обрабатываем нажатие упрвляющих клавиш для отключения таймера
		switch msg.String() {
		case " ", "up", "down", "j", "k", "enter":
			t.stopTimeout()
		case "q", "Q", "esc", "Esc", "ctrl+c", "Ctrl+C":
			// Отмена пользователем
			cancelErr := fmt.Errorf(defaults.ErrorMsgCanceled)
			t.done = true
			t.err = cancelErr
			t.icon = ui.IconCancelled
			t.finalValue = ui.CancelStyle.Render(cancelErr.Error())
			t.SetStopOnError(true)
			return t, nil
		}
	}

	// Делегируем обработку базовому SingleSelectTask
	updatedTask, cmd := t.SingleSelectTask.Update(msg)
	t.SingleSelectTask = updatedTask.(*SingleSelectTask)

	// Если задача завершена, определяем выбранную опцию
	if t.IsDone() {
		selectedIndex := t.GetSelectedIndex()
		switch selectedIndex {
		case 0:
			t.selectedOption = YesOption
		case 1:
			t.selectedOption = NoOption
			if t.noCountsAsError {
				// Выбор "Нет" считается ошибкой для статистики, но не останавливает очередь
				t.SetError(fmt.Errorf("%s \"%s\"", defaults.DefaultSelectedLabel, defaults.DefaultNo))
				t.SetStopOnError(false) // Не останавливаем очередь при выборе "Нет"
			}
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
					if t.noCountsAsError {
						// Выбор "Нет" считается ошибкой для статистики, но не останавливает очередь
						t.SetError(fmt.Errorf("%s \"%s\"", defaults.DefaultSelectedLabel, defaults.DefaultNo))
						t.SetStopOnError(false) // Не останавливаем очередь при выборе "Нет"
					}
				}
			}
		case string:
			// Если это строка (Да, Нет)
			selectedIndex := -1 // Инициализируем значением, указывающим что ничего не найдено
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

			// Устанавливаем выбранную опцию только если найдено соответствие
			if selectedIndex != -1 {
				switch selectedIndex {
				case 0: // Да
					t.selectedOption = YesOption
				case 1: // Нет
					t.selectedOption = NoOption
					if t.noCountsAsError {
						// Выбор "Нет" считается ошибкой для статистики, но не останавливает очередь
						t.SetError(fmt.Errorf("%s \"%s\"", defaults.DefaultSelectedLabel, defaults.DefaultNo))
						t.SetStopOnError(false) // Не останавливаем очередь при выборе "Нет"
					}
				}
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
	// Успешной считается выбор без ошибок и ("Да" или "Нет", если она не считается ошибкой)
	success := !t.HasError() && (t.selectedOption == YesOption || !t.noCountsAsError)
	prefix := t.CompletedPrefix()
	if prefix == "" {
		prefix = ui.GetCompletedTaskPrefix(success)
	}

	// Определяем стиль заголовка в зависимости от результата
	var styledTitle string
	if success {
		styledTitle = t.title
	} else {
		styledTitle = ui.GetErrorStatusStyle().Render(t.title)
	}

	// Сформируем левую часть строки
	left := fmt.Sprintf("%s  %s", prefix, styledTitle)

	// Сформируем правую часть строки
	var right string
	switch t.selectedOption {
	case YesOption:
		right = ui.TaskStatusSuccessStyle.Render(defaults.DefaultYesLabel)
	case NoOption:
		if t.noCountsAsError {
			// Для "Нет" выводим слово ОТКАЗ стилем ошибки
			right = ui.GetErrorStatusStyle().Render(defaults.DefaultNoLabel)
		} else {
			right = ui.TaskStatusSuccessStyle.Render(defaults.DefaultNoLabel)
		}
	}

	// Сформируем строку результата
	var result string
	// Если задача завершилась успешно и есть дополнительные строки для вывода
	if t.showResultLine && t.icon == ui.IconDone && len(t.choices) > 0 {
		result = "\n" + ui.DrawSummaryLine(t.choices[t.cursor]) +
			performance.RepeatEfficient(" ", ui.MainLeftIndent) + ui.VerticalLineSymbol
	}

	// Выравниваем по ширине макета
	return ui.AlignTextToRight(left, right, width) + result
}

// WithResultLine управляет отображением итоговой строки результата (совместимость)
func (t *YesNoTask) WithResultLine(show bool) *YesNoTask {
	t.showResultLine = show
	return t
}

/**
 * WithNoAsError включает интерпретацию ответа "Нет" как успешного результата.
 *
 * @return Указатель на задачу для цепочки вызовов.
 */
func (t *YesNoTask) WithNoAsError() *YesNoTask {
	t.noCountsAsError = false
	if t.selectedOption == NoOption && t.HasError() {
		t.SetError(nil)
	}
	return t
}

/**
 * ResultLineVisible сообщает, отображается ли итоговая строка результата.
 *
 * @return true, если строка результата включена.
 */
func (t *YesNoTask) ResultLineVisible() bool {
	return t.showResultLine
}

// WithCustomLabels позволяет изменить текст опций (перегрузка для 2 параметров)
// Возвращает *YesNoTask для возможности цепочки вызовов
func (t *YesNoTask) WithCustomLabels(yesLabel, noLabel string) *YesNoTask {
	if strings.TrimSpace(yesLabel) != "" {
		t.yesLabel = yesLabel
		t.choices[0] = yesLabel
	}
	if strings.TrimSpace(noLabel) != "" {
		t.noLabel = noLabel
		t.choices[1] = noLabel
	}
	return t
}

// WithDefaultItem позволяет задать опцию, которая будет выбрана по умолчанию при открытии задачи.
// Поддерживает YesNoOption, bool, индекс (int) и строку (string).
func (t *YesNoTask) WithDefaultItem(option interface{}) *YesNoTask {
	if option == nil {
		return t
	}

	selectByIndex := func(index int) {
		t.SingleSelectTask.WithDefaultItem(index)
		switch index {
		case 0:
			t.selectedOption = YesOption
		case 1:
			t.selectedOption = NoOption
		}
	}

	switch v := option.(type) {
	case YesNoOption:
		switch v {
		case YesOption:
			selectByIndex(0)
		case NoOption:
			selectByIndex(1)
		}
	case bool:
		if v {
			selectByIndex(0)
		} else {
			selectByIndex(1)
		}
	case int:
		selectByIndex(v)
	case string:
		for i, choice := range t.choices {
			if strings.EqualFold(choice, v) {
				selectByIndex(i)
				break
			}
		}
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

// WithTimeout устанавливает тайм-аут для задачи Да/Нет
// @param duration Длительность тайм-аута
// @param defaultValue Значение по умолчанию ("Да", "Нет" или индекс 0/1)
// @return Указатель на задачу для цепочки вызовов
func (t *YesNoTask) WithTimeout(duration time.Duration, defaultValue interface{}) *YesNoTask {
	// Нормализуем тип значения по умолчанию до ожидаемых BaseTask/SingleSelectTask типов (int или string)
	var normalized interface{} = nil
	if defaultValue != nil {
		switch v := defaultValue.(type) {
		case YesNoOption:
			if v == YesOption {
				normalized = 0 // индекс "Да"
			} else {
				normalized = 1 // индекс "Нет"
			}
		case bool:
			if v {
				normalized = 0 // Да
			} else {
				normalized = 1 // Нет
			}
		case int:
			normalized = v
		case string:
			// Нормализуем строковое значение к индексу для языко-независимого сравнения
			normalized = t.normalizeStringToIndex(v)
		default:
			// Неизвестный тип — не задаем значение по умолчанию
			normalized = nil
		}
	}

	// Проксируем в SingleSelectTask с нормализованным значением
	t.SingleSelectTask.WithTimeout(duration, normalized)
	return t
}

// normalizeStringToIndex конвертирует строковое значение по умолчанию в соответствующий индекс
// Поддерживает языко-независимое сравнение для "Да"/"Нет" на разных языках
func (t *YesNoTask) normalizeStringToIndex(value string) interface{} {
	// Список известных вариантов для "Да" на разных языках
	yesVariants := []string{"да", "yes", "evet", "так", "так"}
	// Список известных вариантов для "Нет" на разных языках
	noVariants := []string{"нет", "no", "hayır", "не", "ні"}

	valueLower := strings.ToLower(strings.TrimSpace(value))

	// Проверяем среди вариантов "Да"
	for _, variant := range yesVariants {
		if strings.EqualFold(valueLower, variant) {
			return 0 // индекс "Да"
		}
	}

	// Проверяем среди вариантов "Нет"
	for _, variant := range noVariants {
		if strings.EqualFold(valueLower, variant) {
			return 1 // индекс "Нет"
		}
	}

	// Если не найдено соответствие в известных вариантах,
	// пытаемся найти прямое соответствие в текущих choices
	for i, choice := range t.choices {
		if strings.EqualFold(choice, value) {
			return i
		}
	}

	// Если ничего не найдено, возвращаем исходную строку
	return value
}

// WithTimeoutYes устанавливает тайм-аут с выбором "Да" по умолчанию (языко-независимый)
// @param duration Длительность тайм-аута
// @return Указатель на задачу для цепочки вызовов
func (t *YesNoTask) WithTimeoutYes(duration time.Duration) *YesNoTask {
	t.WithTimeout(duration, YesOption)
	return t
}

// WithTimeoutNo устанавливает тайм-аут с выбором "Нет" по умолчанию (языко-независимый)
// @param duration Длительность тайм-аута
// @return Указатель на задачу для цепочки вызовов
func (t *YesNoTask) WithTimeoutNo(duration time.Duration) *YesNoTask {
	t.WithTimeout(duration, NoOption)
	return t
}

// WithDefaultYes устанавливает "Да" как вариант по умолчанию (языко-независимый)
// @return Указатель на задачу для цепочки вызовов
func (t *YesNoTask) WithDefaultYes() *YesNoTask {
	t.WithDefaultItem(YesOption)
	return t
}

// WithDefaultNo устанавливает "Нет" как вариант по умолчанию (языко-независимый)
// @return Указатель на задачу для цепочки вызовов
func (t *YesNoTask) WithDefaultNo() *YesNoTask {
	t.WithDefaultItem(NoOption)
	return t
}

// WithDefaultOption устанавливает вариант по умолчанию и тайм-аут для задачи
// @param defauiltOption Может быть индексом (0 - Да, 1 - Нет) или строкой ("Да", "Нет")
// @param timeout Время ожидания в секундах до автовыбора
// @return Указатель на задачу для цепочки вызовов
func (t *YesNoTask) WithDefaultOption(defauiltOption interface{}, timeout time.Duration) *YesNoTask {
	// Используем метод WithTimeout для единообразия
	t.WithTimeout(timeout, defauiltOption)
	return t
}
