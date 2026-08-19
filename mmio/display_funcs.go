package mmio

import (
	"github.com/belimawr/tonc-tinygo/types"
)

const (
	// The documented screen size
	width  = 240
	height = 160
)

// SetPixelM3 sets the pixel at (x,y) to colour c in display mode 3
func SetPixelM3(x, y int, c types.Colour) {
	// VRAM uses 2-byte (16 bit) words, so we need to double the offset
	offset := uintptr(x+y*width) * 2
	M16(uintptr(MEM_VRAM + offset)).Set(uint16(c))
}
