package main

import (
	"github.com/belimawr/tonc-tinygo/mmio"
	"github.com/belimawr/tonc-tinygo/types"
)

// Compile: tinygo build -o rom.gba -target gameboy-advance .
// Run: mgba rom.gba
// One-shot: tinygo build -o rom.gba -target gameboy-advance . && mgba rom.gba
func main() {
	mmio.REG_DISPCNT.Set(mmio.DCNT_MODE3 | mmio.DCNT_BG2)
	mmio.SetPixelM3(120, 80, types.NewColour(31, 0, 0)) // red
	mmio.SetPixelM3(136, 80, types.NewColour(0, 31, 0)) // blue
	mmio.SetPixelM3(120, 96, types.NewColour(0, 0, 31)) // green

	for {
	}
}
