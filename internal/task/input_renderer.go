package task

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/qzeleza/ziva/internal/defaults"
	terrors "github.com/qzeleza/ziva/internal/errors"
	"github.com/qzeleza/ziva/internal/performance"
	"github.com/qzeleza/ziva/internal/ui"
	"github.com/qzeleza/ziva/internal/validation"
)

// InputRenderer отвечает за отображение задач ввода
// Это отделяет UI логику от бизнес-логики задачи
type InputRenderer struct {
	style       lipgloss.Style
	helpEnabled bool
}

// NewInputRenderer создает новый рендерер для задач ввода
func NewInputRenderer() *InputRenderer {
	return &InputRenderer{
		style:       ui.InputStyle,
		helpEnabled: true,
	}
}

// WithStyle устанавливает пользовательский стиль
func (r *InputRenderer) WithStyle(style lipgloss.Style) *InputRenderer {
	r.style = style
	return r
}

// WithHelp включает или отключает отображение справки
func (r *InputRenderer) WithHelp(enabled bool) *InputRenderer {
	r.helpEnabled = enabled
	return r
}

// RenderInput отображает активное состояние задачи ввода с поддержкой таймера
func (r *InputRenderer) RenderInput(title string, textInput textinput.Model, validator validation.Validator, err error, inputType InputType, prefix string, width int, timerStr ...string) string {
	// Используем переданный префикс или значение по умолчанию
	if strings.TrimSpace(prefix) == "" {
		prefix = ui.GetCurrentTaskPrefix()
	}

	// Формируем заголовок с префиксом
	titleWithPrefix := fmt.Sprintf("%s %s", prefix, ui.ActiveTaskStyle.Render(title))

	// Если передан таймер, выравниваем его справа
	var titleView string
	if len(timerStr) > 0 && timerStr[0] != "" {
		timer := ui.SubtleStyle.Render(timerStr[0])
		titleView = ui.AlignTextToRight(titleWithPrefix, timer, width)
	} else {
		titleView = titleWithPrefix
	}

	// Получаем текст ввода с применением стиля
	inputView := r.style.Render(textInput.View())

	// Отображение ошибки валидации, если есть
	var errView string
	if err != nil {
		errIndent := performance.RepeatEfficient(" ", ui.MainLeftIndent)
		var errText string
		if te, ok := err.(*terrors.TaskError); ok && te.Err != nil {
			errText = ui.CapitalizeFirst(te.Err.Error())
			errText = fmt.Sprintf("%s%s", errIndent, errText)
		} else {
			errText = ui.CapitalizeFirst(err.Error())
			errText = fmt.Sprintf("%s%s", errIndent, errText)
		}
		errView = ui.GetErrorMessageStyle().Render(errText)
	}

	// Текст справки
	var helpText string
	if r.helpEnabled {
		helpIndent := performance.RepeatEfficient(" ", ui.MainLeftIndent)
		helpText = ui.SubtleStyle.Render(fmt.Sprintf("%s%s", helpIndent, defaults.InputConfirmHint))
	}

	// Подсказка о типе ввода
	var typeHint string
	if validator != nil {
		description := validator.Description()
		if description != "" {
			hintIndent := performance.RepeatEfficient(" ", ui.MainLeftIndent)
			typeHint = ui.SubtleStyle.Render(fmt.Sprintf("%s%s %s", hintIndent, defaults.InputFormatLabel, description))
		}
	}

	// Подсказка
	prompt := performance.FastConcat(
		performance.RepeatEfficient(" ", ui.MainLeftIndent),
		ui.CornerDownSymbol,
		ui.HorizontalLineSymbol,
	)

	// Собираем все вместе динамически
	var result strings.Builder

	// Основная часть
	result.WriteString(titleView)
	result.WriteString("\n")
	result.WriteString(prompt + inputView)
	result.WriteString("\n\n")
	lineWidth := width - 1
	if lineWidth < 0 {
		lineWidth = 0
	}
	result.WriteString(ui.DrawLine(lineWidth))

	// Один перевод после линии для первой дополнительной секции
	first := true

	appendSection := func(section string) {
		if section == "" {
			return
		}
		if !first {
			result.WriteString("\n")
		}
		result.WriteString(section)
		first = false
	}

	appendSection(errView)
	appendSection(typeHint)
	appendSection(helpText)

	return result.String()
}

// RenderFinal отображает финальное состояние задачи ввода
func (r *InputRenderer) RenderFinal(title string, value string, hasError bool, err error, prefix string, width int) string {
	var statusStyle lipgloss.Style
	var valueToShow string
	commentLines := make([]string, 0, 2)
	defaultIndent := ui.GetResultIndentWhenNumberingEnabled()

	buildCommentLine := func(indent, text string) string {
		return performance.FastConcat(
			performance.RepeatEfficient(" ", ui.MainLeftIndent),
			ui.VerticalLineSymbol,
			indent,
			text,
		)
	}

	// Если есть ошибка
	if hasError {
		statusStyle = ui.GetErrorStatusStyle()
		if err != nil {
			valueToShow = err.Error()
		} else {
			valueToShow = defaults.DefaultErrorLabel
		}

		if taskErr, ok := err.(*terrors.TaskError); ok && taskErr.Type == terrors.ErrorTypeUserCancel {
			statusStyle = ui.ErrorMessageStyle
			valueToShow = defaults.ErrorTypeUserCancel

			cancelMessage := defaults.ErrorMsgCanceled
			if taskErr.Err != nil {
				// Используем текст исходной ошибки, если он отличается от стандартного
				if msg := strings.TrimSpace(taskErr.Err.Error()); msg != "" {
					cancelMessage = msg
				}
			}
			cancelIndent := "   "
			commentLines = append(commentLines, buildCommentLine(cancelIndent, ui.ErrorMessageStyle.Render(cancelMessage)))
		}
	} else {
		statusStyle = ui.TaskStatusSuccessStyle
		valueToShow = strings.ToUpper(defaults.DefaultSuccessLabel)

		// Для паролей показываем звездочки вместо реального значения
		if r.looksLikePassword(title, value) {
			valueToShow = strings.Repeat("*", len(value))
		}
	}

	// Используем префикс, переданный очередью, либо значение по умолчанию
	if strings.TrimSpace(prefix) == "" {
		prefix = ui.GetCompletedInputTaskPrefix(!hasError)
	}
	leftPart := fmt.Sprintf("%s  %s", prefix, title)
	rightPart := statusStyle.Render(valueToShow)

	if value != "" {
		commentLines = append(commentLines, buildCommentLine(defaultIndent, ui.SubtleStyle.Render(value)))
	}

	var resultBuilder strings.Builder
	resultBuilder.WriteString(ui.AlignTextToRight(leftPart, rightPart, width))
	if len(commentLines) > 0 {
		resultBuilder.WriteString("\n")
		for i, line := range commentLines {
			if i > 0 {
				resultBuilder.WriteString("\n")
			}
			resultBuilder.WriteString(line)
		}
		resultBuilder.WriteString("\n")
	}
	return resultBuilder.String()
}

// looksLikePassword определяет, является ли поле паролем по заголовку
func (r *InputRenderer) looksLikePassword(title, value string) bool {
	lowerTitle := strings.ToLower(title)
	// Переводы и общеупотребимые варианты на поддерживаемых языках (ru, en, tr, be, uk):
	// ru/be/uk: "пароль", "ключ"
	// en:       "password", "pass", "pwd", "key"
	// tr:       "şifre", "sifre", "parola", "anahtar"
	passwordKeywords := []string{
		// ru / be / uk
		"пароль", "ключ",
		// en
		"password", "pass", "pwd", "key",
		// tr
		"şifre", "sifre", "parola", "anahtar",
	}

	for _, keyword := range passwordKeywords {
		if strings.Contains(lowerTitle, keyword) {
			return true
		}
	}
	return false
}

// InputTypeHints предоставляет подсказки для различных типов ввода
var InputTypeHints = map[InputType]string{
	InputTypeText:     "",
	InputTypePassword: defaults.InputHintPassword,
	InputTypeEmail:    defaults.InputHintEmail,
	InputTypeNumber:   defaults.InputHintNumber,
	InputTypeIP:       defaults.InputHintIP,
	InputTypeDomain:   defaults.InputHintDomain,
}

// GetTypeHint возвращает подсказку для типа ввода
func GetTypeHint(inputType InputType) string {
	if hint, exists := InputTypeHints[inputType]; exists {
		return hint
	}
	return ""
}
