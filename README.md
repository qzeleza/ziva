# Термос - текстовый интерфейс для работы на малых маршрутизаторах и устройствах, написанный в виде библиотеки на Go

**Термос** - это мощный и гибкий фреймворк для создания терминальных пользовательских интерфейсов (TUI) на Go. Он построен поверх библиотеки [Bubble Tea](https://github.com/charmbracelet/bubbletea) и предоставляет высокоуровневые компоненты для создания интерактивных консольных приложений.

## ✨ Особенности

- 🎯 **Готовые компоненты** - множественный выбор, одиночный выбор, ввод текста, да/нет
- 🎨 **Гибкая стилизация** - поддержка цветов, стилей и тем оформления
- 📱 **Адаптивность** - оптимизация для embedded устройств с ограниченными ресурсами
- 🔄 **Очереди задач** - система управления последовательностью задач
- ✅ **Валидация** - встроенная система валидации пользовательского ввода
- 🎭 **Анимации** - поддержка анимированных переходов и эффектов
- 🧠 **Умная память** - оптимизация использования памяти для встроенных систем

## 🚀 Быстрый старт

![Внешний вид](docs/images/pic_1.png)

### Установка

```bash
go get github.com/qzeleza/termos
```

## Выполнение функции (FuncTask)

`FuncTask` позволяет выполнить произвольную функцию и отобразить результат (успех/ошибка). Поддерживает дополнительные строки сводки в финальном представлении задачи и настраиваемую метку успеха.

```go
fn := task.NewFuncTaskWithOptions(
    "Проверка соединения",
    func() error {
        // Выполните реальную работу; верните ошибку при сбое
        return nil
    },
    // Дополнительные строки сводки под заголовком при успехе
    task.WithSummaryFunction(func() []string {
        return []string{"Пинг: 12мс", "Потери пакетов: 0%"}
    }),
    // Не останавливать очередь при ошибке (по умолчанию stopOnError=false)
    task.WithStopOnError(false),
    // Переопределить метку успеха (по умолчанию "ГОТОВО")
    task.WithSuccessLabelOption("ЗАВЕРШЕНО"),
)

// Подсказка: q/esc/ctrl+c — отмена задачи пользователем
// Использование: добавьте в очередь задач либо интегрируйте в Bubble Tea модель
```

## Полный пример: очередь всех задач

Ниже приведён практичный пример, который демонстрирует все основные типы задач, стандартные валидаторы и запуск очереди без промежуточного вывода (итог — после завершения всех задач):

```go
package main

import (
    // Встроенные импорты не требуются

    "github.com/qzeleza/termos/common"
    "github.com/qzeleza/termos/examples"
    "github.com/qzeleza/termos/task"
    "github.com/qzeleza/termos/validation"
)

func main() {
    // Заголовок и краткое описание для TUI
    header := "Демонстрация всех типов задач Термос"
    summary := "Мультивыбор, одиночный выбор, ввод с валидаторами, функция, Да/Нет"

    // Формируем очередь задач. ВАЖНО: срез должен быть типа []common.Task
    var tasks []common.Task

    // 1) Задачи мультивыбора (без и с пунктом "Выбрать все")
    //    Пример без "Выбрать все"
    ms1 := task.NewMultiSelectTask(
        "Выберите компоненты установки",
        []string{"CLI", "Сервер", "Агент", "Web UI", "Документация"},
    )
    //    Пример с пунктом "Выбрать все"
    ms2 := task.NewMultiSelectTask(
        "Выберите модули для сборки",
        []string{"auth", "storage", "network", "monitoring"},
    ).WithSelectAll("Выбрать все")

    // 2) Одиночный выбор
    ss := task.NewSingleSelectTask(
        "Выберите среду развертывания",
        []string{"development", "staging", "production"},
    )

    // 3) Ввод с использованием всех стандартных валидаторов
    //    Валидация будет происходить в момент подтверждения (Enter)
    v := validation.DefaultFactory

    inUsername := task.NewInputTaskNew("Имя пользователя", "Введите username:").
        WithValidator(v.Username())

    inEmail := task.NewInputTaskNew("Email", "Введите email:").
        WithInputType(task.InputTypeEmail).WithValidator(v.Email())

    inOptionalEmail := task.NewInputTaskNew("Доп. Email (опционально)", "Введите email или оставьте пустым:").
        WithInputType(task.InputTypeEmail).WithValidator(v.OptionalEmail())

    inPath := task.NewInputTaskNew("Путь к файлу/директории", "Введите путь:").
        WithValidator(v.Path())

    inURL := task.NewInputTaskNew("URL", "Введите URL (http/https):").
        WithValidator(v.URL())

    inPort := task.NewInputTaskNew("Порт", "Введите порт (1-65535):").
        WithInputType(task.InputTypeNumber).WithValidator(v.Port())

    inRange := task.NewInputTaskNew("Число в диапазоне", "Введите число [10..100]:").
        WithInputType(task.InputTypeNumber).WithValidator(v.Range(10, 100))

    inIPv4 := task.NewInputTaskNew("IPv4", "Введите IPv4 адрес:").
        WithInputType(task.InputTypeIP).WithValidator(v.IPv4())

    inIPv6 := task.NewInputTaskNew("IPv6", "Введите IPv6 адрес:").
        WithInputType(task.InputTypeIP).WithValidator(v.IPv6())

    inIPAny := task.NewInputTaskNew("IP (любой)", "Введите IP адрес:").
        WithInputType(task.InputTypeIP).WithValidator(v.IP())

    inDomain := task.NewInputTaskNew("Домен", "Введите доменное имя:").
        WithInputType(task.InputTypeDomain).WithValidator(v.Domain())

    inAlphaNum := task.NewInputTaskNew("Только буквы и цифры", "Введите значение:").
        WithValidator(v.AlphaNumeric())

    inMinLen := task.NewInputTaskNew("Мин. длина", "Минимум 5 символов:").
        WithValidator(v.MinLength(5))

    inMaxLen := task.NewInputTaskNew("Макс. длина", "Не более 10 символов:").
        WithValidator(v.MaxLength(10))

    inExactLen := task.NewInputTaskNew("Точная длина", "Ровно 8 символов:").
        WithValidator(v.Length(8))

    inStdPwd := task.NewInputTaskNew("Пароль (стандарт)", "Введите пароль (>=8):").
        WithInputType(task.InputTypePassword).WithValidator(v.StandardPassword())

    inStrongPwd := task.NewInputTaskNew("Пароль (сильный)", "Введите пароль (>=12):").
        WithInputType(task.InputTypePassword).WithValidator(v.StrongPassword())

    inRequired := task.NewInputTaskNew("Обязательное поле", "Нельзя оставлять пустым:").
        WithValidator(v.Required())

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
    ys := task.NewYesNoTask("Сохранение конфигурации", "Сохранить изменения?")

    // Добавляем задачи в очередь
    tasks = append(tasks,
        ms1, ms2, ss,
        inUsername, inEmail, inOptionalEmail,
        inPath, inURL, inPort, inRange,
        inIPv4, inIPv6, inIPAny, inDomain,
        inAlphaNum, inMinLen, inMaxLen, inExactLen,
        inStdPwd, inStrongPwd, inRequired,
        fn, ys,
    )

    // Запускаем TUI c очередью задач. Результаты отображаются внутри интерфейса;
    // дополнительный вывод через fmt.Print не используется.
    _ = examples.RunTasksWithTUI(header, summary, tasks)
}
```

## 🖼️ Скриншоты

Ниже приведены скриншоты, демонстрирующие работу интерфейса Термос:

![Одиночный выбор](docs/images/pic_2.png)

![Множественный выбор](docs/images/pic_3.png)

![Ввод с валидацией](docs/images/pic_4.png)

![Сводка очереди](docs/images/pic_5.png)

## 📦 Компоненты

### Задачи (Tasks)

- **YesNoTask** - выбор "да/нет" (только 2 опции)
- **SingleSelectTask** - выбор одного элемента из списка
- **MultiSelectTask** - выбор нескольких элементов из списка  
- **InputTaskNew** - ввод текста с валидацией

### Очереди (планируется)

- Оркестрация последовательных задач и сбор статистики в пакете `query/` (в разработке)

### UI компоненты

- **Styles** - система стилизации
- **Colors** - управление цветовой схемой
- **Layout** - компоненты разметки

## 🎨 Кастомизация

### Стили

```go
import "github.com/qzeleza/termos/ui"

// Настройка пользовательских стилей
styles := ui.GetDefaultStyles()
styles.Title = styles.Title.Foreground(lipgloss.Color("#ff6b6b"))
styles.Selected = styles.Selected.Background(lipgloss.Color("#4ecdc4"))
```

### Embedded оптимизации

```go
// 1) Ограничивайте ширину рендера
import "github.com/qzeleza/termos/common"
w := common.CalculateLayoutWidth(terminalWidth)

// 2) Переиспользуйте стили вместо пересоздания в каждом кадре
import "github.com/qzeleza/termos/ui"
ui.TitleStyle = ui.TitleStyle.Foreground(ui.ColorBrightGreen)

// 3) Используйте пулы для сборки строк
import "github.com/qzeleza/termos/performance"
b := performance.GetBuffer()
defer performance.PutBuffer(b)
b.WriteString("...быстрая сборка строки...")
```

> Примечание: оптимизации для embedded теперь включаются автоматически при импортировании модуля (внутренний пакет `internal/autoconfig`). Ручные примеры из каталога `examples/` можно использовать как демонстрацию, но базовая адаптация применяется по умолчанию.

Переменные окружения для управления автодетекцией:

- `TERMOS_EMBEDDED` — принудительное включение/выключение (`true`/`1` или пусто/`0`).
- `TERMOS_MEMORY_LIMIT` — порог памяти для эвристики (например: `64MB`, `128KB`, `1GB`).
- `TERMOS_ASCII_ONLY` — форсировать ASCII-режим (`true`).

## 📚 Документация

- [Руководство разработчика](docs/DEVDOC.md)
- [Примеры использования](docs/EXAMPLES.md)
- [Горячие клавиши](docs/HOTKEYS.md)
- [Настройка стилей](docs/STYLING.md)
- [Оптимизация для embedded](docs/EMBEDDED_OPTIMIZATION_REPORT.md)
- [Решение проблем](docs/TROUBLESHOOTING.md)

## 🏗️ Архитектура

```
termos/
├── task/           # Основные компоненты задач
├── query/          # Система очередей
├── ui/             # UI компоненты и стили  
├── validation/     # Система валидации
├── performance/    # Оптимизации производительности
├── common/         # Общие утилиты
├── errors/         # Обработка ошибок
└── docs/           # Документация
```

## 🎯 Применение

Термос идеально подходит для:

- **CLI утилиты** - интерактивные консольные приложения
- **Инсталляторы** - пошаговые установщики и настройщики
- **Системные утилиты** - управление конфигурациями и сервисами
- **Embedded устройства** - роутеры, IoT устройства с ограниченными ресурсами
- **DevOps инструменты** - развертывание и мониторинг

## 🤝 Совместимость

- **Go версия:** 1.22+
- **Платформы:** Linux, macOS, Windows
- **Терминалы:** все современные терминалы с поддержкой ANSI
- **Embedded:** OpenWrt, Entware, роутеры с ≥32MB RAM

## 📊 Производительность

- **Память:** оптимизировано для работы с ≥16MB RAM
- **CPU:** эффективная работа на ARM, MIPS, x86 архитектурах
- **Анимации:** адаптивная частота кадров в зависимости от мощности системы

## 🐛 Отчеты об ошибках

Если вы нашли ошибку или хотите предложить улучшение:

1. Проверьте [существующие issue](https://github.com/qzeleza/termos/issues)
2. Создайте новый issue с подробным описанием
3. Приложите минимальный пример воспроизведения

## 📄 Лицензия

MIT License - см. файл [LICENSE](LICENSE) для подробностей.

## 🙏 Благодарности

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - основа для TUI
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - стилизация
- Сообщество Go разработчиков

---

**Создано с ❤️ для сообщества Go разработчиков**