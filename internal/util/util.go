package util

import (
	"github.com/spenserblack/go-byteutils"
)

// FillBytes duplicates the last n bytes to fill the whole byte for each byte
// in b.
func FillBytes(b []byte, n byte) {
	for i := range b {
		b[i] = FillByte(b[i], n)
	}
}

// FillByte duplicates the last n bytes to fill the whole byte.
func FillByte(b byte, n byte) byte {
	var bb byte
	bits := byteutils.SliceR(b, 0, n)
	for offset := byte(0); offset < 8; offset += n {
		bb |= bits << offset
	}
	return bb
}
