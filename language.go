package termos

import "github.com/qzeleza/termos/internal/defauilt"

// SetLanguage задаёт язык интерфейса и возвращает фактически установленное значение (например, "ru" или "en").
func SetLanguage(lang string) string {
	return defauilt.SetLanguage(lang)
}

// SetDefaultLanguage задаёт язык по умолчанию, используемый при инициализации пакета.
func SetDefaultLanguage(lang string) string {
	return defauilt.SetDefaultLanguage(lang)
}

// CurrentLanguage возвращает текущий код языка интерфейса.
func CurrentLanguage() string {
	return defauilt.CurrentLanguage()
}

// SupportedLanguages возвращает список поддерживаемых кодов языков.
func SupportedLanguages() []string {
	return defauilt.SupportedLanguages()
}
