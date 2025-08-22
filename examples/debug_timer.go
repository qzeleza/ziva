// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== Тест отображения таймера ===")
	
	// Создаем задачу с таймаутом
	singleTask := task.NewSingleSelectTask(
		"Выберите опцию",
		[]string{"Вариант 1", "Вариант 2", "Вариант 3"},
	)
	
	// Устанавливаем таймаут
	singleTask.WithTimeout(10*time.Second, 1)
	
	// Проверяем состояние таймера
	fmt.Printf("Таймер включен: %v\n", singleTask.BaseTask.GetRemainingTime() > 0)
	fmt.Printf("Оставшееся время: %s\n", singleTask.BaseTask.GetRemainingTimeFormatted())
	fmt.Printf("RenderTimer(): '%s'\n", singleTask.BaseTask.RenderTimer())
	
	// Запускаем команду Run для инициализации таймера
	cmd := singleTask.Run()
	fmt.Printf("Команда Run вернула: %v\n", cmd != nil)
	
	// Вызываем View для проверки отображения
	fmt.Println("\n=== Отображение View(80) ===")
	viewStr := singleTask.View(80)
	fmt.Println(viewStr)
	
	// Проверяем, что строки не пустые
	fmt.Printf("\nДлина View: %d символов\n", len(viewStr))
}
