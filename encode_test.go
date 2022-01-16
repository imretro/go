package imretro

import (
	"bytes"
	"image"
	"image/color"
	"testing"
)

// TestEncode1BitHeader checks that image info would be encoded to a 1-bit imretro
// file.
func TestEncode1BitHeader(t *testing.T) {
	var b bytes.Buffer
	Encode1Bit(t, &b, 320, 240)

	wantHeader := []byte{'I', 'M', 'R', 'E', 'T', 'R', 'O', 0b001_00000}

	for i, actual := range b.Next(len(wantHeader)) {

		if want := wantHeader[i]; actual != want {
			t.Errorf(
				`Header byte %d = %#b (%#x %c), want %#b (%#x %c)`,
				i,
				actual, actual, actual,
				want, want, want,
			)
		}
	}

	FailDimensionHelper(t, &b, "x", "Most", 1)
	FailDimensionHelper(t, &b, "x", "Least", 64)
	FailDimensionHelper(t, &b, "y", "Most", 0)
	FailDimensionHelper(t, &b, "y", "Least", 240)
}

// TestEncodeDimensions checks that the proper bytes are encoded for 2 12-bit
// dimensions.
func TestEncodeDimensions(t *testing.T) {
	var b bytes.Buffer
	Encode1Bit(t, &b, 640, 480)

	// NOTE Skip signature & bit mode
	b.Next(8)

	dimensions := make([]byte, 3)
	if _, err := b.Read(dimensions); err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	tests := []byte{0x28, 0x01, 0xE0}
	for i, want := range tests {
		if actual := dimensions[i]; actual != want {
			t.Errorf(`byte %d = %08b, want %08b`, i, actual, want)
		}
	}
}

// TestEncode1BitPalette checks that a black & white palette would be encoded
// to a 1-bit imretro file.
func TestEncode1BitPalette(t *testing.T) {
	var b bytes.Buffer
	Encode1Bit(t, &b, 320, 240)

	t.Log("Skipping to palette")
	b.Next(11)

	channels := []string{"r", "g", "b", "a"}
	for i, want := range []byte{0, 0, 0, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF} {
		t.Logf(`Checking %s channel of color %d`, channels[i%4], i/4)
		FailByteHelper(t, &b, want)
	}
}

// TestEncode2BitPalette checks that black, white, and 2 shades of gray would
// be encoded to a 2-bit imretro file.
func TestEncode2BitPalette(t *testing.T) {
	var b bytes.Buffer
	Encode2Bit(t, &b, 320, 240)

	t.Log("Skipping to palette")
	b.Next(11)

	channels := []string{"r", "g", "b", "a"}
	bytes := []byte{
		0, 0, 0, 0xFF,
		0x55, 0x55, 0x55, 0xFF,
		0xAA, 0xAA, 0xAA, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
	}
	for i, want := range bytes {
		t.Logf(`Checking %s channel of color %d`, channels[i%4], i/4)
		FailByteHelper(t, &b, want)
	}
}

// TestEncode1BitPixels checks that the pixels have been given the proper indices
// to the palette.
func TestEncode1BitPixels(t *testing.T) {
	var b bytes.Buffer
	Encode1Bit(t, &b, 10, 5)

	t.Log("Skipping to pixels")
	b.Next(11)
	b.Next(8)

	FailByteHelper(t, &b, 0b0000_1111)
	FailByteHelper(t, &b, 0b1010_0000)

	remaining := b.Bytes()

	// NOTE 50 pixels, 8 pixels per byte results in 6 complete bytes (48 pixels)
	// and 1 byte for the 2 remaining pixels. Subtract 2 for the bytes we just
	// tested above.
	if l, want := len(remaining), 5; l != want {
		t.Fatalf(
			`%d remaining pixel bytes (%d total pixel bytes), want %d`,
			l, l+2,
			want,
		)
	}

	t.Logf(`Remaining bytes: %v`, remaining)

	if final, want := remaining[len(remaining)-1], byte(0b0100_0000); final != want {
		t.Errorf(`final byte = %d (%08b), want %d (%08b)`, final, final, want, want)
	}
}

// TestEncode2BitPixels checks that the pixels have been given the proper indices
// to the palette.
func TestEncode2BitPixels(t *testing.T) {
	var b bytes.Buffer
	Encode2Bit(t, &b, 10, 5)

	t.Log("skipping to pixels")
	b.Next(11)
	b.Next(16)

	for i := 0; i < 16/4; i++ {
		FailByteHelper(t, &b, 0b00_01_10_11)
	}

	remaining := b.Bytes()

	// NOTE 50 pixels, 4 pixels per byte results in 12 complete bytes (48 pixels)
	// and 1 byte for the 2 remaining pixels (4 bits in the byte). Subtract 4 for
	// bytes tested above.
	if l, want := len(remaining), 9; l != want {
		t.Fatalf(
			`%d remaining pixel bytes (%d total pixel bytes), want %d`,
			l, l+4,
			want,
		)
	}

	t.Logf(`Remaining bytes: %v`, remaining)

	if final, want := remaining[len(remaining)-1], byte(0b0111_0000); final != want {
		t.Errorf(`final byte = %d (%08b), want %d (%08b)`, final, final, want, want)
	}
}

// TestEncode8BitPixels checks that the pixels have been given the proper indices
// to the palette.
func TestEncode8BitPixels(t *testing.T) {
	var b bytes.Buffer
	Encode8Bit(t, &b, 10, 5)

	t.Log("skipping to pixels")
	b.Next(11)
	b.Next(1024)

	for i := 0; i < 16/4; i++ {
		FailByteHelper(t, &b, 0b1100_0000)
		FailByteHelper(t, &b, 0b1111_1111)
		FailByteHelper(t, &b, 0b1100_0011)
		FailByteHelper(t, &b, 0)
	}

	remaining := b.Bytes()

	// NOTE 50 pixels, 1 pixel per byte results in 50 bytes and 0 remainder. Subtract 16
	// for bytes tested above.
	if l, want := len(remaining), 34; l != want {
		t.Fatalf(
			`%d remaining pixel bytes (%d total pixel bytes), want %d`,
			l, l+16,
			want,
		)
	}

	t.Logf(`Remaining bytes: %v`, remaining)

	if final1, want := remaining[len(remaining)-2], byte(0b11001111); final1 != want {
		t.Errorf(`almost final byte = %d (%08b), want %d (%08b)`, final1, final1, want, want)
	}
	if final2, want := remaining[len(remaining)-1], byte(0b10111100); final2 != want {
		t.Errorf(`almost final byte = %d (%08b), want %d (%08b)`, final2, final2, want, want)
	}
}

// TestEncodeWriteFailures tests that encode will return errors when writing
// fails.
func TestEncodeWriteFailures(t *testing.T) {
	tests := []*struct {
		byteCount int
		expect    string
	}{
		{0, "failure to write signature"},
		{7, "failure to write mode"},
		{8, "failure to write dimensions"},
		{11, "failure to write palette"},
	}
	m := image.NewRGBA(image.Rect(0, 0, 1, 1))
	for i, test := range tests {
		t.Logf(`Test %d`, i)
		b := &cappedWriter{cap: test.byteCount}

		if err := Encode(b, m, OneBit); err == nil {
			t.Error(`err = nil`)
		}
	}
}

// TestEncodeTinyImages tests that extremely small images, which require less
// than a byte for all pixels, won't break.
func TestEncodeTinyImages(t *testing.T) {
	m := image.NewRGBA(image.Rect(0, 0, 1, 1))
	m.Set(0, 0, White)
	tests := []*struct {
		mode PixelMode
		want byte
	}{
		{OneBit, 0b1000_0000},
		{TwoBit, 0b1100_0000},
	}

	for _, test := range tests {
		var b bytes.Buffer
		t.Logf(`Testing mode %08b`, test.mode)
		Encode(&b, m, test.mode)

		byteSlice := b.Bytes()
		lastByte := byteSlice[len(byteSlice)-1]

		if lastByte != test.want {
			t.Errorf(`Last byte = %08b, want %08b`, lastByte, test.want)
		}
	}
}

// TestEncodeTooLargeDimension tests that Encode should fail when the image has
// unsupported dimensions.
func TestEncodeTooLargeDimension(t *testing.T) {
	var b bytes.Buffer
	m := image.NewRGBA(image.Rect(0, 0, 1<<12, 1))
	want := DimensionsTooLargeError(1 << 12)

	if err := Encode(&b, m, EightBit); err != want {
		t.Fatalf(`err = %v, want %v`, err, want)
	}
}

// TestEncodeUnsupportedMode tests that Encode should fail when a bad PixelMode
// is passed.
func TestEncodeUnsupportedMode(t *testing.T) {
	var b bytes.Buffer
	m := image.NewRGBA(image.Rect(0, 0, 1, 1))
	want := UnsupportedBitModeError(1)

	if err := Encode(&b, m, 1); err != want {
		t.Fatalf(`err = %v, want %v`, err, want)
	}
}

// FailDimensionHelper fails if the dimension is not the wanted value.
func FailDimensionHelper(t *testing.T, b *bytes.Buffer, dimension, byteSignificance string, want byte) {
	t.Helper()
	actual, err := b.ReadByte()
	if err != nil {
		panic(err)
	}

	if actual != want {
		t.Errorf(
			`%s significant byte of %s dimension = %d (%08b), want %d (%08b)`,
			byteSignificance, dimension,
			actual, actual,
			want, want,
		)
	}
}

// FailByteHelper fails if the next byte does not match the wanted value.
func FailByteHelper(t *testing.T, b *bytes.Buffer, want byte) {
	t.Helper()
	actual, err := b.ReadByte()
	if err != nil {
		panic(err)
	}

	if actual != want {
		t.Errorf(`byte = %d (%08b), want %d (%08b)`, actual, actual, want, want)
	}
}

// Encode1Bit creates a 1-bit image and encodes it to a buffer.
func Encode1Bit(t *testing.T, b *bytes.Buffer, width, height int) {
	t.Helper()

	m := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < 8; i++ {
		var c color.Color = Black
		if i >= 4 {
			c = White
		}
		m.Set(i, i/width, c)
	}
	for i := 8; i < 16; i++ {
		var c color.Color = Black
		if i%2 == 0 && i < 12 {
			c = White
		}
		m.Set(i%width, i/width, c)
	}
	m.Set(width-2, height-1, DarkerGray)
	m.Set(width-1, height-1, LighterGray)

	Encode(b, m, OneBit)
}

// Encode2Bit creates a 2-bit image and encodes it to a buffer.
func Encode2Bit(t *testing.T, b *bytes.Buffer, width, height int) {
	t.Helper()
	colors := []color.Color{Black, DarkGray, LightGray, White}

	m := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < 16; i++ {
		c := colors[i%len(colors)]
		m.Set(i%width, i/width, c)
	}
	m.Set(width-2, height-1, DarkerGray)
	m.Set(width-1, height-1, LighterGray)

	Encode(b, m, TwoBit)
}

// Encode8Bit creates a 8-bit image and encodes it to a buffer.
func Encode8Bit(t *testing.T, b *bytes.Buffer, width, height int) {
	t.Helper()
	colors := []color.Color{Black, White, color.RGBA{0xFF, 0, 0, 0xFF}, color.RGBA{0, 0, 0, 0}}

	m := image.NewRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < 16; i++ {
		c := colors[i%len(colors)]
		m.Set(i%width, i/width, c)
	}
	m.Set(width-2, height-1, color.RGBA{0xFF, 0xFF, 0, 0xFF})
	m.Set(width-1, height-1, color.RGBA{0, 0xFF, 0xFF, 0x80})

	Encode(b, m, EightBit)
}
