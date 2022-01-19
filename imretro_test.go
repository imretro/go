package imretro

import (
	"image"
	"image/color"
	"testing"
)

// TestIsBitCountSupported tests that true is returned when the bit count is
// supported, false when not supported.
func TestIsBitCountSupported(t *testing.T) {
	if v := IsBitCountSupported(0b0000_0001); v {
		t.Errorf(`IsBitCountSupported(0b0000_0001) = %v, want false`, v)
	}
	if v := IsBitCountSupported(0b1000_0000); !v {
		t.Errorf(`IsBitCountSupported(0b1000_0000) = %v, want true`, v)
	}
}

// TestUnsupportedError tests the error message for unsupported number of bits error.
func TestUnsupportedError(t *testing.T) {
	if actual, want := UnsupportedBitModeError(0b10).Error(), "Unsupported bit count byte: 0b10"; actual != want {
		t.Fatalf(`err = %q, want %q`, actual, want)
	}
}

// TestImagePixelMode tests that an image returns the correct pixel mode.
func TestImagePixelMode(t *testing.T) {
	i := imretroImage{}
	tests := []*struct {
		colorCount int
		want       PixelMode
	}{
		{2, OneBit},
		{4, TwoBit},
		{256, EightBit},
	}
	for _, test := range tests {
		i.config = image.Config{ColorModel: make(ColorModel, test.colorCount)}
		if mode := i.PixelMode(); mode != test.want {
			t.Errorf(
				`pixel mode for %d colors = %08b, want %08b`,
				test.colorCount,
				mode,
				test.want,
			)
		}
	}
}

// TestImagePalette tests that the image returns its color model as a palette.
func TestImagePalette(t *testing.T) {
	i := imretroImage{}
	i.config = image.Config{ColorModel: Default8BitColorModel}
	var palette color.Palette = i.Palette()

	for i := range palette {
		actual := palette[i]
		want := Default8BitColorModel[i]
		CompareColors(t, actual, want)
	}
}

// TestDimensionsTooLargeError tests that the correct string representation of
// the error is returned.
func TestDimensionsTooLargeError(t *testing.T) {
	err := DimensionsTooLargeError(1 << 16)
	want := "Dimensions too large for 16-bit number: 65536"
	if s := err.Error(); s != want {
		t.Fatalf(`Error() = %q, want %q`, s, want)
	}
}
