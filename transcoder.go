package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

const byteTranscoderName = "bytetc"

func init() {
	rand.Seed(time.Now().Unix())
}

type byteChecked struct {
	kind string
	amount uint8
}

func TranscodeBytes(message []byte, bytes []byte) ([]byte, error) {
	newMessage := make([]byte, 0)
	newMessage, err := WriteVersion(byteTranscoderName, newMessage)
	if err != nil {
		return newMessage, err
	}

	byteList := make([]byteGroup, len(message))
	wg := sync.WaitGroup{}
	wg.Add(len(message))

	for i, b := range message {
		go getByteNewBytes(i, b, bytes, byteList, &wg)
	}
	wg.Wait()
	for _, byteGroup := range byteList {
		for _, byte := range byteGroup.bytes {
			newMessage = append(newMessage, byte)
		}
	}
	return newMessage, nil
}

func getByteNewBytes(
	index int,
	char byte,
	bytes []byte,
	byteList []byteGroup,
	group *sync.WaitGroup) {

	byteGroup := byteGroup{
		bytes: make([]byte, 0),
	}
	msg, err := findBytePattern(char, bytes)
	if err != nil {
		log.Println(err.Error())
	}
	for _, b := range msg {
		byteGroup.bytes = append(byteGroup.bytes, b)
	}
	byteList[index] = byteGroup
	group.Done()
}

func findBytePattern(char byte, bytes []byte) ([]byte, error) {
	bounds := len(bytes)
	startX := rand.Intn(bounds)
	firstByte := bytes[startX]

	pattern, err := findBytePartner(
		location{x: startX}, char, firstByte, bytes)
	if err != nil && err.Error() == errorMatchNotFound {
		startX = rand.Intn(bounds)
		firstByte := bytes[startX]

		pattern, err = findBytePartner(
			location{x: startX}, char, firstByte, bytes)
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