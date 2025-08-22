// +build ignore

package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
)

func main() {
	// Создаем очередь задач с демонстрацией таймаутов
	queue := query.New("Демонстрация таймаутов")

	// 1. Задача ввода текста с таймаутом
	inputTask := task.NewInputTaskNew("Введите ваше имя", "Имя")
	inputTask.WithTimeout(5*time.Second, "Гость") // Через 5 сек будет "Гость"

	// 2. Задача выбора Да/Нет с таймаутом
	yesNoTask := task.NewYesNoTask("Продолжить?", "Хотите продолжить выполнение?")
	yesNoTask.WithTimeout(7*time.Second, 0) // Через 7 сек выберет "Да" (индекс 0)

	// 3. Задача одиночного выбора с таймаутом
	singleSelectTask := task.NewSingleSelectTask(
		"Выберите цвет",
		[]string{"Красный", "Зеленый", "Синий", "Желтый"},
	)
	singleSelectTask.WithTimeout(10*time.Second, "Синий") // Через 10 сек выберет "Синий"

	// 4. Задача множественного выбора с таймаутом
	multiSelectTask := task.NewMultiSelectTask(
		"Выберите языки программирования",
		[]string{"Go", "Python", "JavaScript", "Rust", "Java"},
	)
	// Через 8 сек выберет Go и Rust
	multiSelectTask.WithTimeout(8*time.Second, []string{"Go", "Rust"})

	// 5. Задача ввода пароля с таймаутом
	passwordTask := task.NewInputTaskNew("Введите пароль", "Пароль").
		WithInputType(task.InputTypePassword).
		WithAllowEmpty(true) // Разрешаем пустой пароль
	passwordTask.WithTimeout(6*time.Second, "defaultpass123")

	// Добавляем все задачи в очередь
	queue.AddTasks([]task.Task{
		inputTask,
		yesNoTask,
		singleSelectTask,
		multiSelectTask,
		passwordTask,
	})

	// Запускаем программу
	p := tea.NewProgram(queue)
	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Получаем результаты
	if _, ok := finalModel.(*query.Model); ok {
		fmt.Println("\n=== Результаты выполнения задач с таймаутами ===")
		
		// Имя пользователя
		if name := inputTask.GetValue(); name != "" {
			fmt.Printf("Имя: %s\n", name)
		}

		// Выбор продолжения
		if yesNoTask.IsYes() {
			fmt.Println("Пользователь выбрал: Да")
		} else if yesNoTask.IsNo() {
			fmt.Println("Пользователь выбрал: Нет")
		} else {
			fmt.Println("Пользователь выбрал: Выйти")
		}

		// Выбранный цвет
		if color := singleSelectTask.GetSelected(); color != "" {
			fmt.Printf("Выбранный цвет: %s\n", color)
		}

		// Выбранные языки
		if langs := multiSelectTask.GetSelected(); len(langs) > 0 {
			fmt.Printf("Выбранные языки: %v\n", langs)
		}

		// Пароль (маскированный)
		if pass := passwordTask.GetValue(); pass != "" {
			fmt.Printf("Пароль установлен (длина: %d символов)\n", len(pass))
		}

		// Статистика выполнения
		fmt.Printf("\n=== Статистика ===\n")
		fmt.Printf("Задачи выполнены\n")
	}
}
