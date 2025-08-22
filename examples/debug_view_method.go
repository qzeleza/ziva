// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/qzeleza/termos/internal/task"
)

func main() {
	// Создаем простую задачу БЕЗ таймера для отладки
	ss := task.NewSingleSelectTask(
		"Выберите среду",
		[]string{"dev", "prod"},
	)

	fmt.Printf("Задача создана. IsDone(): %t\n", ss.IsDone())
	fmt.Printf("Количество вариантов: %d\n", len([]string{"dev", "prod"}))
	fmt.Printf("Курсор: %d\n", ss.GetSelectedIndex())
	
	// Получаем view напрямую
	view := ss.View(80)
	fmt.Printf("Длина view: %d\n", len(view))
	fmt.Printf("Содержимое view:\n'%s'\n", view)
	
	// Проверим, что будет, если добавить таймер
	fmt.Println("\n--- Добавляем таймер ---")
	ss2 := task.NewSingleSelectTask(
		"Выберите среду с таймером",
		[]string{"dev", "prod"},
	).WithTimeout(10*time.Second, "dev")

	view2 := ss2.View(80)
	fmt.Printf("View с таймером:\n'%s'\n", view2)
}
