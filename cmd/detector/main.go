package main

import (
	"fmt"
	"github.com/qzeleza/termos/examples"
)

// Точка входа для примера детектора embedded окружения
// Запуск: go run ./cmd/detector
func main() {
	fmt.Println("Termos - Детектор embedded окружения")
	fmt.Println("=====================================")

	if examples.IsEmbeddedEnvironment() {
		fmt.Println("Обнаружено embedded окружение: да")
	} else {
		fmt.Println("Обнаружено embedded окружение: нет")
	}
}
