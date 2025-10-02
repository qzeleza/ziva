package task

import (
	"testing"

	"github.com/qzeleza/ziva/internal/defaults"
	"github.com/stretchr/testify/assert"
)

func TestMenuKeywordDetectionByLanguage(t *testing.T) {
	original := defaults.CurrentLanguage()
	t.Cleanup(func() {
		defaults.SetLanguage(original)
	})

	defaults.SetLanguage("ru")
	ruChoices := normalizeItems([]Item{{Name: "Выйти"}, {Name: "Назад"}})
	if assert.Len(t, ruChoices, 2) {
		assert.True(t, isExitChoice(ruChoices[0]), "Русский пункт 'Выйти' должен определяться как выход")
		assert.True(t, isBackChoice(ruChoices[1]), "Русский пункт 'Назад' должен определяться как возврат")
	}

	defaults.SetLanguage("en")
	enChoices := normalizeItems([]Item{{Name: "Exit"}, {Name: "Return"}})
	if assert.Len(t, enChoices, 2) {
		assert.True(t, isExitChoice(enChoices[0]), "Английский пункт 'Exit' должен определяться как выход")
		assert.True(t, isBackChoice(enChoices[1]), "Английский пункт 'Return' должен определяться как возврат")
	}
}

func TestMenuKeywordDetectionIgnoresIrrelevantWords(t *testing.T) {
	original := defaults.CurrentLanguage()
	t.Cleanup(func() {
		defaults.SetLanguage(original)
	})

	defaults.SetLanguage("ru")
	choices := normalizeItems([]Item{{Name: "Выходной"}, {Name: "Назадовой"}})
	if assert.Len(t, choices, 2) {
		assert.False(t, isExitChoice(choices[0]), "Слово 'Выходной' не должно считаться пунктом выхода")
		assert.False(t, isBackChoice(choices[1]), "Слово 'Назадовой' не должно считаться пунктом возврата")
	}
}
