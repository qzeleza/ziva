package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestEmbeddedColorPalette(t *testing.T) {
	// Проверяем, что все цвета в палитре определены
	assert.Equal(t, lipgloss.Color("0"), EmbeddedColorPalette.Black)
	assert.Equal(t, lipgloss.Color("1"), EmbeddedColorPalette.Red)
	assert.Equal(t, lipgloss.Color("2"), EmbeddedColorPalette.Green)
	assert.Equal(t, lipgloss.Color("3"), EmbeddedColorPalette.Yellow)
	assert.Equal(t, lipgloss.Color("4"), EmbeddedColorPalette.Blue)
	assert.Equal(t, lipgloss.Color("5"), EmbeddedColorPalette.Magenta)
	assert.Equal(t, lipgloss.Color("6"), EmbeddedColorPalette.Cyan)
	assert.Equal(t, lipgloss.Color("7"), EmbeddedColorPalette.White)

	assert.Equal(t, lipgloss.Color("8"), EmbeddedColorPalette.BrightBlack)
	assert.Equal(t, lipgloss.Color("9"), EmbeddedColorPalette.BrightRed)
	assert.Equal(t, lipgloss.Color("10"), EmbeddedColorPalette.BrightGreen)
	assert.Equal(t, lipgloss.Color("11"), EmbeddedColorPalette.BrightYellow)
	assert.Equal(t, lipgloss.Color("12"), EmbeddedColorPalette.BrightBlue)
	assert.Equal(t, lipgloss.Color("13"), EmbeddedColorPalette.BrightMagenta)
	assert.Equal(t, lipgloss.Color("14"), EmbeddedColorPalette.BrightCyan)
	assert.Equal(t, lipgloss.Color("15"), EmbeddedColorPalette.BrightWhite)
}

func TestEnableEmbeddedMode(t *testing.T) {
	// Сохраняем исходные значения цветов
	originalBrightGreen := ColorBrightGreen
	originalBrightRed := ColorBrightRed
	originalErrorMessageStyle := ErrorMessageStyle
	originalErrorStatusStyle := ErrorStatusStyle

	// Включаем embedded режим
	EnableEmbeddedMode()

	// Проверяем, что цвета изменились на ANSI
	assert.Equal(t, EmbeddedColorPalette.BrightGreen.(lipgloss.Color), ColorBrightGreen)
	assert.Equal(t, EmbeddedColorPalette.BrightRed.(lipgloss.Color), ColorBrightRed)
	assert.Equal(t, EmbeddedColorPalette.Red.(lipgloss.Color), ColorDarkRed)
	assert.Equal(t, EmbeddedColorPalette.BrightYellow.(lipgloss.Color), ColorBrightYellow)

	// Проверяем, что иконки обновились
	assert.Contains(t, IconDone, "✔")
	assert.Contains(t, IconError, "✕")
	assert.Contains(t, IconCancelled, "⊗")

	// Проверяем, что стили ошибок обновились
	assert.NotEqual(t, originalErrorMessageStyle.GetForeground(), ErrorMessageStyle.GetForeground())
	assert.NotEqual(t, originalErrorStatusStyle.GetForeground(), ErrorStatusStyle.GetForeground())

	// Восстанавливаем исходные значения
	ColorBrightGreen = originalBrightGreen
	ColorBrightRed = originalBrightRed
	ErrorMessageStyle = originalErrorMessageStyle
	ErrorStatusStyle = originalErrorStatusStyle
}

func TestRefreshIconsForEmbedded(t *testing.T) {
	// Сохраняем исходные иконки
	originalIconDone := IconDone
	originalIconError := IconError

	// Вызываем обновление иконок
	refreshIconsForEmbedded()

	// Проверяем, что иконки содержат ожидаемые символы
	assert.Contains(t, IconDone, "✔", "IconDone должен содержать символ галочки")
	assert.Contains(t, IconError, "✕", "IconError должен содержать символ крестика")
	assert.Contains(t, IconCancelled, "⊗", "IconCancelled должен содержать символ отмены")
	assert.Contains(t, IconQuestion, "?", "IconQuestion должен содержать символ вопроса")
	assert.Contains(t, IconSelected, "■", "IconSelected должен содержать символ выбора")
	assert.Contains(t, IconRadioOn, "●", "IconRadioOn должен содержать символ радио")
	assert.Contains(t, IconCursor, "➞", "IconCursor должен содержать символ курсора")
	assert.Contains(t, IconUndone, "◷", "IconUndone должен содержать символ незавершенной задачи")

	// Восстанавливаем исходные иконки
	IconDone = originalIconDone
	IconError = originalIconError
}

func TestIsEmbeddedColorMode(t *testing.T) {
	// Сохраняем исходное состояние
	originalColorBrightGreen := ColorBrightGreen

	// Тестируем до включения embedded режима
	ColorBrightGreen = lipgloss.Color("#00ff00") // Hex цвет
	assert.False(t, IsEmbeddedColorMode(), "Должен вернуть false для hex цветов")

	// Включаем embedded режим
	EnableEmbeddedMode()
	assert.True(t, IsEmbeddedColorMode(), "Должен вернуть true после включения embedded режима")

	// Восстанавливаем исходное состояние
	ColorBrightGreen = originalColorBrightGreen
}

func TestGetEmbeddedMemoryFootprint(t *testing.T) {
	embedded, full := GetEmbeddedMemoryFootprint()

	assert.Equal(t, 32, embedded, "Embedded режим должен использовать 32 байта")
	assert.Equal(t, 144, full, "Полный режим должен использовать 144 байта")
	assert.Less(t, embedded, full, "Embedded режим должен использовать меньше памяти")
}

func TestEmbeddedColorFallbacks(t *testing.T) {
	// Сохраняем исходные цвета
	originalBrightOrange := ColorBrightOrange
	originalDarkOrange := ColorDarkOrange
	originalLightBlue := ColorLightBlue

	EnableEmbeddedMode()

	// Проверяем, что оранжевые цвета fallback на желтые
	assert.Equal(t, EmbeddedColorPalette.BrightYellow.(lipgloss.Color), ColorBrightOrange)
	assert.Equal(t, EmbeddedColorPalette.Yellow.(lipgloss.Color), ColorDarkOrange)

	// Проверяем, что светло-синий fallback на ярко-голубой
	assert.Equal(t, EmbeddedColorPalette.BrightCyan.(lipgloss.Color), ColorLightBlue)

	// Восстанавливаем исходные цвета
	ColorBrightOrange = originalBrightOrange
	ColorDarkOrange = originalDarkOrange
	ColorLightBlue = originalLightBlue
}

func TestEmbeddedModeConsistency(t *testing.T) {
	// Включаем embedded режим
	EnableEmbeddedMode()

	// Проверяем, что все основные цвета соответствуют палитре
	colorMappings := map[lipgloss.Color]lipgloss.Color{
		ColorBrightGreen:   EmbeddedColorPalette.BrightGreen.(lipgloss.Color),
		ColorBrightRed:     EmbeddedColorPalette.BrightRed.(lipgloss.Color),
		ColorDarkRed:       EmbeddedColorPalette.Red.(lipgloss.Color),
		ColorBrightYellow:  EmbeddedColorPalette.BrightYellow.(lipgloss.Color),
		ColorDarkYellow:    EmbeddedColorPalette.Yellow.(lipgloss.Color),
		ColorBrightBlue:    EmbeddedColorPalette.BrightBlue.(lipgloss.Color),
		ColorDarkBlue:      EmbeddedColorPalette.Blue.(lipgloss.Color),
		ColorBrightCyan:    EmbeddedColorPalette.BrightCyan.(lipgloss.Color),
		ColorDarkCyan:      EmbeddedColorPalette.Cyan.(lipgloss.Color),
		ColorBrightMagenta: EmbeddedColorPalette.BrightMagenta.(lipgloss.Color),
		ColorBrightWhite:   EmbeddedColorPalette.BrightWhite.(lipgloss.Color),
		ColorBrightGray:    EmbeddedColorPalette.White.(lipgloss.Color),
		ColorDarkGray:      EmbeddedColorPalette.BrightBlack.(lipgloss.Color),
		ColorBlack:         EmbeddedColorPalette.Black.(lipgloss.Color),
		ColorDarkGreen:     EmbeddedColorPalette.Green.(lipgloss.Color),
	}

	for actual, expected := range colorMappings {
		assert.Equal(t, expected, actual, "Цвет должен соответствовать embedded палитре")
	}
}

func TestEmbeddedModeIdempotency(t *testing.T) {
	// Сохраняем исходное состояние
	originalBrightGreen := ColorBrightGreen

	// Включаем embedded режим дважды
	EnableEmbeddedMode()
	firstResult := ColorBrightGreen

	EnableEmbeddedMode()
	secondResult := ColorBrightGreen

	// Результат должен быть одинаковым
	assert.Equal(t, firstResult, secondResult, "Повторное включение embedded режима не должно изменять цвета")

	// Восстанавливаем исходное состояние
	ColorBrightGreen = originalBrightGreen
}

func TestEmbeddedColorValues(t *testing.T) {
	// Проверяем, что ANSI коды корректные
	ansiValues := []struct {
		color    lipgloss.Color
		expected string
	}{
		{EmbeddedColorPalette.Black.(lipgloss.Color), "0"},
		{EmbeddedColorPalette.Red.(lipgloss.Color), "1"},
		{EmbeddedColorPalette.Green.(lipgloss.Color), "2"},
		{EmbeddedColorPalette.Yellow.(lipgloss.Color), "3"},
		{EmbeddedColorPalette.Blue.(lipgloss.Color), "4"},
		{EmbeddedColorPalette.Magenta.(lipgloss.Color), "5"},
		{EmbeddedColorPalette.Cyan.(lipgloss.Color), "6"},
		{EmbeddedColorPalette.White.(lipgloss.Color), "7"},
		{EmbeddedColorPalette.BrightBlack.(lipgloss.Color), "8"},
		{EmbeddedColorPalette.BrightRed.(lipgloss.Color), "9"},
		{EmbeddedColorPalette.BrightGreen.(lipgloss.Color), "10"},
		{EmbeddedColorPalette.BrightYellow.(lipgloss.Color), "11"},
		{EmbeddedColorPalette.BrightBlue.(lipgloss.Color), "12"},
		{EmbeddedColorPalette.BrightMagenta.(lipgloss.Color), "13"},
		{EmbeddedColorPalette.BrightCyan.(lipgloss.Color), "14"},
		{EmbeddedColorPalette.BrightWhite.(lipgloss.Color), "15"},
	}

	for _, test := range ansiValues {
		assert.Equal(t, test.expected, string(test.color), "ANSI код должен соответствовать ожидаемому")
	}
}
