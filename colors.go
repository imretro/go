package imretro

import (
	"image/color"
)

var (
	// NoColor is "invisible" and signifies a lack of color.
	noColor color.Color = color.Alpha{0}
	black   color.Color = color.Gray{0}
	// DarkerGray is 25% light.
	darkerGray = color.Gray{0x40}
	// DarkGray is 33% light, and can be used for splitting a monochromatic
	// color range into 4 parts (0, 33%, 66%, 100%).
	darkGray = color.Gray{0x55}
	// MediumGray is the exact middle between black and white.
	mediumGray = color.Gray{0x80}
	// LightGray is 66% light, and can be used for splitting a monochromatic
	// color range into 4 parts (0, 33%, 66%, 100%).
	lightGray = color.Gray{0xAA}
	// LighterGray is 75% light.
	lighterGray = color.Gray{0xC0}
	white       = color.Gray{0xFF}
)

// ColorBytes is a color that has a variable number of bytes. It can be
// grayscale, RGB, or RGBA.
type colorBytes []byte

func (c colorBytes) RGBA() (r, g, b, a uint32) {
	return c.AsColor().RGBA()
}

func (c colorBytes) AsColor() color.Color {
	switch len(c) {
	case 1:
		return color.Gray{c[0]}
	case 3:
		return color.RGBA{c[0], c[1], c[2], 0xFF}
	case 4:
		return color.RGBA{c[0], c[1], c[2], c[3]}
	}
	panic("Unreachable")
}
