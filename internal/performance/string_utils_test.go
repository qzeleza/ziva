package performance

import (
	"strings"
	"testing"
)

func BenchmarkTrimSpaceEfficient(b *testing.B) {
	testCases := []string{
		"   hello world   ",
		"no spaces",
		"      ",
		"",
		"   start",
		"end   ",
		"  mixed  spaces  here  ",
	}
	
	b.Run("Efficient", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, s := range testCases {
				_ = TrimSpaceEfficient(s)
			}
		}
	})
	
	b.Run("Standard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, s := range testCases {
				_ = strings.TrimSpace(s)
			}
		}
	})
}

func BenchmarkJoinEfficient(b *testing.B) {
	parts := []string{"one", "two", "three", "four", "five"}
	separator := ", "
	
	b.Run("Efficient", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = JoinEfficient(parts, separator)
		}
	})
	
	b.Run("Standard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = strings.Join(parts, separator)
		}
	})
}

func BenchmarkRepeatEfficient(b *testing.B) {
	text := "test"
	count := 100
	
	b.Run("Efficient", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = RepeatEfficient(text, count)
		}
	})
	
	b.Run("Standard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = strings.Repeat(text, count)
		}
	})
}

func BenchmarkStringPool(b *testing.B) {
	pool := NewStringPool(8)
	
	b.Run("WithPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := pool.Get()
			buf.WriteString("test")
			buf.WriteString(" ")
			buf.WriteString("message")
			_ = buf.String()
			pool.Put(buf)
		}
	})
	
	b.Run("WithoutPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var buf strings.Builder
			buf.WriteString("test")
			buf.WriteString(" ")
			buf.WriteString("message")
			_ = buf.String()
		}
	})
}

func BenchmarkFastConcat(b *testing.B) {
	parts := []string{"part1", "part2", "part3", "part4", "part5"}
	
	b.Run("FastConcat", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = FastConcat(parts...)
		}
	})
	
	b.Run("StandardConcat", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			result := ""
			for _, part := range parts {
				result += part
			}
			_ = result
		}
	})
	
	b.Run("StringsJoin", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = strings.Join(parts, "")
		}
	})
}

func BenchmarkContainsAnyEfficient(b *testing.B) {
	text := "This is a test string with various characters!"
	chars := "xyz123"
	
	b.Run("Efficient", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ContainsAnyEfficient(text, chars)
		}
	})
	
	b.Run("Standard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = strings.ContainsAny(text, chars)
		}
	})
}

func BenchmarkCleanWhitespaceEfficient(b *testing.B) {
	text := "  This   has    multiple    spaces   between   words  "
	
	b.Run("Efficient", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = CleanWhitespaceEfficient(text)
		}
	})
	
	b.Run("Standard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Стандартная обработка с регулярными выражениями
			result := strings.TrimSpace(text)
			parts := strings.Fields(result)
			_ = strings.Join(parts, " ")
		}
	})
}

func TestTrimSpaceEfficient(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"   hello   ", "hello"},
		{"no spaces", "no spaces"},
		{"", ""},
		{"   ", ""},
		{"start   ", "start"},
		{"   end", "end"},
		{"\t\n  hello  \r\n", "hello"},
	}
	
	for _, test := range tests {
		result := TrimSpaceEfficient(test.input)
		if result != test.expected {
			t.Errorf("TrimSpaceEfficient(%q) = %q, want %q", test.input, result, test.expected)
		}
	}
}

func TestJoinEfficient(t *testing.T) {
	tests := []struct {
		parts     []string
		separator string
		expected  string
	}{
		{[]string{"a", "b", "c"}, ",", "a,b,c"},
		{[]string{"one"}, ",", "one"},
		{[]string{}, ",", ""},
		{[]string{"a", "b"}, " - ", "a - b"},
	}
	
	for _, test := range tests {
		result := JoinEfficient(test.parts, test.separator)
		if result != test.expected {
			t.Errorf("JoinEfficient(%v, %q) = %q, want %q", test.parts, test.separator, result, test.expected)
		}
	}
}

func TestRepeatEfficient(t *testing.T) {
	tests := []struct {
		input    string
		count    int
		expected string
	}{
		{"a", 3, "aaa"},
		{"test", 0, ""},
		{"", 5, ""},
		{"x", 1, "x"},
	}
	
	for _, test := range tests {
		result := RepeatEfficient(test.input, test.count)
		if result != test.expected {
			t.Errorf("RepeatEfficient(%q, %d) = %q, want %q", test.input, test.count, result, test.expected)
		}
	}
}

func TestContainsAnyEfficient(t *testing.T) {
	tests := []struct {
		text     string
		chars    string
		expected bool
	}{
		{"hello", "xyz", false},
		{"hello", "el", true},
		{"test", "", false},
		{"", "abc", false},
		{"abc", "c", true},
	}
	
	for _, test := range tests {
		result := ContainsAnyEfficient(test.text, test.chars)
		if result != test.expected {
			t.Errorf("ContainsAnyEfficient(%q, %q) = %v, want %v", test.text, test.chars, result, test.expected)
		}
	}
}

// Тест для embedded-специфичных ограничений
func TestEmbeddedLimits(t *testing.T) {
	// Проверяем, что пул не создает слишком много объектов
	largePool := NewStringPool(100)
	if cap(largePool.buffers) > 16 {
		t.Errorf("StringPool size should be limited to 16 for embedded devices, got %d", cap(largePool.buffers))
	}
	
	// Проверяем ограничение размера результата RepeatEfficient
	longResult := RepeatEfficient("a", 10000)
	if len(longResult) > 2048 {
		t.Errorf("RepeatEfficient should limit result size for embedded devices")
	}
}