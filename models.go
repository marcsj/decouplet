package main

const errorMatchNotFound = "match not found"


type byteGroup struct {
	bytes []byte
}

type decodeGroup struct {
	kind []uint8
	place []string
}

type dictionary struct {
	decoders []decoder
}

type decoder struct {
	character uint8
	amount uint8
}

type location struct {
	x int
	y int
}

type dictionaryChars string

func inDictionary(a byte, chars dictionaryChars) bool {
	for i := range chars {
		if a == chars[i] {
			return true
		}
	}
	return false
}

