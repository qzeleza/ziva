//go:build embedded
// +build embedded

package termos

import (
	"fmt"

	"github.com/qzeleza//logger"
	"github.com/qzeleza/termos/common"
)

// RunTasksWithTUI для embedded сборки - упрощенный CLI интерфейс
func RunTasksWithTUI(log *logger.Logger, header string, summary string, tasks []common.Task) error {
	// Простой текстовый интерфейс без TUI для экономии ресурсов
	fmt.Printf("=== %s ===\n", header)

	for i, task := range tasks {
		fmt.Printf("[%d/%d] %s\n", i+1, len(tasks), task.Title())

		// Простая обработка задач без bubbletea
		if err := runSimpleTask(task); err != nil {
			log.Errorf("Ошибка выполнения задачи %s: %v", task.Title(), err)
			if task.StopOnError() {
				return err
			}
		}
	}

	fmt.Printf("=== %s ===\n", summary)
	return nil
}

// runSimpleTask простая реализация выполнения задачи для embedded
func runSimpleTask(task common.Task) error {
	// Здесь можно добавить упрощенную логику обработки задач
	// без использования TUI интерфейса
	return nil
}
