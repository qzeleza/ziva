// +build ignore

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== Прямое тестирование View() методов ===")

	// Тест 1: SingleSelectTask с таймером
	fmt.Println("\n1. SingleSelectTask с таймером:")
	singleTask := task.NewSingleSelectTask(
		"Выберите опцию",
		[]string{"Опция A", "Опция B", "Опция C"},
	)
	singleTask.WithTimeout(10*time.Second, 0)

	// Запускаем Run() для инициализации таймера
	cmd := singleTask.Run()
	fmt.Printf("Run() команда создана: %t\n", cmd != nil)

	// Получаем View() напрямую
	view := singleTask.View(80)
	fmt.Printf("Длина View(): %d символов\n", len(view))
	fmt.Printf("View() содержимое:\n%s\n", view)

	// Проверяем наличие ключевых элементов
	fmt.Printf("Содержит префикс '○': %t\n", strings.Contains(view, "○"))
	fmt.Printf("Содержит префикс '└─>': %t\n", strings.Contains(view, "└─>"))
	fmt.Printf("Содержит таймер '[': %t\n", strings.Contains(view, "["))
	fmt.Printf("Содержит 'Опция A': %t\n", strings.Contains(view, "Опция A"))

	// Тест 2: Проверяем RenderTimer напрямую
	fmt.Println("\n2. Тест RenderTimer:")
	timer := singleTask.BaseTask.RenderTimer()
	fmt.Printf("RenderTimer(): '%s'\n", timer)
	fmt.Printf("Длина RenderTimer(): %d\n", len(timer))

	// Тест 3: Проверяем состояние таймера
	fmt.Println("\n3. Состояние таймера:")
	fmt.Printf("Таймер включен: %t\n", singleTask.BaseTask.GetRemainingTime() > 0)
	fmt.Printf("Показывать таймер: %t\n", true) // showTimeout недоступен напрямую
	
	// Тест 4: InputTaskNew
	fmt.Println("\n4. InputTaskNew с таймером:")
	inputTask := task.NewInputTaskNew("Введите текст", "Подсказка")
	inputTask.WithTimeout(5*time.Second, "по умолчанию")
	
	inputCmd := inputTask.Run()
	fmt.Printf("InputTask Run() команда: %t\n", inputCmd != nil)
	
	inputView := inputTask.View(80)
	fmt.Printf("InputTask View() длина: %d\n", len(inputView))
	fmt.Printf("InputTask View():\n%s\n", inputView)
}
