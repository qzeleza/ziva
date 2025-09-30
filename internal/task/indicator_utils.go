package task

import (
	"strings"

	"github.com/qzeleza/ziva/internal/ui"
)

// appendIndicatorWithPlainPipe добавляет указатель к индикатору
// Если указатель найден, то он будет добавлен в индикатор
//
// @param sb - билдер строки
// @param indicator - индикатор
func appendIndicatorWithPlainPipe(sb *strings.Builder, indicator string) {
	pipe := ui.VerticalLineSymbol
	idx := strings.Index(indicator, pipe)
	if idx == -1 {
		sb.WriteString(ui.SubtleStyle.Render(indicator))
		return
	}

	prefix := indicator[:idx]
	suffix := indicator[idx+len(pipe):]

	if prefix != "" {
		sb.WriteString(ui.SubtleStyle.Render(prefix))
	}

	sb.WriteString(pipe)

	if suffix != "" {
		sb.WriteString(ui.SubtleStyle.Render(suffix))
	}
}
