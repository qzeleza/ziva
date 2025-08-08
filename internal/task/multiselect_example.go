package task

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/qzeleza/termos/internal/common"
	"github.com/qzeleza/termos/internal/query"
)

// ExampleMultiSelectWithSelectAll демонстрирует использование MultiSelectTask с опцией "Выбрать все"
func ExampleMultiSelectWithSelectAll() error {
	// Создаем задачу множественного выбора с опцией "Выбрать все"
	components := []string{
		"API Server",
		"Frontend App",
		"Database",
		"Worker Queue",
		"Authentication Service",
		"Logging Service",
	}

	// Основная задача с опцией "Выбрать все"
	mainTask := NewMultiSelectTask(
		"Выберите компоненты для установки",
		components,
	).WithSelectAll()

	// Дополнительная задача с кастомным текстом
	additionalComponents := []string{
		"Мониторинг",
		"Бэкап система",
		"CDN",
	}

	additionalTask := NewMultiSelectTask(
		"Выберите дополнительные сервисы",
		additionalComponents,
	).WithSelectAll("Включить все дополнительные сервисы")

	// Создаем список задач
	tasks := []common.Task{
		mainTask,
		additionalTask,
	}

	// Создаем модель очереди
	model := query.New("Мастер установки компонентов")
	model.AddTasks(tasks)

	// Запускаем программу
	_, err := tea.NewProgram(model).Run()
	if err != nil {
		return err
	}

	// Получаем результаты
	selectedMain := mainTask.GetSelected()
	selectedAdditional := additionalTask.GetSelected()

	// В реальном приложении здесь бы была логика обработки выбранных компонентов
	_ = selectedMain
	_ = selectedAdditional

	return nil
}

// ExampleMultiSelectBasic демонстрирует использование обычной MultiSelectTask без "Выбрать все"
func ExampleMultiSelectBasic() error {
	// Создаем обычную задачу множественного выбора (без "Выбрать все")
	features := []string{
		"SSL поддержка",
		"Кеширование",
		"Сжатие данных",
		"Rate limiting",
	}

	featureTask := NewMultiSelectTask(
		"Выберите необходимые функции",
		features,
	)

	// Создаем список задач
	tasks := []common.Task{featureTask}

	// Создаем модель очереди
	model := query.New("Конфигурация функций")
	model.AddTasks(tasks)

	// Запускаем программу
	_, err := tea.NewProgram(model).Run()
	if err != nil {
		return err
	}

	// Получаем результаты
	selectedFeatures := featureTask.GetSelected()
	_ = selectedFeatures

	return nil
}
