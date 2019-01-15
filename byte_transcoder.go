package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type byteChecked struct {
	kind string
	amount uint8
}

type keyBytes []byte

func (keyBytes) KeyType() TranscoderType {
	return TranscoderType("bytetc")
}


func TranscodeBytes(input []byte, key []byte) ([]byte, error) {
	return Transcode(
		input, keyBytes(key), findBytePattern)
}

func findBytePattern(char byte, key key) ([]byte, error) {
	bytes, ok := key.(keyBytes); if !ok {
		return nil, errors.New("failed to cast key")
	}
	bounds := len(bytes)
	startX := rand.Intn(bounds)
	firstByte := bytes[startX]

	pattern, err := findBytePartner(
		location{x: startX}, char, byte(firstByte), bytes)
	if err != nil && err.Error() == errorMatchNotFound {
		startX = rand.Intn(bounds)
		firstByte := bytes[startX]

		pattern, err = findBytePartner(
			location{x: startX}, char, byte(firstByte), bytes)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return pattern, nil
}

func findBytePartner(
	location location,
	difference byte,
	currentByte byte,
	bytes []byte) ([]byte, error) {
	boundary := len(bytes)
	for x := 0; x < boundary; x++ {
		checkedByte := bytes[x]
		if match, firstType, secondType := checkByteMatch(
			difference, currentByte, checkedByte); match {
				return []byte(fmt.Sprintf(
						"%s%v%s%v",
						firstType, location.x,
						secondType, x)), nil
			}
	}
	return nil, errors.New(errorMatchNotFound)
}

func checkByteMatch(
	diff byte,
	current byte,
	checked byte) (bool, string, string) {
	currentBytes := getByteChecks(current)
	checkedBytes := getByteChecks(checked)
	for v := range currentBytes {
		for k := range checkedBytes {
			if checkedBytes[k].amount ==
				currentBytes[v].amount + uint8(diff) {
				return true,
				currentBytes[v].kind,
				currentBytes[k].kind
			}
		}
	}
	return false, "", ""
}

func getByteChecks(current byte) []byteChecked {
	return []byteChecked {
		byteChecked{
			kind: "a",
			amount: current+1,
		},
		byteChecked{
			kind: "b",
			amount: current+2,
		},
		byteChecked{
			kind: "c",
			amount: current+4,
		},
		byteChecked{
			kind: "d",
			amount: current+6,
		},
		byteChecked{
			kind: "e",
			amount: current+8,
		},
		byteChecked{
			kind: "f",
			amount: current+10,
		},
		byteChecked{
			kind: "g",
			amount: current+16,
		},
		byteChecked{
			kind: "h",
			amount: current+32,
		},
		byteChecked{
			kind: "i",
			amount: current+64,
		},
		byteChecked{
			kind: "j",
			amount: current+128,
		},
	}
}