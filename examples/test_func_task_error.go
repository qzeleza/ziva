// +build ignore

package main

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== Тест отображения ошибок в FuncTask ===")
	fmt.Println("Проверяем правильное отображение цветов для текста ошибки и статуса")
	fmt.Println()

	// Создаем задачу с ошибкой
	errorTask := task.NewFuncTask("Задача с ошибкой", func() error {
		return errors.New("тестовая ошибка для проверки стилей")
	})

	// Создаем успешную задачу
	successTask := task.NewFuncTask("Успешная задача", func() error {
		return nil
	}).WithSummary(func() []string {
		return []string{
			"Дополнительная информация",
			"Еще одна строка",
		}
	})

	// Запускаем задачи
	queue := query.New("Тест стилей ошибок")
	queue.AddTasks([]task.Task{errorTask, successTask})

	p := tea.NewProgram(queue)
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
}
