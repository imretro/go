package imretro

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/spenserblack/go-byteutils"
)

// ModelMap maps bit modes to color models.
type ModelMap = map[PixelMode]color.Model

// ErrUnknownModel is raised when an unknown color model is interpreted.
var ErrUnknownModel = errors.New("Color model not recognized")

// MissingModelError is raised when there is no model for the given pixel bit
// mode.
type MissingModelError PixelMode

// Error reports the pixel mode lacking the color model.
func (mode MissingModelError) Error() string {
	return fmt.Sprintf("No model for pixel mode %02b", mode)
}

var (
	Default1BitColorModel = NewOneBitColorModel(Black, White)
	Default2BitColorModel = NewTwoBitColorModel(Black, DarkGray, LightGray, White)
	Default8BitColorModel = ColorModel(Default8BitPalette)
)

// DefaultModelMap maps bit modes to the default color models.
var DefaultModelMap = ModelMap{
	OneBit:   Default1BitColorModel,
	TwoBit:   Default2BitColorModel,
	EightBit: Default8BitColorModel,
}

// ColorModel is color model for imretro images.
type ColorModel Palette

// PixelMode gets the bits-per-pixel according to the color model.
func (model ColorModel) PixelMode() PixelMode {
	l := len(model)
	switch {
	case l <= 2:
		return OneBit
	case l <= 4:
		return TwoBit
	}
	return EightBit
}

// NewOneBitColorModel creates a new color model for 1-bit-pixel images.
func NewOneBitColorModel(off color.Color, on color.Color) ColorModel {
	return ColorModel{off, on}
}

// NewTwoBitColorModel creates a new color model for 2-bit-pixel images.
func NewTwoBitColorModel(off, light, strong, full color.Color) ColorModel {
	return ColorModel{off, light, strong, full}
}

// Index returns the index of the palette color.
func (model ColorModel) Index(c color.Color) uint8 {
	r, g, b, a := ColorAsBytes(c)
	brightness := r | g | b
	isBright := (brightness >= 128) && (a >= 128)
	switch model.PixelMode() {
	case OneBit:
		if isBright {
			return 1
		}
		return 0
	case TwoBit:
		// NOTE Return "off" if <50% opacity
		if a < 0x80 {
			return 0
		}
		// NOTE Two most significant bits of the combined colors.
		return uint8(r|g|b) >> 6
	}
	r = byteutils.SliceL(r, 0, 2)
	g = byteutils.SliceL(g, 0, 2) << 2
	b = byteutils.SliceL(b, 0, 2) << 4
	a = byteutils.SliceL(a, 0, 2) << 6
	return uint8(r | g | b | a)
}

func (model ColorModel) Convert(c color.Color) color.Color {
	index := model.Index(c)
	if int(index) >= len(model) {
		return NoColor
	}
	return model[index]
}
