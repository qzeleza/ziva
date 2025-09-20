// +build ignore

package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("=== Тест выхода по Ctrl+C в FuncTask ===")
	fmt.Println("Нажмите Ctrl+C для выхода из задачи")
	fmt.Println()

	// Создаем задачу, которая выполняется долго
	longTask := task.NewFuncTask("Долгая задача (нажмите Ctrl+C для выхода)", func() error {
		// Имитируем долгую работу
		time.Sleep(30 * time.Second)
		return nil
	})

	// Запускаем задачу
	queue := query.New("Тест выхода по Ctrl+C")
	queue.AddTasks([]task.Task{longTask})

	p := tea.NewProgram(queue)
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
}
