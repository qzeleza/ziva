// +build ignore

package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== Тест выхода по Ctrl+C ===")
	fmt.Println("Нажмите Ctrl+C для выхода из задачи")
	fmt.Println()

	// Создаем простой список для тестирования
	choices := []string{
		"Элемент 1",
		"Элемент 2", 
		"Элемент 3",
		"Элемент 4",
		"Элемент 5",
	}

	singleTask := task.NewSingleSelectTask("Тест выхода по Ctrl+C", choices)

	// Запускаем задачу
	queue := query.New("Тест выхода по Ctrl+C")
	queue.AddTasks([]task.Task{singleTask})

	p := tea.NewProgram(queue)
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	selected := singleTask.GetSelected()
	fmt.Printf("Результат: %s\n", selected)
}
