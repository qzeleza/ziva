// +build ignore

package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
)

func main() {
	// Создаем длинный список из 15 элементов для демонстрации viewport
	longList := []string{
		"Элемент 1: Первый пункт меню",
		"Элемент 2: Второй пункт меню", 
		"Элемент 3: Третий пункт меню",
		"Элемент 4: Четвертый пункт меню",
		"Элемент 5: Пятый пункт меню",
		"Элемент 6: Шестой пункт меню",
		"Элемент 7: Седьмой пункт меню",
		"Элемент 8: Восьмой пункт меню",
		"Элемент 9: Девятый пункт меню",
		"Элемент 10: Десятый пункт меню",
		"Элемент 11: Одиннадцатый пункт меню",
		"Элемент 12: Двенадцатый пункт меню",
		"Элемент 13: Тринадцатый пункт меню",
		"Элемент 14: Четырнадцатый пункт меню",
		"Элемент 15: Пятнадцатый пункт меню",
	}

	fmt.Println("=== Демонстрация Viewport для SingleSelectTask ===")
	fmt.Println("Список из 15 элементов, показываем только 5 одновременно")
	fmt.Println()

	// Создаем SingleSelectTask с viewport размером 5
	singleTask := task.NewSingleSelectTask("Выберите один элемент (viewport=5)", longList).
		WithViewport(5) // Показываем только 5 элементов одновременно

	// Запускаем задачу через query систему
	queue1 := query.New("Демонстрация SingleSelectTask с Viewport")
	queue1.AddTasks([]task.Task{singleTask})

	p1 := tea.NewProgram(queue1)
	_, err := p1.Run()
	if err != nil {
		log.Fatalf("Ошибка выполнения SingleSelectTask: %v", err)
	}

	fmt.Printf("Выбранный элемент: %s\n", singleTask.GetSelected())
	fmt.Println()

	fmt.Println("=== Демонстрация Viewport для MultiSelectTask ===")
	fmt.Println("Тот же список из 15 элементов, показываем только 5 одновременно")
	fmt.Println()

	// Создаем MultiSelectTask с viewport размером 5 и опцией "Выбрать все"
	multiTask := task.NewMultiSelectTask("Выберите несколько элементов (viewport=5)", longList).
		WithViewport(5).        // Показываем только 5 элементов одновременно
		WithSelectAll("Выбрать все элементы") // Добавляем опцию "Выбрать все"

	// Запускаем задачу через query систему
	queue2 := query.New("Демонстрация MultiSelectTask с Viewport")
	queue2.AddTasks([]task.Task{multiTask})

	p2 := tea.NewProgram(queue2)
	_, err = p2.Run()
	if err != nil {
		log.Fatalf("Ошибка выполнения MultiSelectTask: %v", err)
	}

	selected := multiTask.GetSelected()
	fmt.Printf("Выбрано элементов: %d\n", len(selected))
	for i, item := range selected {
		fmt.Printf("  %d. %s\n", i+1, item)
	}
	fmt.Println()

	fmt.Println("=== Демонстрация без Viewport (для сравнения) ===")
	fmt.Println("Короткий список без viewport")
	fmt.Println()

	shortList := []string{
		"Опция A",
		"Опция B", 
		"Опция C",
		"Опция D",
	}

	// Создаем задачу без viewport (по умолчанию показываются все элементы)
	normalTask := task.NewSingleSelectTask("Выберите опцию (без viewport)", shortList)

	// Запускаем задачу через query систему
	queue3 := query.New("Демонстрация без Viewport")
	queue3.AddTasks([]task.Task{normalTask})

	p3 := tea.NewProgram(queue3)
	_, err = p3.Run()
	if err != nil {
		log.Fatalf("Ошибка выполнения обычной задачи: %v", err)
	}

	fmt.Printf("Выбранная опция: %s\n", normalTask.GetSelected())
	fmt.Println()
	fmt.Println("Демонстрация завершена!")
}
