// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/qzeleza/termos/internal/task"
)

func main() {
	// Создаем задачу
	ss := task.NewSingleSelectTask(
		"Выберите среду",
		[]string{"dev", "staging", "prod"},
	)
	
	fmt.Printf("До WithTimeout - IsDone(): %t\n", ss.IsDone())
	
	// Добавляем таймер
	ss = ss.WithTimeout(10*time.Second, "staging")
	
	fmt.Printf("После WithTimeout - IsDone(): %t\n", ss.IsDone())
	fmt.Printf("Выбранный индекс: %d\n", ss.GetSelectedIndex())
	fmt.Printf("Выбранное значение: '%s'\n", ss.GetSelected())
	
	// Проверяем view
	view := ss.View(80)
	fmt.Printf("Длина view: %d\n", len(view))
	fmt.Printf("View content: '%s'\n", view)
	
	// Попробуем с индексом вместо строки
	fmt.Println("\n--- Пробуем с индексом ---")
	ss2 := task.NewSingleSelectTask(
		"Выберите среду (индекс)",
		[]string{"dev", "staging", "prod"},
	).WithTimeout(10*time.Second, 1) // индекс 1 = "staging"
	
	fmt.Printf("С индексом - IsDone(): %t\n", ss2.IsDone())
	view2 := ss2.View(80)
	fmt.Printf("View с индексом: '%s'\n", view2)
}
