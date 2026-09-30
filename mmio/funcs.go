//go:build !tinydebug

package mmio

import (
	"runtime/volatile"
	"unsafe"

	"github.com/belimawr/tonc-tinygo/types"
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

// SetM16 sets a 16-bit memory address
func SetM16Colour(addr int, c types.Colour) {
	M16(uintptr(addr)).Set(uint16(c))
}

// SetM8 sets a 8-bit memory addres
func SetM8(addr uintptr, v uint8) {
	M8(addr).Set(v)
}
