package examples

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/qzeleza/termos/internal/task"
)

// removeANSI удаляет ANSI escape последовательности из строки
func removeANSI(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(s, "")
}

// min возвращает минимальное значение из двух int
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println("=== ФИНАЛЬНЫЙ ТЕСТ КОМПОНЕНТОВ ===")

	success := true

	// Тест 1: SingleSelectTask
	fmt.Println("\n1️⃣ Тест SingleSelectTask:")
	single := task.NewSingleSelectTask("Тест выбора", []string{"A", "B", "C"})
	single.WithTimeout(10*time.Second, 1)
	single.Run()
	view := single.View(80)

	fmt.Printf("   ✓ Длина View: %d символов\n", len(view))
	fmt.Printf("   ✓ Содержит префикс '○': %t\n", strings.Contains(view, "○"))
	fmt.Printf("   ✓ Содержит элементы меню: %t\n", strings.Contains(view, "A"))
	fmt.Printf("   ✓ Содержит таймер '[': %t\n", strings.Contains(view, "["))

	timer := single.BaseTask.RenderTimer()
	fmt.Printf("   ✓ Таймер работает: %t ('%s')\n", len(timer) > 0, timer)

	if !strings.Contains(view, "○") || !strings.Contains(view, "A") || len(timer) == 0 {
		success = false
		fmt.Println("   ❌ ОШИБКА в SingleSelectTask!")
	} else {
		fmt.Println("   ✅ SingleSelectTask OK")
	}

	// Тест 2: InputTaskNew
	fmt.Println("\n2️⃣ Тест InputTaskNew:")
	input := task.NewInputTaskNew("Тест ввода", "подсказка")
	input.WithTimeout(5*time.Second, "по умолчанию")
	input.Run()
	inputView := input.View(80)

	// Убираем ANSI коды для правильной проверки
	cleanInputView := removeANSI(inputView)

	fmt.Printf("   ✓ Длина View: %d символов (сырой), %d символов (очищенный)\n", len(inputView), len(cleanInputView))
	fmt.Printf("   ✓ Содержит префикс '└─>': %t\n", strings.Contains(cleanInputView, "└─>"))
	fmt.Printf("   ✓ Содержит поле ввода '...': %t\n", strings.Contains(cleanInputView, "..."))

	inputTimer := input.BaseTask.RenderTimer()
	cleanTimer := removeANSI(inputTimer)
	fmt.Printf("   ✓ Таймер работает: %t ('%s')\n", len(cleanTimer) > 0, cleanTimer)

	// Проверяем все необходимые элементы в очищенном виде
	hasPrefix := strings.Contains(cleanInputView, "└─>")
	hasInput := strings.Contains(cleanInputView, "...")
	hasTimer := len(cleanTimer) > 0 && strings.Contains(cleanTimer, "[")

	if !hasPrefix || !hasInput || !hasTimer {
		success = false
		fmt.Println("   ❌ ОШИБКА в InputTaskNew!")
		fmt.Printf("     Префикс: %t, Поле ввода: %t, Таймер: %t\n", hasPrefix, hasInput, hasTimer)
		// Отладочная информация
		fmt.Printf("     Первые 100 символов (очищенный): '%s'\n", cleanInputView[:min(100, len(cleanInputView))])
	} else {
		fmt.Println("   ✅ InputTaskNew OK")
	}

	// Тест 3: MultiSelectTask
	fmt.Println("\n3️⃣ Тест MultiSelectTask:")
	multi := task.NewMultiSelectTask("Тест множественного выбора", []string{"X", "Y", "Z"})
	multi.WithTimeout(8*time.Second, []string{"X", "Y"})
	multi.Run()
	multiView := multi.View(80)

	fmt.Printf("   ✓ Длина View: %d символов\n", len(multiView))
	fmt.Printf("   ✓ Содержит элементы: %t\n", strings.Contains(multiView, "X"))
	fmt.Printf("   ✓ Содержит чекбоксы: %t\n", strings.Contains(multiView, "["))

	multiTimer := multi.BaseTask.RenderTimer()
	fmt.Printf("   ✓ Таймер работает: %t ('%s')\n", len(multiTimer) > 0, multiTimer)

	if !strings.Contains(multiView, "X") || len(multiTimer) == 0 {
		success = false
		fmt.Println("   ❌ ОШИБКА в MultiSelectTask!")
	} else {
		fmt.Println("   ✅ MultiSelectTask OK")
	}

	// Тест 4: YesNoTask
	fmt.Println("\n4️⃣ Тест YesNoTask:")
	yesno := task.NewYesNoTask("Тест да/нет", "Вопрос?")
	yesno.WithTimeout(6*time.Second, 0)
	yesno.Run()
	yesnoView := yesno.View(80)

	fmt.Printf("   ✓ Длина View: %d символов\n", len(yesnoView))
	fmt.Printf("   ✓ Содержит варианты: %t\n", strings.Contains(yesnoView, "Да"))

	yesnoTimer := yesno.BaseTask.RenderTimer()
	fmt.Printf("   ✓ Таймер работает: %t ('%s')\n", len(yesnoTimer) > 0, yesnoTimer)

	if !strings.Contains(yesnoView, "Да") || len(yesnoTimer) == 0 {
		success = false
		fmt.Println("   ❌ ОШИБКА в YesNoTask!")
	} else {
		fmt.Println("   ✅ YesNoTask OK")
	}

	// Финальный результат
	fmt.Println("\n" + strings.Repeat("=", 50))
	if success {
		fmt.Println("🎉 ВСЕ ТЕСТЫ ПРОЙДЕНЫ! КОД РАБОТАЕТ КОРРЕКТНО!")
		fmt.Println("✨ Таймеры отображаются, префиксы работают, меню видны!")
		fmt.Println("")
		fmt.Println("Теперь вы можете запустить:")
		fmt.Println("   ./timeout_demo")
		fmt.Println("или")
		fmt.Println("   go run examples/final_timeout_demo.go")
	} else {
		fmt.Println("❌ НАЙДЕНЫ ПРОБЛЕМЫ! Проверьте код выше.")
	}
	fmt.Println(strings.Repeat("=", 50))
}
