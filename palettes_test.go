package imretro

import (
	"image/color"
	"testing"
)

// TestModelPixelMode checks that the correct bit mode is interpreted from the
// color model.
func TestModelPixelMode(t *testing.T) {
	if mode := Default1BitColorModel.PixelMode(); mode != OneBit {
		t.Errorf(
			`mode = %v (%08b), want %v (%08b)`,
			mode, mode,
			OneBit, OneBit,
		)
	}

	if mode := Default2BitColorModel.PixelMode(); mode != TwoBit {
		t.Errorf(
			`mode = %v (%08b), want %v (%08b)`,
			mode, mode,
			TwoBit, TwoBit,
		)
	}

	if mode := Default8BitColorModel.PixelMode(); mode != EightBit {
		t.Errorf(
			`mode = %v (%08b), want %v (%08b)`,
			mode, mode,
			EightBit, EightBit,
		)
	}
}

// Test2BitModelIndex checks that the correct bits (ranging [0b00, 0b11]) are
// returned by colors of varying brightness and opacity.
func Test2BitModelIndex(t *testing.T) {
	model := make(ColorModel, 4)

	if bits := model.Index(color.Alpha{0}); bits != 0 {
		t.Errorf(`bits = %02b, want 0`, bits)
	}
	if bits := model.Index(color.RGBA{0xFF, 0xFF, 0xFF, 0x7F}); bits != 0 {
		t.Errorf(`bits = %02b, want 0`, bits)
	}

	if bits := model.Index(darkerGray); bits != 1 {
		t.Errorf(`bits = %02b, want 1`, bits)
	}
	if bits := model.Index(mediumGray); bits != 2 {
		t.Errorf(`bits = %02b, want 2`, bits)
	}

	if bits := model.Index(color.Gray{0xE0}); bits != 3 {
		t.Errorf(`bits = %02b, want 3`, bits)
	}
}

// Test8BitModelIndex checks that the correct bits (ranging [0x00, 0xFF]) are
// returned by colors of varying brightness and opacity.
func Test8BitModelIndex(t *testing.T) {
	model := make(ColorModel, 256)

	if bits := model.Index(color.Alpha{0}); bits != 0 {
		t.Errorf(`bits = %02b, want 0`, bits)
	}
	if bits := model.Index(color.RGBA{0xFF, 0xFF, 0xFF, 0x80}); bits != 0xBF {
		t.Errorf(`bits = %02b, want 10111111`, bits)
	}

	if bits := model.Index(color.RGBA{0xFF, 0, 0, 0xFF}); bits != 0xC3 {
		t.Errorf(`bits = %02b, want 11000011`, bits)
	}
	if bits := model.Index(color.RGBA{0, 0xFF, 0xFF, 0xFF}); bits != 0xFC {
		t.Errorf(`bits = %02b, want 11111100`, bits)
	}
	if bits := model.Index(color.Gray{0x80}); bits != 0xEA {
		t.Errorf(`bits = %02b, want 11101010`, bits)
	}
}
