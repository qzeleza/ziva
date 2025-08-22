// +build ignore

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== Отладка InputTaskNew ===\n")

	input := task.NewInputTaskNew("Тест ввода", "подсказка")
	input.WithTimeout(5*time.Second, "по умолчанию")
	input.Run()
	
	view := input.View(80)
	lines := strings.Split(view, "\n")
	
	fmt.Printf("Общая длина view: %d символов\n", len(view))
	fmt.Printf("Количество строк: %d\n", len(lines))
	fmt.Println("\nСтроки view:")
	for i, line := range lines {
		fmt.Printf("%2d| '%s'\n", i+1, line)
	}
	
	fmt.Printf("\nПоиск префикса '└─>': %t\n", strings.Contains(view, "└─>"))
	fmt.Printf("Поиск поля ввода '...': %t\n", strings.Contains(view, "..."))
	
	timer := input.BaseTask.RenderTimer()
	fmt.Printf("Таймер: '%s'\n", timer)
}
