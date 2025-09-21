package termos

import "github.com/qzeleza/termos/internal/defaults"

// SetLanguage задаёт язык интерфейса и возвращает фактически установленное значение (например, "ru" или "en").
func SetLanguage(lang string) string {
	return defaults.SetLanguage(lang)
}

// SetDefaultLanguage задаёт язык по умолчанию, используемый при инициализации пакета.
func SetDefaultLanguage(lang string) string {
	return defaults.SetDefaultLanguage(lang)
}

// CurrentLanguage возвращает текущий код языка интерфейса.
func CurrentLanguage() string {
	return defaults.CurrentLanguage()
}

// SupportedLanguages возвращает список поддерживаемых кодов языков.
func SupportedLanguages() []string {
	return defaults.SupportedLanguages()
}
