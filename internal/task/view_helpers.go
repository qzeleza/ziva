package task

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/qzeleza/ziva/internal/performance"
	"github.com/qzeleza/ziva/internal/ui"
)

// renderSelectionSeparator формирует разделитель между заголовком и списком пунктов
func renderSelectionSeparator(width int, enabled bool, inProgressPrefix string) string {
	if !enabled {
		return ""
	}
	if strings.TrimSpace(inProgressPrefix) == "" {
		inProgressPrefix = ui.GetCurrentTaskPrefix()
	}

	basePrefix := performance.FastConcat(
		performance.RepeatEfficient(" ", ui.MainLeftIndent),
		ui.VerticalLineSymbol,
	)

	targetWidth := lipgloss.Width(inProgressPrefix)
	baseWidth := lipgloss.Width(basePrefix)
	if targetWidth < baseWidth {
		targetWidth = baseWidth
	}

	extraSpaces := targetWidth - baseWidth
	if extraSpaces < 0 {
		extraSpaces = 0
	}

	prefix := performance.FastConcat(
		basePrefix,
		performance.RepeatEfficient(" ", extraSpaces),
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
