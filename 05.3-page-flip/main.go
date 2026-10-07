package main

import "github.com/belimawr/tonc-tinygo/mmio"

// Compile: tinygo build -o rom.gba -target gameboy-advance .
// Run: mgba rom.gba
// One-shot: tinygo build -o rom.gba -target gameboy-advance . && mgba rom.gba
func main() {
	mmio.REG_DISPCNT.Set(mmio.DCNT_MODE4 | mmio.DCNT_BG2)

	for {
	}
}

func loadGfx() {
}
