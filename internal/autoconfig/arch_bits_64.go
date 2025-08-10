//go:build amd64 || arm64 || mips64 || mips64le || ppc64 || ppc64le || riscv64 || s390x

package autoconfig

// ArchBits содержит разрядность целевой платформы.
const ArchBits = 64

// Is64Bit возвращает true для 64-битных архитектур.
func Is64Bit() bool { return true }
