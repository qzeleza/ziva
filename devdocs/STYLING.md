# Руководство по стилизации Жива

Настраивайте внешний вид TUI с помощью экспортируемых стилей, цветов, иконок и хелперов из `ui/styles.go`. Это руководство отражает актуальный API.

## Доступные стили и цвета

Из `ui/styles.go` (выборка):
- Цвета: `ui.ColorBrightGreen`, `ui.ColorBrightRed`, `ui.ColorBrightYellow`, `ui.ColorLightBlue`, `ui.ColorBrightWhite`, `ui.ColorBrightGray`, `ui.ColorDarkGray` и др.
- Базовые текстовые стили: `ui.TitleStyle`, `ui.ActiveTitleStyle`, `ui.ActiveTaskStyle`, `ui.SuccessLabelStyle`, `ui.FinishedLabelStyle`
- Стили ошибок: `ui.ErrorMessageStyle`, `ui.ErrorStatusStyle`, `ui.CancelStyle`
- Стили выбора/активного элемента: `ui.SelectionStyle` (Да), `ui.SelectionNoStyle` (Нет), `ui.ActiveStyle`, `ui.InputStyle`, `ui.SpinnerStyle`
- Иконки: `ui.IconDone`, `ui.IconError`, `ui.IconCancelled`, `ui.IconQuestion`, `ui.IconSelected`, `ui.IconRadioOn`, `ui.IconRadioOff`, `ui.IconCursor`, `ui.IconUndone`

Вы можете изменять эти переменные пакета при старте приложения, чтобы подстроить тему.

```go
import (
    "github.com/charmbracelet/lipgloss"
    "github.com/qzeleza/ziva/ui"
)

func init() {
    ui.TitleStyle = ui.TitleStyle.Foreground(ui.ColorBrightGreen).Bold(true)
    ui.SelectionNoStyle = ui.SelectionNoStyle.Foreground(ui.ColorBrightRed).Bold(true)
    ui.ActiveStyle = ui.ActiveStyle.Foreground(ui.ColorLightBlue).Bold(true)
}
```

## Помощники для установки цвета для вывода ошибок

```go
// Переопределение цветов ошибок во время выполнения
ui.SetErrorColor(ui.ColorDarkYellow, ui.ColorBrightYellow)
// ... позже, вернуть значения по умолчанию
ui.ResetErrorColors()
```

## Префиксы задач и визуальное дерево задач

Используйте хелперы префиксов для построения согласованных «деревьев» задач:
- `ui.GetActiveTaskPrefix()`
- `ui.GetTaskBelowPrefix()`
- `ui.GetCompletedTaskPrefix(success bool)`
- `ui.GetCompletedInputTaskPrefix(success bool)`

Пример:

```go
prefix := ui.GetActiveTaskPrefix()
line := prefix + ui.ActiveTitleStyle.Render("Downloading")
```

## Утилиты разметки и текста

- Расчёт ширины: `common.CalculateLayoutWidth(screenWidth int) int`
- Значение по умолчанию ширины: `common.DefaultWidth`
- Выравнивание: `ui.AlignTextToRight(left, right, width int) string`
- Форматирование ошибок с переносами/отступами: `ui.FormatErrorMessage(text string, width int) string`
- ANSI/Unicode-хелперы: `ui.GetPlainTextLength(text)`, `ui.StripANSI(text)`, `ui.WrapText(text, width)`, `ui.GetRuneWidth(r rune)`

```go
import (
    "github.com/qzeleza/ziva/common"
    "github.com/qzeleza/ziva/ui"
)

width := common.CalculateLayoutWidth(120)
left := ui.TitleStyle.Render("Установка")
right := ui.SuccessLabelStyle.Render("Готово")
line := ui.AlignTextToRight(left, right, width)
```

## Рекомендации

- Контраст и читабельность: используйте высококонтрастные сочетания цветов.
- Ограничивайте ширину: применяйте `common.CalculateLayoutWidth` и избегайте «бесконечных» строк.
- Безопасная длина с ANSI: не считайте отступы через `len(coloredText)`, используйте `ui.GetPlainTextLength`.
- Emoji и CJK: используйте `ui.WrapText` и `ui.GetRuneWidth` для предотвращения смещений.
- Последовательность: выводите свои стили из предоставленных, чтобы сохранить целостный вид.
