//go:build embedded
// +build embedded

package ui

import "github.com/charmbracelet/lipgloss"

// Упрощенные ASCII символы для embedded устройств
var (
	IconDone      = "[OK]"
	IconError     = "[ERR]"
	IconCancelled = "[CANCEL]"
	IconSelected  = "X"
	IconRadioOn   = "X"
	IconCursor    = "->"
	IconUndone    = "[ ]"
	IconRadioOff  = "( )"
)

// ASCII символы для рамок
const (
	HorizontalLineSymbol = "-"
	VerticalLineSymbol   = "|"
	CornerDownSymbol     = "+"
	CornerUpSymbol       = "+"
	BranchSymbol         = "+"
	TaskCompletedSymbol  = "X"
	TaskInProgressSymbol = " "
	DownLineSymbol       = "+"
	UpLineSymbol         = "+"
	FinishSymbol         = "#"
)

// Упрощенные ANSI цвета
var (
	ColorBrightGreen   = lipgloss.Color("2") // Зеленый ANSI
	ColorBrightRed     = lipgloss.Color("1") // Красный ANSI
	ColorDarkRed       = lipgloss.Color("1") // Красный ANSI
	ColorBrightYellow  = lipgloss.Color("3") // Желтый ANSI
	ColorDarkYellow    = lipgloss.Color("3") // Желтый ANSI
	ColorBrightOrange  = lipgloss.Color("3") // Fallback на желтый
	ColorDarkOrange    = lipgloss.Color("3") // Fallback на желтый
	ColorBrightBlue    = lipgloss.Color("4") // Синий ANSI
	ColorDarkBlue      = lipgloss.Color("4") // Синий ANSI
	ColorBrightCyan    = lipgloss.Color("6") // Голубой ANSI
	ColorDarkCyan      = lipgloss.Color("6") // Голубой ANSI
	ColorBrightMagenta = lipgloss.Color("5") // Пурпурный ANSI
	ColorBrightWhite   = lipgloss.Color("7") // Белый ANSI
	ColorBrightGray    = lipgloss.Color("7") // Белый ANSI
	ColorDarkGray      = lipgloss.Color("8") // Темно-серый ANSI
	ColorLightBlue     = lipgloss.Color("6") // Fallback на голубой
	ColorBlack         = lipgloss.Color("0") // Черный ANSI
	ColorDarkGreen     = lipgloss.Color("2") // Зеленый ANSI
)

func init() {
	// Автоматически применяем embedded цвета при embedded сборке
	refreshIconsForEmbedded()
}
