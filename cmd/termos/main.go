package main

import (
	// Встроенные импорты не требуются

	"time"

	"github.com/qzeleza/termos/internal/common"
	"github.com/qzeleza/termos/internal/query"
	"github.com/qzeleza/termos/internal/task"
	"github.com/qzeleza/termos/internal/validation"
)

func main() {
	// Заголовок и краткое описание для TUI
	header := "Демонстрация всех типов задач Termos"

	// Формируем очередь задач. ВАЖНО: срез должен быть типа []common.Task
	var tasks []common.Task
	var msel = []string{"CLI", "Сервер", "Агент", "Web UI", "Документация"}
	var ssel = []string{"development", "staging", "production"}
	// 1) Задачи мультивыбора (без и с пунктом "Выбрать все")
	//    Пример без "Выбрать все"
	ms1 := task.NewMultiSelectTask(
		"Выберите компоненты установки",
		msel,
	).WithTimeout(10*time.Second, []string{msel[0], msel[1]})
	//    Пример с пунктом "Выбрать все"
	ms2 := task.NewMultiSelectTask(
		"Выберите модули для сборки",
		ssel,
	).WithSelectAll("Выбрать все").WithTimeout(10*time.Second, []string{ssel[0], ssel[1]})

	// 2) Одиночный выбор
	ss := task.NewSingleSelectTask(
		"Выберите среду развертывания",
		[]string{"development", "staging", "production"},
	).WithTimeout(10*time.Second, "staging")

	// 3) Ввод с использованием всех стандартных валидаторов
	//    Валидация будет происходить в момент подтверждения (Enter)
	v := validation.DefaultFactory

	inUsername := task.NewInputTaskNew("Имя пользователя", "Введите username:").
		WithValidator(v.Username()).
		WithTimeout(10*time.Second, "Alex")

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
	fn := task.NewFuncTaskWithOptions(
		"Проверка соединения",
		func() error {
			// Здесь могла бы быть реальная проверка, для примера считаем, что всё ок
			return nil
		},
		// Выводим краткую сводку под заголовком после успеха
		task.WithSummaryFunction(func() []string {
			return []string{
				"Пинг: 12мс",
				"Потери пакетов: 0%",
			}
		}),
		// Не останавливать очередь при ошибке (для демонстрации поведения)
		task.WithStopOnError(false),
	)

	// 5) Подтверждение Да/Нет (например, для сохранения настроек)
	ys := task.NewYesNoTask("Сохранение конфигурации", "Сохранить изменения?").WithTimeout(5*time.Second, "Нет")

	// Добавляем задачи в очередь
	tasks = append(tasks,

		ss,
		inUsername, ms1, ms2,
		//  inEmail, inOptionalEmail,
		// inPath, inURL, inPort, inRange,
		// inIPv4, inIPv6, inIPAny, inDomain,
		// inAlphaNum, inMinLen, inMaxLen, inExactLen,
		// inStdPwd, inStrongPwd, inRequired,
		fn, ys,
	)

	// Запускаем TUI c очередью задач. Результаты отображаются внутри интерфейса;
	// дополнительный вывод через fmt.Print не используется.
	queue := query.New(header).WithAppName("Термос").WithSummary(true)
	queue.AddTasks(tasks)
	queue.Run()
}
