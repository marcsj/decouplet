package decouplet

import (
	"errors"
)

var errorMatchNotFound = errors.New("match not found")

const partialStart string = ";[&"
const partialEnd string = "&];"

var partialStartBytes = []byte(partialStart)
var partialEndBytes = []byte(partialEnd)

type encoderType string

type dictionaryChars string

type byteGroup struct {
	bytes []byte
}

type decodeGroup struct {
	kind  []uint8
	place []string
}

type dictionary struct {
	decoders []decodeRef
}

type decodeRef struct {
	character uint8
	amount    uint8
}

type splitInfo struct {
	chars  dictionaryChars
	groups int
}

type location struct {
	x int
	y int
}

func (chars dictionaryChars) checkIn(a byte) bool {
	for i := range chars {
		if a == chars[i] {
			return true
		}
	}
	return false
}
