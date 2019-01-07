package main

import (
	"errors"
	"image"
	"image/color"
	"strconv"
	"strings"
)

const colorCodes = "rgbacmyk"
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
	msg, err := CheckTranscoder(imageTranscoderName, msg)
	if err != nil {
		return translation, err
	}
	for {
		pixelRef, end, err := getPixelPair(msg, img.Bounds().Max.X)
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

func pixelNumber(rgb RGBA, checkedColor color.Color) uint8 {
	r, g, b, a := checkedColor.RGBA()
	c, m, y, k := color.RGBToCMYK(uint8(r), uint8(g), uint8(b))
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
	case 'c':
		num = uint8(c)
	case 'm':
		num = uint8(m)
	case 'y':
		num = uint8(y)
	case 'k':
		num = uint8(k)
	}
	return num
}

func getPixelPair(message string, imageWidth int) ([2]pixelRef, int, error) {
	var pair [2]pixelRef
	pixel1, newStart, err := getPixelRef(message, imageWidth)
	if err != nil {
		return pair, 0, err
	}
	pixel2, end, err := getPixelRef(message[newStart:], imageWidth)
	if err != nil {
		return pair, 0, err
	}
	pair[0] = pixel1
	pair[1] = pixel2
	return pair, end+newStart, nil
}

func getPixelRef(message string, imageWidth int) (pixelRef, int, error) {
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
			loc, err := getTextLocation(locString, imageWidth)
			if err != nil {
				return pixelRef, end, err
			}
			pixelRef.location = loc
			break
		}
	}
	if end == 0 {
		locString := message[1:]
		loc, err := getTextLocation(locString, imageWidth)
		if err != nil {
			return pixelRef, end, err
		}
		pixelRef.location = loc
		end = len(message)
	}
	return pixelRef, end, nil
}

func getTextLocation(loc string, imageWidth int) (location, error) {
	location := location{}
	pixelLoc, err := strconv.Atoi(loc)
	if err != nil {
		return location, err
	}
	x, y := GetCoordinates(pixelLoc, imageWidth)
	location.x = x
	location.y = y
	return location, nil
}