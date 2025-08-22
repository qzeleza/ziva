// +build ignore

package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
)

func main() {
	fmt.Println("Запуск демонстрации таймаутов...")
	fmt.Println("Подсказка:")
	fmt.Println("- Таймер отображается справа от заголовка [MM:SS]")  
	fmt.Println("- В задачах выбора нажмите ПРОБЕЛ чтобы остановить таймер")
	fmt.Println("- В задачах ввода любой символ остановит таймер")
	fmt.Println("- Навигация: стрелки ↑/↓, Enter - выбор, Ctrl+C - отмена")
	fmt.Println("")

	// Создаем очередь
	queue := query.New("🕐 Демонстрация таймаутов")

	// 1. Быстрая задача выбора (5 сек)
	quickSelect := task.NewSingleSelectTask(
		"Быстрый выбор (5 сек)",
		[]string{"🚀 Быстро", "🐌 Медленно", "⚡ Мгновенно"},
	)
	quickSelect.WithTimeout(5*time.Second, "⚡ Мгновенно")

	// 2. Задача ввода (8 сек)
	textInput := task.NewInputTaskNew("Введите ваше имя (8 сек)", "имя")
	textInput.WithTimeout(8*time.Second, "Анонимный пользователь")

	// 3. YesNo задача (6 сек)
	yesNo := task.NewYesNoTask("Продолжить демо? (6 сек)", "Согласны ли вы продолжить?")
	yesNo.WithTimeout(6*time.Second, 0) // 0 = "Да"

	// 4. Множественный выбор (10 сек)
	multiSelect := task.NewMultiSelectTask(
		"Выберите технологии (10 сек)",
		[]string{"🔵 Go", "🐍 Python", "⚡ JavaScript", "🦀 Rust", "☕ Java"},
	)
	multiSelect.WithTimeout(10*time.Second, []string{"🔵 Go", "🐍 Python"})

	// Добавляем задачи в очередь
	queue.AddTasks([]task.Task{
		quickSelect,
		textInput,
		yesNo,
		multiSelect,
	})

	// Запускаем интерактивную программу
	p := tea.NewProgram(queue)
	finalModel, err := p.Run()
	if err != nil {
		log.Fatal("Ошибка выполнения:", err)
	}

	// Выводим результаты
	if _, ok := finalModel.(*query.Model); ok {
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println("🎉 РЕЗУЛЬТАТЫ ДЕМОНСТРАЦИИ ТАЙМАУТОВ")
		fmt.Println(strings.Repeat("=", 60))

		fmt.Printf("1️⃣ Быстрый выбор: %s\n", quickSelect.GetSelected())
		fmt.Printf("2️⃣ Ваше имя: %s\n", textInput.GetValue())
		
		if yesNo.IsYes() {
			fmt.Printf("3️⃣ Продолжение: ✅ Да\n")
		} else if yesNo.IsNo() {
			fmt.Printf("3️⃣ Продолжение: ❌ Нет\n")  
		} else {
			fmt.Printf("3️⃣ Продолжение: 🚪 Выход\n")
		}

		selected := multiSelect.GetSelected()
		if len(selected) > 0 {
			fmt.Printf("4️⃣ Технологии: %v\n", selected)
		} else {
			fmt.Printf("4️⃣ Технологии: (ничего не выбрано)\n")
		}

		fmt.Println(strings.Repeat("=", 60))
		fmt.Println("✨ Демонстрация завершена успешно!")
	}
}
