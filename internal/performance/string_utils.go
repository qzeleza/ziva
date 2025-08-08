// Package performance предоставляет оптимизированные утилиты для embedded устройств
package performance

import (
	"bytes"
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

// ToLowerEfficient эффективное приведение к нижнему регистру
func ToLowerEfficient(s string) string {
	if len(s) == 0 {
		return s
	}
	
	// Быстрая проверка - нужно ли изменение
	hasUpper := false
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			hasUpper = true
			break
		}
		if r > 127 { // Не-ASCII символы
			hasUpper = true
			break
		}
	}
	
	if !hasUpper {
		return s
	}
	
	return strings.ToLower(s)
}

// ContainsAnyEfficient проверяет содержание любого из символов
func ContainsAnyEfficient(s string, chars string) bool {
	if len(s) == 0 || len(chars) == 0 {
		return false
	}
	
	// Для коротких строк используем простой поиск
	if len(chars) <= 8 {
		for _, c := range s {
			for _, target := range chars {
				if c == target {
					return true
				}
			}
		}
		return false
	}
	
	// Для длинных строк используем map
	charMap := make(map[rune]bool, len(chars))
	for _, c := range chars {
		charMap[c] = true
	}
	
	for _, c := range s {
		if charMap[c] {
			return true
		}
	}
	
	return false
}

// ReplaceAllEfficient эффективная замена всех вхождений
func ReplaceAllEfficient(s, old, new string) string {
	if len(s) == 0 || len(old) == 0 || old == new {
		return s
	}
	
	// Подсчитываем количество замен
	count := strings.Count(s, old)
	if count == 0 {
		return s
	}
	
	// Для одной замены используем простой способ
	if count == 1 {
		return strings.Replace(s, old, new, 1)
	}
	
	buf := GetBuffer()
	defer PutBuffer(buf)
	
	// Предварительно вычисляем размер результата
	newSize := len(s) + count*(len(new)-len(old))
	if newSize > 0 {
		buf.Grow(newSize)
	}
	
	start := 0
	for {
		idx := strings.Index(s[start:], old)
		if idx == -1 {
			buf.WriteString(s[start:])
			break
		}
		
		buf.WriteString(s[start : start+idx])
		buf.WriteString(new)
		start += idx + len(old)
	}
	
	return buf.String()
}

// CleanWhitespaceEfficient очищает лишние пробелы (оптимизировано для embedded)
func CleanWhitespaceEfficient(s string) string {
	if len(s) == 0 {
		return s
	}
	
	buf := GetBuffer()
	defer PutBuffer(buf)
	
	lastWasSpace := true // Чтобы удалить ведущие пробелы
	
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !lastWasSpace {
				buf.WriteRune(' ')
				lastWasSpace = true
			}
		} else {
			buf.WriteRune(r)
			lastWasSpace = false
		}
	}
	
	result := buf.String()
	// Удаляем завершающий пробел
	if len(result) > 0 && result[len(result)-1] == ' ' {
		result = result[:len(result)-1]
	}
	
	return result
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

// Глобальный пул байтовых буферов
var bytePool = NewByteBufferPool(4)

// GetByteBuffer получает байтовый буфер
func GetByteBuffer() *bytes.Buffer {
	return bytePool.Get()
}

// PutByteBuffer возвращает байтовый буфер
func PutByteBuffer(buf *bytes.Buffer) {
	bytePool.Put(buf)
}

// IntToString эффективно преобразует int в строку для embedded устройств
func IntToString(n int) string {
	if n == 0 {
		return "0"
	}
	
	// Для отрицательных чисел
	negative := n < 0
	if negative {
		n = -n
	}
	
	buf := GetBuffer()
	defer PutBuffer(buf)
	
	// Преобразуем цифры в обратном порядке
	for n > 0 {
		buf.WriteByte(byte('0' + n%10))
		n /= 10
	}
	
	if negative {
		buf.WriteByte('-')
	}
	
	// Переворачиваем строку
	result := buf.String()
	runes := []rune(result)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	
	return string(runes)
}

// EmergencyPoolCleanup выполняет экстренную очистку всех пулов буферов
func EmergencyPoolCleanup() {
	// Очищаем глобальный пул строк
	for {
		select {
		case <-defaultPool.buffers:
		default:
			goto nextPool
		}
	}
	
nextPool:
	// Очищаем глобальный пул байтов
	for {
		select {
		case <-bytePool.buffers:
		default:
			goto cleanupComplete
		}
	}
	
cleanupComplete:
	// Примечание: Для полной очистки также нужно вызвать ui.ClearInternCache()
	// но мы избегаем импорта ui пакета здесь для предотвращения циклических зависимостей
}

// Cleanup для StringPool - экстренная очистка пула строк
func (p *StringPool) Cleanup() {
	targetSize := len(p.buffers) / 2
	if targetSize < 2 {
		targetSize = 2
	}
	
	for len(p.buffers) > targetSize {
		select {
		case <-p.buffers:
		default:
			return
		}
	}
}

// Cleanup для ByteBufferPool - экстренная очистка пула байтов
func (p *ByteBufferPool) Cleanup() {
	targetSize := len(p.buffers) / 2
	if targetSize < 1 {
		targetSize = 1
	}
	
	for len(p.buffers) > targetSize {
		select {
		case <-p.buffers:
		default:
			return
		}
	}
}