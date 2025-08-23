package examples

import (
	"log"

	"github.com/qzeleza/termos"
)

func main____() {
	// Создаем новую очередь задач с заголовком
	queue := termos.NewQueue("Пример очистки экрана")

	// Включаем очистку экрана перед запуском очереди
	queue.WithClearScreen(true)

	// Создаем задачу выбора Да/Нет
	yesNoTask := termos.NewYesNoTask("Подтверждение", "Продолжить выполнение?")

	// Создаем задачу выбора из списка
	options := []string{"Опция 1", "Опция 2", "Опция 3"}
	selectTask := termos.NewSingleSelectTask("Выберите опцию", options)

	// Добавляем задачи в очередь
	queue.AddTasks(yesNoTask, selectTask)

	// Запускаем очередь задач
	// Экран будет очищен перед запуском
	err := queue.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Выводим результаты
	log.Printf("Выбрано: %s", selectTask.GetSelected())
	if yesNoTask.IsYes() {
		log.Println("Пользователь выбрал 'Да'")
	} else {
		log.Println("Пользователь выбрал 'Нет'")
	}
}
