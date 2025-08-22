// +build ignore

package main

import (
	"time"

	"github.com/qzeleza/termos/internal/common"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
)

func main() {
	// Создаем точно такую же задачу как в main.go
	ss := task.NewSingleSelectTask(
		"Выберите среду развертывания",
		[]string{"development", "staging", "production"},
	).WithTimeout(10*time.Second, "staging")

	// Создаем очередь задач
	var tasks []common.Task
	tasks = append(tasks, ss)

	// Запускаем через query систему
	queue := query.New("Тест SingleSelectTask с таймером").WithAppName("Тест").WithSummary(true)
	queue.AddTasks(tasks)
	queue.Run()
}
