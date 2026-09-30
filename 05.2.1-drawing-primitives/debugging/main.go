package main

import (
	"github.com/belimawr/tonc-tinygo/mmio"
	"github.com/belimawr/tonc-tinygo/types"
)

const M3_WIDTH = 240

const MEM_VRAM = 0x0600_0000

// bmp16_line draws a line on a 16bpp canvas
func bmp16_line(x1, y1, x2, y2 int, clr types.Colour, dstBase, dstPitch int) {
	var dx, dy, xstep, ystep int
	dst := dstBase + y1*dstPitch + x1*2
	// dstPitch /= 2

	// --- Normalisation ---
	if x1 > x2 {
		xstep = -2
		dx = x1 - x2
	} else {
		xstep = +2
		dx = x2 - x1
	}

	if y1 > y2 {
		ystep = -dstPitch
		dy = y1 - y2
	} else {
		ystep = +dstPitch
		dy = y2 - y1
	}

	// --- Drawing ---
	if dy == 0 { // Horizontal
		for ii := 0; ii <= dx; ii++ {
			mmio.SetM16Colour(dst+ii*xstep, clr)
		}
	} else if dx == 0 { // Vertical
		for ii := 0; ii <= dy; ii++ {
			mmio.SetM16Colour(dst+ii*ystep, clr)
		}
	} else if dx >= dy { // Diagonal, slope <= 1
		dd := 2*dy - dx
		for ii := 0; ii <= dx; ii++ {
			mmio.SetM16Colour(dst, clr)
			if dd >= 0 {
				dd -= 2 * dx
				dst += ystep
			}
			dd += 2 * dy
			dst += xstep
		}
	} else { // Diagonal, slope > 1
		dd := 2*dx - dy
		for ii := 0; ii <= dy; ii++ {
			mmio.SetM16Colour(dst, clr)
			if dd >= 0 {
				dd -= 2 * dy
				dst += xstep
			}
			dd += 2 * dx
			dst += ystep
		}
	}
}

func main() {
	var CLR_RED = types.NewColour(31, 0, 0)
	// bmp16_line(0, 1, 20, 1, CLR_RED, int(0), M3_WIDTH*2)

	bmp16_line(0, 0, 0, 10, CLR_RED, int(0), M3_WIDTH*2)
}
