// +build ignore

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== Тест отображения таймера ===\n")

	// Тестируем InputTaskNew с активным таймером
	fmt.Println("InputTaskNew с таймером:")
	inputTask := task.NewInputTaskNew("Введите ваше имя", "Имя:")
	inputTask.WithTimeout(10*time.Second, "Пользователь по умолчанию")

	// Симулируем запущенный таймер путем получения строки таймера
	// (в реальности это делается внутренне)
	output := inputTask.View(80)
	fmt.Println(output)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ Проверьте, что:")
	fmt.Println("   - Отображается префикс '└─>'")
	fmt.Println("   - Есть поле для ввода '...'")
	fmt.Println("   - Показана справка внизу")
}
