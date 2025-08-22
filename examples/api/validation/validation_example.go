package main

import (
	"log"
	"time"

	"github.com/qzeleza/termos"
)

func main() {
	// Создаем новую очередь задач
	queue := termos.NewQueue("Пример валидации данных").
		WithAppName("Validation Demo").
		WithSummary(true)

	v := termos.DefaultValidators

	// Email задача
	email := termos.NewInputTask("Email", "Введите email адрес:").
		WithInputType(termos.InputTypeEmail).
		WithValidator(v.Email()).
		WithTimeout(20*time.Second, "user@example.com")

	// Пароль задача
	password := termos.NewInputTask("Пароль", "Введите пароль (минимум 8 символов):").
		WithInputType(termos.InputTypePassword).
		WithValidator(v.StandardPassword()).
		WithTimeout(20*time.Second, "password123")

	// Номер порта
	port := termos.NewInputTask("Порт", "Введите номер порта (1-65535):").
		WithInputType(termos.InputTypeNumber).
		WithValidator(v.Port()).
		WithTimeout(15*time.Second, "8080")

	// IP адрес
	ip := termos.NewInputTask("IP адрес", "Введите IP адрес:").
		WithInputType(termos.InputTypeIP).
		WithValidator(v.IP()).
		WithTimeout(15*time.Second, "192.168.1.1")

	// Домен
	domain := termos.NewInputTask("Доменное имя", "Введите домен:").
		WithInputType(termos.InputTypeDomain).
		WithValidator(v.Domain()).
		WithTimeout(15*time.Second, "example.com")

	// Путь к файлу
	path := termos.NewInputTask("Путь", "Введите путь к файлу/директории:").
		WithValidator(v.Path()).
		WithTimeout(15*time.Second, "/tmp")

	// URL
	url := termos.NewInputTask("URL", "Введите URL (http/https):").
		WithValidator(v.URL()).
		WithTimeout(15*time.Second, "https://example.com")

	// Число в диапазоне
	number := termos.NewInputTask("Число", "Введите число от 10 до 100:").
		WithInputType(termos.InputTypeNumber).
		WithValidator(v.Range(10, 100)).
		WithTimeout(15*time.Second, "50")

	// Только буквы и цифры
	alphaNum := termos.NewInputTask("Логин", "Введите логин (только буквы и цифры):").
		WithValidator(v.AlphaNumeric()).
		WithTimeout(15*time.Second, "user123")

	// Подтверждение
	confirm := termos.NewYesNoTask("Подтверждение", "Сохранить данные?").
		WithTimeout(10*time.Second, "Да")

	// Добавляем все задачи в очередь
	queue.AddTasks(
		email,
		password,
		port,
		ip,
		domain,
		path,
		url,
		number,
		alphaNum,
		confirm,
	)

	// Запускаем очередь
	if err := queue.Run(); err != nil {
		log.Fatal(err)
	}
}
