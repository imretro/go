package util

import "testing"

// TestFillBytes tests that a byte's bits would be duplicated to "fill" the
// bytes.
func TestFillBytes(t *testing.T) {
	tests := []struct {
		b    []byte
		n    byte
		want []byte
	}{
		{[]byte{0b00, 0b01, 0b10, 0b11}, 2, []byte{0, 0x55, 0xAA, 0xFF}},
		{[]byte{0b0000, 0b0001, 0b0010}, 4, []byte{0, 0x11, 0x22}},
	}

	for testCount, tt := range tests {
		t.Logf(`Test %d`, testCount)
		FillBytes(tt.b, tt.n)
		for i, actual := range tt.b {
			want := tt.want[i]
			if actual != want {
				t.Errorf(`byte %d = %02X, want %02X`, i, actual, want)
			}
		}
	}
}
