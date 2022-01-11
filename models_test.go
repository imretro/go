package imretro

import "testing"

// TestModelConvertNoColor tests that a model without enough colors will return
// no color for a color that would be converted to an undefined color.
func TestModelConvertNoColor(t *testing.T) {
	model := ColorModel{Black, Black, Black}

	c := model.Convert(White)
	CompareColors(t, c, NoColor)
}

// TestMissingModelError tests that the correct string representation is
// returned.
func TestMissingModelError(t *testing.T) {
	err := MissingModelError(0b11)
	want := "No model for pixel mode 11"
	if s := err.Error(); s != want {
		t.Fatalf(`Error() = %q, want %q`, s, want)
	}
}
