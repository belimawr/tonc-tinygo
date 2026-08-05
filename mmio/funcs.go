package mmio

import (
	"runtime/volatile"
	"unsafe"
)

// mem16 returns a pointer to a volatile 16-bit register at the given address.
// All VRAM and palette RAM writes go through this to prevent optimizer caching.
func Mem16(addr uintptr) *volatile.Register16 {
	return (*volatile.Register16)(unsafe.Pointer(addr))
}
