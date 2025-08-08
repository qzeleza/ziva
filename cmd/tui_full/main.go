package main

import (
	"fmt"
	"github.com/qzeleza/termos/common"
	"github.com/qzeleza/termos/examples"
	"github.com/qzeleza/termos/task"
)

// Точка входа для полнофункционального TUI
// Запуск: go run ./cmd/tui_full
func main() {
	fmt.Println("Termos - Полнофункциональный TUI")
	fmt.Println("=================================")

	tasks := []common.Task{
		task.NewYesNoTask("Продолжить?", "Подтвердите действие"),
		task.NewInputTaskNew("Введите имя", ""),
		task.NewSingleSelectTask("Выберите вариант", []string{"A", "B", "C"}),
	}

	// Обертка из пакета examples применит embedded-оптимизации при необходимости
	if err := examples.RunTasksWithTUI("Демонстрация Termos", "Итоги будут показаны в конце", tasks); err != nil {
		panic(err)
	}
}
