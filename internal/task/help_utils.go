package task

import (
	"strings"
	"sync"
)

var (
	choiceHelpDelimiter     = "::"
	choiceHelpDelimiterLock sync.RWMutex
)

// SetChoiceHelpDelimiter задаёт глобальный разделитель для встроенных подсказок в строках выбора.
// Пустой разделитель отключает разбор подсказок.
func SetChoiceHelpDelimiter(delim string) {
	choiceHelpDelimiterLock.Lock()
	defer choiceHelpDelimiterLock.Unlock()
	choiceHelpDelimiter = delim
}

func getChoiceHelpDelimiter() string {
	choiceHelpDelimiterLock.RLock()
	defer choiceHelpDelimiterLock.RUnlock()
	return choiceHelpDelimiter
}

func splitChoiceAndHelp(raw string) (string, string) {
	delim := getChoiceHelpDelimiter()
	if delim == "" {
		return raw, ""
	}
	parts := strings.SplitN(raw, delim, 2)
	if len(parts) != 2 {
		return strings.TrimSpace(raw), ""
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func parseChoicesWithHelp(raw []string) ([]string, []string) {
	if len(raw) == 0 {
		return nil, nil
	}
	labels := make([]string, len(raw))
	helps := make([]string, len(raw))
	for i, item := range raw {
		label, help := splitChoiceAndHelp(item)
		labels[i] = label
		helps[i] = help
	}
	return labels, helps
}
