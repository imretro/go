package imretro

import (
	"image"
	"image/color"
	"io"

	"github.com/spenserblack/go-byteutils"

	"github.com/imretro/go/internal/util"
)

// ImretroSignature is the "magic string" used for identifying an imretro file.
const ImretroSignature = "IMRETRO"

// BitsPerPixelIndex is the position of the two bits for the bits-per-pixel
// mode (7 is left-most).
const bitsPerPixelIndex byte = 6

// DecodeError is an error signifying that something unexpected happened when
// decoding the imretro reader.
type DecodeError string

// Decode decodes an image in the imretro format.
//
// Custom color models can be used instead of the default color models. See the
// documentation for the model types for more details. If the decoded image
// contains an in-image palette, the model will be generated from that instead
// of the custom value passed or the default models.
func Decode(r io.Reader, customModels ModelMap) (ImretroImage, error) {
	config, err := DecodeConfig(r, customModels)
	if err != nil {
		return nil, err
	}
	pixels, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return imretroImage{config, pixels}, nil
}

// DecodeConfig returns the color model and dimensions of an imretro image
// without decoding the entire image.
//
// Custom color models can be used instead of the default model.
func DecodeConfig(r io.Reader, customModels ModelMap) (image.Config, error) {
	var buff []byte
	var err error
	modelMap := customModels
	if modelMap == nil {
		modelMap = DefaultModelMap
	}

	buff = make([]byte, len(ImretroSignature)+1)
	mode, err := checkHeader(r, buff)
	if err != nil {
		return image.Config{}, err
	}

	bitsPerPixel := mode & (0b11 << bitsPerPixelIndex)
	hasPalette := byteutils.BitAsBool(byteutils.GetL(mode, PaletteIndex))

	buff = make([]byte, 3)
	_, err = io.ReadFull(r, buff)
	if err != nil {
		return image.Config{}, err
	}

	width, height := util.DimensionsFrom3Bytes(buff[0], buff[1], buff[2])


	var model color.Model
	if !hasPalette {
		var ok bool
		model, ok = modelMap[bitsPerPixel]
		if !ok {
			err = MissingModelError(bitsPerPixel)
		}
	} else {
		switch bitsPerPixel {
		case OneBit:
			model, err = decode1bitModel(r)
		case TwoBit:
			model, err = decode2bitModel(r)
		case EightBit:
			model, err = decode8bitModel(r)
		default:
			err = MissingModelError(bitsPerPixel)
		}
	}

	return image.Config{model, width, height}, err
}

func decode1bitModel(r io.Reader) (color.Model, error) {
	buff := make([]byte, 8)
	if _, err := io.ReadFull(r, buff); err != nil {
		return nil, err
	}
	model := NewOneBitColorModel(ColorFromBytes(buff[:4]), ColorFromBytes(buff[4:]))

	return model, nil
}

func decode2bitModel(r io.Reader) (color.Model, error) {
	buff := make([]byte, 16)
	if _, err := io.ReadFull(r, buff); err != nil {
		return nil, err
	}
	model := NewTwoBitColorModel(
		ColorFromBytes(buff[:4]),
		ColorFromBytes(buff[4:8]),
		ColorFromBytes(buff[8:12]),
		ColorFromBytes(buff[12:]),
	)

	return model, nil
}

func decode8bitModel(r io.Reader) (color.Model, error) {
	colors := make(ColorModel, 0, 256)
	buff := make([]byte, 4)

	for i := 0; i < cap(colors); i++ {
		if _, err := io.ReadFull(r, buff); err != nil {
			return nil, err
		}
		colors = append(colors, ColorFromBytes(buff))
	}

	return colors, nil
}

// CheckHeader confirms the reader is an imretro image by checking the "magic bytes",
// and returns the "mode".
func checkHeader(r io.Reader, buff []byte) (mode byte, err error) {
	_, err = io.ReadFull(r, buff)
	if err != nil {
		return
	}

	for i, b := range buff[:len(buff)-1] {
		if b != ImretroSignature[i] {
			return mode, DecodeError("unexpected signature byte")
		}
	}
	return buff[len(buff)-1], nil
}

// Error reports that the format could not be decoded as imretro.
func (e DecodeError) Error() string {
	return string(e)
}

func init() {
	image.RegisterFormat("imretro", ImretroSignature, globalDecode, globalDecodeConfig)
}

// GlobalDecode returns an image.Image instead of an ImretroImage so that it
// can be registered as a format.
func globalDecode(r io.Reader) (image.Image, error) {
	i, err := Decode(r, nil)
	return i.(image.Image), err
}

// GlobalDecodeConfig has the proper function type to be registered as a
// format.
func globalDecodeConfig(r io.Reader) (image.Config, error) {
	c, err := DecodeConfig(r, nil)
	return c, err
}
