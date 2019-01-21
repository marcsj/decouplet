package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type colorChecked struct {
	color string
	amount uint8
}

type imageKey struct {
	image.Image
}

func (imageKey) KeyType() TranscoderType {
	return TranscoderType("imgtc")
}

func (imageKey) DictionaryChars() dictionaryChars {
	return dictionaryChars("rgbacmyk")
}

func (imageKey) Dictionary() dictionary {
	return dictionary{
		decoders: []decoder{
			{
				character: 'r',
				amount: 0,
			},
			{
				character: 'g',
				amount: 0,
			},
			{
				character: 'b',
				amount: 0,
			},
			{
				character: 'a',
				amount: 0,
			},
			{
				character: 'c',
				amount: 0,
			},
			{
				character: 'm',
				amount: 0,
			},
			{
				character: 'y',
				amount: 0,
			},
			{
				character: 'k',
				amount: 0,
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


func TranscodeImage(input []byte, key image.Image) ([]byte, error) {
	return Transcode(
		input, imageKey{key}, findPixelPattern)
	return nil, nil
}

func TransdecodeImage(input []byte, key image.Image) ([]byte, error) {
	return Transdecode(
		input, imageKey{key}, 2, pixelDiff)
}


func pixelDiff(key key, decodeGroups []decodeGroup) (string, error) {
	img, ok := key.(imageKey); if !ok {
		return "", errors.New("failed to cast key")
	}
	returnString := ""
	for _, dec := range decodeGroups {
		b, err := getImgDefs(img, dec, key.Dictionary())
		if err != nil {
			return "", err
		}
		returnString += string(b)
	}
	return returnString, nil
}

func getImgDefs(img imageKey, group decodeGroup, dict dictionary) (byte, error){
	if len(group.place) < 2 {
		return 0, errors.New("decode group missing locations")
	}

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
	return change2-change1, nil
}

func findPixelPattern(char byte, key key) ([]byte, error) {
	img, ok := key.(imageKey); if !ok {
		return nil, errors.New("failed to cast key")
	}
	bounds := img.Bounds()
	startX := rand.Intn(bounds.Max.X)
	startY := rand.Intn(bounds.Max.Y)
	firstColor := img.At(startX, startY)

	pattern, err := findPixelPartner(
		location{x: startX, y: startY}, char, firstColor, img, key.Dictionary())
	if err != nil && err.Error() == errorMatchNotFound {
		startX = rand.Intn(bounds.Max.X)
		startY = rand.Intn(bounds.Max.Y)
		firstColor = img.At(startX, startY)

		pattern, err = findPixelPartner(
			location{x: startX, y: startY}, char, firstColor, img, key.Dictionary())
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return pattern, nil
}

func findPixelPartner(
	location location,
	difference byte,
	currentColor color.Color,
	img image.Image,
	dict dictionary) ([]byte, error) {
	bounds := img.Bounds()
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			checkedColor := img.At(x, y)
			if match, firstType, secondType := checkColorMatch(
				difference, currentColor, checkedColor, dict); match {
					firstLocation := GetPixelNumber(
						location.x, location.y, bounds.Max.X)
					secondLocation := GetPixelNumber(x, y, bounds.Max.X)
					return []byte(fmt.Sprintf(
						"%s%v%s%v",
						string(firstType), firstLocation,
						string(secondType), secondLocation)), nil
			}
		}
	}
	return nil, errors.New(errorMatchNotFound)
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
				currentColors.decoders[v].amount + uint8(diff) {
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
	x, y := GetCoordinates(loc, imageWidth)
	location.x = x
	location.y = y
	return location, nil
}