package imretro_test

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"log"

	_ "github.com/imretro/go" // register imretro format
)

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

func Example_decode_config() {
	var reader io.Reader = bytes.NewBuffer(ImgBytes)
	config, format, err := image.DecodeConfig(reader)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Format: %s\n", format)
	fmt.Printf("width: %d\nheight: %d\n", config.Width, config.Height)

	// Output:
	// Format: imretro
	// width: 2
	// height: 2
}
