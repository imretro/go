package imretro

import (
	"bytes"
	"image"
	"image/color"
	"io"
	"testing"

	"github.com/imretro/go/internal/util"
	"github.com/spenserblack/go-bitio"
)

// TestPassCheckHeader tests that a reader starting with "IMRETRO" bytes will
// pass.
func TestPassCheckHeader(t *testing.T) {
	buff := make([]byte, 8)
	r := MakeImretroReader(EightBit|WithPalette|RGBA|EightBitColors, nil, 0, 0, nil)
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
	r := MakeImretroReader(OneBit|EightBitColors, [][]byte{}, uint16(320), uint16(240), pixels)

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

	inputAndWant := [][2]color.Color{{darkGray, black}, {lightGray, white}}

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
	r := MakeImretroReader(TwoBit|EightBitColors, [][]byte{}, uint16(320), uint16(240), pixels)

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

	inputAndWant := [][2]color.Color{
		{color.Gray{0x0F}, black},
		{darkGray, darkGray},
		{lightGray, lightGray},
		{color.Gray{0xF0}, white},
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
	r := MakeImretroReader(EightBit|EightBitColors, [][]byte{}, uint16(320), uint16(240), pixels)

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

	inputAndWant := [][2]color.Color{
		{color.Gray{0x0F}, black},
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
	r := MakeImretroReader(OneBit|WithPalette|RGBA|EightBitColors, palette, 2, 2, make([]byte, 1))

	config, err := DecodeConfig(r, nil)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	inputAndWant := [][2]color.Color{
		{black, color.RGBA{0x00, 0xFF, 0x00, 0xFF}},
		{white, color.RGBA{0xEF, 0xFF, 0x00, 0xFF}},
	}

	for _, colors := range inputAndWant {
		input := colors[0]
		want := colors[1]

		t.Logf(`Comparing conversion of %v`, input)
		actual := config.ColorModel.Convert(input)
		CompareColors(t, actual, want)
	}
}

// TestDecode1BitMinGrayscalePalette tests that a 1-bit grayscale palette using
// 2-bit color channels would be properly decoded.
func TestDecode1BitMinGrayscalePalette(t *testing.T) {
	palette := [][]byte{{0b0110 << 4}}
	r := MakeImretroReader(OneBit|WithPalette|Grayscale, palette, 2, 2, make([]byte, 1))

	config, err := DecodeConfig(r, nil)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	tests := []struct {
		input color.Color
		want  color.Color
	}{
		{black, color.Gray{0x55}},
		{white, color.Gray{0xAA}},
	}

	for _, tt := range tests {
		t.Logf(`Converting %#v`, tt.input)
		actual := config.ColorModel.Convert(tt.input)
		CompareColors(t, actual, tt.want)
	}
}

// TestDecode1BitMinRGBPalette tests that a 1-bit RGB palette using
// 2-bit color channels would be properly decoded.
func TestDecode1BitMinRGBPalette(t *testing.T) {
	palette := [][]byte{{0b110100_01}, {0b0011 << 4}}
	r := MakeImretroReader(OneBit|WithPalette|RGB, palette, 2, 2, make([]byte, 1))

	config, err := DecodeConfig(r, nil)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	tests := []struct {
		input color.Color
		want  color.Color
	}{
		{black, color.RGBA{0xFF, 0x55, 0, 0xFF}},
		{white, color.RGBA{0x55, 0, 0xFF, 0xFF}},
	}

	for _, tt := range tests {
		t.Logf(`Converting %#v`, tt.input)
		actual := config.ColorModel.Convert(tt.input)
		CompareColors(t, actual, tt.want)
	}
}

// TestDecode1BitMinRGBAPalette tests that a 1-bit RGB palette using
// 2-bit color channels would be properly decoded.
func TestDecode1BitMinRGBAPalette(t *testing.T) {
	palette := [][]byte{{0}, {0b00011011}}
	r := MakeImretroReader(OneBit|WithPalette|RGBA, palette, 2, 2, make([]byte, 1))

	config, err := DecodeConfig(r, nil)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	tests := []struct {
		input color.Color
		want  color.Color
	}{
		{black, color.Alpha{0}},
		{white, color.RGBA{0, 0x55, 0xAA, 0xFF}},
	}

	for _, tt := range tests {
		t.Logf(`Converting %#v`, tt.input)
		actual := config.ColorModel.Convert(tt.input)
		CompareColors(t, actual, tt.want)
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
	r := MakeImretroReader(TwoBit|WithPalette|RGBA|EightBitColors, palette, 2, 2, make([]byte, 4))

	config, err := DecodeConfig(r, nil)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	inputAndWant := [][2]color.Color{
		{black, color.RGBA{0xFF, 0x00, 0x00, 0xFF}},
		{white, color.RGBA{0x00, 0x00, 0x00, 0x00}},
		{darkGray, color.RGBA{0x00, 0xFF, 0x00, 0xFF}},
		{lightGray, color.RGBA{0x00, 0x00, 0xFF, 0xFF}},
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

	last := len(Default8BitColorModel) - 1
	for i := range Default8BitColorModel {
		c := Default8BitColorModel[last-i]
		r, g, b, a := util.ColorAsBytes(c)
		reversedPalette = append(reversedPalette, []byte{r, g, b, a})
	}

	r := MakeImretroReader(EightBit|WithPalette|RGBA|EightBitColors, reversedPalette, 2, 2, make([]byte, 4))

	config, err := DecodeConfig(r, nil)

	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	inputAndWant := [][2]color.Color{
		{color.Alpha{0}, white},
		{white, color.Alpha{0}},
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
	r := MakeImretroReader(OneBit|EightBitColors, [][]byte{}, 5, 2, []byte{0b10010_100, 0b01_000000})
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
		CompareColors(t, i.At(p.X, p.Y), black)
	}
	for _, p := range whitePoints {
		t.Logf(`Testing point %v`, p)
		CompareColors(t, i.At(p.X, p.Y), white)
	}
	CompareColors(t, i.At(-1, -1), noColor)
	CompareColors(t, i.At(5, 1), noColor)
	CompareColors(t, i.At(5, 2), noColor)
	CompareColors(t, i.At(10, 10), noColor)
}

// TestDecode2BitImage tests that a 2-bit image would be properly decoded.
func TestDecode2BitImage(t *testing.T) {
	pixels := []byte{0b00011011, 0b11_100100, 0b1101_0000}
	r := MakeImretroReader(TwoBit|EightBitColors, nil, 5, 2, pixels)
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
		CompareColors(t, i.At(p.X, p.Y), black)
	}
	for _, p := range lightPoints {
		t.Logf(`Testing point %v`, p)
		CompareColors(t, i.At(p.X, p.Y), darkGray)
	}
	for _, p := range strongPoints {
		t.Logf(`Testing point %v`, p)
		CompareColors(t, i.At(p.X, p.Y), lightGray)
	}
	for _, p := range fullPoints {
		t.Logf(`Testing point %v`, p)
		CompareColors(t, i.At(p.X, p.Y), white)
	}
	CompareColors(t, i.At(-1, -1), noColor)
	CompareColors(t, i.At(5, 1), noColor)
	CompareColors(t, i.At(5, 2), noColor)
	CompareColors(t, i.At(10, 10), noColor)
}

// TestDecode8BitImage tests that an 8-bit image would be properly decoded.
func TestDecode8BitImage(t *testing.T) {
	pixels := []byte{
		0x00, 0xFF, 0xC0, 0xC3, 0xCC, // transparent, white, black, red, green
		0xF0, 0xCF, 0xF3, 0xFC, 0xAA, // blue, yellow, magenta, cyan, 75% light gray
	}
	r := MakeImretroReader(EightBit|EightBitColors, nil, 5, 2, pixels)
	i, err := Decode(r, nil)
	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	wantColors := []color.Color{
		color.Alpha{0}, white, black, color.RGBA{0xFF, 0, 0, 0xFF}, color.RGBA{0, 0xFF, 0, 0xFF},
		color.RGBA{0, 0, 0xFF, 0xFF}, color.RGBA{0xFF, 0xFF, 0, 0xFF}, color.RGBA{0xFF, 0, 0xFF, 0xFF}, color.RGBA{0, 0xFF, 0xFF, 0xFF}, color.RGBA{0xAA, 0xAA, 0xAA, 0xAA},
	}

	for index, want := range wantColors {
		x := index % 5
		y := index / 5
		t.Logf(`Testing point (%d, %d)`, x, y)
		CompareColors(t, i.At(x, y), want)
	}
	CompareColors(t, i.At(-1, -1), noColor)
	CompareColors(t, i.At(5, 1), noColor)
	CompareColors(t, i.At(5, 2), noColor)
	CompareColors(t, i.At(10, 10), noColor)
}

// TestDecodeWithCustomModel tests that an image can be decoded and the custom
// model(s) will be used for the image.
func TestDecodeWithCustomModel(t *testing.T) {
	pixels := []byte{0b0100_0000}
	r := MakeImretroReader(OneBit|EightBitColors, nil, 2, 1, pixels)
	off := color.Alpha{0}
	on := color.RGBA{0, 0xFF, 0, 0xFF}
	i, err := Decode(r, ModelMap{OneBit: NewOneBitColorModel(off, on)})
	if err != nil {
		t.Fatalf(`err = %v, want nil`, err)
	}

	CompareColors(t, i.At(0, 0), off)
	CompareColors(t, i.At(1, 0), on)
}

// TestDecodeMissingModel tests that an image cannot be decoded when the model
// is missing or the pixel mode is not supported (and therefore doesn't have a
// model).
func TestDecodeMissingModel(t *testing.T) {
	var r io.Reader
	var err error

	r = MakeImretroReader(OneBit|EightBitColors, nil, 1, 1, []byte{0})
	_, err = Decode(r, ModelMap{})
	if want := MissingModelError(0); err != want {
		t.Errorf(`err = %v, want %v`, err, want)
	}

	r = MakeImretroReader(0b1110_0001, nil, 1, 1, []byte{0})
	_, err = DecodeConfig(r, nil)
	if want := MissingModelError(0b1100_0000); err != want {
		t.Errorf(`err = %v, want %v`, err, want)
	}
}

// TestDecodeReaderError tests that a reader error would be returned if it
// occurs.
func TestDecodeReaderError(t *testing.T) {
	var r io.Reader
	var err error

	r = bytes.NewBuffer([]byte{})
	if _, err = DecodeConfig(r, nil); err == nil {
		t.Errorf(`err = nil`)
	}

	r = io.LimitReader(MakeImretroReader(EightBit, nil, 1, 1, []byte{0}), 10)
	if _, err = DecodeConfig(r, nil); err == nil {
		t.Errorf(`err = nil`)
	}

	r = io.LimitReader(bytes.NewBuffer([]byte{0xFF, 0xAA, 0x55}), 1)
	if _, _, err = decodeDimensions(r); err != io.EOF {
		t.Errorf(`err = %v, want nil`, err)
	}

	r = errorReader{}
	if _, err = decodeModel(r, 2, true, RGBA); err == nil {
		t.Errorf(`err = nil`)
	}

	r = &errorLimitReader{
		&io.LimitedReader{
			R: MakeImretroReader(EightBit, nil, 10, 10, make([]byte, 100)),
			N: 50,
		},
	}
	if _, err = Decode(r, nil); err == nil {
		t.Errorf(`err = nil`)
	}
}

// TestDecodeError tests that the proper string representation of a failure to
// decode is returned.
func TestDecodeError(t *testing.T) {
	err := DecodeError("Failed!")
	if s := err.Error(); s != "Failed!" {
		t.Fatalf(`Error() = %q, want "Failed!"`, s)
	}
}

// MakeImretroReader makes a 1-bit imretro reader.
func MakeImretroReader(mode byte, palette [][]byte, width, height uint16, pixels []byte) *bytes.Buffer {
	dimensions := (uint(width) << 12) | uint(height)
	b := bytes.NewBuffer([]byte{
		// signature/magic bytes
		'I', 'M', 'R', 'E', 'T', 'R', 'O',
		// Mode byte (8-bit, in-file palette)
		mode,
	})
	{
		w := bitio.NewWriter(b, 3)
		w.WriteBits(dimensions, 24)
	}
	for _, color := range palette {
		b.Write(color)
	}
	b.Write(pixels)
	return b
}
