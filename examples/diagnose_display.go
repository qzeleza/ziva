// +build ignore

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== ДИАГНОСТИКА ПРОБЛЕМ ОТОБРАЖЕНИЯ ===\n")

	// Тестируем InputTaskNew
	fmt.Println("1. InputTaskNew - детальная проверка:")
	input := task.NewInputTaskNew("Введите ваше имя", "Имя:")
	input.WithTimeout(8*time.Second, "Пользователь")
	
	// Проверим, запускается ли таймер
	input.Run()
	
	view := input.View(80)
	fmt.Printf("Сырой вывод View():\n'%s'\n\n", view)
	
	// Разбиваем по строкам для анализа
	lines := strings.Split(view, "\n")
	fmt.Printf("Количество строк: %d\n", len(lines))
	for i, line := range lines {
		fmt.Printf("Строка %d: '%s' (длина: %d)\n", i+1, line, len(line))
	}
	
	// Проверим таймер отдельно
	timer := input.BaseTask.RenderTimer()
	fmt.Printf("\nТаймер отдельно: '%s' (длина: %d)\n", timer, len(timer))
	
	// Проверим состояние таймера (доступ к приватному полю через рефлексию не нужен, просто проверим что таймер работает)
	if len(timer) > 0 {
		fmt.Println("✓ Таймер работает (есть содержимое)")
	} else {
		fmt.Println("✗ Таймер НЕ работает (пустое содержимое)")
	}
	
	// Тестируем SingleSelectTask
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("2. SingleSelectTask - детальная проверка:")
	single := task.NewSingleSelectTask("Выберите опцию", []string{"Вариант A", "Вариант B", "Вариант C"})
	single.WithTimeout(10*time.Second, 1)
	single.Run()
	
	singleView := single.View(80)
	fmt.Printf("Сырой вывод View():\n'%s'\n\n", singleView)
	
	singleLines := strings.Split(singleView, "\n")
	fmt.Printf("Количество строк: %d\n", len(singleLines))
	for i, line := range singleLines {
		fmt.Printf("Строка %d: '%s' (длина: %d)\n", i+1, line, len(line))
	}
	
	singleTimer := single.BaseTask.RenderTimer()
	fmt.Printf("\nТаймер отдельно: '%s' (длина: %d)\n", singleTimer, len(singleTimer))
	
	if len(singleTimer) > 0 {
		fmt.Println("✓ Таймер работает (есть содержимое)")
	} else {
		fmt.Println("✗ Таймер НЕ работает (пустое содержимое)")
	}
}
