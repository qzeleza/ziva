package task

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TimeoutMsg - специальное сообщение, которое отправляется при истечении времени ожидания
type TimeoutMsg struct{}

// TickMsg - сообщение для периодического обновления счетчика
type TickMsg struct{}

// DefaultTimeout - сообщение таймаута по умолчанию
var DefaultTimeout = TimeoutMsg{}

// DefaultTick - сообщение обновления по умолчанию
var DefaultTick = TickMsg{}

// TimeoutManager - менеджер тайм-аутов для задач
type TimeoutManager struct {
	// Длительность тайм-аута
	duration time.Duration
	// Время начала тайм-аута
	startTime time.Time
	// Флаг активности тайм-аута
	active bool
	// Канал для отмены текущего таймера
	cancelCh chan struct{}
}

// NewTimeoutManager создает новый менеджер тайм-аутов
func NewTimeoutManager(duration time.Duration) *TimeoutManager {
	return &TimeoutManager{
		duration: duration,
		active:   false,
		cancelCh: make(chan struct{}, 1),
	}
}

// StartTimeout запускает таймер тайм-аута
func (tm *TimeoutManager) StartTimeout() tea.Cmd {
	// Если таймер уже активен, сначала остановим его
	if tm.active {
		tm.StopTimeout()
	}

	tm.startTime = time.Now()
	tm.active = true

	// Создаем новый канал для отмены
	tm.cancelCh = make(chan struct{}, 1)
	cancelCh := tm.cancelCh

	// Возвращаем команду для таймера
	return func() tea.Msg {
		timer := time.NewTimer(tm.duration)
		select {
		case <-timer.C:
			// Таймер сработал - возвращаем сообщение о тайм-ауте
			return DefaultTimeout
		case <-cancelCh:
			// Таймер был отменен
			timer.Stop()
			return nil
		}
	}
}

// StartTicker запускает периодическое обновление счетчика каждую секунду
func (tm *TimeoutManager) StartTicker() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return DefaultTick
	})
}

// StartTickerAndTimeout запускает одновременно и таймер, и тикер
func (tm *TimeoutManager) StartTickerAndTimeout() tea.Cmd {
	return tea.Batch(
		tm.StartTimeout(),
		tm.StartTicker(),
	)
}

// StopTimeout останавливает текущий таймер, если он активен
func (tm *TimeoutManager) StopTimeout() {
	if tm.active && tm.cancelCh != nil {
		// Отправляем сигнал отмены таймера
		select {
		case tm.cancelCh <- struct{}{}:
		default:
		}
		tm.active = false
	}
}

// IsActive возвращает true, если таймер активен
func (tm *TimeoutManager) IsActive() bool {
	return tm.active
}

// RemainingTime возвращает оставшееся время в секундах
func (tm *TimeoutManager) RemainingTime() int {
	if !tm.active {
		return 0
	}

	elapsed := time.Since(tm.startTime)
	remaining := tm.duration - elapsed
	if remaining < 0 {
		return 0
	}
	return int(remaining.Seconds())
}

// RemainingTimeFormatted возвращает оставшееся время в формате MM:SS
func (tm *TimeoutManager) RemainingTimeFormatted() string {
	seconds := tm.RemainingTime()
	minutes := seconds / 60
	seconds = seconds % 60
	return time.Time{}.Add(time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second).Format("04:05")
}

// SetDuration устанавливает новую длительность тайм-аута
// Не влияет на текущий активный таймер
func (tm *TimeoutManager) SetDuration(duration time.Duration) {
	tm.duration = duration
}

// GetDuration возвращает текущую длительность тайм-аута
func (tm *TimeoutManager) GetDuration() time.Duration {
	return tm.duration
}
