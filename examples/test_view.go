// +build ignore

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== Тест отображения задач с таймаутом ===\n")

	// Тестируем SingleSelectTask
	fmt.Println("1. SingleSelectTask:")
	singleTask := task.NewSingleSelectTask(
		"Выберите опцию",
		[]string{"Вариант 1", "Вариант 2", "Вариант 3"},
	)
	singleTask.WithTimeout(10*time.Second, 1)
	printView(singleTask.View(80))

	// Тестируем InputTaskNew
	fmt.Println("\n2. InputTaskNew:")
	inputTask := task.NewInputTaskNew("Введите ваше имя", "Имя:")
	inputTask.WithTimeout(8*time.Second, "Пользователь")
	printView(inputTask.View(80))

	// Тестируем YesNoTask
	fmt.Println("\n3. YesNoTask:")
	yesNoTask := task.NewYesNoTask("Вопрос для пользователя", "Вы согласны?")
	yesNoTask.WithTimeout(6*time.Second, true)
	printView(yesNoTask.View(80))

	// Тестируем MultiSelectTask
	fmt.Println("\n4. MultiSelectTask:")
	multiTask := task.NewMultiSelectTask(
		"Выберите несколько опций",
		[]string{"Option 1", "Option 2", "Option 3", "Option 4"},
	)
	multiTask.WithTimeout(10*time.Second, []int{0, 2})
	printView(multiTask.View(80))
}

func printView(output string) {
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		fmt.Printf("%2d| %s\n", i+1, line)
	}
}
