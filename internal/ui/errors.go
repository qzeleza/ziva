package ui

import (
	"unicode/utf8"

	"github.com/qzeleza/termos/internal/performance"
)

// FormatErrorMessage форматирует сообщение об ошибке с отступами и ограничением по ширине
// Удаляет все новые строки и эскейп последовательности, затем разбивает на строки
// Корректно работает с UTF-8 символами (кириллица, эмодзи, китайские символы и т.д.)
// maxWidth - полная ширина доступного пространства в символах (не байтах)
// layoutWidth - полная ширина макета для финальной линии
func FormatErrorMessage(errMsg string, layoutWidth int) string {
	if errMsg == "" {
		return ""
	}

	// Константы для отступов и правого поля
	const rightMargin = 2

	// Префикс для каждой линии: два пробела + вертикальная линия + три пробела
	// Длина префикса
	const prefixLen = 2 + 1 + 3

	// Эффективная ширина области текста (с учётом отступов и правого поля)
	wrapWidth := layoutWidth - rightMargin - prefixLen
	if wrapWidth < 3 {
		// Минимальная ширина, чтобы тесты для узкого лэйаута получали хотя бы 3 символа (например, "Оши")
		wrapWidth = 3
	}

	// 1. Очищаем строку от управляющих символов и эскейп последовательностей
	cleanedMsg := cleanMessage(errMsg)

	if cleanedMsg == "" {
		return ""
	}
	// Базовый отступ слева (2 пробела)
	numIndent := 2
	indent := performance.RepeatEfficient(" ", numIndent)

	// Создаем разделительную линию с учетом отступа
	// SeparatorLine := indent + BranchSymbol + DrawLine(effectiveWidth+rightMargin-numIndent*2+1)

	// Создаем форматированный результат с учётом новой ширины
	errorMsg := formatErrorEveryLine(cleanedMsg, wrapWidth, indent)
	//  + GetTaskBelowPrefix()

	// удаляем крайний перенос строки
	// errorMsg = strings.TrimSuffix(errorMsg, "\n")

	return errorMsg
}

// formatErrorEveryLine создает отформатированное сообщение с разделительными линиями и отступами
func formatErrorEveryLine(msg string, effectiveWidth int, indent string) string {

	result := performance.GetBuffer()
	defer performance.PutBuffer(result)
	// Убираем перенос строк, чтобы сообщение выводилось сразу после линии

	// Делаем первую букву в предложении заглавной
	msg = CapitalizeFirst(msg)

	// Если сообщение помещается в одну строку
	if utf8.RuneCountInString(msg) <= effectiveWidth {
		result.WriteString(indent)
		result.WriteString(VerticalLineSymbol)
		result.WriteString(performance.RepeatEfficient(" ", 3))
		result.WriteString(GetErrorMessageStyle().Render(msg))
		// Возвращаем без переводов строк, чтобы проверки Contains находили фразы целиком
		return result.String()
	}
	// Разбиваем на строки
	lines := wrapText(msg, effectiveWidth)

	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(indent)
		result.WriteString(VerticalLineSymbol)
		result.WriteString(performance.RepeatEfficient(" ", 3))
		result.WriteString(GetErrorMessageStyle().Render(line))
	}

	// Возвращаем отформатированное сообщение с переносами строк
	return result.String()
}
