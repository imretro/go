package imretro

import (
	"bytes"
	"image"
	"image/color"
	"io"
	"testing"
)

// TestPassCheckHeader tests that a reader starting with "IMRETRO" bytes will
// pass.
func TestPassCheckHeader(t *testing.T) {
	buff := make([]byte, 8)
	r := MakeImretroReader(0b1010_0000, nil, 0, 0, nil)
	mode, err := checkHeader(r, buff)
	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}
	if pixelMode := mode & (0b1100_0000); pixelMode != EightBit {
		t.Errorf(
			`pixelMode = %d (%08b), want %d (%08b)`,
			pixelMode, pixelMode,
			EightBit, EightBit,
		)
	}
	if hasPalette := mode & (0b0010_0000); hasPalette != 0x20 {
		t.Error("mode does not signify in-file palette")
	}
}

// TestFailCheckHeader tests that a reader with unexpected magic bytes will
// fail.
func TestFailCheckHeader(t *testing.T) {
	buff := make([]byte, 8)
	partialSignature := "IMRET"
	jpgSignature := "\xFF\xD8\xFF\xE0\x00\x10\x4A\x46\x49\x46\x00\x01"

	partialr := bytes.NewBufferString(partialSignature)
	if _, err := checkHeader(partialr, buff); err != io.ErrUnexpectedEOF {
		t.Errorf(`err = %v, want %v`, err, io.ErrUnexpectedEOF)
	}

	jpgr := bytes.NewBufferString(jpgSignature)
	if _, err := checkHeader(jpgr, buff); err != DecodeError("unexpected signature byte") {
		t.Fatalf(`err = %v, want %v`, err, DecodeError("unexpected signature byte"))
	}
}

// TestDecode1BitNoPalette tests that a 1-bit-mode image with no palette can be decoded.
func TestDecode1BitNoPalette(t *testing.T) {
	const width, height int = 320, 240
	var pixels = make([]byte, width*height)
	r := MakeImretroReader(0x00, [][]byte{}, uint16(320), uint16(240), pixels)

	config, err := DecodeConfig(r, nil)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}
	if config.Width != width {
		t.Errorf(`Width = %v, want %v`, config.Width, width)
	}
	if config.Height != height {
		t.Errorf(`Height = %v, want %v`, config.Height, height)
	}
	if _, ok := config.ColorModel.(OneBitColorModel); !ok {
		t.Fatalf(`ColorModel is %T, want OneBitColorModel`, config.ColorModel)
	}

	inputAndWant := [][2]color.Color{{DarkGray, Black}, {LightGray, White}}

	for _, colors := range inputAndWant {
		input := colors[0]
		want := colors[1]

		t.Logf(`Comparing conversion of %v`, input)
		actual := config.ColorModel.Convert(input)
		CompareColors(t, actual, want)
	}
}

// TestDecode2BitNoPalette tests that a 2-bit-mode image with no palette can be decoded.
func TestDecode2BitNoPalette(t *testing.T) {
	const width, height int = 320, 240
	var pixels = make([]byte, width*height)
	r := MakeImretroReader(0x40, [][]byte{}, uint16(320), uint16(240), pixels)

	config, err := DecodeConfig(r, nil)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}
	if config.Width != width {
		t.Errorf(`Width = %v, want %v`, config.Width, width)
	}
	if config.Height != height {
		t.Errorf(`Height = %v, want %v`, config.Height, height)
	}
	if _, ok := config.ColorModel.(TwoBitColorModel); !ok {
		t.Fatalf(`ColorModel is %T, want TwoBitColorModel`, config.ColorModel)
	}

	inputAndWant := [][2]color.Color{
		{color.Gray{0x0F}, Black},
		{DarkGray, DarkGray},
		{LightGray, LightGray},
		{color.Gray{0xF0}, White},
	}

	for _, colors := range inputAndWant {
		input := colors[0]
		want := colors[1]

		t.Logf(`Comparing conversion of %v`, input)
		actual := config.ColorModel.Convert(input)
		CompareColors(t, actual, want)
	}
}

// TestDecode8BitNoPalette tests that an 8-bit-mode image with no palette can be decoded.
func TestDecode8BitNoPalette(t *testing.T) {
	const width, height int = 320, 240
	var pixels = make([]byte, width*height)
	r := MakeImretroReader(0x80, [][]byte{}, uint16(320), uint16(240), pixels)

	config, err := DecodeConfig(r, nil)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}
	if config.Width != width {
		t.Errorf(`Width = %v, want %v`, config.Width, width)
	}
	if config.Height != height {
		t.Errorf(`Height = %v, want %v`, config.Height, height)
	}
	if _, ok := config.ColorModel.(EightBitColorModel); !ok {
		t.Fatalf(`ColorModel is %T, want EightBitColorModel`, config.ColorModel)
	}

	inputAndWant := [][2]color.Color{
		{color.Gray{0x0F}, Black},
		{color.RGBA{0xFF, 0x01, 0xFF, 0xF0}, color.RGBA{0xFF, 0x00, 0xFF, 0xFF}},
	}

	for _, colors := range inputAndWant {
		input := colors[0]
		want := colors[1]

		t.Logf(`Comparing conversion of %v`, input)
		actual := config.ColorModel.Convert(input)
		CompareColors(t, actual, want)
	}
}

// TestDecode1BitPalette tests that a 1-bit palette would be properly decoded.
func TestDecode1BitPalette(t *testing.T) {
	palette := [][]byte{
		{0x00, 0xFF, 0x00, 0xFF},
		{0xEF, 0xFF, 0x00, 0xFF},
	}
	r := MakeImretroReader(0x20, palette, 2, 2, make([]byte, 1))

	config, err := DecodeConfig(r, nil)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	if _, ok := config.ColorModel.(OneBitColorModel); !ok {
		t.Fatalf(`ColorModel is %T, want OneBitColorModel`, config.ColorModel)
	}

	inputAndWant := [][2]color.Color{
		{Black, color.RGBA{0x00, 0xFF, 0x00, 0xFF}},
		{White, color.RGBA{0xEF, 0xFF, 0x00, 0xFF}},
	}

	for _, colors := range inputAndWant {
		input := colors[0]
		want := colors[1]

		t.Logf(`Comparing conversion of %v`, input)
		actual := config.ColorModel.Convert(input)
		CompareColors(t, actual, want)
	}
}

// TestDecode2BitPalette tests that a 2-bit palette would be properly decoded.
func TestDecode2BitPalette(t *testing.T) {
	palette := [][]byte{
		{0xFF, 0x00, 0x00, 0xFF},
		{0x00, 0xFF, 0x00, 0xFF},
		{0x00, 0x00, 0xFF, 0xFF},
		{0x00, 0x00, 0x00, 0x00},
	}
	r := MakeImretroReader(0x60, palette, 2, 2, make([]byte, 4))

	config, err := DecodeConfig(r, nil)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	if _, ok := config.ColorModel.(TwoBitColorModel); !ok {
		t.Fatalf(`ColorModel is %T, want TwoBitColorModel`, config.ColorModel)
	}

	inputAndWant := [][2]color.Color{
		{Black, color.RGBA{0xFF, 0x00, 0x00, 0xFF}},
		{White, color.RGBA{0x00, 0x00, 0x00, 0x00}},
		{DarkGray, color.RGBA{0x00, 0xFF, 0x00, 0xFF}},
		{LightGray, color.RGBA{0x00, 0x00, 0xFF, 0xFF}},
	}

	for _, colors := range inputAndWant {
		input := colors[0]
		want := colors[1]

		t.Logf(`Comparing conversion of %v`, input)
		actual := config.ColorModel.Convert(input)
		CompareColors(t, actual, want)
	}
}

// TestDecode8BitPalette tests that a 2-bit palette would be properly decoded.
func TestDecode8BitPalette(t *testing.T) {
	reversedPalette := make([][]byte, 0, 256)

	last := len(Default8BitPalette) - 1
	for i := range Default8BitPalette {
		c := Default8BitPalette[last-i]
		r, g, b, a := ColorAsBytes(c)
		reversedPalette = append(reversedPalette, []byte{r, g, b, a})
	}

	r := MakeImretroReader(0xA0, reversedPalette, 2, 2, make([]byte, 4))

	config, err := DecodeConfig(r, nil)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	if _, ok := config.ColorModel.(EightBitColorModel); !ok {
		t.Fatalf(`ColorModel is %T, want EightBitColorModel`, config.ColorModel)
	}

	inputAndWant := [][2]color.Color{
		{color.Alpha{0}, White},
		{White, color.Alpha{0}},
	}

	for _, colors := range inputAndWant {
		input := colors[0]
		want := colors[1]

		t.Logf(`Comparing conversion of %v`, input)
		actual := config.ColorModel.Convert(input)
		CompareColors(t, actual, want)
	}
}

// TestDecode1BitImage tests that a 1-bit image would be properly decoded.
func TestDecode1BitImage(t *testing.T) {
	r := MakeImretroReader(0x00, [][]byte{}, 5, 2, []byte{0b10010_100, 0b01_000000})
	i, err := Decode(r, nil)
	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	blackPoints := []image.Point{
		{1, 0}, {2, 0}, {4, 0},
		{1, 1}, {2, 1}, {3, 1},
	}
	whitePoints := []image.Point{
		{0, 0}, {3, 0},
		{0, 1}, {4, 1},
	}
	for _, p := range blackPoints {
		t.Logf(`Testing point %v`, p)
		CompareColors(t, i.At(p.X, p.Y), Black)
	}
	for _, p := range whitePoints {
		t.Logf(`Testing point %v`, p)
		CompareColors(t, i.At(p.X, p.Y), White)
	}
	CompareColors(t, i.At(-1, -1), NoColor)
	CompareColors(t, i.At(5, 1), NoColor)
	CompareColors(t, i.At(5, 2), NoColor)
	CompareColors(t, i.At(10, 10), NoColor)
}

// TestDecode2BitImage tests that a 2-bit image would be properly decoded.
func TestDecode2BitImage(t *testing.T) {
	pixels := []byte{0b00011011, 0b11_100100, 0b1101_0000}
	r := MakeImretroReader(0x40, nil, 5, 2, pixels)
	i, err := Decode(r, nil)
	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	offPoints := []image.Point{{0, 0}, {2, 1}}
	lightPoints := []image.Point{{1, 0}, {1, 1}, {4, 1}}
	strongPoints := []image.Point{{2, 0}, {0, 1}}
	fullPoints := []image.Point{{3, 0}, {4, 0}, {3, 1}}
	for _, p := range offPoints {
		t.Logf(`Testing point %v`, p)
		CompareColors(t, i.At(p.X, p.Y), Black)
	}
	for _, p := range lightPoints {
		t.Logf(`Testing point %v`, p)
		CompareColors(t, i.At(p.X, p.Y), DarkGray)
	}
	for _, p := range strongPoints {
		t.Logf(`Testing point %v`, p)
		CompareColors(t, i.At(p.X, p.Y), LightGray)
	}
	for _, p := range fullPoints {
		t.Logf(`Testing point %v`, p)
		CompareColors(t, i.At(p.X, p.Y), White)
	}
	CompareColors(t, i.At(-1, -1), NoColor)
	CompareColors(t, i.At(5, 1), NoColor)
	CompareColors(t, i.At(5, 2), NoColor)
	CompareColors(t, i.At(10, 10), NoColor)
}

// TestDecode8BitImage tests that an 8-bit image would be properly decoded.
func TestDecode8BitImage(t *testing.T) {
	pixels := []byte{
		0x00, 0xFF, 0xC0, 0xC3, 0xCC, // transparent, white, black, red, green
		0xF0, 0xCF, 0xF3, 0xFC, 0xAA, // blue, yellow, magenta, cyan, 75% light gray
	}
	r := MakeImretroReader(0x80, nil, 5, 2, pixels)
	i, err := Decode(r, nil)
	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	wantColors := []color.Color{
		color.Alpha{0}, White, Black, color.RGBA{0xFF, 0, 0, 0xFF}, color.RGBA{0, 0xFF, 0, 0xFF},
		color.RGBA{0, 0, 0xFF, 0xFF}, color.RGBA{0xFF, 0xFF, 0, 0xFF}, color.RGBA{0xFF, 0, 0xFF, 0xFF}, color.RGBA{0, 0xFF, 0xFF, 0xFF}, color.RGBA{0xAA, 0xAA, 0xAA, 0xAA},
	}

	for index, want := range wantColors {
		x := index % 5
		y := index / 5
		t.Logf(`Testing point (%d, %d)`, x, y)
		CompareColors(t, i.At(x, y), want)
	}
	CompareColors(t, i.At(-1, -1), NoColor)
	CompareColors(t, i.At(5, 1), NoColor)
	CompareColors(t, i.At(5, 2), NoColor)
	CompareColors(t, i.At(10, 10), NoColor)
}

// MakeImretroReader makes a 1-bit imretro reader.
func MakeImretroReader(mode byte, palette [][]byte, width, height uint16, pixels []byte) *bytes.Buffer {
	b := bytes.NewBuffer([]byte{
		// signature/magic bytes
		'I', 'M', 'R', 'E', 'T', 'R', 'O',
		// Mode byte (8-bit, in-file palette)
		mode,
		byte(width >> 8), byte(width & 0xFF),
		byte(height >> 8), byte(height & 0xFF),
	})
	for _, color := range palette {
		b.Write(color)
	}
	b.Write(pixels)
	return b
}

// CompareColors helps compare colors to each other.
func CompareColors(t *testing.T, actual, want color.Color) {
	t.Helper()
	r, g, b, a := actual.RGBA()
	wr, wg, wb, wa := want.RGBA()
	comparisons := [4]channelComparison{
		{"red", r, wr},
		{"green", g, wg},
		{"blue", b, wb},
		{"alpha", a, wa},
	}

	for _, comparison := range comparisons {
		if comparison.actual != comparison.want {
			t.Errorf(
				`%s channel = %v, want %v`,
				comparison.name, comparison.actual,
				comparison.want,
			)
		}
	}
}

// ChannelComparison is used to compare color channels.
type channelComparison struct {
	name         string
	actual, want uint32
}
