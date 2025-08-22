package main

import (
	"log"
	"time"

	"github.com/qzeleza/termos"
)

func main() {
	// Создаем новую очередь задач
	queue := termos.NewQueue("Простой пример использования Termos").
		WithAppName("Пример").
		WithSummary(true)

	// Создаем задачу выбора Да/Нет
	confirm := termos.NewYesNoTask("Подтверждение", "Хотите продолжить?").
		WithTimeout(10*time.Second, "Да")

	// Создаем задачу выбора из списка
	env := termos.NewSingleSelectTask("Выбор среды", []string{
		"development", 
		"staging", 
		"production",
	}).WithTimeout(10*time.Second, "development")

	// Создаем задачу ввода текста
	name := termos.NewInputTask("Имя пользователя", "Введите ваше имя:").
		WithValidator(termos.DefaultValidators.Required()).
		WithTimeout(15*time.Second, "Anonymous")

	// Создаем задачу множественного выбора
	components := termos.NewMultiSelectTask("Выбор компонентов", []string{
		"API Server",
		"Web Interface", 
		"Database",
		"Cache",
		"Monitoring",
	}).WithSelectAll("Выбрать все").
		WithTimeout(15*time.Second, []string{"API Server", "Web Interface"})

	// Создаем задачу выполнения функции
	deploy := termos.NewFuncTaskWithOptions("Развертывание",
		func() error {
			// Имитируем работу
			time.Sleep(2 * time.Second)
			return nil
		},
		termos.WithSummaryFunction(func() []string {
			return []string{
				"Сервисы запущены: 3",
				"Время развертывания: 2.1с",
				"Статус: OK",
			}
		}),
	)

	// Добавляем все задачи в очередь
	queue.AddTasks(confirm, env, name, components, deploy)

	// Запускаем очередь
	if err := queue.Run(); err != nil {
		log.Fatal(err)
	}
}
