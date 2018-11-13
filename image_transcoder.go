package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"
)

const errorMatchNotFound = "match not found"

func init() {
	rand.Seed(time.Now().Unix())
}

func TranscodeImage(message []byte, img image.Image) ([]byte, error){
	newMessage := make([]byte, 0)
	for _, char := range message {
		msg, err := findBytePattern(char, img)
		if err != nil {
			return nil, err
		}
		for _, b := range msg {
			newMessage = append(newMessage, b)
		}
	}
	return newMessage, nil
}

func findBytePattern(char byte, img image.Image) ([]byte, error) {
	bounds := img.Bounds()
	startX := rand.Intn(bounds.Max.X)
	startY := rand.Intn(bounds.Max.Y)
	firstColor := img.At(startX, startY)

	pattern, err := findPixelPartner(
		location{x: startX, y: startY}, char, firstColor, img)
	if err != nil && err.Error() == errorMatchNotFound {
		startX = rand.Intn(bounds.Max.X)
		startY = rand.Intn(bounds.Max.Y)
		firstColor = img.At(startX, startY)

		pattern, err = findPixelPartner(
			location{x: startX, y: startY}, char, firstColor, img)
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
	img image.Image) ([]byte, error) {
		bounds := img.Bounds()
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			checkedColor := img.At(x, y)
			if match, firstColor, secondColor := checkColorMatch(
				difference, currentColor, checkedColor); match {
					return []byte(fmt.Sprintf(
						"%s%v,%v%s%v%v",
						firstColor, location.x, location.y,
						secondColor, x, y)), nil
			}
		}
	}
	return nil, errors.New(errorMatchNotFound)
}

func checkColorMatch(
	diff byte,
	current color.Color,
	checked color.Color) (bool, string, string) {
	or, og, ob, oa := current.RGBA()
	r, g, b, a := checked.RGBA()
	currentColors := map[string]uint8{
		"r": uint8(or),
		"g": uint8(og),
		"b": uint8(ob),
		"a": uint8(oa),
	}
	checkedColors := map[string]uint8{
		"r": uint8(r),
		"g": uint8(g),
		"b": uint8(b),
		"a": uint8(a),
	}
	for v := range currentColors {
		for k := range checkedColors {
			if uint8(currentColors[k]) ==
				uint8(currentColors[v]) + uint8(diff) {
				return true, v, k
			}
		}
	}
	return false, "", ""
}