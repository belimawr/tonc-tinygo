package main

import (
	"runtime/volatile"
	"unsafe"
)

// Compile: tinygo build -o rom.gba -target gameboy-advance .
// Run: mgba rom.gba

func main() {
	// IO registers: 0400:0000 to 0400:03FFh
	// Size: 1kb
	// Word: 16 bits
	lcdControll := (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000000)))

	// Register explanation: https://problemkaputt.de/gbatek.htm#gbalcdvideocontroller
	// 4000000h - DISPCNT - LCD Control (Read/Write)
	// b000010000000011
	// Video mode: 3
	// Screen display BG2: on
	lcdControll.Set(0x403)

	// VRAM	0600:0000h to 0601:7FFFh
	// Size: 96kb
	// Word: 16 bit
	// Screen: 240x160

	// Address: 0x06004B78
	dot1 := (*volatile.Register16)(unsafe.Pointer(uintptr(0x06000000 + 120 + 80*240)))
	dot1.Set(0x001F)

	// Address: 0x06004B88
	dot2 := (*volatile.Register16)(unsafe.Pointer(uintptr(0x06000000 + 136 + 80*240)))
	dot2.Set(0x03E0)

	// Address: 0x06005A78
	dot3 := (*volatile.Register16)(unsafe.Pointer(uintptr(0x06000000 + 120 + 96*240)))
	dot3.Set(0x7C00)

	// Go/TinyGo does not seem to need a forever loop o.O
}

// C code
//
// int main()
// {
//     *(unsigned int*)0x04000000 = 0x0403;

//     ((unsigned short*)0x06000000)[120+80*240] = 0x001F;
//     ((unsigned short*)0x06000000)[136+80*240] = 0x03E0;
//     ((unsigned short*)0x06000000)[120+96*240] = 0x7C00;

//     while(1);

//     return 0;
// }
