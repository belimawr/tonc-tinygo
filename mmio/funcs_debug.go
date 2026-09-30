//go:build tinydebug

package mmio

import (
	"fmt"

	"github.com/belimawr/tonc-tinygo/types"
)

func SetM16Colour(addr int, c types.Colour) {
	fmt.Printf("addr: %04d, colour: %X\n", addr, c)
}
