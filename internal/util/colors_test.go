package util

import (
	"image/color"
	"testing"
)

// TestColorAsBytes tests that a color would be converted to 4 bytes.
func TestColorAsBytes(t *testing.T) {
	white := color.Gray{255}
	gray := color.Gray{127}
	invisible := color.RGBA{0, 0, 0, 0}

	if r, g, b, a := ColorAsBytes(white); r != 255 || g != 255 || b != 255 || a != 255 {
		t.Fatalf(`r, g, b, a = %d %d %d %d, want 255, 255, 255, 255`, r, g, b, a)
	}
	if r, g, b, a := ColorAsBytes(gray); r != 127 || g != 127 || b != 127 || a != 255 {
		t.Fatalf(`r, g, b, a = %d %d %d %d, want 127, 127, 127, 255`, r, g, b, a)
	}
	if _, _, _, a := ColorAsBytes(invisible); a != 0 {
		t.Fatalf(`a = %d, want 0`, a)
	}
}

// TestColorFromBytes tests that a color can be created from 4 bytes.
func TestColorFromBytes(t *testing.T) {
	tests := []struct {
		color string
		b     [4]byte
		want  map[rune]uint32
	}{
		{
			"white", [4]byte{0xFF, 0xFF, 0xFF, 0xFF},
			map[rune]uint32{'r': 0xFFFF, 'g': 0xFFFF, 'b': 0xFFFF, 'a': 0xFFFF},
		},
		{
			"black", [4]byte{0, 0, 0, 0},
			map[rune]uint32{'r': 0, 'g': 0, 'b': 0, 'a': 0},
		},
	}

	for _, tt := range tests {
		r, g, b, a := ColorFromBytes(tt.b[:]).RGBA()
		actual := map[rune]uint32{'r': r, 'g': g, 'b': b, 'a': a}
		for k, v := range actual {
			if w := tt.want[k]; v != w {
				t.Errorf(`%s color channel %c = %04X, want %04X`, tt.color, k, v, w)
			}
		}
	}
}

// TestColorFromBytesSingleByte tests that a single byte can be expanded to 4
// colors.
func TestColorFromBytesSingleByte(t *testing.T) {
	tests := []struct {
		b    byte
		want map[rune]uint32
	}{
		{
			0b00011011,
			map[rune]uint32{'r': 0, 'g': 0x5555, 'b': 0xAAAA, 'a': 0xFFFF},
		},
	}

	for _, tt := range tests {
		r, g, b, a := ColorFromBytes([]byte{tt.b}).RGBA()
		actual := map[rune]uint32{'r': r, 'g': g, 'b': b, 'a': a}
		for k, v := range actual {
			if w := tt.want[k]; v != w {
				t.Errorf(`color channel %c = %04X, want %04X`, k, v, w)
			}
		}
	}
}
