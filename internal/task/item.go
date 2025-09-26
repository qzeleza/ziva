package task

import (
	"fmt"
	"strings"
)

/// Item описывает элемент списка для задач выбора.
type Item struct {
	Key         string
	Name        string
	Description string
}

type choice struct {
	key         string
	name        string
	description string
}

func (c choice) displayName() string {
	return c.name
}

func (c choice) helpText() string {
	return c.description
}

func (c choice) valueKey() string {
	return c.key
}

/// normalizeItems подготавливает элементы к использованию внутри задач.
func normalizeItems(source []Item) []choice {
	normalized := make([]choice, len(source))
	for i, it := range source {
		key := strings.TrimSpace(it.Key)
		name := strings.TrimSpace(it.Name)
		if key == "" && name != "" {
			key = name
		}
		if name == "" && key != "" {
			name = key
		}
		if key == "" && name == "" {
			key = fmt.Sprintf("item_%d", i+1)
			name = key
		}
		desc := strings.TrimSpace(it.Description)
		normalized[i] = choice{
			key:         key,
			name:        name,
			description: desc,
		}
	}
	return normalized
}
