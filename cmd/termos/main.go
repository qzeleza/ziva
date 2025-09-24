package main

import (
	// Встроенные импорты не требуются

	"errors"
	"flag"
	"fmt"
	"log"
	"os"
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
	// Настраиваем язык интерфейса и предупреждаем о возможных ограничениях терминала
	activeLang := configureLanguage()
	warnTerminalCapabilities(activeLang)
	termos.SetDefaultLanguage("ru")
	// Используем собственный разделитель для подсказок в пунктах меню
	termos.SetChoiceHelpDelimiter("||")

	// Заголовок и краткое описание для TUI
	header := "Демонстрация всех типов задач Termos"

	// Создаем очередь и добавляем задачи
	queue := termos.NewQueue(header)
	queue.WithAppName("Термос")
	// queue.WithOutResultLine()
	// queue.WithOutSummary()
	// queue.WithTasksNumbered(false, "[%d]")

	// Формируем очередь задач
	var msel = []string{
		"CLI||Командный интерфейс для администрирования",
		"Сервер||Бэкенд сервисы",
		"Агент",
		"Web UI||Веб-интерфейс для пользователей",
		"Документация||Автоматическая генерация документации",
		"Компилировать||Сборка исполняемых файлов",
		"Выходные данные||Архивация результатов",
		"Область просмотра",
		"Поле ввода||Пример текстовой задачи",
		"Мультивыбор||Дополнительные параметры",
		"Одиночный выбор||Переключатель режимов",
		"Проверка ввода||Встроенные валидаторы"}

	var ssel = []string{
		"development||Среда разработки",
		"staging||Промежуточная среда",
		"production||Боевая среда",
		"другое||Пользовательское значение",
		"отмена||Отмена выбора",
		"выход||Выход из программы",
	}
	// 1) Задачи мультивыбора (без и с пунктом "Выбрать все")
	//    Пример без "Выбрать все"
	ms1 := termos.NewMultiSelectTask("Выберите компоненты установки", msel).
		WithViewport(5, false).
		WithTimeout(3*time.Second, []string{msel[0], msel[1]}).
		WithItemsDisabled([]string{msel[2], msel[3]})

	// //    Пример с пунктом "Выбрать все"
	// ms2 := termos.NewMultiSelectTask("Выберите модули для сборки", ssel).
	// 	WithViewport(3, false).
	// 	WithSelectAll("Выбрать все").
	// 	WithTimeout(10*time.Second, []string{ssel[0], ssel[1]}).
	// 	WithDefaultItems([]string{ssel[0], ssel[1]}).
	// 	WithItemsDisabled([]string{ssel[2], ssel[3]})

	// 2) Одиночный выбор
	ss := termos.NewSingleSelectTask(
		"Выберите среду развертывания",
		ssel,
	).WithViewport(3).
		WithTimeout(3*time.Second, "staging").
		WithDefaultItem("production")

	// 3) Ввод с использованием всех стандартных валидаторов
	//    Валидация будет происходить в момент подтверждения (Enter)
	// v := termos.DefaultValidators

	// inPath := task.NewInputTaskNew("Путь к файлу/директории", "Введите путь:").
	// 	WithValidator(v.Path())

	queue.AddTasks(
		ss,
		// inUsername,
		ms1,
		// ms2,
		// inPath,
	)

	// 4) Задача-выполнение функции (FuncTask)
	//    Выполняет полезную работу и выводит результат в финальном представлении задачи (без fmt.Print)
	errorTaskRun := false
	var fn *termos.FuncTask
	if errorTaskRun {
		data := pingResult{}
		fn = termos.NewFuncTask(
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
		queue.AddTasks(fn)
	}
	// 5) Подтверждение Да/Нет (например, для сохранения настроек)
	// Используем языко-независимый метод вместо строки "Да"
	ys := termos.NewYesNoTask("Сохранение конфигурации", "Сохранить изменения?").WithTimeoutYes(2 * time.Second)

	// inUsername := termos.NewInputTask("Имя пользователя", "Введите username:").
	// 	WithValidator(v.Username()).
	// 	WithTimeout(10*time.Second, "Alex")

	// inEmail := task.NewInputTaskNew("Email", "Введите email:").
	// 	WithInputType(task.InputTypeEmail).WithValidator(v.Email()).
	// 	WithTimeout(10*time.Second, "defauilt@example.com")

	// inOptionalEmail := task.NewInputTaskNew("Доп. Email (опционально)", "Введите email или оставьте пустым:").
	// 	WithInputType(task.InputTypeEmail).WithValidator(v.OptionalEmail())

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

	queue.AddTasks(

		//  inEmail, inOptionalEmail,
		// inPath, inURL, inPort, inRange,
		// inIPv4, inIPv6, inIPAny, inDomain,
		// inAlphaNum, inMinLen, inMaxLen, inExactLen,
		// inStdPwd, inStrongPwd, inRequired,
		ys,
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

// configureLanguage разбирает язык из CLI/конфига, проверяет доступность локали и применяет его.
func configureLanguage() string {
	if defLang := strings.TrimSpace(os.Getenv("TERMOS_DEFAULT_LANG")); defLang != "" {
		termos.SetDefaultLanguage(defLang)
	}
	langFlag := flag.String("lang", "", "язык интерфейса Termos (например, ru или en)")
	flag.Parse()

	lang := strings.TrimSpace(*langFlag)
	if lang == "" {
		lang = strings.TrimSpace(os.Getenv("TERMOS_LANG"))
	}
	if lang == "" {
		lang = "ru"
	}
	lang = strings.ToLower(lang)

	if strings.HasPrefix(lang, "ru") && !localeAvailable("ru_RU") {
		printRussianLocaleHint()
		lang = "en"
	}

	return termos.SetLanguage(lang)
}

// localeAvailable проверяет присутствие заданной локали в системе.
func localeAvailable(locale string) bool {
	if strings.HasPrefix(locale, "en") {
		return true
	}
	cmd := exec.Command("locale", "-a")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	list := strings.Split(strings.ToLower(string(output)), "\n")
	targets := []string{"ru_ru.utf-8", "ru_ru.utf8", "ru_ru", "ru"}
	for _, line := range list {
		line = strings.TrimSpace(line)
		for _, target := range targets {
			if line == target {
				return true
			}
		}
	}
	return false
}

// printRussianLocaleHint выводит рекомендации по установке русской локали и шрифта.
func printRussianLocaleHint() {
	fmt.Println("⚠️ Не удалось найти локаль ru_RU.UTF-8. Переключаю интерфейс на английский.")
	fmt.Println("   Установите русскую локаль командой (Debian/Ubuntu): sudo locale-gen ru_RU.UTF-8 && sudo update-locale LANG=ru_RU.UTF-8")
	fmt.Println("   Для Entware/BusyBox: opkg install locale-full glibc-binary-locales && export LANG=ru_RU.UTF-8")
	fmt.Println("   При необходимости настройте шрифт: setterm -reset && setterm -store")
}

// warnTerminalCapabilities предупреждает о возможных ограничениях терминала.
func warnTerminalCapabilities(lang string) {
	langEnv := strings.ToLower(os.Getenv("LANG"))
	term := strings.ToLower(os.Getenv("TERM"))
	colorTerm := strings.TrimSpace(os.Getenv("COLORTERM"))

	if !strings.Contains(langEnv, "utf") {
		fmt.Println("⚠️ Текущая локаль не содержит UTF-8. Псевдографика может отображаться некорректно.")
		fmt.Println("   Совет: export LANG=ru_RU.UTF-8 && export LC_ALL=ru_RU.UTF-8")
	}
	if colorTerm == "" {
		fmt.Println("⚠️ Терминал не сообщает о поддержке цвета (COLORTERM пуст). Включите цветной режим или используйте современный эмулятор.")
	}
	if term == "linux" || strings.Contains(term, "vt100") || strings.Contains(term, "busybox") {
		fmt.Println("ℹ️ Для Entware/BusyBox установите шрифт UTF-8: setterm -store && setterm -font latarcyrheb-sun32")
	}
}
