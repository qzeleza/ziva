package examples

import (
	"fmt"
	"log"

	"github.com/qzeleza/termos/task"
)

func RunBasicUsage() {
	fmt.Println("Termos - Пример базового использования")
	fmt.Println("=====================================")

	// Пример 1: Задача Yes/No
	fmt.Println("\n1. Пример задачи Yes/No:")

	yesNoTask := task.NewYesNoTask(
		"Хотите продолжить установку?",
		"Подтвердите ваше действие",
	)

	// Запускаем задачу (в реальном приложении здесь был бы TUI)
	// Для примера просто симулируем выбор
	fmt.Println("   Заголовок: " + yesNoTask.Title())
	fmt.Println("   Результат: Да (симуляция)")

	fmt.Println("\n2. Пример задачи ввода текста:")

	inputTask := task.NewInputTaskNew(
		"Введите ваше имя",
		"",
	)

	fmt.Println("   Заголовок: " + inputTask.Title())
	fmt.Println("   Результат: Иван Иванов (симуляция)")

	fmt.Println("\n3. Пример задачи выбора из списка:")

	options := []string{
		"Установить все компоненты",
		"Выборочная установка",
		"Только базовые компоненты",
	}

	selectTask := task.NewSingleSelectTask(
		"Выберите тип установки",
		options,
	)

	fmt.Println("   Заголовок: " + selectTask.Title())
	fmt.Println("   Опции:")
	for i, option := range options {
		fmt.Printf("     %d. %s\n", i+1, option)
	}
	fmt.Println("   Результат: Установить все компоненты (симуляция)")

	fmt.Println("\n✅ Все примеры завершены!")
	fmt.Println("\nВ реальном приложении эти задачи запускались бы через:")
	fmt.Println("  result, err := task.Run()")
	fmt.Println("  if err != nil { ... }")

	// Пример с обработкой ошибок
	fmt.Println("\n4. Пример обработки ошибок:")

	// Симуляция ошибки
	errorTask := &task.BaseTask{}
	errorTask.SetError(fmt.Errorf("не удалось подключиться к серверу"))

	if errorTask.HasError() {
		fmt.Printf("   Ошибка: %v\n", errorTask.Error())
		fmt.Printf("   Остановка очереди: %v\n", errorTask.StopOnError())
	}

	log.Println("\nДля запуска интерактивных примеров используйте терминал с поддержкой TUI")
}
