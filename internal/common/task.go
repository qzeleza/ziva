// Package common содержит общие интерфейсы и константы для всего приложения.
package common

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Task представляет собой интерфейс для выполнения задач в очереди.
// Этот интерфейс используется как в пакете task, так и в пакете query.
type Task interface {
	// Title возвращает заголовок задачи.
	Title() string
	
	// Run запускает выполнение задачи и возвращает команду bubbletea.
	Run() tea.Cmd
	
	// Update обновляет состояние задачи на основе полученного сообщения.
	Update(msg tea.Msg) (Task, tea.Cmd)
	
	// View отображает текущее состояние задачи с учетом указанной ширины.
	View(width int) string
	
	// IsDone возвращает true, если задача завершена.
	IsDone() bool
	
	// FinalView отображает финальное состояние задачи с учетом указанной ширины.
	FinalView(width int) string
	
	// HasError возвращает true, если при выполнении задачи произошла ошибка.
	HasError() bool
	
	// Error возвращает ошибку, если она есть.
	Error() error
	
	// StopOnError возвращает true, если при возникновении ошибки в этой задаче
	// нужно остановить выполнение всей очереди задач.
	StopOnError() bool
	
	// SetStopOnError устанавливает флаг остановки очереди при ошибке.
	SetStopOnError(stop bool)
}