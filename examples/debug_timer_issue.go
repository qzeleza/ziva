// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
)

func main() {
	// Создаем лог файл для отладки
	logFile, err := os.Create("/tmp/timer_debug.log")
	if err != nil {
		fmt.Printf("Ошибка создания лог файла: %v\n", err)
		return
	}
	defer logFile.Close()
	
	log.SetOutput(logFile)
	log.Println("=== НАЧАЛО ОТЛАДКИ ТАЙМЕРА ===")

	fmt.Println("=== Тест проблемы с таймером в MultiSelectTask ===")
	fmt.Println("Таймер: 10 секунд")
	fmt.Println("Попробуйте навигацию стрелками - таймер НЕ должен сбрасываться")
	fmt.Println("Лог отладки записывается в /tmp/timer_debug.log")
	fmt.Println()

	choices := []string{
		"Элемент 1",
		"Элемент 2", 
		"Элемент 3",
		"Элемент 4",
		"Элемент 5",
	}

	// Создаем задачу с коротким таймером для тестирования
	multiTask := task.NewMultiSelectTask("Тест таймера", choices).
		WithTimeout(10 * time.Second).
		WithViewport(3).
		WithSelectAll("Выбрать все")

	log.Printf("Создана задача с таймером %v", 10*time.Second)

	// Запускаем задачу
	queue := query.New("Тест таймера")
	queue.AddTasks([]task.Task{multiTask})

	log.Println("Запускаем программу...")
	p := tea.NewProgram(queue)
	_, err = p.Run()
	if err != nil {
		log.Printf("Ошибка выполнения: %v", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}

	selected := multiTask.GetSelected()
	log.Printf("Результат: выбрано %d элементов", len(selected))
	fmt.Printf("Результат: выбрано %d элементов\n", len(selected))
	
	log.Println("=== КОНЕЦ ОТЛАДКИ ТАЙМЕРА ===")
}
