//go:build !embedded
// +build !embedded

package examples

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/common"
	"github.com/qzeleza/termos/query"
	"github.com/qzeleza/termos/ui"
)

// RunTasksWithTUI для полной сборки - полнофункциональный TUI
func RunTasksWithTUI(header string, summary string, tasks []common.Task) error {
	// Проверяем, не нужно ли переключиться в embedded режим автоматически
	if IsEmbeddedEnvironment() {
		// Применяем embedded оптимизации
		config := OptimizedEmbeddedConfig()
		ApplyEmbeddedConfig(config)
		ui.EnableEmbeddedMode()

		// Обнаружено embedded окружение, применены оптимизации
	}

	queue := query.New(header)
	queue.AddTasks(tasks)

	if _, err := tea.NewProgram(queue).Run(); err != nil {
		return err
	}
	return nil
}
