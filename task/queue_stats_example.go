package task

import (
	"errors"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/common"
	"github.com/qzeleza/termos/query"
	"github.com/qzeleza/termos/ui"
)

// ExampleQueueWithStats демонстрирует использование очереди задач с подсчетом статистики
func ExampleQueueWithStats() error {
	// Создаем задачи с разными результатами выполнения
	tasks := []common.Task{
		// Успешные задачи
		NewFuncTask("Инициализация системы", func() error {
			time.Sleep(500 * time.Millisecond)
			return nil // Успешно
		}),

		NewFuncTask("Загрузка конфигурации", func() error {
			time.Sleep(300 * time.Millisecond)
			return nil // Успешно
		}),

		// Задача с ошибкой, но не прерывающая выполнение
		NewFuncTask("Подключение к базе данных", func() error {
			time.Sleep(400 * time.Millisecond)
			return errors.New("не удалось подключиться к базе данных")
		}, false), // false означает "не останавливать очередь при ошибке"

		// Еще одна успешная задача
		NewFuncTask("Запуск веб-сервера", func() error {
			time.Sleep(600 * time.Millisecond)
			return nil // Успешно
		}),

		// Задача с ошибкой
		NewFuncTask("Отправка уведомлений", func() error {
			time.Sleep(200 * time.Millisecond)
			return errors.New("сервис уведомлений недоступен")
		}, false), // Не останавливать очередь

		// Финальная успешная задача
		NewFuncTask("Завершение инициализации", func() error {
			time.Sleep(300 * time.Millisecond)
			return nil // Успешно
		}),
	}

	// Создаем модель очереди с настроенными стилями
	model := query.New("Запуск приложения").
		WithSummary(true). // Включаем отображение сводки
		WithTitleColor(ui.ActiveStyle.GetForeground(), true).
		WithAppName("KvasPro v1.0").
		WithAppNameColor(ui.SubtleStyle.GetForeground(), false)

	// Добавляем задачи в очередь
	model.AddTasks(tasks)

	// Запускаем выполнение
	_, err := tea.NewProgram(model).Run()
	if err != nil {
		return err
	}

	return nil
}

// ExampleQueueAllSuccessful демонстрирует очередь, где все задачи выполняются успешно
func ExampleQueueAllSuccessful() error {
	tasks := []common.Task{
		NewFuncTask("Проверка системы", func() error {
			time.Sleep(300 * time.Millisecond)
			return nil
		}),

		NewFuncTask("Подготовка данных", func() error {
			time.Sleep(400 * time.Millisecond)
			return nil
		}),

		NewFuncTask("Выполнение операции", func() error {
			time.Sleep(500 * time.Millisecond)
			return nil
		}),
	}

	model := query.New("Успешный сценарий").
		WithSummary(true) // Включаем отображение сводки

	model.AddTasks(tasks)

	// В этом случае сводка покажет: "Обработка операций прошла (3/3)" и статус "УСПЕШНО"
	_, err := tea.NewProgram(model).Run()
	return err
}

// ExampleQueueWithCriticalError демонстрирует очередь, прерванную критической ошибкой
func ExampleQueueWithCriticalError() error {
	tasks := []common.Task{
		NewFuncTask("Проверка прав доступа", func() error {
			time.Sleep(300 * time.Millisecond)
			return nil
		}),

		// Критическая ошибка, останавливающая выполнение
		NewFuncTask("Критическая операция", func() error {
			time.Sleep(400 * time.Millisecond)
			return errors.New("критическая ошибка системы")
		}), // По умолчанию stopOnError = true

		// Эта задача не будет выполнена из-за критической ошибки выше
		NewFuncTask("Заключительные операции", func() error {
			time.Sleep(200 * time.Millisecond)
			return nil
		}),
	}

	model := query.New("Сценарий с критической ошибкой").
		WithSummary(true) // Включаем отображение сводки

	model.AddTasks(tasks)

	// В этом случае сводка покажет: "Обработка операций прошла (1/3)" и статус "С ОШИБКАМИ"
	_, err := tea.NewProgram(model).Run()
	return err
}

// ExampleQueueWithoutSummary демонстрирует очередь без отображения сводки (по умолчанию)
func ExampleQueueWithoutSummary() error {
	tasks := []common.Task{
		NewFuncTask("Подготовка системы", func() error {
			time.Sleep(300 * time.Millisecond)
			return nil
		}),

		NewFuncTask("Выполнение операции", func() error {
			time.Sleep(400 * time.Millisecond)
			return nil
		}),

		NewFuncTask("Завершение", func() error {
			time.Sleep(200 * time.Millisecond)
			return nil
		}),
	}

	// Создаем модель без включения сводки (по умолчанию showSummary = false)
	model := query.New("Быстрое выполнение").
		WithTitleColor(ui.ActiveStyle.GetForeground(), true)

	model.AddTasks(tasks)

	// В этом случае сводка НЕ будет отображена после завершения всех задач
	_, err := tea.NewProgram(model).Run()
	return err
}
