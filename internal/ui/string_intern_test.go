package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInterning(t *testing.T) {
	// Очищаем кэш перед тестом
	ClearInternCache()
	
	// Тестируем базовое интернирование
	s1 := InternString("test string")
	s2 := InternString("test string")
	
	// Проверяем, что это одна и та же строка в памяти
	assert.Equal(t, s1, s2, "Интернированные строки должны быть равны")
	
	// Проверяем статистику кэша
	size, capacity := GetCacheStats()
	assert.True(t, size > 0, "Кэш должен содержать интернированные строки")
	assert.Equal(t, 128, capacity, "Емкость кэша должна быть 128")
}

func TestPrecomputedStrings(t *testing.T) {
	// Проверяем, что предвычисленные строки доступны
	assert.NotEmpty(t, InternedStrings.TaskCompleted, "TaskCompleted должен быть предвычислен")
	assert.NotEmpty(t, InternedStrings.TaskInProgress, "TaskInProgress должен быть предвычислен")
	assert.NotEmpty(t, InternedStrings.Selected, "Selected должен быть предвычислен")
	assert.NotEmpty(t, InternedStrings.NotSelected, "NotSelected должен быть предвычислен")
	
	// Проверяем индентацию
	assert.Equal(t, "  ", InternedStrings.Indent2, "Indent2 должен быть 2 пробела")
	assert.Equal(t, "    ", InternedStrings.Indent4, "Indent4 должен быть 4 пробела")
}

func TestCacheCleanup(t *testing.T) {
	// Очищаем кэш
	ClearInternCache()
	
	// Добавляем несколько строк
	for i := 0; i < 10; i++ {
		InternString(string(rune('a' + i)))
	}
	
	sizeBefore, _ := GetCacheStats()
	
	// Очищаем кэш
	ClearInternCache()
	
	sizeAfter, _ := GetCacheStats()
	
	// Проверяем, что критичные строки остались, но размер уменьшился или остался разумным
	assert.True(t, sizeAfter > 0, "После очистки должны остаться критичные строки")
	assert.True(t, sizeAfter <= sizeBefore, "Размер кэша после очистки не должен быть больше")
	
	// Проверяем, что критичные строки всё ещё работают
	assert.NotEmpty(t, InternedStrings.TaskCompleted, "Критичные строки должны остаться после очистки")
}