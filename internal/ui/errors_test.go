package ui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFormatErrorMessageWithPreserveNewLines проверяет функцию форматирования сообщений об ошибках
// с разными значениями параметра preserveNewLines
func TestFormatErrorMessageWithPreserveNewLines(t *testing.T) {
	tests := []struct {
		name             string
		errMsg           string
		layoutWidth      int
		preserveNewLines bool
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name:             "Сообщение с переносами строк (сохранение)",
			errMsg:           "Строка 1\nСтрока 2\nСтрока 3",
			layoutWidth:      50,
			preserveNewLines: true,
			shouldContain:    []string{"Строка 1", "Строка 2", "Строка 3"},
			shouldNotContain: []string{"Строка 1 Строка 2 Строка 3"},
		},
		{
			name:             "Сообщение с переносами строк (без сохранения)",
			errMsg:           "Строка 1\nСтрока 2\nСтрока 3",
			layoutWidth:      50,
			preserveNewLines: false,
			shouldContain:    []string{"Строка 1 Строка 2 Строка 3"},
			shouldNotContain: []string{},
		},
		{
			name:             "Сообщение с табуляцией и возвратом каретки (сохранение)",
			errMsg:           "Строка 1\tТабуляция\rВозврат",
			layoutWidth:      50,
			preserveNewLines: true,
			shouldContain:    []string{"Строка 1", "Табуляция", "Возврат"},
			shouldNotContain: []string{"Строка 1 Табуляция Возврат"},
		},
		{
			name:             "Сообщение с эскейп-последовательностями (сохранение)",
			errMsg:           "Строка\\nС\\tЭскейпами",
			layoutWidth:      50,
			preserveNewLines: true,
			shouldContain:    []string{"Строка", "С", "Эскейпами"},
			shouldNotContain: []string{},
		},
		{
			name:             "Пустое сообщение",
			errMsg:           "",
			layoutWidth:      50,
			preserveNewLines: true,
			shouldContain:    []string{},
			shouldNotContain: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatErrorMessage(tt.errMsg, tt.layoutWidth, tt.preserveNewLines)

			if tt.errMsg == "" {
				assert.Empty(t, result, "Результат должен быть пустым для пустого сообщения")
				return
			}

			// Проверяем наличие ожидаемых строк
			for _, expected := range tt.shouldContain {
				assert.Contains(t, result, expected, "Результат должен содержать: %s", expected)
			}

			// Проверяем отсутствие нежелательных строк
			for _, notExpected := range tt.shouldNotContain {
				assert.NotContains(t, result, notExpected, "Результат не должен содержать: %s", notExpected)
			}

			// Проверяем наличие отступов
			assert.True(t, strings.Contains(result, MessageIndent), "Результат должен содержать отступ")
		})
	}
}

// TestTaskErrorFormatting проверяет форматирование ошибок в задачах с разными настройками сохранения переносов строк
func TestTaskErrorFormatting(t *testing.T) {
	// Этот тест является заглушкой для будущей реализации
	// Здесь можно будет проверить, как работает форматирование ошибок в реальных задачах
	// с использованием метода PreserveErrorNewLines
	t.Skip("Этот тест будет реализован позже")
}
