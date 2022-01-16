package util

import (
	"github.com/spenserblack/go-byteutils"
)

// DimensionsAs3Bytes converts 2 16 bit dimensions into a set of 3 bytes.
func DimensionsAs3Bytes(width uint16, height uint16) (dimensions [3]byte) {
	lowerWidth := byte(width & 0xFF)
	dimensions[0] |= byte((width & 0xF00) >> 4)
	dimensions[0] |= byteutils.SliceL(lowerWidth, 0, 4)
	dimensions[1] |= byteutils.SliceL(lowerWidth, 4, 8) << 4
	dimensions[1] |= byte((height & 0xF00) >> 8)
	dimensions[2] = byte(height & 0xFF)
	return
}

// DimensionsFrom3Bytes converts 3 bytes to width and height, where width
// and height are both 12 bit numbers.
func DimensionsFrom3Bytes(upper, middle, low byte) (width, height int) {
	width |= int(upper) << 4
	width |= int(middle) >> 4
	height |= int(middle & 0x0F) << 8
	height |= int(low)
	return
}
