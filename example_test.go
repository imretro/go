package imretro_test

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"log"

	_ "github.com/spenserblack/imretro" // register imretro format
)

// ImgBytes declares a 2x2 image with no in-file palette, 1 bit per pixel, and
// an alternating white/black checkerboard pattern.
var ImgBytes = []byte{'I', 'M', 'R', 'E', 'T', 'R', 'O', 0x00, 0, 2, 0, 2, 0b1001_0000}

func Example_decode() {
	var reader io.Reader = bytes.NewBuffer(ImgBytes)
	img, format, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Format: %s\n", format)

	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			fmt.Printf("r = %04X, g = %04X, b = %04X\n", r, g, b)
		}
	}

	// Output:
	// Format: imretro
	// r = FFFF, g = FFFF, b = FFFF
	// r = 0000, g = 0000, b = 0000
	// r = 0000, g = 0000, b = 0000
	// r = FFFF, g = FFFF, b = FFFF
}
