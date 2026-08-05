package mmio

import (
	"github.com/belimawr/tonc-tinygo/types"
)

const (
	// The documented screen size
	// width  = 240
	// height = 160

	// What works for me (real device && emulator) in video mode 3
	width  = 480
	height = 320
)

// SetPixelM3 sets the pixel at (x,y) to colour c in display mode 3
func SetPixelM3(x, y int, c types.Colour) {
	Mem16(uintptr(VRAM + uintptr(x+y*width))).Set(uint16(c))
}
