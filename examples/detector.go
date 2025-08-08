package examples

import (
	"os"
	"runtime"
	"strconv"
	"strings"
)

// IsEmbeddedEnvironment автоматически определяет, является ли окружение embedded
func IsEmbeddedEnvironment() bool {
	// Принудительная установка через переменную окружения
	if embedded := os.Getenv("TERMOS_EMBEDDED"); embedded != "" {
		return embedded == "true" || embedded == "1"
	}

	// Автоматическое определение
	return isMemoryConstrained() ||
		isLimitedTerminal() ||
		isKnownEmbeddedEnvironment()
}

// isMemoryConstrained проверяет ограничения памяти
func isMemoryConstrained() bool {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// Принудительный лимит из переменной окружения
	if limitStr := os.Getenv("TERMOS_MEMORY_LIMIT"); limitStr != "" {
		if limit, err := parseMemoryLimit(limitStr); err == nil {
			return m.Sys < limit
		}
	}
	
	const memoryThreshold = 512 * 1024 * 1024 // 512MB
	return m.Sys < memoryThreshold
}

// parseMemoryLimit парсит строку с лимитом памяти (например "64MB", "128KB")
func parseMemoryLimit(limitStr string) (uint64, error) {
	limitStr = strings.ToUpper(strings.TrimSpace(limitStr))
	
	var multiplier uint64 = 1
	var numStr string
	
	if strings.HasSuffix(limitStr, "MB") {
		multiplier = 1024 * 1024
		numStr = strings.TrimSuffix(limitStr, "MB")
	} else if strings.HasSuffix(limitStr, "KB") {
		multiplier = 1024
		numStr = strings.TrimSuffix(limitStr, "KB")
	} else if strings.HasSuffix(limitStr, "GB") {
		multiplier = 1024 * 1024 * 1024
		numStr = strings.TrimSuffix(limitStr, "GB")
	} else {
		numStr = limitStr
	}
	
	num, err := strconv.ParseUint(numStr, 10, 64)
	if err != nil {
		return 0, err
	}
	
	return num * multiplier, nil
}

// isLimitedTerminal проверяет ограничения терминала  
func isLimitedTerminal() bool {
	term := os.Getenv("TERM")
	
	// Принудительное отключение UTF-8
	if os.Getenv("TERMOS_ASCII_ONLY") == "true" {
		return true
	}
	
	// Ограниченные терминалы
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

// isKnownEmbeddedEnvironment проверяет известные embedded окружения
func isKnownEmbeddedEnvironment() bool {
	// Проверяем наличие характерных файлов OpenWrt/Entware
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