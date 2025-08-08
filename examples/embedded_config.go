// Package examples - конфигурация для embedded устройств
package examples

import (
	"github.com/qzeleza/termos/internal/ui"
)

// EmbeddedConfig содержит настройки оптимизации для embedded устройств
type EmbeddedConfig struct {
	// Ограничения памяти
	MaxCompletedTasks  int // Максимальное количество завершенных задач в памяти
	MemoryThresholdMB  int // Порог для запуска очистки памяти (в МБ)
	StringPoolSize     int // Размер пула строк
	ByteBufferPoolSize int // Размер пула байтовых буферов

	// Цветовые настройки
	UseSimpleColors bool // Использовать упрощенную ANSI палитру
	MaxColors       int  // Максимальное количество одновременно используемых цветов

	// Битсет оптимизации
	Use32BitBitset    bool // Использовать 32-битный битсет вместо 64-битного
	MaxBitsetElements int  // Максимальное количество элементов в битсете

	// Интернирование строк
	StringCacheSize    int  // Размер кэша для интернирования строк
	EnableStringIntern bool // Включить интернирование часто используемых строк
}

// DefaultEmbeddedConfig возвращает стандартную конфигурацию для embedded устройств
func DefaultEmbeddedConfig() *EmbeddedConfig {
	return &EmbeddedConfig{
		// Ограничения памяти для устройств с 128MB RAM
		MaxCompletedTasks:  50,
		MemoryThresholdMB:  64, // 64MB порог
		StringPoolSize:     4,  // Уменьшенный размер пула
		ByteBufferPoolSize: 2,  // Минимальный размер

		// Простые цвета для максимальной совместимости
		UseSimpleColors: true,
		MaxColors:       16, // Только стандартные ANSI цвета

		// 32-битная оптимизация
		Use32BitBitset:    true,
		MaxBitsetElements: 32,

		// Интернирование строк
		StringCacheSize:    64, // Уменьшенный кэш для embedded
		EnableStringIntern: true,
	}
}

// OptimizedEmbeddedConfig возвращает максимально оптимизированную конфигурацию
// для устройств с очень ограниченными ресурсами (например, 64MB RAM)
func OptimizedEmbeddedConfig() *EmbeddedConfig {
	return &EmbeddedConfig{
		MaxCompletedTasks:  25, // Еще меньше задач в памяти
		MemoryThresholdMB:  32, // Более агрессивная очистка
		StringPoolSize:     2,
		ByteBufferPoolSize: 1,

		UseSimpleColors: true,
		MaxColors:       8, // Только базовые цвета

		Use32BitBitset:    true,
		MaxBitsetElements: 16, // Еще более ограниченный битсет

		StringCacheSize:    32, // Минимальный кэш строк
		EnableStringIntern: true,
	}
}

// ApplyEmbeddedConfig применяет embedded конфигурацию к модулю Термос
func ApplyEmbeddedConfig(config *EmbeddedConfig) {
	if config == nil {
		config = DefaultEmbeddedConfig()
	}

	// Применяем цветовые настройки
	if config.UseSimpleColors {
		ui.EnableEmbeddedMode()
	}

	// Настраиваем интернирование строк
	if config.EnableStringIntern {
		// Кэш уже инициализирован при импорте ui пакета
		// Здесь мы могли бы настроить его размер, если бы была такая функция
	}

	// Для пулов буферов нужно было бы пересоздать их с новыми размерами,
	// но в текущей реализации они инициализированы как глобальные переменные
	// В производственном коде стоило бы сделать их настраиваемыми
}

// GetMemoryFootprintEstimate возвращает оценку потребления памяти
func GetMemoryFootprintEstimate(config *EmbeddedConfig) (estimateMB int) {
	if config == nil {
		config = DefaultEmbeddedConfig()
	}

	estimate := 0

	// Базовый модуль: ~1-2MB с оптимизациями
	estimate += 1500 // 1.5MB в килобайтах

	// Задачи в памяти: ~0.5KB на задачу с оптимизацией битсета
	estimate += config.MaxCompletedTasks / 2 // в KB

	// Буферные пулы: зависят от размера
	estimate += config.StringPoolSize * 4     // ~4KB на буфер строк
	estimate += config.ByteBufferPoolSize * 2 // ~2KB на байтовый буфер

	// Кэш интернирования строк
	if config.EnableStringIntern {
		estimate += config.StringCacheSize / 4 // примерно 4 строки на KB
	}

	// Цвета: ANSI vs Hex
	if config.UseSimpleColors {
		estimate += 1 // ~1KB для ANSI цветов
	} else {
		estimate += 4 // ~4KB для hex цветов и дополнительных структур
	}

	// Конвертируем в мегабайты
	return estimate / 1024
}

// ValidateEmbeddedConfig проверяет корректность конфигурации
func ValidateEmbeddedConfig(config *EmbeddedConfig) error {
	if config.MaxCompletedTasks < 5 {
		config.MaxCompletedTasks = 5
	}

	if config.MemoryThresholdMB < 16 {
		config.MemoryThresholdMB = 16
	}

	if config.StringPoolSize < 1 {
		config.StringPoolSize = 1
	}

	if config.ByteBufferPoolSize < 1 {
		config.ByteBufferPoolSize = 1
	}

	if config.MaxBitsetElements < 8 {
		config.MaxBitsetElements = 8
	}

	if config.StringCacheSize < 16 {
		config.StringCacheSize = 16
	}

	return nil
}
