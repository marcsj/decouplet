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
	color color.Color,
	img image.Image) ([]byte, error) {
		bounds := img.Bounds()
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			newColor := img.At(x, y)
			or, og, ob, oa := color.RGBA()
			r, g, b, a := newColor.RGBA()
			if uint8(r) == uint8(or) + uint8(difference) {
				return []byte(fmt.Sprintf("r%v,%vr%v,%v",
					location.x, location.y, x, y)), nil
			} else if uint8(g) == uint8(og) + uint8(difference) {
				return []byte(fmt.Sprintf("g%v,%vg%v,%v",
					location.x, location.y, x, y)), nil
			} else if uint8(b) == uint8(ob) + uint8(difference) {
				return []byte(fmt.Sprintf("b%v,%vb%v,%v",
					location.x, location.y, x, y)), nil
			} else if uint8(a) == uint8(oa) + uint8(difference) {
				return []byte(fmt.Sprintf("a%v,%va%v,%v",
					location.x, location.y, x, y)), nil
			}
		}
	}
	return nil, errors.New(errorMatchNotFound)
}