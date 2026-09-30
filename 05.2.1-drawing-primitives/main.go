package main

import (
	"github.com/belimawr/tonc-tinygo/mmio"
	"github.com/belimawr/tonc-tinygo/types"
)

// bmp16_line draws a line on a 16bpp canvas
func bmp16_line(x1, y1, x2, y2 int, clr types.Colour, dstBase, dstPitch int) {
	var dx, dy, xstep, ystep int
	dst := dstBase + y1*dstPitch + x1*2

	// This is based on a C algorithm dstPitch is the array offset
	// for the 16 bits type, in Go we're dealing directly with the 8 bits
	// memory addresses and 16 bits colurs, so we need to multiply it by 2.
	// Which is the same as skipping the divide by 2 step of the
	// original algorithm.
	// dstPitch /= 2

	// --- Normalisation ---
	// Same thing here for xstep and ystep we need to advance:
	// 16 bits = 2 bytes = 2 memory addresses so use +2/-2 instead of +1/-1
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

// bmp16_rec draws a rectangle on a 16bpp canvas
func bmp16_rec(left, top, right, bottom int, clr types.Colour, dstBase, dstPitch int) {
	width, height := uint(right-left), uint(bottom-top)
	dst := dstBase + top*dstPitch + left*2
	dstPitch /= 2

	// --- Draw ---
	for iy := 0; iy < int(height); iy++ {
		for ix := 0; ix < int(width); ix++ {
			// This is based on a C algorithm ix and iy are the array offset
			// for the 16 bits type, in Go we're dealing directly with the 8 bits
			// memory addresses, so multiply dx and dy by 2
			mmio.SetM16Colour(dst+(iy*2)*dstPitch+(ix*2), clr)
		}
	}
}

// bmp16_frame draws a frame on a 16bpp canvas
func bmp16_frame(left, top, right, bottom int, clr types.Colour, dstBase, dstPitch int) {
	// Frame is RB exclusive
	right--
	bottom--

	bmp16_line(left, top, right, top, clr, dstBase, dstPitch)
	bmp16_line(left, bottom, right, bottom, clr, dstBase, dstPitch)

	bmp16_line(left, top, left, bottom, clr, dstBase, dstPitch)
	bmp16_line(right, top, right, bottom, clr, dstBase, dstPitch)
}

const M3_WIDTH = 240

func m3_plot(x, y int, clr types.Colour) {
	mmio.SetPixelM3(x, y, clr)
}

func m3_line(x1, y1, x2, y2 int, clr types.Colour) {
	bmp16_line(x1, y1, x2, y2, clr, int(mmio.MEM_VRAM), M3_WIDTH*2)
}

func m3_rect(left, top, right, bottom int, clr types.Colour) {
	bmp16_rec(left, top, right, bottom, clr, int(mmio.MEM_VRAM), M3_WIDTH*2)
}

func m3_frame(left, top, right, bottom int, clr types.Colour) {
	bmp16_frame(left, top, right, bottom, clr, int(mmio.MEM_VRAM), M3_WIDTH*2)
}

// m3_fill fills the whole screen with one colour
// This is not the most efficient way. Ideally we'd write 32 bit colours and
// do pointer arithmetic, however, we're writing Go, not C.
func m3_fill(clr types.Colour) {
	const width = 240
	const height = 160
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			mmio.SetPixelM3(x, y, clr)
		}
	}
}

var CLR_BLACK = types.NewColour(0, 0, 0)
var CLR_RED = types.NewColour(31, 0, 0)
var CLR_LIME = types.NewColour(0, 31, 0)
var CLR_YELLOW = types.NewColour(31, 31, 0)
var CLR_BLUE = types.NewColour(0, 0, 31)
var CLR_MAGENTA = types.NewColour(31, 0, 31)
var CLR_CYAN = types.NewColour(0, 31, 31)
var CLR_WHITE = types.NewColour(31, 31, 31)

// Compile: tinygo build -o rom.gba -target gameboy-advance .
// Run: mgba rom.gba
// One-shot: tinygo build -o rom.gba -target gameboy-advance . && mgba rom.gba
func main() {
	mmio.REG_DISPCNT.Set(mmio.DCNT_MODE3 | mmio.DCNT_BG2)

	// Fill screen with grey colour
	m3_fill(types.NewColour(12, 12, 14))

	// Rectangles:
	m3_rect(12, 8, 108, 72, CLR_RED)
	m3_rect(108, 72, 132, 88, CLR_LIME)
	m3_rect(132, 88, 228, 152, CLR_BLUE)

	// Rectangle frames
	m3_frame(132, 8, 228, 72, CLR_CYAN)
	m3_frame(109, 73, 131, 87, CLR_BLACK)
	m3_frame(12, 88, 108, 152, CLR_YELLOW)

	// Lines in top right frame
	for ii := 0; ii <= 8; ii++ {
		jj := 3*ii + 7
		m3_line(132+11*ii, 9, 226, 12+7*ii, types.NewColour(uint16(jj), 0, uint16(jj)))
		m3_line(226-11*ii, 70, 133, 69-7*ii, types.NewColour(uint16(jj), 0, uint16(jj)))
	}

	// Lines in bottom left frame
	for ii := 0; ii <= 8; ii++ {
		jj := 3*ii + 7
		m3_line(15+11*ii, 88, 104-11*ii, 150, types.NewColour(0, uint16(jj), uint16(jj)))
	}

	for {
	}
}
