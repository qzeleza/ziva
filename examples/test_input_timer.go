// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== Тест таймера для input задач ===")

	// Создаем input задачу с 10-секундным таймером
	input := task.NewInputTaskNew("Введите ваше имя", "Имя:")
	input.WithTimeout(10*time.Second, "Пользователь по умолчанию")

	fmt.Printf("Задача создана. IsDone(): %t\n", input.IsDone())
	
	// Запускаем задачу
	cmd := input.Run()
	fmt.Printf("Команда Run() выполнена: %v\n", cmd != nil)
	
	// Показываем начальное view
	view := input.View(80)
	fmt.Printf("Начальное view:\n%s\n", view)
	
	// Проверяем таймер каждую секунду
	for i := 0; i < 12; i++ {
		timer := input.BaseTask.RenderTimer()
		fmt.Printf("Секунда %d: Таймер = '%s', IsDone = %t\n", i, timer, input.IsDone())
		
		if input.IsDone() {
			fmt.Printf("Задача завершилась на %d секунде\n", i)
			fmt.Printf("Финальное значение: '%s'\n", input.GetValue())
			break
		}
		
		time.Sleep(1 * time.Second)
	}
}
