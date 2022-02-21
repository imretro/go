package imretro_test

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	_ "github.com/imretro/go"
)

// TestDecodedImage decodes an image and tests its pixels.
func TestDecodedImage(t *testing.T) {
	contents := []byte("IMRETRO")
	pixels := []byte{
		0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
		0b1000_0000,
	}
	contents = append(
		contents,
		0,
		0x00, 0x90, 0x09, // dimensions
	)
	contents = append(contents, pixels...)
	r := bytes.NewBuffer(contents)

	m, _, err := image.Decode(r)
	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	bounds := m.Bounds()

	boundTests := []struct {
		d    int
		name rune
	}{{bounds.Dx(), 'x'}, {bounds.Dy(), 'y'}}
	for _, tt := range boundTests {
		if tt.d != 9 {
			t.Fatalf(`dimension %c = %d, want 9`, tt.name, tt.d)
		}
	}

	colors := []color.Color{color.Gray{0xFF}, color.Gray{0}}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			want := colors[(x+y)%2]

			wr, wg, wb, wa := want.RGBA()
			ar, ag, ab, aa := m.At(x, y).RGBA()

			tests := []struct {
				want    uint32
				actual  uint32
				channel rune
			}{
				{wr, ar, 'r'}, {wg, ag, 'g'}, {wb, ab, 'b'}, {wa, aa, 'a'},
			}

			for _, tt := range tests {
				if tt.actual != tt.want {
					t.Fatalf(
						`%c color channel of pixel (%d, %d) = %v, want %v`,
						tt.channel, x, y, tt.actual, tt.want,
					)
				}
			}
		}
	}
}
