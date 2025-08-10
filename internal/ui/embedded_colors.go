// Package ui - оптимизированная цветовая схема для embedded устройств
package ui

import "github.com/charmbracelet/lipgloss"

// EmbeddedColorPalette - упрощенная палитра для embedded устройств
// Использует только стандартные ANSI цвета для максимальной совместимости
var EmbeddedColorPalette = struct {
	// Основные цвета (ANSI 0-7)
	Black   lipgloss.TerminalColor
	Red     lipgloss.TerminalColor
	Green   lipgloss.TerminalColor
	Yellow  lipgloss.TerminalColor
	Blue    lipgloss.TerminalColor
	Magenta lipgloss.TerminalColor
	Cyan    lipgloss.TerminalColor
	White   lipgloss.TerminalColor

	// Яркие цвета (ANSI 8-15)
	BrightBlack   lipgloss.TerminalColor
	BrightRed     lipgloss.TerminalColor
	BrightGreen   lipgloss.TerminalColor
	BrightYellow  lipgloss.TerminalColor
	BrightBlue    lipgloss.TerminalColor
	BrightMagenta lipgloss.TerminalColor
	BrightCyan    lipgloss.TerminalColor
	BrightWhite   lipgloss.TerminalColor
}{
	// Используем строковые коды ANSI вместо hex для оптимизации
	Black:   lipgloss.Color("0"),
	Red:     lipgloss.Color("1"),
	Green:   lipgloss.Color("2"),
	Yellow:  lipgloss.Color("3"),
	Blue:    lipgloss.Color("4"),
	Magenta: lipgloss.Color("5"),
	Cyan:    lipgloss.Color("6"),
	White:   lipgloss.Color("7"),

	BrightBlack:   lipgloss.Color("8"),
	BrightRed:     lipgloss.Color("9"),
	BrightGreen:   lipgloss.Color("10"),
	BrightYellow:  lipgloss.Color("11"),
	BrightBlue:    lipgloss.Color("12"),
	BrightMagenta: lipgloss.Color("13"),
	BrightCyan:    lipgloss.Color("14"),
	BrightWhite:   lipgloss.Color("15"),
}

// EnableEmbeddedMode переключает модуль на использование embedded-совместимых цветов
func EnableEmbeddedMode() {
	// Заменяем тяжелые hex-цвета на легкие ANSI (с приведением типов)
	ColorBrightGreen = EmbeddedColorPalette.BrightGreen.(lipgloss.Color)
	ColorBrightRed = EmbeddedColorPalette.BrightRed.(lipgloss.Color)
	ColorDarkRed = EmbeddedColorPalette.Red.(lipgloss.Color)
	ColorBrightYellow = EmbeddedColorPalette.BrightYellow.(lipgloss.Color)
	ColorDarkYellow = EmbeddedColorPalette.Yellow.(lipgloss.Color)
	ColorBrightBlue = EmbeddedColorPalette.BrightBlue.(lipgloss.Color)
	ColorDarkBlue = EmbeddedColorPalette.Blue.(lipgloss.Color)
	ColorBrightCyan = EmbeddedColorPalette.BrightCyan.(lipgloss.Color)
	ColorDarkCyan = EmbeddedColorPalette.Cyan.(lipgloss.Color)
	ColorBrightMagenta = EmbeddedColorPalette.BrightMagenta.(lipgloss.Color)
	ColorBrightWhite = EmbeddedColorPalette.BrightWhite.(lipgloss.Color)
	ColorBrightGray = EmbeddedColorPalette.White.(lipgloss.Color)
	ColorDarkGray = EmbeddedColorPalette.BrightBlack.(lipgloss.Color)
	ColorBlack = EmbeddedColorPalette.Black.(lipgloss.Color)
	ColorDarkGreen = EmbeddedColorPalette.Green.(lipgloss.Color)

	// Упрощаем оранжевые цвета (не все терминалы поддерживают)
	ColorBrightOrange = EmbeddedColorPalette.BrightYellow.(lipgloss.Color) // Fallback на желтый
	ColorDarkOrange = EmbeddedColorPalette.Yellow.(lipgloss.Color)
	ColorLightBlue = EmbeddedColorPalette.BrightCyan.(lipgloss.Color) // Fallback на голубой

	// Пересоздаем иконки с новыми цветами
	refreshIconsForEmbedded()
}

// refreshIconsForEmbedded обновляет иконки для embedded режима
func refreshIconsForEmbedded() {
	IconDone = lipgloss.NewStyle().SetString("✔").Foreground(ColorBrightGreen).String()
	IconError = lipgloss.NewStyle().SetString("✕").Foreground(ColorBrightRed).String()
	IconCancelled = lipgloss.NewStyle().SetString("⊗").Foreground(ColorBrightYellow).String()
	IconQuestion = lipgloss.NewStyle().SetString("?").Foreground(ColorBrightGreen).String()
	IconSelected = lipgloss.NewStyle().SetString("■").Foreground(ColorBrightGreen).String()
	IconRadioOn = lipgloss.NewStyle().SetString("●").Foreground(ColorLightBlue).String()
	IconCursor = lipgloss.NewStyle().SetString("➞").Foreground(ColorLightBlue).String()
	IconUndone = lipgloss.NewStyle().SetString("◷").Foreground(ColorLightBlue).Bold(true).String()

	// Обновляем стили ошибок (с приведением типов)
	ErrorMessageStyle = ErrorMessageStyle.Foreground(EmbeddedColorPalette.Yellow)
	ErrorStatusStyle = ErrorStatusStyle.Foreground(EmbeddedColorPalette.BrightYellow)
}

// EnableASCIIMode включает максимально совместимый ASCII-набор иконок
func EnableASCIIMode() {
    // Устанавливаем простые ASCII-иконки без цвета
    IconDone = "*"
    IconError = "x"
    IconCancelled = "!"
    IconQuestion = "?"
    IconSelected = ">"
    IconRadioOn = "(x)"
    IconCursor = ">"
    IconUndone = "."
}

// IsEmbeddedColorMode возвращает true если включен embedded режим
func IsEmbeddedColorMode() bool {
	// Проверяем, используется ли ANSI код вместо hex
	return ColorBrightGreen == EmbeddedColorPalette.BrightGreen
}

// GetEmbeddedMemoryFootprint возвращает примерную оценку потребления памяти цветами
func GetEmbeddedMemoryFootprint() (embedded int, full int) {
	// ANSI коды: 1-2 байта на цвет x 16 цветов = ~32 байта
	embeddedFootprint := 32

	// Hex коды: ~8 байт на цвет x 18 цветов = ~144 байта
	fullFootprint := 144

	return embeddedFootprint, fullFootprint
}
