package ui

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/qzeleza/termos/performance"
)

// Яркая цветовая палитра для UI
var (
	ColorBrightGreen   = lipgloss.Color("#00ff00") // Ярко-зелёный
	ColorBrightRed     = lipgloss.Color("#FF2104") // Ярко-красный
	ColorDarkRed       = lipgloss.Color("#B10C01") // Темно-красный
	ColorBrightYellow  = lipgloss.Color("#ffff00") // Ярко-жёлтый
	ColorDarkYellow    = lipgloss.Color("#D2BE88") // Темно-жёлтый
	ColorBrightOrange  = lipgloss.Color("#FFA500") // Ярко-оранжевый
	ColorDarkOrange    = lipgloss.Color("#D2691E") // Темно-оранжевый
	ColorBrightBlue    = lipgloss.Color("#0000ff") // Ярко-синий
	ColorDarkBlue      = lipgloss.Color("#01189B") // Темно-синий
	ColorBrightCyan    = lipgloss.Color("#00ffff") // Ярко-голубой
	ColorDarkCyan      = lipgloss.Color("#006875") // Темно-голубой
	ColorBrightMagenta = lipgloss.Color("#ff00ff") // Ярко-фиолетовый
	ColorBrightWhite   = lipgloss.Color("#fff")    // Ярко-белый
	ColorBrightGray    = lipgloss.Color("#777")    // Светло-серый
	ColorDarkGray      = lipgloss.Color("#333")    // Темно-серый
	ColorLightBlue     = lipgloss.Color("#5DA9E9") // Светло-голубой
	ColorBlack         = lipgloss.Color("#000")    // Черный
	ColorDarkGreen     = lipgloss.Color("#008000") // Темно-зелёный
)

// Яркие иконки
var (
	IconDone      = lipgloss.NewStyle().SetString("✔").Foreground(ColorBrightGreen).String()          // Галочка (ярко-зелёная)
	IconError     = lipgloss.NewStyle().SetString("✕").Foreground(ColorBrightYellow).String()         // Крестик (ярко-красный)
	IconCancelled = lipgloss.NewStyle().SetString("⊗").Foreground(ColorBrightYellow).String()         // Отмена (ярко-жёлтый)
	IconQuestion  = lipgloss.NewStyle().SetString("?").Foreground(ColorBrightGreen).String()          // Вопрос (ярко-жёлтый)
	IconSelected  = lipgloss.NewStyle().SetString("■").Foreground(ColorBrightGreen).String()          // Галочка выбора (ярко-синий)
	IconRadioOn   = lipgloss.NewStyle().SetString("●").Foreground(ColorLightBlue).String()            // Радио включено (ярко-синий)
	IconCursor    = lipgloss.NewStyle().SetString("➞").Foreground(ColorLightBlue).String()            // Курсор (ярко-синий)
	IconUndone    = lipgloss.NewStyle().SetString("◷").Foreground(ColorLightBlue).Bold(true).String() // Неактивный элемент (ярко-белый)
	IconRadioOff  = "○"
)

// Стили для текста
// Стили для текста с яркой палитрой
var (
	TitleStyle = lipgloss.NewStyle().Bold(true)

	// Динамические стили ошибок (могут быть изменены через SetErrorColor)
	ErrorMessageStyle = lipgloss.NewStyle().Foreground(ColorDarkYellow)              // Ошибка: ярко-красный (без курсива)
	ErrorStatusStyle  = lipgloss.NewStyle().Foreground(ColorBrightYellow).Bold(true) // Вывод сообщений об ошибках: ярко-красный (курсив)
	CancelStyle       = lipgloss.NewStyle().Foreground(ColorBrightYellow)            // Вывод статуса ошибки: ярко-жёлтый
	SubtleStyle       = lipgloss.NewStyle().Foreground(ColorBrightGray)              // Подписи: светло-серый
	SelectionStyle    = lipgloss.NewStyle().Foreground(ColorBrightGreen)             // Выделение (Да): ярко-зелёный
	SelectionNoStyle  = lipgloss.NewStyle().Foreground(ColorBrightRed).Bold(true)    // Выделение (Нет): ярко-красный
	ActiveStyle       = lipgloss.NewStyle().Foreground(ColorLightBlue).Bold(true)    // Активный элемент: ярко-синий
	InputStyle        = lipgloss.NewStyle().Foreground(ColorLightBlue).Bold(true)    // Стиль для активного ввода
	SpinnerStyle      = lipgloss.NewStyle().Foreground(ColorLightBlue).Bold(true)    // Стиль для спиннера
	ActiveTitleStyle  = lipgloss.NewStyle().Foreground(ColorBrightGreen).Bold(true)  // Активный заголовок ввода
	ActiveTaskStyle   = lipgloss.NewStyle().Foreground(ColorBrightGreen)             // Стиль для активной задачи
	SuccessLabelStyle = lipgloss.NewStyle().Foreground(ColorBrightGreen).Bold(true)  // Успешное завершение

	FinishedLabelStyle = lipgloss.NewStyle().Foreground(ColorBrightWhite).Bold(true) // Завершение
)

// Константы отступов (в пробелах)
const (
	// Основная система отступов с вертикальными линиями
	MainLeftIndent      = 2 // Основной отступ от левого края для всех задач
	MessageIndentSpaces = 3 // Отступ для вывода ошибок и сообщений при завершении задачи
)

// Строковые отступы
var (
	MessageIndent = strings.Repeat(" ", MessageIndentSpaces) // Отступ для сообщений об ошибках
)

// Символы для новой системы отображения
const (
	HorizontalLineSymbol = "─" // Символ горизонтальной линии
	VerticalLineSymbol   = "│" // Символ вертикальной линии
	CornerDownSymbol     = "└" // Угловой символ вправо (устарело)
	CornerUpSymbol       = "┌" // Угловой символ вправо (устарело)
	ArrowSymbol          = ">" // Стрелка вправо (устарело)
	BranchSymbol         = "├" // Символ ветки для активных задач
	TaskCompletedSymbol  = "●" // Символ включенного радио
	TaskInProgressSymbol = "○" // Символ выключенного радио
	DownLineSymbol       = "┬" // Символ специальной линии
	UpLineSymbol         = "┴" // Символ специальной линии
	FinishSymbol         = "■" // Префикс для задач
)

// Функции для создания префиксов с новой системой отображения
// GetTaskBelowPrefix возвращает префикс для задачи ниже текущей выполняющейся задачи
// Формат: "   │ " (отступ + ветка + линия + пробел)
func GetTaskBelowPrefix() string {
	return strings.Repeat(" ", MainLeftIndent) + VerticalLineSymbol + " "
}

// GetActiveTaskPrefix возвращает префикс для активной (выполняющейся) задачи
// Формат: "   ├ " (отступ + разветвление + горизонтальная линия + пробел)
func GetActiveTaskPrefix() string {
	return strings.Repeat(" ", MainLeftIndent) + BranchSymbol + " "
}

// GetCompletedTaskPrefix возвращает префикс для завершенной задачи (успешно или с ошибкой)
// Формат: "   └ " для задач выбора
func GetCompletedTaskPrefix(success bool) string {
	prefix := strings.Repeat(" ", MainLeftIndent) + CornerDownSymbol + " "
	if success {
		return prefix
	}
	return prefix
}

// GetCompletedInputTaskPrefix возвращает префикс для завершенной задачи ввода текста
// Формат: "   ■ " для текстовых задач
func GetCompletedInputTaskPrefix(success bool) string {
	symbol := FinishSymbol
	if success {
		return strings.Repeat(" ", MainLeftIndent) + ActiveTitleStyle.Render(symbol) + " "
	}
	return strings.Repeat(" ", MainLeftIndent) + ErrorStatusStyle.Render(symbol) + " "
}

// GetErrorMessageStyle возвращает текущий стиль для сообщений об ошибках
func GetErrorMessageStyle() lipgloss.Style {
	return ErrorMessageStyle
}

// GetErrorStatusStyle возвращает текущий стиль для статуса ошибки
func GetErrorStatusStyle() lipgloss.Style {
	return ErrorStatusStyle
}

// SetErrorColor устанавливает пользовательский цвет для ошибок
func SetErrorColor(messageColor, statusColor lipgloss.Color) {
	ErrorMessageStyle = ErrorMessageStyle.Foreground(messageColor)
	ErrorStatusStyle = ErrorStatusStyle.Foreground(statusColor)
}

// ResetErrorColors сбрасывает цвета ошибок к значениям по умолчанию
func ResetErrorColors() {
	ErrorMessageStyle = lipgloss.NewStyle().Foreground(ColorDarkYellow)
	ErrorStatusStyle = lipgloss.NewStyle().Foreground(ColorBrightYellow).Bold(true)
}

// FormatErrorMessage форматирует сообщение об ошибке с отступом и переносом строк
// Используется в FinalView для отображения подробной информации об ошибке
func FormatErrorMessage(errorText string, width int) string {
	if errorText == "" {
		return ""
	}

	// Убираем лишние пробелы в начале и конце
	errorText = strings.TrimSpace(errorText)

	// Применяем стиль к тексту ошибки
	styledError := GetErrorMessageStyle().Render(errorText)

	// Вычисляем доступную ширину для текста с учетом отступа
	availableWidth := width - MessageIndentSpaces
	if availableWidth < 10 { // Минимальная ширина
		availableWidth = 10
	}

	// Разбиваем текст на строки с учетом доступной ширины
	lines := WrapText(styledError, availableWidth)

	// Добавляем отступ к каждой строке
	var result strings.Builder
	for _, line := range lines {
		result.WriteString(MessageIndent)
		result.WriteString(line)
		result.WriteString("\n")
	}

	// Убираем последний символ перевода строки
	return strings.TrimSuffix(result.String(), "\n")
}

// AlignTextToRight выравнивает текст по правому краю с учетом ширины
// Используется для отображения результатов задач справа от заголовка
func AlignTextToRight(left, right string, width int) string {
	// Очищаем ANSI escape последовательности для подсчета реальной длины
	leftLen := GetPlainTextLength(left)
	rightLen := GetPlainTextLength(right)

	// Вычисляем количество пробелов для выравнивания
	spaces := width - leftLen - rightLen
	if spaces < 1 {
		spaces = 1 // Минимум один пробел
	}

	return left + strings.Repeat(" ", spaces) + right
}

// GetPlainTextLength возвращает длину текста без ANSI escape последовательностей
func GetPlainTextLength(text string) int {
	// Используем оптимизированную функцию из пакета performance
	return performance.StripANSILength(text)
}

// WrapText разбивает текст на строки с учетом заданной ширины
// Поддерживает Unicode символы и ANSI escape последовательности
func WrapText(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	// Удаляем ANSI коды для корректного подсчета длины
	plainText := StripANSI(text)

	var lines []string
	var currentLine strings.Builder
	var currentWidth int

	for _, r := range plainText {
		runeWidth := GetRuneWidth(r)

		// Если добавление символа превысит ширину, переходим на новую строку
		if currentWidth+runeWidth > width && currentLine.Len() > 0 {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentWidth = 0
		}

		// Обработка символов перевода строки
		if r == '\n' {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentWidth = 0
			continue
		}

		currentLine.WriteRune(r)
		currentWidth += runeWidth
	}

	// Добавляем последнюю строку, если она не пустая
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// StripANSI удаляет ANSI escape последовательности из текста
func StripANSI(text string) string {
	// Регулярное выражение для поиска ANSI escape последовательностей
	re := regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)
	return re.ReplaceAllString(text, "")
}

// GetRuneWidth возвращает ширину отображения Unicode символа
func GetRuneWidth(r rune) int {
	// Для большинства символов ширина равна 1
	// Для широких символов (например, китайские иероглифы) ширина равна 2
	// Для управляющих символов ширина равна 0

	if unicode.IsControl(r) {
		return 0
	}

	// Простая проверка на широкие символы
	// В реальном приложении лучше использовать библиотеку runewidth
	if r >= 0x1100 && (r <= 0x115F || // Hangul Jamo
		(r >= 0x2E80 && r <= 0x9FFF) || // CJK
		(r >= 0xAC00 && r <= 0xD7A3) || // Hangul Syllables
		(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility Ideographs
		(r >= 0xFE10 && r <= 0xFE19) || // Vertical forms
		(r >= 0xFE30 && r <= 0xFE6F) || // CJK Compatibility Forms
		(r >= 0xFF00 && r <= 0xFF60) || // Fullwidth Forms
		(r >= 0xFFE0 && r <= 0xFFE6) || // Fullwidth Forms
		(r >= 0x20000 && r <= 0x2FFFD) || // CJK Extension B
		(r >= 0x30000 && r <= 0x3FFFD)) { // CJK Extension C
		return 2
	}

	return utf8.RuneLen(r) // Возвращаем количество байт в UTF-8 представлении
}