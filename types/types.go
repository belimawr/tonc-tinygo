package types

type Colour uint16

// NewColour returns a new colour in the 5.5.5 BRG Blue-Green-Red format.
//
// Each colour in the 0 to 31 range
// Red  : bits 0 to 4
// Blue : bits 5 to 9
// Green: bits 10 to 14
// NoUse: bit 15, not used
// 0 G G G G G B B B B B R R R R R
func NewColour(r, g, b uint16) Colour {
	return Colour(r | g<<5 | b<<10)
}
