// +build ignore

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== Отладка таймера для input задач ===")

	// Создаем input задачу с таймером
	input := task.NewInputTaskNew("Введите ваше имя", "Имя:")
	input.WithTimeout(10*time.Second, "Пользователь")

	fmt.Printf("1. После создания: IsDone = %t\n", input.IsDone())

	// Проверяем состояние таймера до запуска
	timer := input.BaseTask.RenderTimer()
	fmt.Printf("2. Таймер до Run(): '%s' (длина: %d)\n", timer, len(timer))

	// Запускаем задачу
	cmd := input.Run()
	fmt.Printf("3. Run() выполнен, есть команда: %t\n", cmd != nil)

	// Проверяем таймер сразу после запуска
	timer = input.BaseTask.RenderTimer()
	fmt.Printf("4. Таймер после Run(): '%s' (длина: %d)\n", timer, len(timer))

	// Ждем немного и проверяем снова
	time.Sleep(100 * time.Millisecond)
	timer = input.BaseTask.RenderTimer()
	fmt.Printf("5. Таймер через 100ms: '%s' (длина: %d)\n", timer, len(timer))

	// Проверяем полный View
	view := input.View(80)
	fmt.Printf("6. View content:\n%s\n", view)
	
	// Проверяем, есть ли таймер в view
	if len(view) > 0 {
		lines := strings.Split(view, "\n")
		for i, line := range lines {
			if strings.Contains(line, "[") && strings.Contains(line, "]") {
				fmt.Printf("7. ✅ Таймер найден в строке %d: '%s'\n", i+1, line)
			}
		}
	}
}
