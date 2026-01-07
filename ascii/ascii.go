// Package ascii can convert a image pixel to a raw char
// base on it's RGBA value, in another word, input a image pixel
// output a raw char ascii.
package ascii

import (
	"image/color"
	"math"
	"reflect"

	"github.com/aybabtme/rgbterm"
)

// CharPixel is converted pixel ascii
type CharPixel struct {
	Char rune
	R    uint8
	G    uint8
	B    uint8
	A    uint8
}

// Options convert pixel to raw char
type Options struct {
	Pixels   []rune
	Reversed bool
	Colored  bool
	Block    bool
}

// DefaultOptions that contains the default pixels
var DefaultOptions = Options{
	// Pixels:   []rune(" .,:;i1tfLCG08@"),
	Pixels:   []rune(" .,:;i1tfLCG08@"),
	Reversed: false,
	Colored:  true,
	Block:    true,
}

// NewOptions create a new convert option
func NewOptions() Options {
	newOptions := Options{}
	newOptions.mergeOptions(&DefaultOptions)
	return newOptions
}

// mergeOptions merge two options
func (options *Options) mergeOptions(newOptions *Options) {
	options.Pixels = append([]rune{}, newOptions.Pixels...)
	options.Reversed = newOptions.Reversed
	options.Colored = newOptions.Colored
	options.Block = newOptions.Block
}

// NewPixelConverter create a new pixel converter
func NewPixelConverter() PixelConverter {
	return PixelASCIIConverter{}
}

// PixelConverter define the convert pixel operation
type PixelConverter interface {
	ConvertPixelToASCII(pixel color.Color, options *Options) string
	ConvertPixelToPixelASCII(pixel color.Color, options *Options) CharPixel
}

// PixelASCIIConverter responsible for pixel ascii conversion
type PixelASCIIConverter struct {
}

// ConvertPixelToPixelASCII convert a image pixel to CharPixel
func (converter PixelASCIIConverter) ConvertPixelToPixelASCII(pixel color.Color, options *Options) CharPixel {
	convertOptions := NewOptions()
	convertOptions.mergeOptions(options)

	if convertOptions.Reversed {
		convertOptions.Pixels = converter.reverse(convertOptions.Pixels)
	}

	var rawChar rune
	r := reflect.ValueOf(pixel).FieldByName("R").Uint()
	g := reflect.ValueOf(pixel).FieldByName("G").Uint()
	b := reflect.ValueOf(pixel).FieldByName("B").Uint()
	a := reflect.ValueOf(pixel).FieldByName("A").Uint()

	if convertOptions.Block {
		rawChar = '█'
	} else {
		value := converter.intensity(r, g, b, a)

		precision := float64(255*3) / float64(len(convertOptions.Pixels)-1)
		rawChar = convertOptions.Pixels[converter.roundValue(float64(value)/precision)]
	}

	return CharPixel{
		Char: rawChar,
		R:    uint8(r),
		G:    uint8(g),
		B:    uint8(b),
		A:    uint8(a),
	}
}

// ConvertPixelToASCII converts a pixel to a ASCII char string
func (converter PixelASCIIConverter) ConvertPixelToASCII(pixel color.Color, options *Options) string {
	convertOptions := NewOptions()
	convertOptions.mergeOptions(options)

	pixelASCII := converter.ConvertPixelToPixelASCII(pixel, options)
	rawChar, r, g, b := pixelASCII.Char, pixelASCII.R, pixelASCII.G, pixelASCII.B
	if convertOptions.Colored {
		return converter.decorateWithColor(r, g, b, rawChar)
	}
	return string([]rune{rawChar})
}

func (converter PixelASCIIConverter) roundValue(value float64) int {
	return int(math.Floor(value + 0.5))
}

func (converter PixelASCIIConverter) reverse(numbers []rune) []rune {
	for i := 0; i < len(numbers)/2; i++ {
		j := len(numbers) - i - 1
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}

func (converter PixelASCIIConverter) intensity(r, g, b, a uint64) uint64 {
	return (r + g + b) * a / 255
}

// decorateWithColor decorate the raw char with the color base on r,g,b value
func (converter PixelASCIIConverter) decorateWithColor(r, g, b uint8, rawChar rune) string {
	coloredChar := rgbterm.FgString(string([]rune{rawChar}), uint8(r), uint8(g), uint8(b))
	return coloredChar
}
