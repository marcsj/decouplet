package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"sync"
	"time"
)

const imageTranscoderName = "imgtc"
const errorMatchNotFound = "match not found"

func init() {
	rand.Seed(time.Now().Unix())
}

type colorChecked struct {
	color string
	amount uint8
}

type byteGroup struct {
	bytes []byte
}

func TranscodeImage(message []byte, img image.Image) ([]byte, error) {
	newMessage := make([]byte, 0)
	newMessage, err := WriteVersion(imageTranscoderName, newMessage)
	if err != nil {
		return newMessage, err
	}

	byteList := make([]byteGroup, len(message))
	wg := sync.WaitGroup{}
	wg.Add(len(message))

	for i, b := range message {
		go getNewBytes(i, b, img, byteList, &wg)
	}
	wg.Wait()
	for _, byteGroup := range byteList {
		for _, byte := range byteGroup.bytes {
			newMessage = append(newMessage, byte)
		}
	}
	return newMessage, nil
}

func getNewBytes(
	index int,
	char byte,
	img image.Image,
	byteList []byteGroup,
	group *sync.WaitGroup) {

	byteGroup := byteGroup{
		bytes: make([]byte, 0),
	}
	msg, err := findBytePattern(char, img)
	if err != nil {
		log.Println(err.Error())
	}
	for _, b := range msg {
		byteGroup.bytes = append(byteGroup.bytes, b)
	}
	byteList[index] = byteGroup
	group.Done()
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
			if match, firstType, secondType := checkColorMatch(
				difference, currentColor, checkedColor); match {
					firstLocation := GetPixelNumber(
						location.x, location.y, bounds.Max.X)
					secondLocation := GetPixelNumber(x, y, bounds.Max.X)
					return []byte(fmt.Sprintf(
						"%s%v%s%v",
						firstType, firstLocation,
						secondType, secondLocation)), nil
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
	oc, om, oy, ok := color.RGBToCMYK(uint8(or), uint8(og), uint8(ob))
	c, m, y, k := color.RGBToCMYK(uint8(r), uint8(g), uint8(b))
	currentColors := []colorChecked{
		colorChecked{
			color: "r",
			amount: uint8(or),
		},
		colorChecked{
			color: "g",
			amount: uint8(og),
		},
		colorChecked{
			color: "b",
			amount: uint8(ob),
		},
		colorChecked{
			color: "a",
			amount: uint8(oa),
		},
		colorChecked{
			color: "c",
			amount: uint8(oc),
		},
		colorChecked{
			color: "m",
			amount: uint8(om),
		},
		colorChecked{
			color: "y",
			amount: uint8(oy),
		},
		colorChecked{
			color: "k",
			amount: uint8(ok),
		},
	}
	checkedColors := []colorChecked{
		colorChecked{
			color: "r",
			amount: uint8(r),
		},
		colorChecked{
			color: "g",
			amount: uint8(g),
		},
		colorChecked{
			color: "b",
			amount: uint8(b),
		},
		colorChecked{
			color: "a",
			amount: uint8(a),
		},
		colorChecked{
			color: "c",
			amount: uint8(c),
		},
		colorChecked{
			color: "m",
			amount: uint8(m),
		},
		colorChecked{
			color: "y",
			amount: uint8(y),
		},
		colorChecked{
			color: "k",
			amount: uint8(k),
		},
	}
	for v := range currentColors {
		for k := range checkedColors {
			if checkedColors[k].amount ==
				currentColors[v].amount + uint8(diff) {
				return true,
				currentColors[v].color,
				currentColors[k].color
			}
		}
	}
	return false, "", ""
}