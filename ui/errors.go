package ui

import (
	"strings"
	"unicode/utf8"

	"github.com/qzeleza/termos/performance"
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

	// Константы для отступов при выводе сообщений
	const (
		rightMargin = 2
	)

	// Эффективная ширина для текста
	effectiveWidth := layoutWidth - rightMargin - 3
	if effectiveWidth <= 0 {
		return errMsg // Возвращаем как есть, если ширина недостаточна
	}

	// 1. Очищаем строку от управляющих символов и эскейп последовательностей
	cleanedMsg := cleanMessage(errMsg)

	if cleanedMsg == "" {
		return ""
	}
	numIndent := 2
	indent := performance.RepeatEfficient(" ", numIndent)

	// Создаем разделительную линию с учетом отступа
	// SeparatorLine := indent + BranchSymbol + DrawLine(effectiveWidth+rightMargin-numIndent*2+1)

	// Создаем форматированный результат
	errorMsg := formatErrorEveryLine(cleanedMsg, effectiveWidth-numIndent*2, indent)
	//  + GetTaskBelowPrefix()

	// удаляем крайний перенос строки
	// errorMsg = strings.TrimSuffix(errorMsg, "\n")

	return errorMsg
}

// formatErrorEveryLine создает отформатированное сообщение с разделительными линиями и отступами
func formatErrorEveryLine(msg string, effectiveWidth int, indent string) string {

	result := performance.GetBuffer()
	defer performance.PutBuffer(result)
	// Убираем перенос строки, чтобы сообщение выводилось сразу после линии

	// Если сообщение помещается в одну строку (считаем символы, а не байты)
	if utf8.RuneCountInString(msg) <= effectiveWidth {
		result.WriteString(indent + msg)
		result.WriteString("\n")
		return result.String()
	}

	// Делаем первую букву в предложении заглавной
	msg = CapitalizeFirst(msg)

	// Разбиваем на строки
	lines := wrapText(msg, effectiveWidth)

	for i, line := range lines {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(indent +
			VerticalLineSymbol +
			performance.RepeatEfficient(" ", 3) +
			GetErrorMessageStyle().Render(line))
	}

	// Заменяем все переносы строк более одного в конце строки result
	styledMsg := strings.TrimSuffix(result.String(), "\n")

	// Возвращаем отформатированное сообщение
	return styledMsg
}
