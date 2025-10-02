package task

import (
	"strings"
	"unicode"

	"github.com/qzeleza/ziva/internal/defaults"
)

// isExitChoice проверяет, является ли выбор пунктом выхода.
// @param c - выбор
//
// @return true, если выбор является пунктом выхода
func isExitChoice(c choice) bool {
	return matchesMenuKeywords(c, defaults.MenuExitKeywords)
}

// isBackChoice проверяет, является ли выбор пунктом возврата.
// @param c - выбор
//
// @return true, если выбор является пунктом возврата
func isBackChoice(c choice) bool {
	return matchesMenuKeywords(c, defaults.MenuBackKeywords)
}

// matchesMenuKeywords проверяет, соответствует ли выбор ключевым словам.
// @param c - выбор
// @param keywords - ключевые слова
//
// @return true, если выбор соответствует ключевым словам
func matchesMenuKeywords(c choice, keywords []string) bool {
	if len(keywords) == 0 {
		return false
	}

	variants := []string{
		c.displayName(),
		c.valueKey(),
	}

	for _, variant := range variants {
		normalized := normalizeMenuCandidate(variant)
		if normalized == "" {
			continue
		}
		for _, keyword := range keywords {
			if keywordMatches(normalized, keyword) {
				return true
			}
		}
	}
	return false
}

// normalizeMenuCandidate нормализует строку для сравнения с ключевыми словами.
// @param value - строка для нормализации
//
// @return нормализованная строка
func normalizeMenuCandidate(value string) string {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	if trimmed == "" {
		return ""
	}
	trimmed = strings.Trim(trimmed, " .:;!?…\"'«»")
	return trimmed
}

// keywordMatches проверяет, соответствует ли строка ключевому слову.
// @param value - строка для сравнения
// @param keyword - ключевое слово
//
// @return true, если строка соответствует ключевому слову
func keywordMatches(value, keyword string) bool {
	normKeyword := normalizeMenuCandidate(keyword)
	if normKeyword == "" {
		return false
	}
	if value == normKeyword {
		return true
	}
	if strings.Contains(normKeyword, " ") {
		return strings.Contains(value, normKeyword)
	}

	tokens := strings.FieldsFunc(value, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	for _, token := range tokens {
		if token == normKeyword {
			return true
		}
	}
	return false
}
