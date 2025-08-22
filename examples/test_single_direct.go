// +build ignore

package main

import (
	"fmt"
	"time"

	"github.com/qzeleza/termos/internal/task"
)

func main() {
	// Создаем точно такую же задачу как в main.go
	ss := task.NewSingleSelectTask(
		"Выберите среду развертывания",
		[]string{"development", "staging", "production"},
	).WithTimeout(10*time.Second, "staging")

	fmt.Printf("Задача создана. Состояние IsDone(): %t\n", ss.IsDone())
	
	// Запускаем задачу
	cmd := ss.Run()
	fmt.Printf("Команда Run() выполнена. Команда: %v\n", cmd)
	fmt.Printf("Состояние IsDone() после Run(): %t\n", ss.IsDone())

	// Проверяем View() напрямую
	view := ss.View(80)
	fmt.Printf("Длина View(): %d\n", len(view))
	fmt.Printf("Содержимое View():\n%s\n", view)
	
	// Ждем немного и проверим снова
	time.Sleep(1 * time.Second)
	fmt.Printf("Состояние IsDone() через 1 секунду: %t\n", ss.IsDone())
	view2 := ss.View(80)
	fmt.Printf("View() через 1 секунду:\n%s\n", view2)
}
