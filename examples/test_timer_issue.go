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
	fmt.Println("=== Тест проблемы с таймером в MultiSelectTask ===")
	fmt.Println("Таймер: 10 секунд, попробуйте навигацию ↑/↓")
	fmt.Println("Ожидается: таймер НЕ должен сбрасываться при навигации")
	fmt.Println()

	// Создаем MultiSelectTask с коротким таймером для тестирования
	choices := []string{
		"Элемент 1",
		"Элемент 2", 
		"Элемент 3",
		"Элемент 4",
		"Элемент 5",
		"Элемент 6",
	}

	multiTask := task.NewMultiSelectTask("Тест таймера (10 сек)", choices).
		WithViewport(3).
		WithSelectAll("Выбрать все").
		WithTimeout(10*time.Second, []string{"Элемент 1", "Элемент 2"})

	// Запускаем задачу
	queue := query.New("Тест таймера MultiSelectTask")
	queue.AddTasks([]task.Task{multiTask})

	p := tea.NewProgram(queue)
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	selected := multiTask.GetSelected()
	fmt.Printf("Результат: выбрано %d элементов\n", len(selected))
	for i, item := range selected {
		fmt.Printf("  %d. %s\n", i+1, item)
	}
}
