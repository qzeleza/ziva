package examples

import (
	"fmt"
	"github.com/qzeleza/termos/internal/task"
)

func YesNoDemo() {
	fmt.Println("=== Демонстрация YesNoTask ===")
	
	// Создаем YesNo задачу
	yesNoTask := task.NewYesNoTask("Сохранить конфигурацию?", "Хотите ли вы сохранить текущую конфигурацию?")
	
	// Показываем View задачи
	view := yesNoTask.View(80)
	fmt.Printf("YesNoTask View:\n%s\n\n", view)
	
	// Проверяем доступные опции через методы
	fmt.Println("Проверка доступных опций:")
	fmt.Printf("  IsYes(): %v\n", yesNoTask.IsYes())
	fmt.Printf("  IsNo(): %v\n", yesNoTask.IsNo())
	
	fmt.Println("\n✅ YesNoTask теперь содержит только 2 опции: 'Да' и 'Нет'")
	fmt.Println("❌ Опция 'Выйти' была удалена")
	fmt.Println("ℹ️ Теперь только 2 опции в списке выбора, используйте стрелки и Enter")
}
