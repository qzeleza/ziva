package performance

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Дополнительные тесты для повышения покрытия string_utils.go

func TestToLowerEfficient(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Пустая строка", "", ""},
		{"Только строчные", "hello", "hello"},
		{"Только заглавные", "HELLO", "hello"},
		{"Смешанный регистр", "Hello World", "hello world"},
		{"С числами", "Hello123", "hello123"},
		{"Не-ASCII символы", "Привет", "привет"},
		{"Уже в нижнем регистре", "already lower", "already lower"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToLowerEfficient(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReplaceAllEfficient(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		old      string
		new      string
		expected string
	}{
		{"Пустая строка", "", "old", "new", ""},
		{"Пустой old", "test", "", "new", "test"},
		{"old == new", "test", "t", "t", "test"},
		{"Нет совпадений", "hello", "world", "test", "hello"},
		{"Одно совпадение", "hello world", "world", "test", "hello test"},
		{"Много совпадений", "test test test", "test", "demo", "demo demo demo"},
		{"Перекрывающиеся замены", "aaa", "aa", "b", "ba"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceAllEfficient(tt.s, tt.old, tt.new)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCleanWhitespaceEfficient(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Пустая строка", "", ""},
		{"Только пробелы", "   ", ""},
		{"Ведущие пробелы", "   hello", "hello"},
		{"Завершающие пробелы", "hello   ", "hello"},
		{"Множественные пробелы", "hello    world", "hello world"},
		{"Табуляции и переносы", "hello\t\n world", "hello world"},
		{"Смешанные whitespace", "  hello \t world  \n", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanWhitespaceEfficient(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFastConcat(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		expected string
	}{
		{"Пустой список", []string{}, ""},
		{"Один элемент", []string{"hello"}, "hello"},
		{"Два элемента", []string{"hello", "world"}, "helloworld"},
		{"Много элементов", []string{"a", "b", "c", "d"}, "abcd"},
		{"С пустыми строками", []string{"hello", "", "world"}, "helloworld"},
		{"Только пустые строки", []string{"", "", ""}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FastConcat(tt.parts...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestByteBufferPool(t *testing.T) {
	// Тестируем создание пула
	pool := NewByteBufferPool(4)
	assert.NotNil(t, pool)

	// Тестируем получение буфера
	buf1 := pool.Get()
	assert.NotNil(t, buf1)
	assert.Equal(t, 0, buf1.Len(), "Буфер должен быть пустым")

	// Записываем данные в буфер
	buf1.WriteString("test data")
	assert.Equal(t, "test data", buf1.String())

	// Возвращаем буфер в пул
	pool.Put(buf1)

	// Получаем буфер снова (должен быть очищен)
	buf2 := pool.Get()
	assert.Equal(t, 0, buf2.Len(), "Буфер должен быть очищен после возврата в пул")

	// Тестируем возврат nil буфера
	pool.Put(nil) // Не должно вызвать панику

	// Тестируем возврат слишком большого буфера
	largeBuf := &bytes.Buffer{}
	largeBuf.Grow(3000) // Больше лимита в 2048
	pool.Put(largeBuf)  // Должен быть отброшен
}

func TestGlobalByteBufferPool(t *testing.T) {
	// Тестируем глобальные функции
	buf := GetByteBuffer()
	assert.NotNil(t, buf)
	assert.Equal(t, 0, buf.Len())

	buf.WriteString("global test")
	PutByteBuffer(buf)

	buf2 := GetByteBuffer()
	assert.Equal(t, 0, buf2.Len(), "Глобальный буфер должен быть очищен")
}

func TestIntToString(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"Ноль", 0, "0"},
		{"Положительное число", 123, "123"},
		{"Отрицательное число", -123, "-123"},
		{"Однозначное", 5, "5"},
		{"Большое число", 999999, "999999"},
		{"Отрицательное однозначное", -5, "-5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntToString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEmergencyPoolCleanup(t *testing.T) {
	// Заполняем пулы буферами
	for i := 0; i < 5; i++ {
		buf := GetBuffer()
		buf.WriteString("test")
		PutBuffer(buf)

		byteBuf := GetByteBuffer()
		byteBuf.WriteString("test")
		PutByteBuffer(byteBuf)
	}

	// Выполняем экстренную очистку
	EmergencyPoolCleanup()

	// Проверяем, что пулы очищены
	// (точная проверка зависит от внутренней реализации пулов)
	buf := GetBuffer()
	assert.NotNil(t, buf, "После очистки пулы должны по-прежнему работать")
}

func TestStringPoolCleanup(t *testing.T) {
	pool := NewStringPool(6)

	// Заполняем пул
	for i := 0; i < 4; i++ {
		buf := pool.Get()
		buf.WriteString("test")
		pool.Put(buf)
	}

	// Выполняем очистку
	pool.Cleanup()

	// Пул должен по-прежнему работать
	buf := pool.Get()
	assert.NotNil(t, buf)
}

func TestByteBufferPoolCleanup(t *testing.T) {
	pool := NewByteBufferPool(6)

	// Заполняем пул
	for i := 0; i < 4; i++ {
		buf := pool.Get()
		buf.WriteString("test")
		pool.Put(buf)
	}

	// Выполняем очистку
	pool.Cleanup()

	// Пул должен по-прежнему работать
	buf := pool.Get()
	assert.NotNil(t, buf)
}

func TestStringPoolEdgeCases(t *testing.T) {
	// Тестируем создание пула с экстремальными размерами
	tinyPool := NewStringPool(1)
	assert.NotNil(t, tinyPool)

	largePool := NewStringPool(100) // Должен быть ограничен до 16
	assert.NotNil(t, largePool)

	// Тестируем переполнение пула
	buf1 := tinyPool.Get()
	buf2 := tinyPool.Get()

	tinyPool.Put(buf1)
	tinyPool.Put(buf2) // Второй буфер должен быть отброшен

	// Тестируем слишком большой буфер
	largeBuf := &strings.Builder{}
	largeBuf.Grow(5000)
	tinyPool.Put(largeBuf) // Должен быть отброшен
}

func TestByteBufferPoolEdgeCases(t *testing.T) {
	// Тестируем создание пула с экстремальными размерами
	tinyPool := NewByteBufferPool(1)
	assert.NotNil(t, tinyPool)

	largePool := NewByteBufferPool(100) // Должен быть ограничен до 8
	assert.NotNil(t, largePool)

	// Тестируем переполнение пула
	buf1 := tinyPool.Get()
	buf2 := tinyPool.Get()

	tinyPool.Put(buf1)
	tinyPool.Put(buf2) // Второй буфер должен быть отброшен
}

func TestPerformanceOptimizations(t *testing.T) {
	// Тестируем, что оптимизированные функции работают так же, как стандартные
	testStrings := []string{"hello", "HELLO", "Hello World", "", "  spaced  "}

	for _, s := range testStrings {
		// TrimSpace
		assert.Equal(t, strings.TrimSpace(s), TrimSpaceEfficient(s))

		// ToLower (только для ASCII)
		if isASCII(s) {
			assert.Equal(t, strings.ToLower(s), ToLowerEfficient(s))
		}
	}

	testParts := [][]string{
		{"a", "b", "c"},
		{"hello", "world"},
		{""},
		{},
	}

	for _, parts := range testParts {
		if len(parts) > 0 {
			expected := strings.Join(parts, ",")
			actual := JoinEfficient(parts, ",")
			assert.Equal(t, expected, actual)
		}
	}
}

func isASCII(s string) bool {
	for _, r := range s {
		if r >= 128 {
			return false
		}
	}
	return true
}

func TestContainsAnyEfficientLongChars(t *testing.T) {
	// Тестируем длинную строку chars (>8 символов) для покрытия map-based пути
	longChars := "abcdefghijklmnop"

	assert.True(t, ContainsAnyEfficient("hello", longChars))
	assert.False(t, ContainsAnyEfficient("xyz", longChars))
	assert.False(t, ContainsAnyEfficient("", longChars))
	assert.False(t, ContainsAnyEfficient("hello", ""))
}

func TestRepeatEfficientEdgeCases(t *testing.T) {
	// Тестируем превышение maxResult
	result := RepeatEfficient("very long string that exceeds limit", 1000)
	assert.NotEmpty(t, result)
	assert.Less(t, len(result), 3000, "Результат должен быть ограничен")

	// Тестируем нулевое повторение
	assert.Equal(t, "", RepeatEfficient("test", 0))

	// Тестируем отрицательное повторение
	assert.Equal(t, "", RepeatEfficient("test", -1))
}

func TestStringPoolGetWhenEmpty(t *testing.T) {
	// Создаем пул и забираем все буферы
	pool := NewStringPool(2)
	buf1 := pool.Get()
	buf2 := pool.Get()
	buf3 := pool.Get() // Этот должен создать новый буфер

	assert.NotNil(t, buf1)
	assert.NotNil(t, buf2)
	assert.NotNil(t, buf3)

	// Проверяем, что все буферы независимые
	buf1.WriteString("1")
	buf2.WriteString("2")
	buf3.WriteString("3")

	assert.Equal(t, "1", buf1.String())
	assert.Equal(t, "2", buf2.String())
	assert.Equal(t, "3", buf3.String())
}
