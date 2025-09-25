package query

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/qzeleza/ziva/internal/common"
	"github.com/qzeleza/ziva/internal/defaults"
	"github.com/qzeleza/ziva/internal/performance"
	"github.com/qzeleza/ziva/internal/ui"
)

type ErrorColor int

const (
	Yellow ErrorColor = iota
	Red
	Orange
)

// Константы для управления памятью на embedded устройствах
var (
	// Значения могут быть переопределены через переменные окружения
	MaxCompletedTasks       int    = 50               // Максимальное количество завершенных задач в памяти
	MemoryPressureThreshold uint64 = 64 * 1024 * 1024 // 64MB - порог для запуска очистки памяти
)

func init() {
	// ZIVA_MAX_COMPLETED_TASKS=int
	if v := strings.TrimSpace(os.Getenv("ZIVA_MAX_COMPLETED_TASKS")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			MaxCompletedTasks = n
		}
	}

	// ZIVA_MEMORY_PRESSURE_THRESHOLD=64MB/65536KB/67108864/64MiB
	if v := strings.TrimSpace(os.Getenv("ZIVA_MEMORY_PRESSURE_THRESHOLD")); v != "" {
		if bytes, err := parseMemoryEnv(v); err == nil && bytes > 0 {
			MemoryPressureThreshold = bytes
		}
	} else if v := strings.TrimSpace(os.Getenv("GOMEMLIMIT")); v != "" {
		if bytes, err := parseMemoryEnv(v); err == nil && bytes > 0 {
			// По умолчанию порог в 0.8 от лимита памяти рантайма
			MemoryPressureThreshold = uint64(float64(bytes) * 0.8)
		}
	}
}

// parseMemoryEnv — локальный парсер для значений памяти с суффиксами
func parseMemoryEnv(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "_", "")
	up := strings.ToUpper(s)
	var mult uint64 = 1
	num := s
	switch {
	case strings.HasSuffix(up, "GIB"):
		mult = 1024 * 1024 * 1024
		num = s[:len(s)-3]
	case strings.HasSuffix(up, "MIB"):
		mult = 1024 * 1024
		num = s[:len(s)-3]
	case strings.HasSuffix(up, "KIB"):
		mult = 1024
		num = s[:len(s)-3]
	case strings.HasSuffix(up, "GB"):
		mult = 1024 * 1024 * 1024
		num = s[:len(s)-2]
	case strings.HasSuffix(up, "MB"):
		mult = 1024 * 1024
		num = s[:len(s)-2]
	case strings.HasSuffix(up, "KB"):
		mult = 1024
		num = s[:len(s)-2]
	case strings.HasSuffix(up, "B"):
		mult = 1
		num = s[:len(s)-1]
	}
	n, err := strconv.ParseUint(strings.TrimSpace(num), 10, 64)
	if err != nil {
		return 0, err
	}
	return n * mult, nil
}

// Используем константы из пакета common для расчета ширины макета

// Модель очереди задач.
// Модель представляет собой очередь задач, которые выполняются последовательно.
type Model struct {
	tasks      []common.Task  // Список задач.
	current    int            // Индекс текущей задачи.
	title      string         // Заголовок очереди.
	titleStyle lipgloss.Style // Стиль заголовка.
	summary    string         // Сводка по выполненным задачам.
	// summaryStyle   lipgloss.Style // Стиль сводки.
	width          int            // Ширина экрана.
	quitting       bool           // Флаг завершения работы.
	stoppedOnError bool           // Флаг прерывания очереди из-за ошибки.
	appName        string         // Название приложения.
	appNameStyle   lipgloss.Style // Стиль названия приложения.
	errorTask      common.Task    // Задача, вызвавшая прерывание очереди.
	clearScreen    bool           // Флаг очистки экрана перед запуском

	// Счетчики для подсчета результатов выполнения
	successCount int  // Количество успешно выполненных задач
	errorCount   int  // Количество задач с ошибками
	showSummary  bool // Флаг отображения сводки

	// Параметры отображения префиксов завершённых задач
	numberCompletedTasks bool   // Включает отображение номеров вместо символа завершения
	keepFirstSymbol      bool   // Если true, первая завершённая задача сохраняет символ
	numberFormat         string // Строка формата для отображения номера задачи

	// Параметры форматирования вывода результатов
	resultFormattingEnabled bool   // Включает форматирование результатов с разделительными линиями
	resultLinePrefix        string // Префикс для разделительной линии (по умолчанию "  │  ")
	resultLineLength        int    // Количество символов "─" в разделительной линии
}

const defauiltNumberFormat = "[%02d]" // формат по умолчанию для отображения номеров задач

// New создает новую модель очереди с заданным заголовком и задачами.
func New(title string) *Model {
	return &Model{
		title:        title,
		summary:      defaults.SummaryCompleted,
		width:        common.DefaultWidth, // Начальная ширина
		showSummary:  true,                // По умолчанию сводка отображается
		clearScreen:  false,               // По умолчанию экран не очищается
		appNameStyle: lipgloss.NewStyle().Foreground(ui.ColorDarkGray).Background(ui.ColorBrightWhite).Bold(false),
		numberFormat: defauiltNumberFormat,
		// Инициализация параметров форматирования результатов
		resultFormattingEnabled: true,                           // По умолчанию отключено
		resultLinePrefix:        "  │  ",                        // Префикс по умолчанию с символом │
		resultLineLength:        common.DefaultWidth * 93 / 100, // Длина линии перед выводом результатов задачи по умолчанию
	}
}

// Добавляет список задач для выполнения.
func (m *Model) AddTasks(tasks []common.Task) {
	m.tasks = append(m.tasks, tasks...)
}

// WithTasksNumbered включает нумерацию и задаёт формат представления числа в префиксе.
// Формат передаётся как строка для fmt.Sprintf (например, "[%02d]", "(%d)", "[0%d]").
// Если формат пустой, используется значение по умолчанию.
func (m *Model) WithTasksNumbered(enable bool, keepFirstSymbol bool, numberFormat string) *Model {
	m.numberCompletedTasks = enable
	ui.NumberingEnabled = enable
	m.keepFirstSymbol = keepFirstSymbol
	if strings.TrimSpace(numberFormat) == "" {
		m.numberFormat = defauiltNumberFormat
	} else {
		m.numberFormat = numberFormat
	}
	return m
}

// updateTaskStats обновляет статистику выполнения задач
func (m *Model) updateTaskStats() {
	m.successCount = 0
	m.errorCount = 0

	// Подсчитываем все задачи - просматриваем все до текущей позиции или все задачи если завершены
	tasksToCheck := m.current
	if m.current >= len(m.tasks) {
		tasksToCheck = len(m.tasks)
	}

	for i := 0; i < tasksToCheck; i++ {
		task := m.tasks[i]
		if task.IsDone() {
			if task.HasError() {
				m.errorCount++
			} else {
				m.successCount++
			}
		}
	}

	// Если очередь была остановлена из-за ошибки, и задача с ошибкой сохранена,
	// проверяем, учтена ли она в счетчике ошибок
	if m.stoppedOnError && m.errorTask != nil {
		// Проверяем, является ли задача с ошибкой текущей задачей
		if m.current < len(m.tasks) && m.tasks[m.current] == m.errorTask && m.errorTask.HasError() {
			// Если текущая задача с ошибкой не была учтена в цикле выше, учитываем её
			if m.current >= tasksToCheck {
				m.errorCount++
			}
		}
	}
}

// formatSummaryWithStats форматирует сводку с учетом статистики
func (m *Model) formatSummaryWithStats() (string, string) {
	totalTasks := len(m.tasks)
	completedTasks := m.successCount + m.errorCount

	// Формируем левую часть: summary + (успешных/всего)
	leftSummary := performance.FastConcat(
		m.summary,
		" ",
		performance.FastConcat(
			performance.IntToString(m.successCount),
			" ",
			defaults.DefaultFromSummaryLabel,
			" ",
			performance.IntToString(totalTasks),
		),
		" ",
		defaults.DefaultTasksSummaryLabel,
	)

	// Формируем правую часть: УСПЕШНО или С ОШИБКАМИ
	var rightStatus string
	if m.errorCount > 0 {
		rightStatus = defaults.StatusProblem
	} else if completedTasks == totalTasks && completedTasks > 0 {
		rightStatus = defaults.StatusSuccess
	} else {
		// Для состояния "В ПРОЦЕССЕ" или когда нет завершенных задач
		rightStatus = defaults.StatusInProgress
	}

	return leftSummary, rightStatus
}

// Запускает очередь задач
func (m *Model) Run() error {
	// Если установлен флаг очистки экрана, очищаем экран перед запуском
	if m.clearScreen {
		// Используем ANSI-последовательность для очистки экрана
		fmt.Print(defaults.ClearScreen)
	}

	_, err := tea.NewProgram(m).Run()
	return err
}

// WithTitleColor устанавливает цвет заголовка.
func (m *Model) WithTitleColor(titleColor lipgloss.TerminalColor, bold bool) *Model {
	m.titleStyle = lipgloss.NewStyle().Foreground(titleColor).Bold(bold)
	return m
}

// WithAppName устанавливает название приложения.
func (m *Model) WithAppName(appName string) *Model {
	m.appName = "  " + appName + "  "
	return m
}

// WithAppNameStyle устанавливает стиль названия приложения.
func (m *Model) WithAppNameColor(textColor lipgloss.TerminalColor, bold bool) *Model {
	m.appNameStyle = lipgloss.NewStyle().Foreground(textColor).Bold(bold).Background(ui.ColorBrightWhite)
	return m
}

// WithSummary устанавливает флаг отображения сводки.
func (m *Model) WithSummary(show bool) *Model {
	m.showSummary = show
	return m
}

// WithClearScreen устанавливает флаг очистки экрана перед запуском очереди задач.
func (m *Model) WithClearScreen(clear bool) *Model {
	m.clearScreen = clear
	return m
}

// WithResultFormatting включает форматирование результатов задач с разделительными линиями.
// Если enabled=true, то перед каждым результатом задачи будет добавляться разделительная линия
// из префикса и указанного количества символов "─".
// При enabled=false поведение остается как и раньше - результаты выводятся сразу после строки с задачей.
// @param enabled - включить/выключить форматирование
func (m *Model) WithResultFormatting(enabled bool) *Model {
	m.resultFormattingEnabled = enabled
	return m
}

// SetErrorColor устанавливает цвет для отображения ошибок в очереди.
func (m *Model) SetErrorColor(color ErrorColor) *Model {

	switch color {
	case Yellow:
		ui.SetErrorColor(ui.ColorDarkYellow, ui.ColorBrightYellow)
	case Red:
		ui.SetErrorColor(ui.ColorDarkRed, ui.ColorBrightRed)
	case Orange:
		ui.SetErrorColor(ui.ColorDarkOrange, ui.ColorBrightOrange)
	}
	return m
}

// layoutWidth вычисляет ширину для рендеринга задач.
// Использует функцию из пакета common.
func (m *Model) layoutWidth() int {
	return common.CalculateLayoutWidth(m.width)
}

// Init запускает первую задачу.
func (m *Model) Init() tea.Cmd {
	if len(m.tasks) > 0 {
		return m.tasks[0].Run()
	}
	return tea.Quit
}

// setTitle прорисовывает заголовок
func (m *Model) setTitle(width int) string {
	var result string
	if m.appName != "" {
		appName := m.appNameStyle.Render(m.appName)
		title := m.titleStyle.Render(m.title)
		result = ui.AlignTextToRight("  "+title, appName, width) + "\n"
	} else {
		result = "  " + m.title + "\n"
	}
	return result
}

// Update обрабатывает сообщения и делегирует их задачам.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Обработка обновлений размера окна.
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = size.Width
	}

	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "Ctrl+c" {
		m.quitting = true
		return m, tea.Quit
	}

	if m.current >= len(m.tasks) {
		return m, tea.Quit
	}

	currentTask := m.tasks[m.current]
	var cmd tea.Cmd

	// Обновляем текущую задачу и сохраняем её обратно
	updatedTask, cmd := currentTask.Update(msg)
	m.tasks[m.current] = updatedTask

	// Проверяем завершение уже ОБНОВЛЁННОЙ задачи
	if m.tasks[m.current].IsDone() {
		// Обновляем статистику выполненных задач
		m.updateTaskStats()

		// Проверяем давление памяти и выполняем очистку если необходимо
		m.checkMemoryPressure()

		// Проверяем, есть ли ошибка в текущей задаче
		if updatedTask.HasError() && updatedTask.StopOnError() {
			// Если есть ошибка и флаг StopOnError установлен,
			// прекращаем выполнение очереди и сохраняем информацию об ошибке
			m.stoppedOnError = true
			m.errorTask = updatedTask

			// Явно увеличиваем счетчик ошибок для текущей задачи
			m.errorCount++

			// Финальное обновление статистики перед остановкой из-за ошибки
			m.updateTaskStats()
			return m, tea.Quit
		}

		// Иначе переходим к следующей задаче
		m.current++
		if m.current < len(m.tasks) {
			nextCmd := m.tasks[m.current].Run()
			return m, tea.Batch(cmd, nextCmd)
		}
		// Финальное обновление статистики при завершении всех задач
		m.updateTaskStats()
		return m, tea.Quit
	}

	return m, cmd
}

// View отображает список задач.
func (m *Model) View() string {

	// Если очередь завершена, отображаем просто надпись о завершении
	// без прорисовки задач
	// if m.quitting {
	// 	return ui.CancelStyle.Render(task.DefaultCancelLabel) + "\n"
	// }

	var sb strings.Builder

	layoutWidth := m.layoutWidth()

	// Используем настроенный заголовок.
	sb.WriteString(ui.DrawLine(layoutWidth))
	sb.WriteString(m.setTitle(layoutWidth))
	sb.WriteString(ui.DrawLine(layoutWidth) + "\n")

	for i, t := range m.tasks {
		if i < m.current {
			hasError := t.HasError()
			m.applyCompletedTaskPrefix(t, i, hasError)
			// Завершенные задачи: отображаем их с форматированием (или без, если отключено)
			sb.WriteString(m.formatTaskResult(t, layoutWidth) + "\n")
		} else if i == m.current {
			hasError := t.HasError()
			if t.IsDone() {
				m.applyCompletedTaskPrefix(t, i, hasError)
				sb.WriteString(m.formatTaskResult(t, layoutWidth) + "\n")
			} else {
				m.applyInProgressTaskPrefix(t, i, hasError)
				// Активная задача: отображаем ее интерактивный вид.
				// Обрезаем только завершающие символы новой строки, сохраняя ведущие пробелы
				// view := strings.TrimRight(t.View(layoutWidth), "\n")
				view := t.View(layoutWidth)
				sb.WriteString(view + "\n")
			}
			// Добавляем разделитель, если есть ожидающие задачи и нет ошибки
			if i+1 < len(m.tasks) && (!m.stoppedOnError && !hasError) {
				sb.WriteString(ui.DrawLine(layoutWidth))
			}
		}
	}

	// Убираем крайнюю линию, если она есть
	removeDuplicateLines(&sb)

	// Добавляем финальную разделительную линию
	// Если есть активная задача, добавляем обычную линию
	// Иначе добавляем специальную линию
	if m.current < len(m.tasks) && !m.stoppedOnError {
		sb.WriteString(ui.DrawLine(layoutWidth))
	} else {

		// Отображаем сводку только если включен флаг showSummary
		if m.showSummary {
			// Получаем форматированную сводку с статистикой
			leftSummary, rightStatus := m.formatSummaryWithStats()

			// Определяем стиль для правой части в зависимости от статуса
			var rightStyle lipgloss.Style
			switch rightStatus {
			case defaults.StatusSuccess:
				rightStyle = ui.SuccessLabelStyle
			case defaults.StatusProblem:
				rightStyle = ui.GetErrorStatusStyle()
			case defaults.StatusInProgress:
				rightStyle = ui.SubtleStyle
			default:
				rightStyle = ui.SubtleStyle
			}

			// Определяем стиль для левой части (summary) - используем те же стили что и для правой части
			var summaryStyle lipgloss.Style
			switch rightStatus {
			case defaults.StatusSuccess:
				summaryStyle = ui.SuccessLabelStyle // Тот же стиль что и для правой части "SUCCESS"
			case defaults.StatusProblem:
				summaryStyle = ui.GetErrorStatusStyle() // Тот же стиль что и для правой части при ошибках
			case defaults.StatusInProgress:
				// Для состояния "В ПРОЦЕССЕ" используем тот же стиль что и справа
				summaryStyle = ui.SubtleStyle
			default:
				summaryStyle = ui.SubtleStyle
			}

			// Создаем левую часть футера
			summaryPrefix := performance.FastConcat(
				performance.RepeatEfficient(" ", ui.MainLeftIndent),
				ui.VerticalLineSymbol, "\n", "  ", // лишняя строка при выводе результатов выполнения задач
			)
			leftPart := performance.FastConcat(
				summaryPrefix,
				ui.FinishedLabelStyle.Render(ui.TaskCompletedSymbol),
				"  ",
				summaryStyle.Render(leftSummary),
			)

			// Создаем правую часть футера
			rightPart := rightStyle.Render(rightStatus)

			// Выравниваем по ширине макета
			footer := ui.AlignTextToRight(leftPart, rightPart, layoutWidth) + "\n\n"
			footer += ui.DrawLine(layoutWidth) + "\n"
			sb.WriteString(footer)
		} else {
			// Заменяем вертикальные линии перед символами задач ПЕРЕД добавлением финальных элементов
			removeVerticalLinesBeforeTaskSymbols(&sb)
			// Если сводка отключена, добавляем только финальную линию
			// sb.WriteString("\n")
			sb.WriteString(ui.DrawLine(layoutWidth) + "\n")
		}
	}

	return sb.String()
}

// applyCompletedTaskPrefix настраивает префикс завершённой задачи в зависимости от параметров модели
func (m *Model) applyCompletedTaskPrefix(task common.Task, index int, hasError bool) {
	setter, ok := task.(interface{ SetCompletedPrefix(string) })
	if !ok {
		return
	}

	if !m.numberCompletedTasks {
		setter.SetCompletedPrefix("")
		return
	}

	if m.keepFirstSymbol && index == 0 && !hasError {
		setter.SetCompletedPrefix("")
		return
	}

	number := index + 1
	if m.keepFirstSymbol && !hasError {
		number = index
		if number <= 0 {
			number = 1
		}
	}
	setter.SetCompletedPrefix(buildCompletedPrefix(number, m.numberFormat))
}

func (m *Model) applyInProgressTaskPrefix(task common.Task, index int, hasError bool) {
	setter, ok := task.(interface{ SetInProgressPrefix(string) })
	if !ok {
		return
	}

	if !m.numberCompletedTasks {
		setter.SetInProgressPrefix("")
		return
	}

	if m.keepFirstSymbol && index == 0 && !hasError {
		setter.SetInProgressPrefix("")
		return
	}

	number := index + 1
	if m.keepFirstSymbol && !hasError {
		number = index
		if number <= 0 {
			number = 1
		}
	}

	setter.SetInProgressPrefix(buildInProgressPrefix(number, m.numberFormat))
}

// calculateResultLinePrefix вычисляет префикс для разделительной линии и строк результата
// в зависимости от настроек нумерации
func (m *Model) calculateResultLinePrefix() string {
	if m.numberCompletedTasks {
		// При включенной нумерации нужно учесть ширину номера
		// Например, для "[1]  " нужно "     │  " (5 пробелов + │ + 2 пробела)
		// Вычисляем ширину номера для максимального номера задачи
		maxNumber := len(m.tasks)
		sampleNumber := fmt.Sprintf(m.numberFormat, maxNumber)
		numberWidth := len(sampleNumber)

		// Добавляем 2 пробела после номера (как в основном префиксе) + пробел перед │
		totalSpaces := numberWidth - 1 // +2 после номера, +1 перед │
		return performance.FastConcat(
			performance.RepeatEfficient(" ", totalSpaces),
			ui.VerticalLineSymbol,
			performance.RepeatEfficient(" ", 3), // 2 пробела после │
		)
	} else {
		// При отключенной нумерации используем стандартный префикс
		return m.resultLinePrefix
	}
}

// formatTaskResult форматирует результат задачи с разделительной линией
// Создает линию из префикса и символов "─", затем выводит результат задачи с новой строки
func (m *Model) formatTaskResult(task common.Task, width int) string {
	if !m.resultFormattingEnabled {
		// Если форматирование отключено, возвращаем обычное представление
		return task.FinalView(width)
	}

	var result strings.Builder

	// Получаем обычное представление задачи
	taskView := task.FinalView(width)

	// Разделяем на строки
	lines := strings.Split(taskView, "\n")
	if len(lines) == 0 {
		return taskView
	}

	// Первая строка - это заголовок задачи, добавляем её как есть
	result.WriteString(lines[0])

	// Если есть дополнительные строки (результаты), добавляем разделительную линию
	if len(lines) > 1 {
		result.WriteString("\n")

		// Определяем стиль линии на основе типа результата задачи
		var lineStyle lipgloss.Style
		if task.HasError() {
			// Для ошибок используем очень приглушенный желтый цвет (более приглушенный чем текст ошибки)
			lineStyle = ui.VerySubtleErrorStyle
		} else {
			// Для успешных результатов используем очень приглушенный стиль (едва заметный)
			lineStyle = ui.VerySubtleStyle
		}

		// Создаем стилизованную разделительную линию
		separatorContent := performance.RepeatEfficient(ui.HorizontalLineSymbol, m.resultLineLength)
		styledSeparator := lineStyle.Render(separatorContent)

		// Вычисляем динамический префикс для разделительной линии
		dynamicPrefix := m.calculateResultLinePrefix()
		separatorLine := performance.FastConcat(dynamicPrefix, styledSeparator)
		result.WriteString(separatorLine + "\n")

		// Добавляем остальные строки результата
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) != "" { // Пропускаем пустые строки
				// Строки результата уже содержат префикс, просто добавляем их
				result.WriteString(lines[i])
				if i < len(lines)-1 { // Добавляем перенос строки, кроме последней строки
					result.WriteString("\n")
				}
			}
		}
	}

	return result.String()
}

// buildCompletedPrefix возвращает префикс завершённой задачи, учитывая формат номера.
func buildCompletedPrefix(number int, format string) string {
	indent := ui.MainLeftIndent - 1
	if indent < 0 {
		indent = 0
	}
	return performance.FastConcat(
		performance.RepeatEfficient(" ", indent),
		formatTaskNumber(number, format),
	)
}

// buildInProgressPrefix возвращает префикс активной задачи с учётом формата номера.
func buildInProgressPrefix(number int, format string) string {
	indent := ui.MainLeftIndent - 1
	if indent < 0 {
		indent = 0
	}
	return performance.FastConcat(
		performance.RepeatEfficient(" ", indent),
		formatTaskNumber(number, format),
		" ",
	)
}

// formatTaskNumber форматирует номер задачи на основе заданного шаблона fmt.Sprintf.
// Если переданный формат пустой, используется формат по умолчанию.
func formatTaskNumber(number int, format string) string {
	if number <= 0 {
		number = 1
	}
	if strings.TrimSpace(format) == "" {
		format = defauiltNumberFormat
	}
	if !strings.Contains(format, "%") {
		format = format + "%d"
	}
	return fmt.Sprintf(format, number)
}

// Убираем из потока вывода крайнюю линию, если она есть
func removeDuplicateLines(sb *strings.Builder) {

	// Обработка финальной разделительной линии
	// Получаем текущее содержимое Builder для анализа
	content := sb.String()

	// Находим последний перенос строки для выделения последней строки
	lastNewlineIndex := strings.LastIndex(content, "\n")
	if lastNewlineIndex != -1 {
		// Проверяем, не является ли последняя строка пустой (только перенос строки)
		if lastNewlineIndex == len(content)-1 {
			// Если последняя строка пустая, ищем предыдущий перенос строки
			lastNewlineIndex = strings.LastIndex(content[:lastNewlineIndex], "\n")
		}

		if lastNewlineIndex != -1 && lastNewlineIndex < len(content)-1 {
			// Получаем последнюю строку (или предпоследнюю, если последняя пустая)
			lastLine := content[lastNewlineIndex+1:]
			// Убираем возможный перенос строки в конце
			lastLine = strings.TrimSuffix(lastLine, "\n")

			// Проверяем, является ли строка горизонтальной линией
			// (состоит только из символов HorizontalLineSymbol)
			isHorizontalLine := strings.Contains(lastLine, ui.HorizontalLineSymbol)

			// Если строка - горизонтальная линия и не пустая, удаляем её
			if isHorizontalLine && len(lastLine) > 0 {
				// Очищаем Builder и записываем контент без последней линии-разделителя
				sb.Reset()
				sb.WriteString(content[:lastNewlineIndex+1])
			}
		}
	}
}

// removeVerticalLinesBeforeTaskSymbols убирает вертикальные линии, ведущие к последнему (самому нижнему)
// символу задачи (TaskCompletedSymbol/TaskInProgressSymbol). Это позволяет «обрубить» вертикальные
// ветки, чтобы они не тянулись до конца списка задач.
//
// Алгоритм:
//  1. Ищем строку и колонку (в рунах) последнего символа задачи
//  2. Ищем самый нижний вертикальный сегмент в той же колонке и меняем его на пробел
func removeVerticalLinesBeforeTaskSymbols(sb *strings.Builder) {
	// Конвертируем буфер в строки
	content := sb.String()
	lines := strings.Split(content, "\n")

	// 1. Находим строку и колонку (в рунах) последнего символа задачи
	lastLine := -1
	col := -1
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		completed := strings.Index(line, ui.TaskCompletedSymbol)
		progress := strings.Index(line, ui.TaskInProgressSymbol)
		if completed == -1 && progress == -1 {
			continue
		}
		bytePos := completed
		if bytePos == -1 || (progress != -1 && progress < completed) {
			bytePos = progress
		}
		lastLine = i
		col = utf8.RuneCountInString(line[:bytePos])
		break
	}

	if lastLine == -1 {
		return // символы задач не найдены
	}

	// 2. Ищем все вертикальные сегменты в той же колонке после последней задачи и заменяем их на пробелы
	for i := lastLine + 1; i < len(lines); i++ {
		runes := []rune(lines[i])
		if col < len(runes) && string(runes[col]) == ui.VerticalLineSymbol {
			// Заменяем вертикальную линию на пробел
			runes[col] = ' '
			lines[i] = string(runes)
		}
	}

	// 3. Обновляем builder
	sb.Reset()
	sb.WriteString(strings.Join(lines, "\n"))
}

// DrawSpecialLine создает горизонтальную линию заданной ширины c угловой линией вверху
// типа ──┴─
func DrawFooterLine(width int) string {
	return performance.FastConcat(
		ui.HorizontalLineSymbol,
		ui.FinishedLabelStyle.Render(" "+ui.TaskCompletedSymbol+" "),
		performance.RepeatEfficient(ui.HorizontalLineSymbol, width-4), "\n")
}

// checkMemoryPressure проверяет использование памяти и выполняет очистку при необходимости
func (m *Model) checkMemoryPressure() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	// Если использование памяти превышает порог, выполняем экстренную очистку
	if ms.Sys > MemoryPressureThreshold {
		m.emergencyCleanup()
		// Принудительный сборщик мусора для embedded устройств
		runtime.GC()
	}
}

// emergencyCleanup выполняет экстренную очистку памяти
func (m *Model) emergencyCleanup() {
	// Ограничиваем количество завершенных задач
	m.cleanupOldTasks()

	// Очищаем буферные пулы через performance пакет
	performance.EmergencyPoolCleanup()

	// Очищаем кэш интернирования строк
	ui.ClearInternCache()
}

// cleanupOldTasks ограничивает количество завершенных задач в памяти
func (m *Model) cleanupOldTasks() {
	if m.current <= MaxCompletedTasks {
		return // Нет необходимости в очистке
	}

	// Сохраняем только последние MaxCompletedTasks завершенных задач
	// плюс все активные/незавершенные задачи
	keepFrom := m.current - MaxCompletedTasks
	if keepFrom < 0 {
		keepFrom = 0
	}

	// Создаем новый срез с ограниченным количеством задач
	newTasks := make([]common.Task, len(m.tasks)-keepFrom)
	copy(newTasks, m.tasks[keepFrom:])

	// Обновляем индекс текущей задачи
	m.current -= keepFrom
	m.tasks = newTasks
}
