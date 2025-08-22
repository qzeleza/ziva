// +build ignore

package main

import (
	"github.com/qzeleza/termos/internal/common"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
)

func main() {
	// Создаем задачу БЕЗ таймера
	ss := task.NewSingleSelectTask(
		"Выберите среду развертывания",
		[]string{"development", "staging", "production"},
	)

	// Создаем очередь задач
	var tasks []common.Task
	tasks = append(tasks, ss)

	// Запускаем через query систему
	queue := query.New("Тест SingleSelectTask БЕЗ таймера").WithAppName("Тест").WithSummary(true)
	queue.AddTasks(tasks)
	queue.Run()
}
