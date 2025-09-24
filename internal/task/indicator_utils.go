package task

import (
	"strings"

	"github.com/qzeleza/termos/internal/ui"
)

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
