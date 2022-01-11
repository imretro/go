package imretro

import (
	"image"
	"image/color"
	"io"

	"github.com/spenserblack/go-byteutils"
)

// Encode writes the image m to w in imretro format.
func Encode(w io.Writer, m image.Image, pixelMode PixelMode) error {
	if _, err := w.Write([]byte("IMRETRO")); err != nil {
		return err
	}
	if _, err := w.Write([]byte{pixelMode | WithPalette}); err != nil {
		return err
	}

	bounds := m.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	for _, d := range []int{width, height} {
		if d > MaximumDimension {
			return DimensionsTooLargeError(d)
		}
		if _, err := w.Write(byteutils.BytesFromUint16(uint16(d), byteutils.LittleEndian)); err != nil {
			return err
		}
	}

	if err := writePalette(w, DefaultModelMap[pixelMode].(ColorModel)); err != nil {
		return err
	}

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
			bit := byte(Default1BitColorModel.Index(c))
			byteutils.ChangeL(&buffer[len(buffer)-1], bitIndex, bit)
			bitIndex++
		}
	}
	w.Write(buffer)
	return nil
}

func encodeTwoBit(w io.Writer, m image.Image) error {
	// NOTE Write the pixels
	bounds := m.Bounds()
	buffer := make([]byte, 1, (bounds.Dx()*bounds.Dy())/4)
	var bitIndex byte = 0

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if bitIndex >= 8 {
				bitIndex = 0
				buffer = append(buffer, 0)
			}
			c := m.At(x, y)
			bits := byte(Default2BitColorModel.Index(c))
			buffer[len(buffer)-1] |= bits << (6 - bitIndex)
			// NOTE Each index is 2 bits
			bitIndex += 2
		}
	}
	w.Write(buffer)
	return nil
}

func encodeEightBit(w io.Writer, m image.Image) error {
	bounds := m.Bounds()
	buffer := make([]byte, 0, bounds.Dx()*bounds.Dy())
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := m.At(x, y)
			bits := byte(Default8BitColorModel.Index(c))
			buffer = append(buffer, bits)
		}
	}
	w.Write(buffer)
	return nil
}

// WriteColor writes a color as 4 bytes to a Writer.
func writeColor(w io.Writer, c color.Color) error {
	r, g, b, a := ColorAsBytes(c)
	_, err := w.Write([]byte{r, g, b, a})
	return err
}

// WritePalette writes all the colors of the palette, where each color is 4
// bytes, to a Writer.
func writePalette(w io.Writer, p ColorModel) error {
	for _, c := range p {
		if err := writeColor(w, c); err != nil {
			return err
		}
	}
	return nil
}
