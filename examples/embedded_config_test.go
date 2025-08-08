package examples

import (
	"testing"

	"github.com/qzeleza/termos/ui"
	"github.com/stretchr/testify/assert"
)

func TestDefaultEmbeddedConfig(t *testing.T) {
	config := DefaultEmbeddedConfig()

	assert.NotNil(t, config, "DefaultEmbeddedConfig должен возвращать конфигурацию")
	assert.Equal(t, 50, config.MaxCompletedTasks, "MaxCompletedTasks должен быть 50")
	assert.Equal(t, 64, config.MemoryThresholdMB, "MemoryThresholdMB должен быть 64")
	assert.Equal(t, 4, config.StringPoolSize, "StringPoolSize должен быть 4")
	assert.Equal(t, 2, config.ByteBufferPoolSize, "ByteBufferPoolSize должен быть 2")
	assert.True(t, config.UseSimpleColors, "UseSimpleColors должен быть true")
	assert.Equal(t, 16, config.MaxColors, "MaxColors должен быть 16")
	assert.True(t, config.Use32BitBitset, "Use32BitBitset должен быть true")
	assert.Equal(t, 32, config.MaxBitsetElements, "MaxBitsetElements должен быть 32")
	assert.Equal(t, 64, config.StringCacheSize, "StringCacheSize должен быть 64")
	assert.True(t, config.EnableStringIntern, "EnableStringIntern должен быть true")
}

func TestOptimizedEmbeddedConfig(t *testing.T) {
	config := OptimizedEmbeddedConfig()

	assert.NotNil(t, config, "OptimizedEmbeddedConfig должен возвращать конфигурацию")
	assert.Equal(t, 25, config.MaxCompletedTasks, "Оптимизированная конфигурация должна иметь меньше задач")
	assert.Equal(t, 32, config.MemoryThresholdMB, "Оптимизированная конфигурация должна иметь меньший порог памяти")
	assert.Equal(t, 8, config.MaxColors, "Оптимизированная конфигурация должна использовать меньше цветов")
	assert.Equal(t, 16, config.MaxBitsetElements, "Оптимизированная конфигурация должна иметь меньший битсет")
	assert.Equal(t, 32, config.StringCacheSize, "Оптимизированная конфигурация должна иметь меньший кэш строк")
}

func TestApplyEmbeddedConfig(t *testing.T) {
	// Сохраняем исходное состояние
	originalColorMode := ui.IsEmbeddedColorMode()

	// Тестируем применение конфигурации
	config := DefaultEmbeddedConfig()
	ApplyEmbeddedConfig(config)

	// Проверяем, что embedded режим включился
	assert.True(t, ui.IsEmbeddedColorMode(), "Embedded режим должен быть включен")

	// Тестируем применение с nil конфигурацией
	ApplyEmbeddedConfig(nil)
	// Должно работать без паники

	// Восстанавливаем состояние если нужно
	if !originalColorMode {
		// При необходимости можно добавить функцию отключения embedded режима
	}
}

func TestGetMemoryFootprintEstimate(t *testing.T) {
	// Тестируем с дефолтной конфигурацией
	config := DefaultEmbeddedConfig()
	estimate := GetMemoryFootprintEstimate(config)

	assert.Greater(t, estimate, 0, "Оценка потребления памяти должна быть больше 0")
	assert.Less(t, estimate, 100, "Оценка не должна быть чрезмерно большой")

	// Тестируем с оптимизированной конфигурацией
	optimizedConfig := OptimizedEmbeddedConfig()
	optimizedEstimate := GetMemoryFootprintEstimate(optimizedConfig)

	assert.Less(t, optimizedEstimate, estimate, "Оптимизированная конфигурация должна потреблять меньше памяти")

	// Тестируем с nil
	nilEstimate := GetMemoryFootprintEstimate(nil)
	assert.Equal(t, estimate, nilEstimate, "nil конфигурация должна использовать дефолтную")
}

func TestValidateEmbeddedConfig(t *testing.T) {
	// Создаем конфигурацию с некорректными значениями
	config := &EmbeddedConfig{
		MaxCompletedTasks:  1, // Слишком мало
		MemoryThresholdMB:  5, // Слишком мало
		StringPoolSize:     0, // Неверное значение
		ByteBufferPoolSize: 0, // Неверное значение
		MaxBitsetElements:  3, // Слишком мало
		StringCacheSize:    5, // Слишком мало
	}

	err := ValidateEmbeddedConfig(config)
	assert.NoError(t, err, "ValidateEmbeddedConfig не должен возвращать ошибку")

	// Проверяем, что значения исправлены
	assert.Equal(t, 5, config.MaxCompletedTasks, "MaxCompletedTasks должен быть исправлен до минимального значения")
	assert.Equal(t, 16, config.MemoryThresholdMB, "MemoryThresholdMB должен быть исправлен")
	assert.Equal(t, 1, config.StringPoolSize, "StringPoolSize должен быть исправлен до 1")
	assert.Equal(t, 1, config.ByteBufferPoolSize, "ByteBufferPoolSize должен быть исправлен до 1")
	assert.Equal(t, 8, config.MaxBitsetElements, "MaxBitsetElements должен быть исправлен")
	assert.Equal(t, 16, config.StringCacheSize, "StringCacheSize должен быть исправлен")
}

func TestEmbeddedConfigComparison(t *testing.T) {
	defaultConfig := DefaultEmbeddedConfig()
	optimizedConfig := OptimizedEmbeddedConfig()

	// Оптимизированная конфигурация должна иметь меньшие значения для экономии памяти
	assert.Less(t, optimizedConfig.MaxCompletedTasks, defaultConfig.MaxCompletedTasks)
	assert.Less(t, optimizedConfig.MemoryThresholdMB, defaultConfig.MemoryThresholdMB)
	assert.Less(t, optimizedConfig.StringPoolSize, defaultConfig.StringPoolSize)
	assert.Less(t, optimizedConfig.ByteBufferPoolSize, defaultConfig.ByteBufferPoolSize)
	assert.Less(t, optimizedConfig.MaxColors, defaultConfig.MaxColors)
	assert.Less(t, optimizedConfig.MaxBitsetElements, defaultConfig.MaxBitsetElements)
	assert.Less(t, optimizedConfig.StringCacheSize, defaultConfig.StringCacheSize)

	// Общие настройки должны быть одинаковыми
	assert.Equal(t, defaultConfig.UseSimpleColors, optimizedConfig.UseSimpleColors)
	assert.Equal(t, defaultConfig.Use32BitBitset, optimizedConfig.Use32BitBitset)
	assert.Equal(t, defaultConfig.EnableStringIntern, optimizedConfig.EnableStringIntern)
}

func TestMemoryFootprintCalculation(t *testing.T) {
	// Создаем конфигурацию с известными значениями для проверки расчетов
	config := &EmbeddedConfig{
		MaxCompletedTasks:  10,
		StringPoolSize:     2,
		ByteBufferPoolSize: 1,
		StringCacheSize:    32,
		EnableStringIntern: true,
		UseSimpleColors:    true,
	}

	estimate := GetMemoryFootprintEstimate(config)

	// Проверяем, что оценка разумная (базовые 1.5MB + дополнительные компоненты)
	assert.Greater(t, estimate, 1, "Оценка должна быть больше 1MB")
	assert.Less(t, estimate, 10, "Оценка должна быть меньше 10MB для embedded конфигурации")
}

func TestEmbeddedConfigEdgeCases(t *testing.T) {
	// Тестируем экстремально маленькие значения
	tinyConfig := &EmbeddedConfig{
		MaxCompletedTasks:  0,
		MemoryThresholdMB:  0,
		StringPoolSize:     -1,
		ByteBufferPoolSize: -1,
		MaxBitsetElements:  0,
		StringCacheSize:    0,
	}

	ValidateEmbeddedConfig(tinyConfig)

	// Все значения должны быть исправлены до минимальных
	assert.GreaterOrEqual(t, tinyConfig.MaxCompletedTasks, 5)
	assert.GreaterOrEqual(t, tinyConfig.MemoryThresholdMB, 16)
	assert.GreaterOrEqual(t, tinyConfig.StringPoolSize, 1)
	assert.GreaterOrEqual(t, tinyConfig.ByteBufferPoolSize, 1)
	assert.GreaterOrEqual(t, tinyConfig.MaxBitsetElements, 8)
	assert.GreaterOrEqual(t, tinyConfig.StringCacheSize, 16)
}
