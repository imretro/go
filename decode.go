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
		var modelSize int
		switch bitsPerPixel {
		case OneBit:
			modelSize = 1 << 1
		case TwoBit:
			modelSize = 1 << 2
		case EightBit:
			modelSize = 1 << 8
		default:
			return image.Config{}, MissingModelError(bitsPerPixel)
		}
		model, err = decodeModel(r, modelSize, mode&EightBitColors != 0)
	}

	return image.Config{model, width, height}, err
}

// DecodeModel will decode bytes into a ColorModel. The bytes decoded depend on
// the length of the ColorModel.
func decodeModel(r io.Reader, size int, accurateColors bool) (color.Model, error) {
	model := make(ColorModel, size)
	buffSize := 1
	if accurateColors {
		buffSize = 4
	}
	buff := make([]byte, buffSize)
	for i := range model {
		if _, err := io.ReadFull(r, buff); err != nil {
			return nil, err
		}
		model[i] = util.ColorFromBytes(buff)
	}
	return model, nil
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
