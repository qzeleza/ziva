package task

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/qzeleza/ziva/internal/performance"
	"github.com/qzeleza/ziva/internal/ui"
)

// renderSelectionSeparator формирует разделитель между заголовком и списком пунктов
func renderSelectionSeparator(width int, enabled bool) string {
	if !enabled {
		return ""
	}
	// Получаем отступ для результата, в зависимости от включенной нумерации задач
	additionalIndent := ui.GetResultIndentWhenNumberingEnabled()
	if len(additionalIndent) >= 2 {
		// Убираем последние два символа отступа
		additionalIndent = additionalIndent[:len(additionalIndent)-2]
	} else {
		additionalIndent = ""
	}

	// Формируем префикс для разделителя
	prefix := performance.FastConcat(
		performance.RepeatEfficient(" ", ui.MainLeftIndent),
		ui.VerticalLineSymbol,
		additionalIndent,
	)

	// Вычисляем доступную ширину для горизонтальной линии
	available := width - lipgloss.Width(prefix)
	if available > 0 {
		available--
	}
	if available < 0 {
		available = 0
	}

	// Формируем горизонтальную линию с бледным серым оттенком
	horizontal := ui.VerySubtleStyle.Render(performance.RepeatEfficient(ui.HorizontalLineSymbol, available))

	// Формируем разделитель
	return performance.FastConcat(
		prefix,
		horizontal,
		"\n",
	)
}
