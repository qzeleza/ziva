// +build ignore

package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== Тест MultiSelectTask с опцией \"Выбрать все\" и viewport ===")
	fmt.Println("Проверяем отображение индикаторов прокрутки при наличии опции \"Выбрать все\"")
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
		"Элемент 9",
		"Элемент 10",
		"Элемент 11",
		"Элемент 12",
	}

	multiTask := task.NewMultiSelectTask("Тест индикаторов прокрутки с опцией \"Выбрать все\"", choices).
		WithViewport(4).
		WithSelectAll("Выбрать все")

	// Запускаем задачу
	queue := query.New("Тест индикаторов прокрутки")
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
