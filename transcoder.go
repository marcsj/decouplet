package main

import (
	"errors"
	"log"
	"strings"
	"sync"
)

type TranscoderType string

type key interface{
	KeyType() TranscoderType
	DictionaryChars() dictionaryChars
	Dictionary() dictionary
}

func Transcode(
	input []byte,
	key key,
	encoder func(byte, key) ([]byte, error),
	) (output []byte, err error) {
	bytes, err := WriteVersion(key.KeyType())
	if err != nil {
		return nil, err
	}

	byteGroups := make([]byteGroup, len(bytes))
	wg := &sync.WaitGroup{}
	wg.Add(len(input))

	for i := range input {
		go bytesRelay(i, input, byteGroups, key, encoder, wg)
	}
	wg.Wait()

	for _, byteGroup := range byteGroups {
		for _, byte := range byteGroup.bytes {
			bytes = append(bytes, byte)
		}
	}
	return bytes, nil
}

func bytesRelay(
	index int,
	input []byte,
	bytes []byteGroup,
	key key,
	encoder func(byte, key) ([]byte, error),
	wg *sync.WaitGroup) {
	byteGroup := byteGroup{
		bytes: make([]byte, 0),
	}
	msg, err := encoder(input[index], key)
	if err != nil {
		wg.Done()
		log.Fatal(err)
		return
	}
	for _, b := range msg {
		byteGroup.bytes = append(byteGroup.bytes, b)
	}
	bytes[index] = byteGroup
	wg.Done()
}

func Transdecode(
	input []byte,
	key key,
	groups int,
	decoder func(key, []decodeGroup) (string, error),
	) (output []byte, err error) {
	msg := string(input)
	err = CheckTranscoder(key.KeyType(), &msg)
	if err != nil {
		return nil, err
	}
	decodeGroups, err := findDecodeGroups(msg, key.DictionaryChars(), groups)
	decoded, err := decoder(key, decodeGroups)
	return []byte(decoded), err
}

func findDecodeGroups(
	input string,
	characters dictionaryChars,
	numGroups int,
	) (decodeGroups []decodeGroup, err error) {
	if !strings.ContainsAny(input[0:1], string(characters)) {
		return decodeGroups, errors.New("no decode characters found")
	}
	start := false
	decode := decodeGroup{
		kind: []uint8{},
		place: []string{},
	}
	buffer := make([]uint8, 0)
	numberAdded := 0
	for i := range input {
		if !start {
			if strings.ContainsAny(string(input[i]), string(characters)) {
				start = true
				decode.place = append(decode.place, string(buffer))
			} else {
				buffer = append(buffer, input[i])
			}
		}
		if start {
			decode.kind = append(decode.kind, input[i])
			numberAdded ++
			if numberAdded == numGroups-1 {
				numberAdded = 0
				decodeGroups = append(decodeGroups, decode)
				decode = decodeGroup{
					kind: []uint8{},
					place: []string{},
				}
			}
			start = false
		}
	}
	return decodeGroups, nil
}