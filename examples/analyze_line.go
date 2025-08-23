//go:build ignore
// +build ignore

package examples

import (
	"fmt"
	"strings"
	"time"

	"github.com/qzeleza/termos/internal/task"
)

func main___() {
	fmt.Println("=== Анализ строки InputTaskNew ===\n")

	input := task.NewInputTaskNew("Тест ввода", "подсказка")
	input.WithTimeout(5*time.Second, "по умолчанию")

	view := input.View(80)
	lines := strings.Split(view, "\n")

	if len(lines) > 0 {
		firstLine := lines[0]
		fmt.Printf("Первая строка: '%s'\n", firstLine)
		fmt.Printf("Длина: %d символов\n", len(firstLine))

		// Ищем символы по отдельности
		for i, char := range firstLine {
			if i < 10 { // Первые 10 символов
				fmt.Printf("  [%d]: '%c' (\\u%04X)\n", i, char, char)
			}
		}

		// Ищем подстроку "└─>" в строке
		if strings.Contains(firstLine, "└─>") {
			fmt.Println("✅ Префикс '└─>' НАЙДЕН!")
		} else {
			fmt.Println("❌ Префикс '└─>' НЕ найден")
		}

		// Проверим символы Unicode
		if strings.Contains(firstLine, "└") {
			fmt.Println("✅ Символ '└' найден")
		}
		if strings.Contains(firstLine, "─") {
			fmt.Println("✅ Символ '─' найден")
		}
		if strings.Contains(firstLine, ">") {
			fmt.Println("✅ Символ '>' найден")
		}

		// Попробуем найти с разными начальными пробелами
		trimmed := strings.TrimLeft(firstLine, " ")
		fmt.Printf("После удаления пробелов: '%s'\n", trimmed)
		if strings.HasPrefix(trimmed, "└─>") {
			fmt.Println("✅ Найден префикс после удаления пробелов!")
		}
	}

	if len(lines) > 1 {
		secondLine := lines[1]
		fmt.Printf("\nВторая строка: '%s'\n", secondLine)

		if strings.Contains(secondLine, "...") {
			fmt.Println("✅ Поле ввода '...' НАЙДЕНО!")
		} else {
			fmt.Println("❌ Поле ввода '...' НЕ найдено")
		}
	}
}
