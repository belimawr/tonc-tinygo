package mmio

import (
	"runtime/volatile"
	"unsafe"
)

// M16 returns a pointer to a volatile 16-bit register at the given address.
// All VRAM and palette RAM writes go through this to prevent optimizer caching.
func M16(addr uintptr) *volatile.Register16 {
	return (*volatile.Register16)(unsafe.Pointer(addr))
}

// M8 returns a pointer to a volatile 8-bit register at the given address.
// All VRAM and palette RAM writes go through this to prevent optimizer caching.
func M8(addr uintptr) *volatile.Register8 {
	return (*volatile.Register8)(unsafe.Pointer(addr))
}
