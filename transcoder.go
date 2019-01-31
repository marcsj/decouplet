package main

import (
	"errors"
	"log"
	"sync"
)

type Key interface {
	GetKeyType() TranscoderType
	GetDictionaryChars() DictionaryChars
	GetDictionary() Dictionary
}

func Transcode(
	input []byte,
	key Key,
	encoder func(byte, Key) ([]byte, error),
) (output []byte, err error) {
	bytes, err := WriteVersion(key.GetKeyType())
	if err != nil {
		return nil, err
	}

	byteGroups := make([]ByteGroup, len(input))
	wg := &sync.WaitGroup{}
	wg.Add(len(input))

	for i := range input {
		bytesRelay(i, input, byteGroups, key, encoder, wg)
	}

	for _, byteGroup := range byteGroups {
		for _, b := range byteGroup.bytes {
			bytes = append(bytes, b)
		}
	}
	return bytes, nil
}

func TranscodeConcurrent(
	input []byte,
	key Key,
	encoder func(byte, Key) ([]byte, error),
) (output []byte, err error) {
	bytes, err := WriteVersion(key.GetKeyType())
	if err != nil {
		return nil, err
	}

	byteGroups := make([]ByteGroup, len(input))
	wg := &sync.WaitGroup{}
	wg.Add(len(input))

	for i := range input {
		go bytesRelay(i, input, byteGroups, key, encoder, wg)
	}
	wg.Wait()

	for _, byteGroup := range byteGroups {
		for _, b := range byteGroup.bytes {
			bytes = append(bytes, b)
		}
	}
	return bytes, nil
}

func bytesRelay(
	index int,
	input []byte,
	bytes []ByteGroup,
	key Key,
	encoder func(byte, Key) ([]byte, error),
	wg *sync.WaitGroup) {
	byteGroup := ByteGroup{
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
	key Key,
	groups int,
	decodeFunc func(Key, DecodeGroup) (byte, error),
) (output []byte, err error) {
	err = CheckTranscoder(key.GetKeyType(), &input)
	if err != nil {
		return nil, err
	}
	decodeGroups, err := findDecodeGroups(input, key.GetDictionaryChars(), groups)
	decoded, err := decodeBytes(key, decodeGroups, decodeFunc)
	return decoded, err
}

func findDecodeGroups(
	input []byte,
	characters DictionaryChars,
	numGroups int,
) (decodeGroups []DecodeGroup, err error) {
	if !characters.CheckIn(input[0]) {
		return decodeGroups, errors.New("no decode characters found")
	}
	decode := DecodeGroup{
		kind:  []uint8{},
		place: []string{},
	}
	buffer := make([]uint8, 0)
	numberAdded := 0

	for i := range input {
		if characters.CheckIn(input[i]) {
			if len(buffer) > 0 {
				decode.place = append(decode.place, string(buffer))
				buffer = make([]uint8, 0)
				if numberAdded == numGroups {
					numberAdded = 0
					decodeGroups = append(decodeGroups, decode)
					decode = DecodeGroup{
						kind:  []uint8{},
						place: []string{},
					}
				}
			}
			if i != len(input)-1 {
				decode.kind = append(decode.kind, input[i])
				numberAdded++
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
	key Key,
	decodeGroups []DecodeGroup,
	decodeFunc func(Key, DecodeGroup) (byte, error),
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
