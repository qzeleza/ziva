package examples

import (
	"regexp"
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

// Удалена неиспользуемая функция main
// func main() {
//     // пример использования
// }
