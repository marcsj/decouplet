package main

import (
	"errors"
	"log"
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

	byteGroups := make([]byteGroup, len(input))
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
	decodeFunc func(key, decodeGroup) (byte, error),
	) (output []byte, err error) {
	err = CheckTranscoder(key.KeyType(), &input)
	if err != nil {
		return nil, err
	}
	decodeGroups, err := findDecodeGroups(input, key.DictionaryChars(), groups)
	decoded, err := decodeBytes(key, decodeGroups, decodeFunc)
	return decoded, err
}

func findDecodeGroups(
	input []byte,
	characters dictionaryChars,
	numGroups int,
	) (decodeGroups []decodeGroup, err error) {
	if !inDictionary(input[0], characters) {
		return decodeGroups, errors.New("no decode characters found")
	}
	decode := decodeGroup{
		kind: []uint8{},
		place: []string{},
	}
	buffer := make([]uint8, 0)
	numberAdded := 0

	for i := range input {
		if inDictionary(input[i], characters) {
			if len(buffer) > 0 {
				decode.place = append(decode.place, string(buffer))
				buffer = make([]uint8, 0)
				if numberAdded == numGroups {
					numberAdded = 0
					decodeGroups = append(decodeGroups, decode)
					decode = decodeGroup{
						kind:  []uint8{},
						place: []string{},
					}
				}
			}
			if i != len(input)-1 {
				decode.kind = append(decode.kind, input[i])
				numberAdded ++
			}
		} else {
			buffer = append(buffer, input[i])
			if i == len(input)-1 {
				decode.place = append(decode.place, string(buffer))
				decodeGroups = append(decodeGroups, decode)
			}
		}
	}
	return decodeGroups, nil
}

func decodeBytes(
	key key,
	decodeGroups []decodeGroup,
	decodeFunc func(key, decodeGroup) (byte, error),
	) ([]byte, error) {
	returnBytes := make([]byte, 0)
	for _, dec := range decodeGroups {
		b, err := decodeFunc(key, dec)
		if err != nil {
			return nil, err
		}
		returnBytes = append(returnBytes, b)
	}
	return returnBytes, nil
}