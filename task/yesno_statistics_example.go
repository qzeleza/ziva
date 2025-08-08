package task

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/common"
	"github.com/qzeleza/termos/query"
	"github.com/qzeleza/termos/ui"
)

// ExampleYesNoStatistics демонстрирует как YesNoTask теперь учитывается в статистике ошибок
func ExampleYesNoStatistics() error {
	// Создаем различные задачи для демонстрации статистики
	tasks := []common.Task{
		// Обычная функциональная задача
		NewFuncTask("Инициализация системы", func() error {
			time.Sleep(300 * time.Millisecond)
			return nil // Успешно
		}),

		// YesNoTask - пользователь может выбрать "Да" или "Нет"
		NewYesNoTask("Включить автозапуск", "Включить автоматический запуск при загрузке системы?"),

		// Еще одна функциональная задача
		NewFuncTask("Настройка сети", func() error {
			time.Sleep(400 * time.Millisecond)
			return nil // Успешно
		}),

		// Еще один YesNoTask
		NewYesNoTask("Отправлять статистику", "Разрешить отправку анонимной статистики использования?"),

		// Финальная задача
		NewFuncTask("Завершение настройки", func() error {
			time.Sleep(200 * time.Millisecond)
			return nil // Успешно
		}),
	}

	// Создаем модель очереди с включенной статистикой
	model := query.New("Настройка приложения").
		WithSummary(true). // Включаем отображение сводки со статистикой
		WithTitleColor(ui.ActiveStyle.GetForeground(), true).
		WithAppName("KvasPro Setup").
		WithAppNameColor(ui.SubtleStyle.GetForeground(), false)

	model.AddTasks(tasks)

	// Запускаем выполнение
	_, err := tea.NewProgram(model).Run()
	if err != nil {
		return err
	}

	return nil
}

// ExampleYesNoWithMixedResults демонстрирует очередь с смешанными результатами
func ExampleYesNoWithMixedResults() error {
	// Предварительно настроенные задачи для демонстрации
	task1 := NewYesNoTask("Согласие на обработку данных", "Согласны на обработку персональных данных?")
	task2 := NewYesNoTask("Подписка на рассылку", "Подписаться на рассылку новостей?")
	task3 := NewYesNoTask("Участие в программе улучшений", "Участвовать в программе улучшения продукта?")

	// Предустанавливаем ответы для демонстрации:
	// - task1: "Да" (успешно)
	// - task2: "Нет" (ошибка для статистики)
	// - task3: "Да" (успешно)

	// Симулируем пользовательский ввод
	// task1: выбираем "Да" (по умолчанию)
	task1.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// task2: переходим к "Нет" и выбираем
	task2.Update(tea.KeyMsg{Type: tea.KeyDown})
	task2.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// task3: выбираем "Да" (по умолчанию)
	task3.Update(tea.KeyMsg{Type: tea.KeyEnter})

	tasks := []common.Task{task1, task2, task3}

	model := query.New("Настройка согласий").
		WithSummary(true).
		WithTitleColor(ui.ActiveStyle.GetForeground(), true)

	model.AddTasks(tasks)

	// В результате статистика покажет:
	// "Обработка операций прошла (2/3)" и статус "С ОШИБКАМИ"
	// потому что task2 (отказ от рассылки) считается "ошибкой" для статистики

	_, err := tea.NewProgram(model).Run()
	return err
}

// ExampleYesNoAllAccepted демонстрирует полностью успешную очередь
func ExampleYesNoAllAccepted() error {
	task1 := NewYesNoTask("Условия использования", "Принимаете условия использования?")
	task2 := NewYesNoTask("Политика конфиденциальности", "Согласны с политикой конфиденциальности?")

	// Оба ответа "Да"
	task1.Update(tea.KeyMsg{Type: tea.KeyEnter})
	task2.Update(tea.KeyMsg{Type: tea.KeyEnter})

	tasks := []common.Task{task1, task2}

	model := query.New("Принятие соглашений").
		WithSummary(true).
		WithTitleColor(ui.ActiveStyle.GetForeground(), true)

	model.AddTasks(tasks)

	// В результате статистика покажет:
	// "Обработка операций прошла (2/2)" и статус "УСПЕШНО"

	_, err := tea.NewProgram(model).Run()
	return err
}

// ExampleYesNoAllRejected демонстрирует очередь, где все отвечают "Нет"
func ExampleYesNoAllRejected() error {
	task1 := NewYesNoTask("Дополнительные функции", "Установить дополнительные функции?")
	task2 := NewYesNoTask("Телеметрия", "Включить отправку телеметрии?")

	// Оба ответа "Нет"
	task1.Update(tea.KeyMsg{Type: tea.KeyDown})
	task1.Update(tea.KeyMsg{Type: tea.KeyEnter})
	task2.Update(tea.KeyMsg{Type: tea.KeyDown})
	task2.Update(tea.KeyMsg{Type: tea.KeyEnter})

	tasks := []common.Task{task1, task2}

	model := query.New("Отказ от дополнений").
		WithSummary(true).
		WithTitleColor(ui.ActiveStyle.GetForeground(), true)

	model.AddTasks(tasks)

	// В результате статистика покажет:
	// "Обработка операций прошла (0/2)" и статус "С ОШИБКАМИ"
	// потому что оба отказа считаются "ошибками" для статистики

	_, err := tea.NewProgram(model).Run()
	return err
}
