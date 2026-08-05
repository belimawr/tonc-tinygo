package mmio

import (
	"device/gba"
	"runtime/volatile"
	"unsafe"
)

// ========================= Display Control Register =========================
// Address: 0x4000000 - REG_DISPCNT (The display control register)
// var REG_DISPCNT = gba.DISP.DISPCNT
var REG_DISPCNT = (*volatile.Register16)(unsafe.Add(unsafe.Pointer(gba.REG_BASE), uintptr(0x0)))

// Source: https://www.cs.rit.edu/~tjh8300/CowBite/CowBiteSpec.htm#REG_DISPCNT
/*
   F E D C  B A 9 8  7 6 5 4  3 2 1 0
   W V U S  L K J I  F D B A  C M M M
   0-2 (M) = The video mode. See video modes list above.
   3   (C) = Game Boy Color mode. Read only - should stay at 0.
   4   (A) = This bit controls the starting address of the bitmap in bitmapped modes
             and is used for page flipping. See the description of the specific
             video mode for details.
   5   (B) = Force processing during hblank. Setting this causes the display
             controller to process data earlier and longer, beginning from the end of
             the previous scanline up to the end of the current one. This added
             processing time can help prevent flickering when there are too many
             sprites on a scanline.
   6   (D) = Sets whether sprites stored in VRAM use 1 dimension or 2.
             1 - 1d: tiles are are stored sequentially
             0 - 2d: each row of tiles is stored 32 x 64 bytes in from the start of the
             previous row.
   7   (F) = Force the display to go blank when set. This can be used to save power
             when the display isn't needed, or to blank the screen when it is being
             built up (such as in mode 3, which has only one framebuffer). On the SNES,
             transfers rates to VRAM were improved during a forced blank; it is logical
             to assume that this would also hold true on the GBA.

   8   (I) = If set, enable display of BG0.
   9   (J) = If set, enable display of BG1.
   A   (K) = If set, enable display of BG2.
   B   (L) = If set, enable display of BG3.
   C   (S) = If set, enable display of OAM (sprites).
   D   (U) = Enable Window 0
   E   (V) = Enable Window 1
   F   (W) = Enable Sprite Windows
*/

const (
	// Video modes/BG modes
	DCNT_MODE0 = 0x0000 //M
	DCNT_MODE1 = 0x0001 //M
	DCNT_MODE2 = 0x0002 //M
	// Display mode 3: Bitmap, requires background 2
	DCNT_MODE3 = 0x0003 //M
	DCNT_MODE4 = 0x0004 //M
	DCNT_MODE5 = 0x0005 //M

	// Layers
	DCNT_BG0 = 0x0100 // I
	DCNT_BG1 = 0x0200 // J
	// Background 2
	DCNT_BG2 = 0x0400 // K
	DCNT_BG3 = 0x0800 // L
	DCNT_OBJ = 0x1000 // S
)

// VRAM Video RAM, only 16 bits write
var VRAM = gba.MEM_VRAM
