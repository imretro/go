package util

import "image/color"

// ColorAsBytes converts a color to a 4-byte (one byte for each channel)
// representation.
func ColorAsBytes(c color.Color) (r, g, b, a byte) {
	rchan, gchan, bchan, achan := c.RGBA()
	return ChannelAsByte(rchan), ChannelAsByte(gchan), ChannelAsByte(bchan), ChannelAsByte(achan)
}

// ColorFromBytes converts 4 bytes into a color. Panics if the slice has less
// than 4 bytes.
func ColorFromBytes(bs []byte) color.Color {
	return color.RGBA{bs[0], bs[1], bs[2], bs[3]}
}

// ChannelAsByte converts a uint32 color channel ranging within [0, 0xFFFF] to
// a byte.
func ChannelAsByte(channel uint32) byte {
	return byte(channel >> 8)
}
