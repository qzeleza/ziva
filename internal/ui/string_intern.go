// Package ui - оптимизация строк для embedded устройств
package ui

// stringInterningCache кэш интернирования строк для частых UI элементов
// Ограничен размером для embedded устройств
var stringInterningCache = make(map[string]string, 128)

// Часто используемые строки заранее интернированы
var (
	InternedStrings = struct {
		// Префиксы задач
		TaskCompleted  string
		TaskInProgress string
		TaskCancelled  string
		TaskError      string

		// Общие UI элементы
		Selected    string
		NotSelected string
		Loading     string
		Done        string
		Cancel      string
		Error       string
		Success     string

		// Сообщения
		PressEnter   string
		PressSpace   string
		PressQ       string
		NavigateKeys string

		// Индентация (предвычислено)
		Indent2 string
		Indent4 string
		Indent6 string
		Indent8 string

		// Символы
		VerticalLine   string
		HorizontalLine string
		Corner         string
		Branch         string
	}{
		// Инициализируем частые строки
		TaskCompleted:  intern("✓"),
		TaskInProgress: intern("○"),
		TaskCancelled:  intern("✗"),
		TaskError:      intern("!"),

		Selected:    intern("●"),
		NotSelected: intern(" "),
		Loading:     intern("..."),
		Done:        intern("ГОТОВО"),
		Cancel:      intern("ОТМЕНА"),
		Error:       intern("ОШИБКА"),
		Success:     intern("УСПЕШНО"),

		PressEnter:   intern("Enter - подтвердить"),
		PressSpace:   intern("Space - выбрать/снять"),
		PressQ:       intern("Q - отмена"),
		NavigateKeys: intern("↑/↓ - навигация"),

		Indent2: intern("  "),
		Indent4: intern("    "),
		Indent6: intern("      "),
		Indent8: intern("        "),

		VerticalLine:   intern("│"),
		HorizontalLine: intern("─"),
		Corner:         intern("└"),
		Branch:         intern("├"),
	}
)

// intern интернирует строку для экономии памяти на embedded устройствах
func intern(s string) string {
	if cached, exists := stringInterningCache[s]; exists {
		return cached
	}

	// Ограничиваем размер кэша для embedded устройств
	if len(stringInterningCache) >= 128 {
		// При переполнении очищаем старые записи (простая стратегия)
		clearOldEntries()
	}

	stringInterningCache[s] = s
	return s
}

// InternString публичная функция для интернирования пользовательских строк
func InternString(s string) string {
	return intern(s)
}

// clearOldEntries очищает половину записей при переполнении кэша
func clearOldEntries() {
	count := 0
	target := len(stringInterningCache) / 2

	for key := range stringInterningCache {
		if count >= target {
			break
		}
		delete(stringInterningCache, key)
		count++
	}
}

// GetCacheStats возвращает статистику кэша интернирования для мониторинга
func GetCacheStats() (size int, capacity int) {
	return len(stringInterningCache), 128
}

// ClearInternCache полностью очищает кэш (для emergency cleanup)
func ClearInternCache() {
	// Сохраняем только предвычисленные строки
	newCache := make(map[string]string, 128)

	// Восстанавливаем критичные интернированные строки
	criticalStrings := []string{
		"✓", "○", "✗", "!", "●", " ", "...", "ГОТОВО", "ОТМЕНА", "ОШИБКА", "УСПЕШНО",
		"Enter - подтвердить", "Space - выбрать/снять", "Q - отмена", "↑/↓ - навигация",
		"  ", "    ", "      ", "        ", "│", "─", "└", "├",
	}

	for _, s := range criticalStrings {
		newCache[s] = s
	}

	stringInterningCache = newCache
}
