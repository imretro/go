package imretro

import (
	"image"
	"image/color"
	"io"

	"github.com/spenserblack/go-bitio"

	"github.com/imretro/go/internal/util"
)

// Encode writes the image m to w in imretro format.
func Encode(w io.Writer, m image.Image, pixelMode PixelMode) error {
	var helper encoderHelper
	switch pixelMode {
	case OneBit:
		helper = encodeOneBit
	case TwoBit:
		helper = encodeTwoBit
	case EightBit:
		helper = encodeEightBit
	default:
		return UnsupportedBitModeError(pixelMode)
	}

	if _, err := w.Write([]byte("IMRETRO")); err != nil {
		return err
	}
	if _, err := w.Write([]byte{pixelMode | WithPalette | EightBitColors}); err != nil {
		return err
	}

	bounds := m.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	for _, d := range []int{width, height} {
		if d > MaximumDimension {
			return DimensionsTooLargeError(d)
		}
	}
	dimensions := uint(width<<12 | height)
	{
		writer := bitio.NewWriter(w, 3)
		if _, err := writer.WriteBits(dimensions, 24); err != nil {
			return err
		}
	}

	if err := writePalette(w, DefaultModelMap[pixelMode].(ColorModel)); err != nil {
		return err
	}
	return helper(w, m)
}

// EncoderHelper is a unifying type for the specialized pixel encoding
// functions.
type encoderHelper = func(io.Writer, image.Image) error

func encodeOneBit(w io.Writer, m image.Image) error {
	// NOTE Write the pixels
	bounds := m.Bounds()
	pixels := bitio.NewWriter(w, 1)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := m.At(x, y)
			// NOTE If at least 1 color is bright and not transparent, it is bright
			bit := bitio.Bit(Default1BitColorModel.Index(c))
			if _, err := pixels.WriteBit(bit); err != nil {
				return err
			}
		}
	}
	_, err := pixels.CommitPending()
	return err
}

func encodeTwoBit(w io.Writer, m image.Image) error {
	// NOTE Write the pixels
	bounds := m.Bounds()
	pixels := bitio.NewWriter(w, 1)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := m.At(x, y)
			bits := bitio.Bits(Default2BitColorModel.Index(c))
			if _, err := pixels.WriteBits(bits, 2); err != nil {
				return err
			}
		}
	}
	_, err := pixels.CommitPending()
	return err
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
	r, g, b, a := util.ColorAsBytes(c)
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
