package mmio

// Main memory sections
const (
	// MEM_VRAM is the video RAM. Note: no 8bit write!
	MEM_VRAM uintptr = 0x06000000

	// MEM_IO Starting address of the IO registers
	MEM_IO uintptr = 0x04000000
)
