package imretro_test

import imretro "github.com/imretro/go"

// ImgBytes declares a 2x2 image with no in-file palette, 1 bit per pixel, and
// an alternating white/black checkerboard pattern.
var ImgBytes = []byte{
	'I', 'M', 'R', 'E', 'T', 'R', 'O', // Signature
	imretro.OneBit,   // Mode
	0x00, 0x20, 0x02, // Width & Height (2 12-bit numbers)
	0b1001_0000, // Pixels (on, off, off, on, ignored)
}
