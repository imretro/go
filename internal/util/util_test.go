package util

import "testing"

// TestDimensionsAs3Bytes tests that 2 16 bit dimensions would be converted to
// a 24-bit set.
func TestDimensionsAs3Bytes(t *testing.T) {
	dimensions := DimensionsAs3Bytes(0x248, 0xABC)
	for i, want := range []byte{0x24, 0x8A, 0xBC} {
		if actual := dimensions[i]; actual != want {
			t.Errorf(`dimensions[%d] = %02X, want %02X`, i, actual, want)
		}
	}
}

// TestDimensionsFrom3Bytes tests that 3 bytes can be converted to width and
// height, where both dimensions are treated as 12 bit numbers.
func TestDimensionsFrom3Bytes(t *testing.T) {
	width, height := DimensionsFrom3Bytes(0x24, 0x8A, 0xBC)
	tests := []struct{
		dim    string
		actual int
		want   int
	}{
		{"width", width, 0x248},
		{"height", height, 0xABC},
	}
	for _, test := range tests {
		if test.actual != test.want {
			t.Errorf(`%s = %03X, want %03X`, test.dim, test.actual, test.want)
		}
	}
}
