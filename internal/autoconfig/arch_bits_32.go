//go:build 386 || arm || mips || mipsle || ppc || wasm

package autoconfig

// ArchBits содержит разрядность целевой платформы.
const ArchBits = 32

// Is64Bit возвращает false для 32-битных архитектур.
func Is64Bit() bool { return false }
