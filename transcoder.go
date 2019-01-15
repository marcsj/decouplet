package main

import (
	"log"
	"sync"
)

type TranscoderType string

type key interface{
	KeyType() TranscoderType
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

func Transdecode(input []byte, decoder interface{}) (output []byte, err error) {
	return nil, nil
}