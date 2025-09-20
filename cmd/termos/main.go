package main

import (
	// Встроенные импорты не требуются

	"errors"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/qzeleza/termos"
)

// pingResult представляет результат проверки подключения
type pingResult struct {
	Ping   string
	Loss   string
	Status string
}

func main() {
	// Заголовок и краткое описание для TUI
	header := "Демонстрация всех типов задач Termos"

	// Формируем очередь задач
	var msel = []string{
		"CLI",
		"Сервер",
		"Агент",
		"Web UI",
		"Документация",
		"Компилировать",
		"Выходные данные",
		"Область просмотра",
		"Поле ввода",
		"Мультивыбор",
		"Одиночный выбор",
		"Проверка ввода"}

	var ssel = []string{"development", "staging", "production"}
	// 1) Задачи мультивыбора (без и с пунктом "Выбрать все")
	//    Пример без "Выбрать все"
	ms1 := termos.NewMultiSelectTask(
		"Выберите компоненты установки",
		msel,
	).WithViewport(5).WithTimeout(3*time.Second, []string{msel[0], msel[1]})
	//    Пример с пунктом "Выбрать все"
	ms2 := termos.NewMultiSelectTask(
		"Выберите модули для сборки",
		ssel,
	).WithViewport(3).WithSelectAll("Выбрать все").WithTimeout(10*time.Second, []string{ssel[0], ssel[1]}).WithDefaultItems([]string{ssel[0], ssel[1]})

	// 2) Одиночный выбор
	ss := termos.NewSingleSelectTask(
		"Выберите среду развертывания",
		[]string{"development", "staging", "production", "другое", "отмена", "выход"},
	).WithViewport(3).
		// WithTimeout(3*time.Second, "staging")
		WithDefaultItem("production")

	// 3) Ввод с использованием всех стандартных валидаторов
	//    Валидация будет происходить в момент подтверждения (Enter)
	// v := termos.DefaultValidators

	// inUsername := termos.NewInputTask("Имя пользователя", "Введите username:").
	// 	WithValidator(v.Username()).
	// 	WithTimeout(10*time.Second, "Alex")

	// inEmail := task.NewInputTaskNew("Email", "Введите email:").
	// 	WithInputType(task.InputTypeEmail).WithValidator(v.Email()).
	// 	WithTimeout(10*time.Second, "default@example.com")

	// inOptionalEmail := task.NewInputTaskNew("Доп. Email (опционально)", "Введите email или оставьте пустым:").
	// 	WithInputType(task.InputTypeEmail).WithValidator(v.OptionalEmail())

	// inPath := task.NewInputTaskNew("Путь к файлу/директории", "Введите путь:").
	// 	WithValidator(v.Path())

	// inURL := task.NewInputTaskNew("URL", "Введите URL (http/https):").
	// 	WithValidator(v.URL())

	// inPort := task.NewInputTaskNew("Порт", "Введите порт (1-65535):").
	// 	WithInputType(task.InputTypeNumber).WithValidator(v.Port())

	// inRange := task.NewInputTaskNew("Число в диапазоне", "Введите число [10..100]:").
	// 	WithInputType(task.InputTypeNumber).WithValidator(v.Range(10, 100))

	// inIPv4 := task.NewInputTaskNew("IPv4", "Введите IPv4 адрес:").
	// 	WithInputType(task.InputTypeIP).WithValidator(v.IPv4())

	// inIPv6 := task.NewInputTaskNew("IPv6", "Введите IPv6 адрес:").
	// 	WithInputType(task.InputTypeIP).WithValidator(v.IPv6())

	// inIPAny := task.NewInputTaskNew("IP (любой)", "Введите IP адрес:").
	// 	WithInputType(task.InputTypeIP).WithValidator(v.IP())

	// inDomain := task.NewInputTaskNew("Домен", "Введите доменное имя:").
	// 	WithInputType(task.InputTypeDomain).WithValidator(v.Domain())

	// inAlphaNum := task.NewInputTaskNew("Только буквы и цифры", "Введите значение:").
	// 	WithValidator(v.AlphaNumeric())

	// inMinLen := task.NewInputTaskNew("Мин. длина", "Минимум 5 символов:").
	// 	WithValidator(v.MinLength(5))

	// inMaxLen := task.NewInputTaskNew("Макс. длина", "Не более 10 символов:").
	// 	WithValidator(v.MaxLength(10))

	// inExactLen := task.NewInputTaskNew("Точная длина", "Ровно 8 символов:").
	// 	WithValidator(v.Length(8))

	// inStdPwd := task.NewInputTaskNew("Пароль (стандарт)", "Введите пароль (>=8):").
	// 	WithInputType(task.InputTypePassword).WithValidator(v.StandardPassword())

	// inStrongPwd := task.NewInputTaskNew("Пароль (сильный)", "Введите пароль (>=12):").
	// 	WithInputType(task.InputTypePassword).WithValidator(v.StrongPassword())

	// inRequired := task.NewInputTaskNew("Обязательное поле", "Нельзя оставлять пустым:").
	// 	WithValidator(v.Required())

	// 4) Задача-выполнение функции (FuncTask)
	//    Выполняет полезную работу и выводит результат в финальном представлении задачи (без fmt.Print)
	data := pingResult{}
	fn := termos.NewFuncTask(
		"Проверка соединения",
		func() error {
			// return checkConnection(&data)
			return errors.New("симуляция ошибки в середине выполнения очереди\nне ясная причина стимуляции проблемы\nдополнительная информация")
		},
		// Выводим краткую сводку под заголовком после успеха
		termos.WithSummaryFunction(func() []string {
			return []string{
				"Пинг: " + data.Ping,
				"Потери пакетов: " + data.Loss,
			}
		}),
		// Не останавливать очередь при ошибке (для демонстрации поведения)
		termos.WithStopOnError(true),
	)

	// 5) Подтверждение Да/Нет (например, для сохранения настроек)
	ys := termos.NewYesNoTask("Сохранение конфигурации", "Сохранить изменения?").WithTimeout(5*time.Second, "Нет")

	// Создаем очередь и добавляем задачи
	queue := termos.NewQueue(header).WithAppName("Термос").WithSummary(true)
	queue.AddTasks(
		ss,
		// inUsername,
		ms1,
		ms2,
		//  inEmail, inOptionalEmail,
		// inPath, inURL, inPort, inRange,
		// inIPv4, inIPv6, inIPAny, inDomain,
		// inAlphaNum, inMinLen, inMaxLen, inExactLen,
		// inStdPwd, inStrongPwd, inRequired,
		fn, ys,
	)

	// Запускаем TUI c очередью задач. Результаты отображаются внутри интерфейса;
	// дополнительный вывод через fmt.Print не используется.
	if err := queue.Run(); err != nil {
		// Обработка ошибки
		log.Fatalf("Ошибка при запуске очереди: %v", err)
	}
}

func checkConnection(result *pingResult) error {
	time.Sleep(2 * time.Second)
	cmd := exec.Command("ping", "-c", "1", "google.com")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	// Парсим результат команды ping
	outputStr := string(output)

	// Извлекаем время пинга
	if strings.Contains(outputStr, "time=") {
		parts := strings.Split(outputStr, "time=")
		if len(parts) > 1 {
			timePart := strings.Fields(parts[1])[0]
			result.Ping = timePart
		}
	}

	// Извлекаем потери пакетов
	if strings.Contains(outputStr, "packet loss") {
		parts := strings.Split(outputStr, " packet loss")
		if len(parts) > 0 {
			lossFields := strings.Fields(parts[0])
			if len(lossFields) > 0 {
				result.Loss = lossFields[len(lossFields)-1]
			}
		}
	}

	// Определяем статус на основе успешности команды
	if !strings.Contains(outputStr, "1 received") {
		result.Status = "FAILED"
	} else {
		result.Status = "OK"
	}

	return nil
}
