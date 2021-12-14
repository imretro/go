package imretro

import (
	"errors"
	"image"
	"image/color"
	"io"

	"github.com/spenserblack/go-byteutils"
)

// Encode writes the image m to w in imretro format.
func Encode(w io.Writer, m image.Image, pixelMode byte) error {
	w.Write([]byte("IMRETRO"))
	w.Write([]byte{pixelMode | WithPalette})

	bounds := m.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	for _, d := range []int{width, height} {
		if d > MaximumDimension {
			return DimensionsTooLargeError(d)
		}
		w.Write(byteutils.BytesFromUint16(uint16(d), byteutils.LittleEndian))
	}

	writePalette(w, DefaultPaletteMap[pixelMode])

	switch pixelMode {
	case OneBit:
		return encodeOneBit(w, m)
	case TwoBit:
		return encodeTwoBit(w, m)
	case EightBit:
		return encodeEightBit(w, m)
	}
	return UnsupportedBitModeError(pixelMode)
}

func encodeOneBit(w io.Writer, m image.Image) error {
	// NOTE Write the pixels
	bounds := m.Bounds()
	buffer := make([]byte, 1, (bounds.Dx()*bounds.Dy())/8)
	var bitIndex byte = 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if bitIndex >= 8 {
				bitIndex = 0
				// NOTE Next byte is being written
				buffer = append(buffer, 0)
			}
			c := m.At(x, y)
			// NOTE If at least 1 color is bright and not transparent, it is bright
			bit := Default1BitColorModel.Bit(c)
			byteutils.ChangeL(&buffer[len(buffer)-1], bitIndex, bit)
			bitIndex++
		}
	}
	w.Write(buffer)
	return nil
}

func encodeTwoBit(w io.Writer, m image.Image) error {
	return errors.New("Not implemented")
}

func encodeEightBit(w io.Writer, m image.Image) error {
	return errors.New("Not implemented")
}

// WriteColor writes a color as 4 bytes to a Writer.
func writeColor(w io.Writer, c color.Color) error {
	r, g, b, a := ColorAsBytes(c)
	_, err := w.Write([]byte{r, g, b, a})
	return err
}

// WritePalette writes all the colors of the palette, where each color is 4
// bytes, to a Writer.
func writePalette(w io.Writer, p Palette) error {
	for _, c := range p {
		if err := writeColor(w, c); err != nil {
			return err
		}
	}
	return nil
}
