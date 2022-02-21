package imretro

import (
	"image"
	"image/color"
	"io"

	"github.com/spenserblack/go-bitio"
	"github.com/spenserblack/go-byteutils"

	"github.com/imretro/go/internal/util"
)

// ImretroSignature is the "magic string" used for identifying an imretro file.
const ImretroSignature = "IMRETRO"

// DecodeError is an error signifying that something unexpected happened when
// decoding the imretro reader.
type DecodeError string

// Decode decodes an image in the imretro format.
//
// Custom color models can be used instead of the default color models. For
// simplicity's sake, a single ColorModel can be passed as the CustomModel. If
// multiple PixelMode values are expected, it is recommended to use a ModelMap
// for the CustomModel. See the documentation for the model types for more
// details. If the decoded image contains an in-image palette, the model will be
// generated from that instead of the custom value passed or the default models.
func Decode(r io.Reader, customModels CustomModel) (Image, error) {
	config, err := DecodeConfig(r, customModels)
	if err != nil {
		return nil, err
	}
	area := config.Width * config.Height
	pixelsForByte := 8 / config.ColorModel.(ColorModel).BitsPerPixel()
	bytesNeeded := area / pixelsForByte
	if area%pixelsForByte != 0 {
		bytesNeeded++
	}
	pixels := make([]byte, bytesNeeded)
	if _, err := io.ReadFull(r, pixels); err != nil {
		return nil, err
	}

	return imretroImage{config, pixels}, nil
}

// DecodeConfig returns the color model and dimensions of an imretro image
// without decoding the entire image.
//
// Custom color models can be used instead of the default model.
func DecodeConfig(r io.Reader, customModels CustomModel) (image.Config, error) {
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

	bitsPerPixel := mode & (0b11 << pixelBitsIndex)
	hasPalette := byteutils.BitAsBool(byteutils.GetR(mode, paletteIndex))

	width, height, err := decodeDimensions(r)
	if err != nil {
		return image.Config{}, err
	}

	var model color.Model
	if !hasPalette {
		var ok bool
		model, ok = modelMap.ColorModel(bitsPerPixel)
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
		model, err = decodeModel(r, modelSize, mode&EightBitColors != 0, mode&(0b11<<colorChannelIndex))
	}

	return image.Config{ColorModel: model, Width: width, Height: height}, err
}

// DecodeDimensions gets the dimensions from a reader.
func decodeDimensions(r io.Reader) (width, height int, err error) {
	var w, h uint
	reader := bitio.NewReader(r, 3)
	w, _, err = reader.ReadBits(12)
	if err != nil {
		return
	}
	h, _, err = reader.ReadBits(12)
	width = int(w)
	height = int(h)
	return
}

// DecodeModel will decode bytes into a ColorModel. The bytes decoded depend on
// the length of the ColorModel.
func decodeModel(r io.Reader, size int, accurateColors bool, colorChannels ModeFlag) (color.Model, error) {
	model := make(ColorModel, size)
	channelCount := -1
	chunkSize := 1
	bitsPerChannel := 2
	switch colorChannels {
	case Grayscale:
		channelCount = 1
	case RGB:
		channelCount = 3
	case RGBA:
		channelCount = 4
	}
	if accurateColors {
		chunkSize = channelCount
		bitsPerChannel = 8
	}
	reader := bitio.NewReader(r, chunkSize)
	for i := range model {
		channels := make(colorBytes, channelCount)
		for i := range channels {
			bits, _, err := reader.ReadBits(bitsPerChannel)
			if err != nil {
				return nil, err
			}
			filledBits := util.FillByte(byte(bits), byte(bitsPerChannel))
			channels[i] = filledBits
		}
		model[i] = channels
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
	return i, err
}

// GlobalDecodeConfig has the proper function type to be registered as a
// format.
func globalDecodeConfig(r io.Reader) (image.Config, error) {
	c, err := DecodeConfig(r, nil)
	return c, err
}
