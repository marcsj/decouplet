package decouplet

import (
	"errors"
)

var errorMatchNotFound = errors.New("match not found")
var errorDecodeNotFound = errors.New("valid decode character not found")
var errorKeyCastFailed = errors.New("failed to cast key")
var errorDecodeGeneric = errors.New("decode error")
var errorDecodeGroup = errors.New("decode groups missing locations")

const partialStart string = ";[&"
const partialEnd string = "&];"

var partialStartBytes = []byte(partialStart)
var partialEndBytes = []byte(partialEnd)

type dictionarySet string

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
	chars  dictionarySet
	groups int
}

type location struct {
	x int
	y int
}

func (chars dictionarySet) checkIn(a byte) bool {
	for i := range chars {
		if a == chars[i] {
			return true
		}
	}
	return false
}
