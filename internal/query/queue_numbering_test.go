package query

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qzeleza/termos/internal/common"
	"github.com/qzeleza/termos/internal/performance"
	"github.com/qzeleza/termos/internal/ui"
)

// numberedMockTask расширяет mockTask, добавляя поддержку установки префикса завершения.
type numberedMockTask struct {
	*mockTask
	prefix string
}

// newNumberedMockTask создаёт мок с поддержкой установки префикса завершения.
func newNumberedMockTask(title string) *numberedMockTask {
	return &numberedMockTask{mockTask: newMockTask(title)}
}

// SetCompletedPrefix переопределяет префикс завершения для мока.
func (t *numberedMockTask) SetCompletedPrefix(prefix string) {
	t.prefix = prefix
}

// FinalView возвращает отображение завершённой задачи.
func (t *numberedMockTask) FinalView(width int) string {
	prefix := t.prefix
	if prefix == "" {
		prefix = performance.FastConcat(
			performance.RepeatEfficient(" ", ui.MainLeftIndent),
			ui.TaskCompletedSymbol,
		)
	}
	return fmt.Sprintf("%s %s", prefix, t.title)
}

// prepareNumberedTasks помогает быстро собрать набор моков с фиксированным состоянием.
func prepareNumberedTasks(titles ...string) ([]common.Task, []*numberedMockTask) {
	var tasks []common.Task
	var mocks []*numberedMockTask
	for _, title := range titles {
		mt := newNumberedMockTask(title)
		mt.done = true
		tasks = append(tasks, mt)
		mocks = append(mocks, mt)
	}
	return tasks, mocks
}

// expectedLine возвращает строку представления задачи с учётом формата номера.
func expectedLine(n int, title string) string {
	return performance.FastConcat(
		performance.RepeatEfficient(" ", ui.MainLeftIndent-1),
		fmt.Sprintf(defaultNumberFormat, n),
		" ",
		title,
	)
}

func TestQueueNumberedCompletions(t *testing.T) {
	tasks, mocks := prepareNumberedTasks("Task 1", "Task 2", "Task 3")

	model := New("Нумерация задач")
	model.AddTasks(tasks)
	model.current = len(model.tasks)
	model.WithTasksNumbered(true, false, defaultNumberFormat)

	view := model.View()

	assert.Contains(t, view, expectedLine(1, "Task 1"), "Первая задача должна отображаться с номером 1")
	assert.Contains(t, view, expectedLine(2, "Task 2"), "Вторая задача должна отображаться с номером 2")
	assert.Contains(t, view, expectedLine(3, "Task 3"), "Третья задача должна отображаться с номером 3")

	assert.Equal(t, performance.FastConcat(performance.RepeatEfficient(" ", ui.MainLeftIndent-1), fmt.Sprintf(defaultNumberFormat, 1)), mocks[0].prefix)
	assert.Equal(t, performance.FastConcat(performance.RepeatEfficient(" ", ui.MainLeftIndent-1), fmt.Sprintf(defaultNumberFormat, 2)), mocks[1].prefix)
}

func TestQueueNumberedCompletionsKeepFirstSymbol(t *testing.T) {
	tasks, mocks := prepareNumberedTasks("Task A", "Task B")

	model := New("Нумерация с символом")
	model.AddTasks(tasks)
	model.current = len(model.tasks)
	model.WithTasksNumbered(true, true, defaultNumberFormat)

	view := model.View()

	assert.Contains(t, view, performance.FastConcat(
		performance.RepeatEfficient(" ", ui.MainLeftIndent),
		ui.TaskCompletedSymbol,
		" Task A",
	), "Первая задача должна сохранять символ завершения")
	assert.Contains(t, view, expectedLine(1, "Task B"), "Вторая задача должна отображаться с номером 1")

	assert.Equal(t, "", mocks[0].prefix, "Для первой задачи префикс должен оставаться стандартным")
	assert.Equal(t, performance.FastConcat(performance.RepeatEfficient(" ", ui.MainLeftIndent-1), fmt.Sprintf(defaultNumberFormat, 1)), mocks[1].prefix)
}
