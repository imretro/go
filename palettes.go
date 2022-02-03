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
