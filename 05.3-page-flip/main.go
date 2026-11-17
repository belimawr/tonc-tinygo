package main

import (
	_ "embed"
	"unsafe"

	"github.com/belimawr/tonc-tinygo/mmio"
)

const HalfKByte = 0x0200

var front32Start = (*uint32)(unsafe.Pointer(mmio.MEM_VRAM))
var back32Start = (*uint32)(unsafe.Pointer(mmio.MEM_VRAM + 0xA000))
var pallet32Start = (*uint32)(unsafe.Pointer(mmio.MEM_PALRAM))

var FRONT32 = unsafe.Slice(front32Start, 40*(HalfKByte/2))
var BACK32 = unsafe.Slice(back32Start, 40*(HalfKByte/2))
var PALLET32 = unsafe.Slice(pallet32Start, (HalfKByte / 2))

// Compile: tinygo build -o rom.gba -target gameboy-advance .
// Run: mgba rom.gba
// One-shot: tinygo build -o rom.gba -target gameboy-advance . && mgba rom.gba
func main() {
	loadgfx()

	mmio.REG_DISPCNT.Set(mmio.DCNT_MODE4 | mmio.DCNT_BG2)
	vptr := (*uint16)(unsafe.Pointer(uintptr(0x0400_0000)))
	for {
		// mmio.REG_DISPCNT.Set(mmio.DCNT_MODE4 | mmio.DCNT_BG2)
		*vptr ^= 0x0010
		for range 60 {
			vidVsync()
		}
		*vptr ^= 0x0010

		// mmio.REG_DISPCNT.Set(mmio.DCNT_MODE4 | mmio.DCNT_BG2 | 0x0010)
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

func loadgfx() {
	for i := range len(frontPal) {
		PALLET32[i] = frontPal[i]
	}

	for ii := range 16 {
		for i := range 144 / 4 {
			FRONT32[ii*60+i] = frontBitmap[(ii*144/4)+i]
			BACK32[ii*60+i] = backBitmap[(ii*144/4)+i]
		}
	}
}
