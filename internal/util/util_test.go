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
