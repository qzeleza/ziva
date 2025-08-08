// query/query_model_test.go

package query

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/qzeleza/termos/internal/common"
)

// mockTask - мок задачи для тестирования модели Query.
// Реализует интерфейс common.Task
type mockTask struct {
	title     string
	done      bool
	err       error
	value     string
	stopOnErr bool // Флаг остановки при ошибке
}

func newMockTask(title string) *mockTask {
	return &mockTask{
		title: title,
		done:  false,
	}
}

func (t *mockTask) Title() string {
	return t.title
}

func (t *mockTask) IsDone() bool {
	return t.done
}

func (t *mockTask) Error() error {
	return t.err
}

func (t *mockTask) HasError() bool {
	return t.err != nil
}

func (t *mockTask) StopOnError() bool {
	return t.stopOnErr
}

func (t *mockTask) SetStopOnError(stop bool) {
	t.stopOnErr = stop
}

// Run запускает выполнение задачи и возвращает команду bubbletea.
func (t *mockTask) Run() tea.Cmd {
	// В моке просто возвращаем nil, так как нам не нужно реальное выполнение
	return nil
}

// Update обновляет состояние задачи на основе полученного сообщения.
func (t *mockTask) Update(msg tea.Msg) (common.Task, tea.Cmd) {
	// Для тестирования считаем, что любое сообщение завершает задачу
	t.done = true
	return t, nil
}

// View отображает текущее состояние задачи с учетом указанной ширины.
func (t *mockTask) View(width int) string {
	if t.done {
		return ""
	}
	return t.title
}

// FinalView отображает финальное состояние задачи с учетом указанной ширины.
func (t *mockTask) FinalView(width int) string {
	if t.HasError() {
		return "ОШИБКА: " + t.title + ": " + t.err.Error()
	}
	return t.title + ": " + t.value
}

// TestQueryModelCreation проверяет корректность создания модели Query
func TestQueryModelCreation(t *testing.T) {
	// Создаем список задач для модели
	tasks := []common.Task{
		newMockTask("Задача 1"),
		newMockTask("Задача 2"),
	}

	// Создаем модель Query с заголовком
	model := New("Тестовый заголовок")
	model.AddTasks(tasks)

	// Проверяем, что модель создана корректно
	assert.NotNil(t, model, "Модель не должна быть nil")
	assert.Equal(t, "Тестовый заголовок", model.title, "Заголовок модели должен соответствовать заданному")
	assert.Equal(t, 2, len(model.tasks), "Количество задач должно соответствовать заданному")
	assert.Equal(t, 0, model.current, "Индекс текущей задачи должен быть 0")
	assert.False(t, model.quitting, "Флаг завершения должен быть false")
	assert.False(t, model.stoppedOnError, "Флаг остановки из-за ошибки должен быть false")
}

// TestQueryModelUpdate проверяет обновление модели при получении сообщения
func TestQueryModelUpdate(t *testing.T) {
	// Создаем список задач для модели
	task1 := newMockTask("Задача 1")
	task2 := newMockTask("Задача 2")
	tasks := []common.Task{task1, task2}

	// Создаем модель Query
	model := New("Тестовая модель")
	model.AddTasks(tasks)

	// Проверяем обработку сообщения для первой задачи
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Проверяем, что модель обновлена корректно
	queryModel, ok := updatedModel.(*Model)
	assert.True(t, ok, "Обновленная модель должна быть типа *Model")

	// После обновления первая задача должна быть завершена и индекс должен перейти к следующей задаче
	assert.True(t, task1.IsDone(), "Первая задача должна быть завершена")
	assert.Equal(t, 1, queryModel.current, "Индекс текущей задачи должен быть 1")

	// Обрабатываем сообщение для второй задачи
	updatedModel, _ = queryModel.Update(tea.KeyMsg{Type: tea.KeyEnter})
	queryModel, _ = updatedModel.(*Model)

	// Проверяем, что вторая задача завершена и модель перешла к следующей задаче
	assert.True(t, task2.IsDone(), "Вторая задача должна быть завершена")
	assert.Equal(t, 2, queryModel.current, "Индекс текущей задачи должен быть равен количеству задач, что означает завершение всех задач")
}

// TestQueryModelView проверяет отображение модели
func TestQueryModelView(t *testing.T) {
	// Создаем список задач для модели
	task1 := newMockTask("Задача 1")
	task2 := newMockTask("Задача 2")
	tasks := []common.Task{task1, task2}

	// Создаем модель Query с заголовком
	model := New("Тестовый заголовок")
	model.AddTasks(tasks)

	// Проверяем, что View содержит заголовок и текущую задачу
	view := model.View()
	assert.Contains(t, view, "Тестовый заголовок", "View должен содержать заголовок модели")
	assert.Contains(t, view, "Задача 1", "View должен содержать название текущей задачи")

	// Завершаем первую задачу и проверяем отображение второй
	model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	view = model.View()
	assert.Contains(t, view, "Задача 2", "View должен содержать название текущей задачи")
}

// TestQueryModelErrorHandling проверяет обработку ошибок в модели
func TestQueryModelErrorHandling(t *testing.T) {
	// Создаем задачу с ошибкой
	task1 := newMockTask("Задача с ошибкой")
	task1.err = errors.New("тестовая ошибка")
	task1.SetStopOnError(true)

	// Создаем вторую задачу, которая не должна выполниться
	task2 := newMockTask("Задача 2")

	// Создаем модель Query
	model := New("Тестовая модель с ошибкой")
	model.AddTasks([]common.Task{task1, task2})

	// Обновляем модель
	updatedModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	queryModel, _ := updatedModel.(*Model)

	// Проверяем, что модель остановлена из-за ошибки
	assert.True(t, queryModel.stoppedOnError, "Модель должна быть остановлена из-за ошибки")
	assert.Equal(t, task1, queryModel.errorTask, "Задача с ошибкой должна быть сохранена в errorTask")

	// Проверяем, что вторая задача не выполнена
	assert.False(t, task2.IsDone(), "Вторая задача не должна быть выполнена")
}

// TestQueryModelSetTitle проверяет установку заголовка модели
func TestQueryModelSetTitle(t *testing.T) {
	// Создаем модель Query
	tasks := []common.Task{
		newMockTask("Задача 1"),
	}

	model := New("Тестовый заголовок")
	model.AddTasks(tasks)

	// Устанавливаем заголовок
	title := "Тестовый заголовок"
	model.title = title

	// Проверяем, что заголовок установлен корректно
	assert.Equal(t, title, model.title, "Заголовок модели должен соответствовать переданному значению")
}

// TestQueryModelResults проверяет получение результатов выполнения задач
func TestQueryModelResults(t *testing.T) {
	// Создаем задачи с предопределенными значениями
	task1 := &mockTask{
		title: "Задача 1",
		done:  true,
		value: "Результат 1",
	}

	task2 := &mockTask{
		title: "Задача 2",
		done:  true,
		value: "Результат 2",
	}

	// Создаем модель Query с завершенными задачами
	tasks := []common.Task{task1, task2}
	model := New("Тестовый заголовок")
	model.AddTasks(tasks)

	// Сначала проверяем финальное представление задач без флага quitting
	model.current = len(tasks) // Индекс за пределами массива задач означает, что все задачи выполнены

	// Проверяем финальные представления задач
	assert.Equal(t, "Задача 1: Результат 1", task1.FinalView(80), "Финальное представление первой задачи должно содержать её результат")
	assert.Equal(t, "Задача 2: Результат 2", task2.FinalView(80), "Финальное представление второй задачи должно содержать её результат")

	// Теперь проверяем состояние завершения всех задач
	// Включаем сводку для этого теста
	model.WithSummary(true)
	// Обновляем статистику задач
	model.updateTaskStats()
	view := model.View()

	// Проверяем, что View содержит финальные представления задач и статистику
	assert.Contains(t, view, "Задача 1: Результат 1", "View должен содержать финальное представление первой задачи")
	assert.Contains(t, view, "Задача 2: Результат 2", "View должен содержать финальное представление второй задачи")
	assert.Contains(t, view, "Обработка операций прошла", "View должен содержать сводку")
	assert.Contains(t, view, "(2/2)", "View должен содержать статистику выполнения")
	assert.Contains(t, view, "УСПЕШНО", "View должен показывать статус УСПЕШНО")
}
