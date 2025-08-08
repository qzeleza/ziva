// Package performance предоставляет оптимизированные утилиты для embedded устройств
package performance

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"
)

// StringPool пул для переиспользования строк и буферов (критично для embedded)
type StringPool struct {
	buffers chan *strings.Builder
}

// NewStringPool создает новый пул с указанным размером для embedded устройств
func NewStringPool(size int) *StringPool {
	// Для embedded устройств используем меньший размер пула
	if size > 16 {
		size = 16
	}
	if size < 4 {
		size = 4
	}
	
	pool := &StringPool{
		buffers: make(chan *strings.Builder, size),
	}
	
	// Предварительно заполняем пул
	for i := 0; i < size; i++ {
		pool.buffers <- &strings.Builder{}
	}
	
	return pool
}

// Get получает буфер из пула
func (p *StringPool) Get() *strings.Builder {
	select {
	case buf := <-p.buffers:
		buf.Reset()
		return buf
	default:
		return &strings.Builder{}
	}
}

// Put возвращает буфер в пул
func (p *StringPool) Put(buf *strings.Builder) {
	if buf == nil {
		return
	}
	
	// Для embedded устройств ограничиваем размер буфера
	if buf.Cap() > 4096 {
		return // Не возвращаем слишком большие буферы
	}
	
	select {
	case p.buffers <- buf:
	default:
		// Пул полон, отбрасываем буфер
	}
}

// Глобальный пул для всего приложения
var defaultPool = NewStringPool(8)

// GetBuffer получает буфер из глобального пула
func GetBuffer() *strings.Builder {
	return defaultPool.Get()
}

// PutBuffer возвращает буфер в глобальный пул
func PutBuffer(buf *strings.Builder) {
	defaultPool.Put(buf)
}

// StripANSILength возвращает длину строки без ANSI escape последовательностей
var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)

func StripANSILength(s string) int {
	cleaned := ansiRegex.ReplaceAllString(s, "")
	return len([]rune(cleaned))
}

// TrimSpaceEfficient эффективная замена strings.TrimSpace для embedded
func TrimSpaceEfficient(s string) string {
	if len(s) == 0 {
		return s
	}
	
	start := 0
	end := len(s)
	
	// Ищем начало без пробелов
	for start < end && unicode.IsSpace(rune(s[start])) {
		start++
	}
	
	// Ищем конец без пробелов
	for end > start && unicode.IsSpace(rune(s[end-1])) {
		end--
	}
	
	if start == 0 && end == len(s) {
		return s // Нет изменений
	}
	
	return s[start:end]
}

// JoinEfficient эффективная замена strings.Join для embedded
func JoinEfficient(parts []string, separator string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	
	buf := GetBuffer()
	defer PutBuffer(buf)
	
	buf.WriteString(parts[0])
	for i := 1; i < len(parts); i++ {
		buf.WriteString(separator)
		buf.WriteString(parts[i])
	}
	
	return buf.String()
}

// RepeatEfficient эффективная замена strings.Repeat для embedded
func RepeatEfficient(s string, count int) string {
	if count <= 0 || len(s) == 0 {
		return ""
	}
	if count == 1 {
		return s
	}
	
	// Для embedded устройств ограничиваем размер результата
	const maxResult = 2048
	if len(s)*count > maxResult {
		count = maxResult / len(s)
		if count <= 0 {
			count = 1
		}
	}
	
	buf := GetBuffer()
	defer PutBuffer(buf)
	
	buf.Grow(len(s) * count)
	for i := 0; i < count; i++ {
		buf.WriteString(s)
	}
	
	return buf.String()
}

// FastConcat быстрая конкатенация строк для embedded
func FastConcat(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}
	
	// Вычисляем общую длину
	totalLen := 0
	for _, part := range parts {
		totalLen += len(part)
	}
	
	if totalLen == 0 {
		return ""
	}
	
	buf := GetBuffer()
	defer PutBuffer(buf)
	
	buf.Grow(totalLen)
	for _, part := range parts {
		buf.WriteString(part)
	}
	
	return buf.String()
}

// ByteBufferPool пул байтовых буферов (еще более эффективен для embedded)
type ByteBufferPool struct {
	buffers chan *bytes.Buffer
}

// NewByteBufferPool создает пул байтовых буферов
func NewByteBufferPool(size int) *ByteBufferPool {
	if size > 8 {
		size = 8 // Для embedded еще меньше
	}
	if size < 2 {
		size = 2
	}
	
	pool := &ByteBufferPool{
		buffers: make(chan *bytes.Buffer, size),
	}
	
	for i := 0; i < size; i++ {
		pool.buffers <- &bytes.Buffer{}
	}
	
	return pool
}

// Get получает байтовый буфер
func (p *ByteBufferPool) Get() *bytes.Buffer {
	select {
	case buf := <-p.buffers:
		buf.Reset()
		return buf
	default:
		return &bytes.Buffer{}
	}
}

// Put возвращает буфер в пул
func (p *ByteBufferPool) Put(buf *bytes.Buffer) {
	if buf == nil || buf.Cap() > 2048 {
		return
	}
	
	select {
	case p.buffers <- buf:
	default:
	}
}

// EmergencyPoolCleanup выполняет экстренную очистку всех пулов буферов
func EmergencyPoolCleanup() {
	// Очищаем глобальный пул строк
	for {
		select {
		case <-defaultPool.buffers:
		default:
			return
		}
	}
}