package imretro_test

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"log"

	imretro "github.com/imretro/go"
)

// Custom models can be used when decoding to define your own palette.
// In this example, instead of using the default black & white 1-bit-pixel
// palette, black and green palettes are passed.
func ExampleCustomModel_model_map() {
	black := color.Gray{0}
	green := color.RGBA{0, 0xFF, 0, 0xFF}
	custom := imretro.ModelMap{
		imretro.OneBit: imretro.ColorModel{black, green},
		imretro.TwoBit: imretro.ColorModel{
			black, color.RGBA{0, 0x55, 0, 0}, color.RGBA{0, 0xAA, 0, 0}, black,
		},
		imretro.EightBit: make(imretro.ColorModel, 256),
	}
	var reader io.Reader = bytes.NewBuffer(ImgBytes)
	img, err := imretro.Decode(reader, custom)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			fmt.Printf("r = %04X, g = %04X, b = %04X\n", r, g, b)
		}
	}

	// Output:
	// r = 0000, g = FFFF, b = 0000
	// r = 0000, g = 0000, b = 0000
	// r = 0000, g = 0000, b = 0000
	// r = 0000, g = FFFF, b = 0000
}

// If the pixel mode is predictable, a single ColorModel can be passed. In this
// example, it is assumed that the image will always be in 1-bit-pixel mode
// (two colors).
func ExampleCustomModel_single_model() {
	black := color.Gray{0}
	green := color.RGBA{0, 0xFF, 0, 0xFF}
	custom := imretro.ColorModel{black, green}
	var reader io.Reader = bytes.NewBuffer(ImgBytes)
	img, err := imretro.Decode(reader, custom)
	if err != nil {
		log.Fatal(err)
	}

	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			fmt.Printf("r = %04X, g = %04X, b = %04X\n", r, g, b)
		}
	}

	// Output:
	// r = 0000, g = FFFF, b = 0000
	// r = 0000, g = 0000, b = 0000
	// r = 0000, g = 0000, b = 0000
	// r = 0000, g = FFFF, b = 0000
}
