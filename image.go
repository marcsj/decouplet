package decouplet

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func LoadImage(filename string) (image.Image, error) {
	imageFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer imageFile.Close()

	img, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func getPixelNumber(x int, y int, width int) int {
	return y*width + x
}

func getCoordinates(pixelNumber int, imageWidth int) (int, int) {
	x := pixelNumber % imageWidth
	y := pixelNumber / imageWidth
	return x, y
}
