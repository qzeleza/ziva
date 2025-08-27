// +build ignore

package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== Тест проблемы с выравниванием в MultiSelectTask ===")
	fmt.Println("Список из 8 элементов, viewport=3")
	fmt.Println("Проверяем выравнивание индикаторов '...N выше/ниже'")
	fmt.Println()

	// Создаем список с достаточным количеством элементов для тестирования viewport
	choices := []string{
		"Элемент 1",
		"Элемент 2", 
		"Элемент 3",
		"Элемент 4",
		"Элемент 5",
		"Элемент 6",
		"Элемент 7",
		"Элемент 8",
	}

	multiTask := task.NewMultiSelectTask("Тест выравнивания viewport", choices).
		WithViewport(3).
		WithSelectAll("Выбрать все")

	// Запускаем задачу
	queue := query.New("Тест выравнивания")
	queue.AddTasks([]task.Task{multiTask})

	p := tea.NewProgram(queue)
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	selected := multiTask.GetSelected()
	fmt.Printf("Результат: выбрано %d элементов\n", len(selected))
}
