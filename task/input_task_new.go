package task

import (
	"errors"

	terrors "github.com/qzeleza/termos/errors"
	"github.com/qzeleza/termos/performance"
	"github.com/qzeleza/termos/ui"
	"github.com/qzeleza/termos/validation"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InputType определяет тип ввода для InputTask
type InputType int

// Константы для различных типов ввода
const (
	InputTypeText     InputType = iota // Обычный текстовый ввод
	InputTypePassword                  // Ввод пароля (скрытый текст)
	InputTypeEmail                     // Ввод email
	InputTypeNumber                    // Ввод числа
	InputTypeIP                        // Ввод IP-адреса
	InputTypeDomain                    // Ввод доменного имени
)

// InputTaskNew представляет улучшенную задачу ввода с разделенными обязанностями
type InputTaskNew struct {
	BaseTask

	// Компоненты с разделенными обязанностями
	textInput    textinput.Model       // UI компонент ввода
	validator    validation.Validator  // Валидатор данных
	renderer     *InputRenderer        // Рендерер для отображения
	errorHandler *terrors.ErrorHandler // Обработчик ошибок

	// Состояние задачи
	inputType     InputType // Тип ввода
	validationErr error     // Ошибка валидации
	value         string    // Введенное значение
	prompt        string    // Подсказка для ввода
	placeholder   string    // Текст-заполнитель

	// Настройки
	width      int  // Ширина поля ввода
	maskInput  bool // Маскировать ввод (для паролей)
	allowEmpty bool // Разрешить пустое значение
}

// NewInputTaskNew создает новую улучшенную задачу ввода
func NewInputTaskNew(title, prompt string) *InputTaskNew {
	// Создаем базовую задачу
	baseTask := NewBaseTask(title)
	baseTask.SetStopOnError(true)

	// Инициализируем компонент ввода
	ti := textinput.New()
	ti.Placeholder = DefaultPlaceholder
	ti.Focus()
	ti.CharLimit = MaxInputLength
	ti.Width = DefaultInputWidth

	return &InputTaskNew{
		BaseTask:     baseTask,
		textInput:    ti,
		renderer:     NewInputRenderer(),
		errorHandler: terrors.DefaultErrorHandler,
		inputType:    InputTypeText,
		prompt:       prompt,
		placeholder:  DefaultPlaceholder,
		width:        DefaultInputWidth,
		maskInput:    false,
		allowEmpty:   false,
	}
}

// WithValidator устанавливает валидатор для задачи
func (t *InputTaskNew) WithValidator(validator validation.Validator) *InputTaskNew {
	t.validator = validator
	return t
}

// WithInputType устанавливает тип ввода
func (t *InputTaskNew) WithInputType(inputType InputType) *InputTaskNew {
	t.inputType = inputType

	// Автоматически настраиваем параметры в зависимости от типа
	switch inputType {
	case InputTypePassword:
		t.maskInput = true
		t.textInput.EchoMode = textinput.EchoPassword
		t.textInput.EchoCharacter = PasswordMask

		// Если валидатор не установлен, используем стандартный валидатор паролей
		if t.validator == nil {
			t.validator = validation.StandardPassword()
		}

	case InputTypeEmail:
		if t.validator == nil {
			t.validator = validation.Email()
		}

	case InputTypeNumber:
		if t.validator == nil {
			t.validator = validation.NewNumberValidator(DefaultNumberMin, DefaultNumberMax)
		}

	case InputTypeIP:
		if t.validator == nil {
			t.validator = validation.IP()
		}

	case InputTypeDomain:
		if t.validator == nil {
			t.validator = validation.Domain()
		}
	}

	return t
}

// WithWidth устанавливает ширину поля ввода
func (t *InputTaskNew) WithWidth(width int) *InputTaskNew {
	if width < MinInputWidth {
		width = MinInputWidth
	} else if width > MaxInputWidth {
		width = MaxInputWidth
	}

	t.width = width
	t.textInput.Width = width
	return t
}

// WithPlaceholder устанавливает текст-заполнитель
func (t *InputTaskNew) WithPlaceholder(placeholder string) *InputTaskNew {
	t.placeholder = placeholder
	t.textInput.Placeholder = placeholder
	return t
}

// WithAllowEmpty разрешает пустые значения
func (t *InputTaskNew) WithAllowEmpty(allow bool) *InputTaskNew {
	t.allowEmpty = allow
	return t
}

// WithRenderer устанавливает пользовательский рендерер
func (t *InputTaskNew) WithRenderer(renderer *InputRenderer) *InputTaskNew {
	t.renderer = renderer
	return t
}

// WithStyle устанавливает стиль для рендерера
func (t *InputTaskNew) WithStyle(style lipgloss.Style) *InputTaskNew {
	t.renderer.WithStyle(style)
	return t
}

// GetValue возвращает введенное значение
func (t *InputTaskNew) GetValue() string {
	return t.value
}

// Run запускает задачу ввода
func (t *InputTaskNew) Run() tea.Cmd {
	return textinput.Blink
}

// Update обновляет состояние задачи
func (t *InputTaskNew) Update(msg tea.Msg) (Task, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			// Отмена ввода
			return t.handleCancel()

		case "enter":
			// Подтверждение ввода
			return t.handleSubmit()
		}

	case error:
		// Обработка ошибки валидации с контекстом
		taskErr := t.errorHandler.Handle(t.title, msg)
		t.validationErr = taskErr
		t.SetError(taskErr)
		return t, nil
	}

	// Обновляем компонент ввода
	t.textInput, cmd = t.textInput.Update(msg)

	// Валидируем введенное значение в реальном времени
	t.validateInput()

	return t, cmd
}

// validateInput выполняет валидацию текущего ввода
func (t *InputTaskNew) validateInput() {
	currentValue := t.textInput.Value()

	// Проверяем пустое значение
	if currentValue == "" && !t.allowEmpty {
		if t.validator != nil {
			// Даем валидатору шанс проверить пустое значение
			if err := t.validator.Validate(currentValue); err != nil {
				t.validationErr = err
				return
			}
		}
		t.validationErr = nil
		return
	}

	// Выполняем валидацию если есть валидатор
	if t.validator != nil {
		if err := t.validator.Validate(currentValue); err != nil {
			// Создаем TaskError с контекстом валидации
			t.validationErr = terrors.NewValidationError(t.title, err).
				WithContext("input_type", t.inputType).
				WithContext("value_length", len(currentValue))
		} else {
			// Если валидация успешна — очищаем ошибку
			t.validationErr = nil
			t.SetError(nil)
		}
	} else {
		t.validationErr = nil
		t.SetError(nil)
	}
}

// handleSubmit обрабатывает подтверждение ввода
func (t *InputTaskNew) handleSubmit() (Task, tea.Cmd) {
	currentValue := t.textInput.Value()

	// Финальная валидация
	if !t.allowEmpty && performance.TrimSpaceEfficient(currentValue) == "" {
		emptyErr := terrors.NewValidationError(t.title, errors.New("поле не может быть пустым")).
			WithContext("required", true)
		t.validationErr = emptyErr
		t.SetError(emptyErr)
		t.textInput.Focus()
		return t, nil
	}

	if t.validator != nil {
		if err := t.validator.Validate(currentValue); err != nil {
			validationErr := terrors.NewValidationError(t.title, err).
				WithContext("final_validation", true).
				WithContext("input_type", t.inputType)
			t.validationErr = validationErr
			t.SetError(validationErr)
			t.textInput.Focus()
			return t, nil
		}
	}

	// Ввод успешен — очищаем ошибки и фиксируем значение
	t.validationErr = nil
	t.SetError(nil)
	t.value = currentValue
	t.done = true
	t.icon = ui.IconDone
	t.finalValue = ui.SuccessLabelStyle.Render(t.getDisplayValue())

	return t, nil
}

// handleCancel обрабатывает отмену ввода
func (t *InputTaskNew) handleCancel() (Task, tea.Cmd) {
	cancelErr := terrors.NewCancelError(t.title).
		WithContext("input_type", t.inputType).
		WithContext("partial_value", t.textInput.Value())

	t.SetError(cancelErr)
	t.done = true
	t.icon = ui.IconCancelled
	t.finalValue = ui.CancelStyle.Render("Отменено")

	return t, nil
}

// getDisplayValue возвращает значение для отображения (маскирует пароли)
func (t *InputTaskNew) getDisplayValue() string {
	if t.maskInput {
		return performance.RepeatEfficient(string(PasswordMask), len(t.value))
	}
	return t.value
}

// View отображает текущее состояние задачи
func (t *InputTaskNew) View(width int) string {
	if t.IsDone() {
		return t.FinalView(width)
	}

	return t.renderer.RenderInput(
		t.title,
		t.textInput,
		t.validator,
		t.validationErr,
		t.inputType,
		width,
	)
}

// FinalView отображает финальное состояние задачи
func (t *InputTaskNew) FinalView(width int) string {
	if t.HasError() {
		return t.renderer.RenderFinal(t.title, "", true, t.Error(), width)
	}

	return t.renderer.RenderFinal(t.title, t.getDisplayValue(), false, nil, width)
}

// InputTaskBuilder предоставляет fluent API для создания InputTask
type InputTaskBuilder struct {
	task *InputTaskNew
}

// NewInputTaskBuilder создает новый построитель задач ввода
func NewInputTaskBuilder(title, prompt string) *InputTaskBuilder {
	return &InputTaskBuilder{
		task: NewInputTaskNew(title, prompt),
	}
}

// Password настраивает задачу для ввода пароля
func (b *InputTaskBuilder) Password() *InputTaskBuilder {
	b.task.WithInputType(InputTypePassword)
	return b
}

// Email настраивает задачу для ввода email
func (b *InputTaskBuilder) Email() *InputTaskBuilder {
	b.task.WithInputType(InputTypeEmail)
	return b
}

// Number настраивает задачу для ввода числа
func (b *InputTaskBuilder) Number(min, max int) *InputTaskBuilder {
	b.task.WithInputType(InputTypeNumber).WithValidator(validation.NewNumberValidator(min, max))
	return b
}

// IP настраивает задачу для ввода IP адреса
func (b *InputTaskBuilder) IP() *InputTaskBuilder {
	b.task.WithInputType(InputTypeIP)
	return b
}

// Domain настраивает задачу для ввода домена
func (b *InputTaskBuilder) Domain() *InputTaskBuilder {
	b.task.WithInputType(InputTypeDomain)
	return b
}

// Required делает поле обязательным
func (b *InputTaskBuilder) Required() *InputTaskBuilder {
	b.task.WithAllowEmpty(false)
	return b
}

// Optional делает поле опциональным
func (b *InputTaskBuilder) Optional() *InputTaskBuilder {
	b.task.WithAllowEmpty(true)
	return b
}

// Width устанавливает ширину поля
func (b *InputTaskBuilder) Width(width int) *InputTaskBuilder {
	b.task.WithWidth(width)
	return b
}

// Placeholder устанавливает текст-заполнитель
func (b *InputTaskBuilder) Placeholder(placeholder string) *InputTaskBuilder {
	b.task.WithPlaceholder(placeholder)
	return b
}

// Validator устанавливает пользовательский валидатор
func (b *InputTaskBuilder) Validator(validator validation.Validator) *InputTaskBuilder {
	b.task.WithValidator(validator)
	return b
}

// Build возвращает готовую задачу
func (b *InputTaskBuilder) Build() *InputTaskNew {
	return b.task
}
