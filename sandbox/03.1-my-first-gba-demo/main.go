package main

import (
	"machine"
	"runtime/interrupt"
	"runtime/volatile"
	"unsafe"
)

var width = uint(240)

type point struct {
	x         uint8
	y         uint8
	colour    uint16
	colourReg *volatile.Register16
}

func (p *point) setPos(x, y uint8) {
	p.x = x
	p.y = y

	// VRAM	0600:0000h to 0601:7FFFh
	// Size: 96kb
	// Word: 16 bit
	// Screen: 240x160
	p.colourReg = (*volatile.Register16)(unsafe.Pointer(uintptr(0x06000000 + uint(x) + uint(y)*width)))
	p.colourReg.Set(p.colour)
}

func (p *point) pos() (x, y uint8) {
	return p.x, p.y
}

func (p *point) clear() {
	p.colourReg.Set(0)
}

var (
	p1 = point{x: 120, y: 80, colour: 0x001F}
	p2 = point{x: 136, y: 80, colour: 0x03E0}
	p3 = point{x: 120, y: 96, colour: 0x7C00}
)

// Compile: tinygo build -o rom.gba -target gameboy-advance .
// Run: mgba rom.gba

var saveGame *volatile.Register16

func main() {
	// IO registers: 0400:0000 to 0400:03FFh
	// Size: 1kb
	// Word: 16 bits
	lcdControll := (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000000)))

	// Register explanation: https://problemkaputt.de/gbatek.htm#gbalcdvideocontroller
	// https://www.cs.rit.edu/~tjh8300/CowBite/CowBiteSpec.htm#REG_DISPCNT
	// 4000000h - DISPCNT - LCD Control (Read/Write)
	// F E D C  B A 9 8  7 6 5 4  3 2 1 0
	// 0 0 0 0  0 1 0 0  0 0 0 0  0 0 1 1
	// Video mode: 3
	// Screen display BG2: on
	lcdControll.Set(0x403)

	// VRAM	0600:0000h to 0601:7FFFh
	// Size: 96kb
	// Word: 16 bit
	// Screen: 240x160

	// Address: 0x06004B78
	p1.setPos(120, 80)
	current = &p1

	// Address: 0x06004B88
	p2.setPos(136, 80)

	// Address: 0x06005A78
	p3.setPos(120, 96)

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
	saveGame = (*volatile.Register16)(unsafe.Pointer(uintptr(0x0E000000)))
	saveGame.Set(0xBA) // Something that is easy to see in a hex dump and memory viewer

	// Enable V and H blank interruptions
	// Interrupt code source:
	// https://dev.to/aurelievache/learning-go-by-examples-part-5-create-a-game-boy-advance-gba-game-in-go-5944
	regDISPSTAT := (*volatile.Register16)(unsafe.Pointer(uintptr(0x4000004)))
	regDISPSTAT.SetBits(1<<3 | 1<<4)

	interrupt.New(machine.IRQ_VBLANK, update).Enable()

	for {
	}
}

// ============================================================
// Let's try to add some movement using the keys.
var regKEYPAD = (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000130)))

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
	case keyDOWN:
		y++
	case keyLEFT:
		x--
	case keyRIGHT:
		x++
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
}
