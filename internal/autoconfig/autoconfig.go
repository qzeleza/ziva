// Package autoconfig — автоматическая детекция окружения и применение embedded-конфигурации
package autoconfig

import (
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/qzeleza/termos/internal/ui"
)

// init вызывается при импортировании пакета и запускает авто-конфигурацию
func init() {
	AutoConfigure()
}

// AutoConfigure выполняет авто-детекцию embedded-окружения и применяет настройки по умолчанию
func AutoConfigure() {
	if isEmbeddedEnvironment() {
		applyEmbeddedDefaults()
	}
}

// applyEmbeddedDefaults включает минимально инвазивные embedded-настройки ядра
// На текущем этапе это переключение упрощённой палитры и иконок для терминалов
// с ограниченной поддержкой цветов/UTF-8. По мере необходимости сюда можно
// добавить другие безопасные оптимизации по умолчанию.
func applyEmbeddedDefaults() {
	ui.EnableEmbeddedMode()
}

// isEmbeddedEnvironment автоматически определяет, является ли окружение embedded
func isEmbeddedEnvironment() bool {
	// Принудительная установка через переменную окружения
	if embedded := os.Getenv("TERMOS_EMBEDDED"); embedded != "" {
		return embedded == "true" || embedded == "1"
	}

	// Автоматическое определение по набору эвристик
	return isMemoryConstrained() ||
		isLimitedTerminal() ||
		isKnownEmbeddedEnvironment()
}

// isMemoryConstrained проверяет ограничения памяти
func isMemoryConstrained() bool {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Принудительный лимит из переменной окружения (например, 64MB, 128KB, 1GB)
	if limitStr := os.Getenv("TERMOS_MEMORY_LIMIT"); limitStr != "" {
		if limit, err := parseMemoryLimit(limitStr); err == nil {
			return m.Sys < limit
		}
	}

	// Консервативный порог по умолчанию
	const memoryThreshold = 512 * 1024 * 1024 // 512MB
	return m.Sys < memoryThreshold
}

// parseMemoryLimit парсит строку с лимитом памяти (например "64MB", "128KB", "1GB")
func parseMemoryLimit(limitStr string) (uint64, error) {
	limitStr = strings.ToUpper(strings.TrimSpace(limitStr))

	var multiplier uint64 = 1
	var numStr string

	switch {
	case strings.HasSuffix(limitStr, "MB"):
		multiplier = 1024 * 1024
		numStr = strings.TrimSuffix(limitStr, "MB")
	case strings.HasSuffix(limitStr, "KB"):
		multiplier = 1024
		numStr = strings.TrimSuffix(limitStr, "KB")
	case strings.HasSuffix(limitStr, "GB"):
		multiplier = 1024 * 1024 * 1024
		numStr = strings.TrimSuffix(limitStr, "GB")
	default:
		numStr = limitStr
	}

	num, err := strconv.ParseUint(numStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return num * multiplier, nil
}

// isLimitedTerminal проверяет ограничения терминала и локали
func isLimitedTerminal() bool {
	term := os.Getenv("TERM")

	// Принудительное отключение UTF-8
	if os.Getenv("TERMOS_ASCII_ONLY") == "true" {
		return true
	}

	// Распространённые ограниченные терминалы
	if term == "linux" || term == "console" || term == "vt100" || term == "vt102" {
		return true
	}

	// Проверяем поддержку UTF-8
	lang := strings.ToLower(os.Getenv("LANG"))
	lcAll := strings.ToLower(os.Getenv("LC_ALL"))
	lcCtype := strings.ToLower(os.Getenv("LC_CTYPE"))

	return !strings.Contains(lang, "utf") &&
		!strings.Contains(lcAll, "utf") &&
		!strings.Contains(lcCtype, "utf")
}

// isKnownEmbeddedEnvironment проверяет известные признаки embedded-окружений
func isKnownEmbeddedEnvironment() bool {
	// Характерные файлы/пути OpenWrt/Entware и др.
	embeddedPaths := []string{
		"/opt/etc/init.d",      // Entware
		"/etc/openwrt_release", // OpenWrt
	}
	for _, path := range embeddedPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	// Проверяем архитектуру CPU
	return isEmbeddedCPU()
}

// isEmbeddedCPU проверяет архитектуру CPU для embedded устройств
func isEmbeddedCPU() bool {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return false
	}
	cpuInfo := strings.ToLower(string(data))

	// Характерные признаки embedded процессоров
	embeddedCPUs := []string{
		"arm", "mips", "aarch64", "armv7", "cortex",
		"mediatek", "qualcomm", "broadcom", "rockchip",
	}
	for _, cpu := range embeddedCPUs {
		if strings.Contains(cpuInfo, cpu) {
			return true
		}
	}
	return false
}
