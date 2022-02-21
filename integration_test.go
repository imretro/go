package imretro_test

import (
	"bytes"
	"image"
	"image/color"
	"testing"

	imretro "github.com/imretro/go"
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

// TestDecodedMissingPixels decodes an image with missing pixels with
// image.Decode and ensures that it would not panic.
func TestDecodedMissingPixels(t *testing.T) {
	contents := []byte("IMRETRO")
	pixels := []byte{
		0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
		0b1000_0000,
	}
	contents = append(
		contents,
		0,
		0o000, 0o132, 0o011, // dimensions
	)
	contents = append(contents, pixels...)
	r := bytes.NewBuffer(contents)

	if _, _, err := image.Decode(r); err == nil {
		t.Fatalf(`err = nil`)
	}
}

// TestEncode1BitImage encodes an image and ensures the expected bytes are
// written.
//
// Issue #18
func TestEncode1BitImage(t *testing.T) {
	width := 2
	height := 2
	m := image.NewRGBA(image.Rect(0, 0, width, height))
	colors := []color.Color{color.Gray{0}, color.Gray{0xFF}}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			c := colors[(x+y)%len(colors)]
			m.Set(x, y, c)
		}
	}

	var b bytes.Buffer

	imretro.Encode(&b, m, imretro.OneBit)

	wantSignature := []byte("IMRETRO")
	actualSignature := b.Next(7)

	for i, want := range wantSignature {
		actual := actualSignature[i]
		if actual != want {
			t.Fatalf(`signature byte %d = %c, want %c`, i, actual, want)
		}
	}
	wantModeByte := imretro.OneBit | imretro.WithPalette | imretro.EightBitColors
	if actual, err := b.ReadByte(); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	} else if actual != wantModeByte {
		t.Fatalf(`mode byte = %08b, want %08b`, actual, wantModeByte)
	}

	dimensionTests := []byte{0x00, 0x20, 0x02}

	for i, want := range dimensionTests {
		actual, err := b.ReadByte()
		if err != nil {
			t.Fatalf(`err = %v, want nil`, err)
		}
		if actual != want {
			t.Fatalf(`dimension byte %d = %v, want %v`, i, actual, want)
		}
	}

	for i, want := range []byte{0, 0xFF} {
		for j := 0; j < 4; j++ {
			index := j + (i * 4)
			actual, err := b.ReadByte()
			if err != nil {
				t.Fatalf(`err = %v, want nil`, err)
			}
			if actual != want {
				t.Fatalf(`palette byte %d = %v, want %v`, index, actual, want)
			}
		}
	}

	if actual, err := b.ReadByte(); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	} else if actual != 0b0110_0000 {
		t.Fatalf(`pixel byte = 0b%08b, want 0b01100000`, actual)
	}
}
