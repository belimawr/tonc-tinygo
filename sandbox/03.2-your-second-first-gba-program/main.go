package main

import (
	"machine"
	"runtime/interrupt"
	"runtime/volatile"

	"github.com/belimawr/tonc-tinygo/mmio"
	"github.com/belimawr/tonc-tinygo/types"
)

const (
	size   = 10
	width  = 240
	height = 160
)

type point struct {
	x      int
	y      int
	colour types.Colour
}

func (p *point) setPos(x, y int) {
	p.x = x
	p.y = y
}

func (p *point) draw() {
	for x := range size {
		for y := range size {
			mmio.SetPixelM3(p.x+x, p.y+y, p.colour)
		}
	}
}

func (p *point) pos() (x, y int) {
	return p.x, p.y
}

func (p *point) clear() {
	for x := range size {
		for y := range size {
			mmio.SetPixelM3(p.x+x, p.y+y, 0x0)
		}
	}
}

var (
	p1 = point{x: 120, y: 80, colour: types.NewColour(31, 0, 0)}
	p2 = point{x: 136, y: 80, colour: types.NewColour(0, 31, 0)}
	p3 = point{x: 120, y: 96, colour: types.NewColour(0, 0, 31)}
)

// Compile: tinygo build -o rom.gba -target gameboy-advance .
// Run: mgba rom.gba

var saveGame *volatile.Register16

func main() {
	mmio.REG_DISPCNT.Set(mmio.DCNT_MODE3 | mmio.DCNT_BG2)

	p1.setPos(120, 80)
	p2.setPos(136, 80)
	p3.setPos(120, 96)

	p1.draw()
	p2.draw()
	p3.draw()

	current = &p1

	// Write to the Cart RAM memory section, this is where saved data is stored.
	// mGBA creates a .sav file with this data.
	// That works as an interesting way to make the game "send data" directly
	// to my computer.
	// 0xBA is easy to spot in a hex dump:
	//
	// % hexdump rom.sav
	// 0000000 ffba ffff ffff ffff ffff ffff ffff ffff
	// 0000010 ffff ffff ffff ffff ffff ffff ffff ffff
	// *
	// 0008000
	saveGame = mmio.Mem16(uintptr(0x0E000000))
	saveGame.Set(0xBA) // Something that is easy to see in a hex dump and memory viewer

	// Enable V and H blank interruptions
	// Interrupt code source:
	// https://dev.to/aurelievache/learning-go-by-examples-part-5-create-a-game-boy-advance-gba-game-in-go-5944
	regDISPSTAT := mmio.Mem16(uintptr(0x4000004))
	regDISPSTAT.SetBits(1<<3 | 1<<4)

	interrupt.New(machine.IRQ_VBLANK, update).Enable()

	for {
	}
}

// ============================================================
// Let's try to add some movement using the keys.
var regKEYPAD = mmio.Mem16(uintptr(0x04000130))

var (
	keyUP        = uint16(959)
	keyDOWN      = uint16(895)
	keyLEFT      = uint16(991)
	keyRIGHT     = uint16(1007)
	keyA         = uint16(1022)
	keyB         = uint16(1021)
	keyLSHOULDER = uint16(511)
	keyRSHOULDER = uint16(767)

	current = &point{}
)

func update(evt interrupt.Interrupt) {
	x, y := current.pos()
	switch keyValue := regKEYPAD.Get(); keyValue {
	case keyUP:
		y--
		if y < 0 {
			y = 0
		}
	case keyDOWN:
		y++
		if y > height-size {
			y = height - 1 - size
		}
	case keyLEFT:
		x--
		if x < 0 {
			x = 0
		}
	case keyRIGHT:
		x++
		if x > width-size {
			x = width - 1 - size
		}
	case keyA:
		current = &p1
		x, y = current.pos()
		saveGame.Set(uint16(1))
	case keyB:
		current = &p2
		x, y = current.pos()
		saveGame.Set(uint16(2))
	case keyLSHOULDER, keyRSHOULDER:
		current = &p3
		x, y = current.pos()
		saveGame.Set(uint16(3))
	}

	current.clear()
	current.setPos(x, y)
	current.draw()
}
