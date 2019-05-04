package decouplet

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type byteChecked struct {
	kind   string
	amount uint8
}

type keyBytes []byte

func (keyBytes) GetKeyType() TranscoderType {
	return TranscoderType("bytetc")
}

func (keyBytes) GetDictionaryChars() DictionaryChars {
	return DictionaryChars("abcdefghijk")
}

func (keyBytes) GetDictionary() Dictionary {
	return Dictionary{
		decoders: []Decoder{
			{
				character: 'a',
				amount:    0,
			},
			{
				character: 'b',
				amount:    1,
			},
			{
				character: 'c',
				amount:    2,
			},
			{
				character: 'd',
				amount:    4,
			},
			{
				character: 'e',
				amount:    6,
			},
			{
				character: 'f',
				amount:    8,
			},
			{
				character: 'g',
				amount:    10,
			},
			{
				character: 'h',
				amount:    16,
			},
			{
				character: 'i',
				amount:    32,
			},
			{
				character: 'j',
				amount:    64,
			},
			{
				character: 'k',
				amount:    128,
			},
		},
	}
}

func TranscodeBytes(input []byte, key []byte) ([]byte, error) {
	return Transcode(
		input, keyBytes(key), findBytePattern)
}

func TranscodeBytesConcurrent(input []byte, key []byte) ([]byte, error) {
	return TranscodeConcurrent(
		input, keyBytes(key), findBytePattern)
}

func TranscodeBytesStream(input io.Reader, key []byte) (io.Reader, error) {
	return TranscodeStream(
		input, keyBytes(key), findBytePattern)
}

func TranscodeBytesStreamPartial(input io.Reader, key []byte, take int, skip int) (io.Reader, error) {
	return TranscodeStreamPartial(
		input, keyBytes(key), take, skip, findBytePattern)
}

func TransdecodeBytes(input []byte, key []byte) ([]byte, error) {
	return Transdecode(
		input, keyBytes(key), 2, getByteDefs)
}

func TransdecodeBytesStream(input io.Reader, key []byte) (io.Reader, error) {
	return TransdecodeStream(
		input, keyBytes(key), 2, getByteDefs)
}

func TransdecodeBytesStreamPartial(input io.Reader, key []byte) (io.Reader, error) {
	return TransdecodeStreamPartial(
		input, keyBytes(key), 2, getByteDefs)
}

func getByteDefs(key Key, group DecodeGroup) (byte, error) {
	if len(group.place) < 2 {
		return 0, errors.New("decode group missing locations")
	}
	bytes, ok := key.(keyBytes)
	if !ok {
		return 0, errors.New("failed to cast Key")
	}
	dict := key.GetDictionary()

	loc1, err := strconv.Atoi(group.place[0])
	if err != nil {
		return 0, err
	}
	loc2, err := strconv.Atoi(group.place[1])
	if err != nil {
		return 0, err
	}

	var change1 uint8
	var change2 uint8
	for _, g := range dict.decoders {
		if g.character == group.kind[0] {
			change1 = bytes[loc1] + g.amount
		}
	}
	for _, g := range dict.decoders {
		if g.character == group.kind[1] {
			change2 = bytes[loc2] + g.amount
		}
	}
	return change2 - change1, nil
}

func findBytePattern(char byte, key Key) ([]byte, error) {
	bytes, ok := key.(keyBytes)
	if !ok {
		return nil, errors.New("failed to cast Key")
	}
	bounds := len(bytes)
	startX := rand.Intn(bounds)
	firstByte := bytes[startX]

	pattern, err := findBytePartner(
		location{x: startX}, char, byte(firstByte), bytes, key.GetDictionary())
	if err != nil && err.Error() == errorMatchNotFound {
		startX = rand.Intn(bounds)
		firstByte := bytes[startX]

		pattern, err = findBytePartner(
			location{x: startX}, char, byte(firstByte), bytes, key.GetDictionary())
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
	bytes []byte,
	dict Dictionary) ([]byte, error) {
	boundary := len(bytes)
	for x := 0; x < boundary; x++ {
		checkedByte := bytes[x]
		if match, firstType, secondType := checkByteMatch(
			difference, currentByte, checkedByte, dict); match {
			return []byte(fmt.Sprintf(
				"%s%v%s%v",
				string(firstType), location.x,
				string(secondType), x)), nil
		}
	}
	return nil, errors.New(errorMatchNotFound)
}

func checkByteMatch(
	diff byte,
	current byte,
	checked byte,
	dict Dictionary) (bool, uint8, uint8) {
	for v := range dict.decoders {
		for k := range dict.decoders {
			if checked+dict.decoders[k].amount ==
				current+dict.decoders[v].amount+uint8(diff) {
				return true,
					dict.decoders[v].character,
					dict.decoders[k].character
			}
		}
	}
	return false, 0, 0
}
