// +build ignore

package main

import (
	"time"

	"github.com/qzeleza/termos/internal/common"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
)

func main() {
	// Создаем простую задачу для тестирования
	single := task.NewSingleSelectTask("Тест выбора", []string{"A", "B", "C"})
	single.WithTimeout(10*time.Second, 1)

	input := task.NewInputTaskNew("Введите имя", "Имя:")
	input.WithTimeout(8*time.Second, "Пользователь")

	// Создаем очередь задач
	var tasks []common.Task
	tasks = append(tasks, single, input)

	// Запускаем через query систему
	queue := query.New("Диагностика задач").WithAppName("Тест").WithSummary(true)
	queue.AddTasks(tasks)
	queue.Run()
}
