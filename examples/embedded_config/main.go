package main

import (
	"fmt"
	"github.com/qzeleza/termos/examples"
)

// Точка входа для демонстрации embedded-конфигураций
// Запуск: go run ./cmd/embedded_config
func main() {
	fmt.Println("Termos - Embedded конфигурации")
	fmt.Println("================================")

	isEmbedded := examples.IsEmbeddedEnvironment()
	fmt.Printf("Обнаружено embedded окружение: %v\n", isEmbedded)

	var cfg *examples.EmbeddedConfig
	if isEmbedded {
		cfg = examples.OptimizedEmbeddedConfig()
		fmt.Println("Применяется оптимизированная конфигурация для embedded")
	} else {
		cfg = examples.DefaultEmbeddedConfig()
		fmt.Println("Применяется стандартная конфигурация")
	}

	examples.ApplyEmbeddedConfig(cfg)
	est := examples.GetMemoryFootprintEstimate(cfg)
	fmt.Printf("Оценка потребления памяти: ~%d MB\n", est)
}
