package main

import (
	_ "embed"
	"unsafe"

	"github.com/belimawr/tonc-tinygo/mmio"
)

//go:embed gopher01.palette
var palette []byte

//go:embed gopher01.dat
var gopher01 []byte

//go:embed gopher02.dat
var gopher02 []byte

func Memcp16(dst uintptr, src []byte) {
	dstPtr := (*uint16)(unsafe.Pointer(dst))
	src16 := unsafe.Slice((*uint16)(unsafe.Pointer(&src[0])), len(src)/2)
	dst16 := unsafe.Slice(dstPtr, len(src)/2)
	copy(dst16, src16)
}

func Copy16(dst []uint16, src []byte) {
	src16 := unsafe.Slice((*uint16)(unsafe.Pointer(&src[0])), len(src)/2)
	copy(dst, src16)
}

// Compile: tinygo build -o rom.gba -target gameboy-advance .
// Run: mgba rom.gba
// One-shot: tinygo build -o rom.gba -target gameboy-advance . && mgba rom.gba
func main() {
	mmio.REG_DISPCNT.Set(mmio.DCNT_MODE4 | mmio.DCNT_BG2)

	Memcp16(mmio.MEM_PALRAM, palette)
	Memcp16(mmio.MEM_VRAM, gopher01)
	Memcp16(mmio.MEM_VRAM+0xA000, gopher02)

	vptr := (*uint16)(unsafe.Pointer(uintptr(0x0400_0000)))
	for {
		// mmio.REG_DISPCNT.Set(mmio.DCNT_MODE4 | mmio.DCNT_BG2)
		*vptr ^= 0x0010
		for range 60 {
			vidVsync()
		}

		// mmio.REG_DISPCNT.Set(mmio.DCNT_MODE4 | mmio.DCNT_BG2 | 0x0010)
		*vptr ^= 0x0010
		for range 60 {
			vidVsync()
		}
	}
}

func vidVsync() {
	for mmio.REG_VCOUNT.Get() >= 160 {
	}
	for mmio.REG_VCOUNT.Get() < 160 {
	}
}
