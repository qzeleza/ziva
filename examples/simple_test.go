// +build ignore

package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== Тест простой задачи с таймером ===")

	// Создаем ОДНУ простую задачу с таймером
	singleTask := task.NewSingleSelectTask(
		"Выберите вариант",
		[]string{"Вариант А", "Вариант Б", "Вариант В"},
	)
	singleTask.WithTimeout(30*time.Second, 0) // 30 секунд, по умолчанию первый

	// Создаем очередь с одной задачей
	queue := query.New("Простой тест")
	queue.AddTasks([]task.Task{singleTask})

	// Запускаем
	p := tea.NewProgram(queue)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	// Показываем результаты
	if _, ok := finalModel.(*query.Model); ok {
		fmt.Printf("\nВыбрано: %s\n", singleTask.GetSelected())
	}
}
