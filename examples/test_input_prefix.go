// +build ignore

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== Тест префикса для input задачи ===")

	// Создаем input задачу
	input := task.NewInputTaskNew("Введите ваше имя", "Имя:")
	input.WithTimeout(10*time.Second, "Пользователь")

	// Показываем view
	view := input.View(80)
	fmt.Printf("View content:\n%s\n", view)
	
	// Анализируем первую строку
	lines := strings.Split(view, "\n")
	if len(lines) > 0 {
		firstLine := lines[0]
		fmt.Printf("\nПервая строка: '%s'\n", firstLine)
		
		// Проверяем префиксы
		if strings.Contains(firstLine, "└─>") {
			fmt.Println("❌ Префикс НЕПРАВИЛЬНЫЙ: найден '└─>' (должен быть '○ ')")
		} else if strings.Contains(firstLine, "○") {
			fmt.Println("✅ Префикс ПРАВИЛЬНЫЙ: найден '○'")
		} else {
			fmt.Println("❓ Префикс неопределен")
		}
		
		// Проверяем таймер
		if strings.Contains(firstLine, "[") && strings.Contains(firstLine, "]") {
			fmt.Println("✅ Таймер присутствует")
		} else {
			fmt.Println("❌ Таймер отсутствует")
		}
	}
}
