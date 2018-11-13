package main

import (
	"errors"
	"image"
	"image/color"
	"strconv"
	"strings"
)

const colorCodes = "rgba"
const errorColorCode = "message does not contain color code"
const errorLocation = "incorrect text for pixel location"

type RGBA string

type pixelRef struct {
	color    RGBA
	location location
}

type location struct {
	x int
	y int
}


func TransdecodeImage(message []byte, img image.Image) ([]byte, error){
	msg := string(message)
	var translation []byte
	for {
		pixelRef, end, err := getPixelPair(msg)
		if err != nil {
			return nil, err
		}
		firstColor := img.At(
			pixelRef[0].location.x, pixelRef[0].location.y)
		secondColor := img.At(
			pixelRef[1].location.x, pixelRef[1].location.y)
		difference :=
			pixelNumber(pixelRef[1].color, secondColor) -
			pixelNumber(pixelRef[0].color, firstColor)
		translation = append(translation, byte(difference))
		if len(msg[end:]) < 2 {
			return translation, nil
		}
		msg = msg[end:]
	}
}

func pixelNumber(rgb RGBA, color color.Color) uint8 {
	r, g, b, a := color.RGBA()
	num := uint8(0)
	switch rgb[0] {
	case 'r':
		num = uint8(r)
	case 'g':
		num = uint8(g)
	case 'b':
		num = uint8(b)
	case 'a':
		num = uint8(a)
	}
	return num
}

func getPixelPair(message string) ([2]pixelRef, int, error) {
	var pair [2]pixelRef
	pixel1, newStart, err := getPixelRef(message)
	if err != nil {
		return pair, 0, err
	}
	pixel2, end, err := getPixelRef(message[newStart:])
	if err != nil {
		return pair, 0, err
	}
	pair[0] = pixel1
	pair[1] = pixel2
	return pair, end+newStart, nil
}

func getPixelRef(message string) (pixelRef, int, error) {
	pixelRef := pixelRef{}
	end := 0
	if !strings.ContainsAny(message[0:1], colorCodes) {
		return pixelRef, end, errors.New(errorColorCode)
	}
	pixelRef.color = RGBA(message[0:1])
	for i := 1; i < len(message); i++ {
		if strings.ContainsAny(message[i:i+1], colorCodes) {
			locString := message[1:i]
			end = i
			loc, err := getTextLocation(locString)
			if err != nil {
				return pixelRef, end, err
			}
			pixelRef.location = loc
			break
		}
	}
	if end == 0 {
		locString := message[1:]
		loc, err := getTextLocation(locString)
		if err != nil {
			return pixelRef, end, err
		}
		pixelRef.location = loc
		end = len(message)
	}
	return pixelRef, end, nil
}

func getTextLocation(loc string) (location, error) {
	location := location{}
	locs := strings.Split(loc, ",")
	if len(locs) < 2 {
		return location, errors.New(errorLocation)
	}
	x, err := strconv.Atoi(locs[0])
	if err != nil {
		return location, err
	}
	y, err := strconv.Atoi(locs[1])
	if err != nil {
		return location, err
	}
	location.x = x
	location.y = y
	return location, nil
}