package main

import (
	"unsafe"

	_ "embed"

	"github.com/belimawr/tonc-tinygo/mmio"
	"github.com/belimawr/tonc-tinygo/types"
)

//go:embed bin.dat
var myImg []byte

var CLR_BLACK = types.NewColour(0, 0, 0)
var CLR_BLUE = types.NewColour(0, 0, 31)
var CLR_LIME = types.NewColour(0, 31, 0)
var CLR_RED = types.NewColour(31, 0, 0)
var CLR_WHITE = types.NewColour(31, 31, 31)

func memcpy(dst uintptr, src []byte) {
	copy(unsafe.Slice((*byte)(unsafe.Pointer(dst)), len(src)), src)
}

// Compile: tinygo build -o rom.gba -target gameboy-advance .
// Run: mgba rom.gba
// One-shot: tinygo build -o rom.gba -target gameboy-advance . && mgba rom.gba
func main() {
	mmio.REG_DISPCNT.Set(mmio.DCNT_MODE3 | mmio.DCNT_BG2)
	memcpy(mmio.MEM_VRAM, myImg)

	// Manually write the same data aligning it with the second line
	// of the memory inspector from mgba making it easier to debug
	mmio.SetM16Colour(int(mmio.MEM_VRAM+16+0), CLR_BLUE)
	mmio.SetM16Colour(int(mmio.MEM_VRAM+16+2), CLR_RED)
	mmio.SetM16Colour(int(mmio.MEM_VRAM+16+4), CLR_LIME)
	mmio.SetM16Colour(int(mmio.MEM_VRAM+16+6), CLR_WHITE)
	mmio.SetM16Colour(int(mmio.MEM_VRAM+16+8), CLR_BLACK)

	for {
	}
}

func loadGfx() {
}
