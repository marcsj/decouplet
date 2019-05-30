package decouplet

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type imageKey struct {
	image.Image
}

const matchFindRetriesImage = 4
const imageKeySize = 300

var errorImageKeyTooSmall = errors.New("key needs to be larger than 300x300")

func (imageKey) GetVersion() EncoderInfo {
	return EncoderInfo{
		Name:    "imgec",
		Version: "0.2",
	}
}

func (k imageKey) CheckValid() (bool, error) {
	if k.Image.Bounds().Max.X < imageKeySize || k.Image.Bounds().Max.Y < imageKeySize {
		return false, errorImageKeyTooSmall
	}
	return true, nil
}

func (imageKey) GetDictionaryChars() dictionaryChars {
	return dictionaryChars("rgbacmyk")
}

func (imageKey) GetDictionary() dictionary {
	return dictionary{
		decoders: []decodeRef{
			{
				character: 'r',
				amount:    0,
			},
			{
				character: 'g',
				amount:    0,
			},
			{
				character: 'b',
				amount:    0,
			},
			{
				character: 'a',
				amount:    0,
			},
			{
				character: 'c',
				amount:    0,
			},
			{
				character: 'm',
				amount:    0,
			},
			{
				character: 'y',
				amount:    0,
			},
			{
				character: 'k',
				amount:    0,
			},
		},
	}
}

func dictionaryRGBACMYK(col color.Color, dict dictionary) dictionary {
	r, g, b, a := col.RGBA()
	c, m, y, k := color.RGBToCMYK(uint8(r), uint8(g), uint8(b))
	for i := range dict.decoders {
		switch dict.decoders[i].character {
		case 'r':
			dict.decoders[i].amount = uint8(r)
		case 'g':
			dict.decoders[i].amount = uint8(g)
		case 'b':
			dict.decoders[i].amount = uint8(b)
		case 'a':
			dict.decoders[i].amount = uint8(a)
		case 'c':
			dict.decoders[i].amount = uint8(c)
		case 'm':
			dict.decoders[i].amount = uint8(m)
		case 'y':
			dict.decoders[i].amount = uint8(y)
		case 'k':
			dict.decoders[i].amount = uint8(k)
		}
	}
	return dict
}

// EncodeImage encodes a slice of bytes against an image key.
func EncodeImage(input []byte, key image.Image) ([]byte, error) {
	return encode(
		input, imageKey{key}, findPixelPattern)
	return nil, nil
}

// EncodeImageStream encodes a stream of bytes against an image key.
func EncodeImageStream(input io.Reader, key image.Image) (*io.PipeReader, error) {
	return encodeStream(
		input, imageKey{key}, findPixelPattern)
}

// EncodeImageStreamPartial encodes a byte stream partially against an image key.
// Arguments take and skip are used to determine how many bytes to take, and skip along a stream.
func EncodeImageStreamPartial(input io.Reader, key image.Image, take int, skip int) (*io.PipeReader, error) {
	return encodePartialStream(
		input, imageKey{key}, take, skip, findPixelPattern)
}

// DecodeImage encodes a slice of bytes against an image key.
func DecodeImage(input []byte, key image.Image) ([]byte, error) {
	return decode(
		input, imageKey{key}, 2, getImgDefs)
}

// DecodeImageStream decodes a stream of bytes against an image key.
func DecodeImageStream(input io.Reader, key image.Image) (*io.PipeReader, error) {
	return decodeStream(
		input, imageKey{key}, 2, getImgDefs)
}

// DecodeImageStreamPartial decodes a byte stream with delimiters against an image key.
func DecodeImageStreamPartial(input io.Reader, key image.Image) (*io.PipeReader, error) {
	return decodePartialStream(
		input, imageKey{key}, 2, getImgDefs)
}

func getImgDefs(key encodingKey, group decodeGroup) (byte, error) {
	if len(group.place) < 2 {
		return 0, errors.New("decode group missing locations")
	}
	img, ok := key.(imageKey)
	if !ok {
		return 0, errors.New("failed to cast key")
	}
	dict := key.GetDictionary()

	loc1, err := strconv.Atoi(group.place[0])
	if err != nil {
		return 0, err
	}
	loc2, err := strconv.Atoi(group.place[1])
	if err != nil {
		return 0, err
	}
	location1, err := getXYLocation(loc1, img.Bounds().Max.X)
	if err != nil {
		return 0, err
	}
	location2, err := getXYLocation(loc2, img.Bounds().Max.X)
	if err != nil {
		return 0, err
	}

	var change1 uint8
	var change2 uint8
	changeColor1 := img.At(location1.x, location1.y)
	changeColor2 := img.At(location2.x, location2.y)
	dict1 := dictionaryRGBACMYK(changeColor1, dict)
	dict2 := dictionaryRGBACMYK(changeColor2, dict)

	for _, g := range dict1.decoders {
		if g.character == group.kind[0] {
			change1 = g.amount
		}
	}
	for _, g := range dict2.decoders {
		if g.character == group.kind[1] {
			change2 = g.amount
		}
	}
	return change2 - change1, nil
}

func findPixelPattern(char byte, key encodingKey) ([]byte, error) {
	imageKey, ok := key.(imageKey)
	if !ok {
		return nil, errorKeyCastFailed
	}
	var pattern []byte
	var err error

	for i := 0; i < matchFindRetriesImage; i++ {
		pattern, err = getPixelPattern(char, imageKey)
		if err == nil {
			return pattern, nil
		}
	}

	return nil, err
}

func getPixelPattern(char byte, key imageKey) ([]byte, error) {
	bounds := key.Bounds()
	currentX := rand.Intn(bounds.Max.X)
	currentY := rand.Intn(bounds.Max.Y)
	startX := rand.Intn(bounds.Max.X)
	startY := rand.Intn(bounds.Max.Y)
	dictionary := key.GetDictionary()

	pattern := make([]byte, 0)
	var err error

	currentLocation := location{x: currentX, y: currentY}
	currentColor := key.At(currentX, currentY)

	if startX > bounds.Max.X/2 {
		for x := startX; x >= 0; x-- {
			if startY > bounds.Max.Y/2 {
				for y := startY; y >= 0; y-- {
					checkedLocation := location{x, y}
					checkedColor := key.At(x, y)
					pattern, err = findPixelPartner(
						currentLocation, checkedLocation, char, currentColor, checkedColor, key, dictionary)
					if err == nil {
						return pattern, nil
					}
				}
			} else {
				for y := startY; y < bounds.Max.Y; y++ {
					checkedLocation := location{x, y}
					checkedColor := key.At(x, y)
					pattern, err = findPixelPartner(
						currentLocation, checkedLocation, char, currentColor, checkedColor, key, dictionary)
					if err == nil {
						return pattern, nil
					}
				}
			}
		}
	} else {
		for x := startX; x < bounds.Max.X; x++ {
			if startY > bounds.Max.Y/2 {
				for y := startY; y >= 0; y-- {
					checkedLocation := location{x, y}
					checkedColor := key.At(x, y)
					pattern, err = findPixelPartner(
						currentLocation, checkedLocation, char, currentColor, checkedColor, key, dictionary)
					if err == nil {
						return pattern, nil
					}
				}
			} else {
				for y := startY; y < bounds.Max.Y; y++ {
					checkedLocation := location{x, y}
					checkedColor := key.At(x, y)
					pattern, err = findPixelPartner(
						currentLocation, checkedLocation, char, currentColor, checkedColor, key, dictionary)
					if err == nil {
						return pattern, nil
					}
				}
			}
		}
	}

	return nil, err
}

func findPixelPartner(
	currentLocation location,
	checkedLocation location,
	difference byte,
	currentColor color.Color,
	checkedColor color.Color,
	key image.Image,
	dict dictionary) ([]byte, error) {
	bounds := key.Bounds()
	if match, firstType, secondType := checkColorMatch(
		difference, currentColor, checkedColor, dict); match {
		firstLocation := getPixelNumber(
			currentLocation.x, currentLocation.y, bounds.Max.X)
		secondLocation := getPixelNumber(
			checkedLocation.x, checkedLocation.y, bounds.Max.X)
		return []byte(fmt.Sprintf(
			"%s%v%s%v",
			string(firstType), firstLocation,
			string(secondType), secondLocation)), nil
	}

	return nil, errorMatchNotFound
}

func checkColorMatch(
	diff byte,
	current color.Color,
	checked color.Color,
	dict dictionary) (bool, uint8, uint8) {
	currentColors := dictionaryRGBACMYK(current, dict)
	checkedColors := dictionaryRGBACMYK(checked, dict)
	for v := range currentColors.decoders {
		for k := range checkedColors.decoders {
			if checkedColors.decoders[k].amount ==
				currentColors.decoders[v].amount+uint8(diff) {
				return true,
					currentColors.decoders[v].character,
					currentColors.decoders[k].character
			}
		}
	}
	return false, 0, 0
}

func getXYLocation(loc int, imageWidth int) (location, error) {
	location := location{}
	x, y := getCoordinates(loc, imageWidth)
	location.x = x
	location.y = y
	return location, nil
}
