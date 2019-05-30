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

type bytesKey []byte

const matchFindRetriesByte = 16
const minByteKeySize = 64

var errorByteKeyTooShort = errors.New("key is smaller than minimum length of 64 bytes")

func (bytesKey) GetVersion() EncoderInfo {
	return EncoderInfo{
		Name:    "byteec",
		Version: "0.2",
	}
}

func (k bytesKey) CheckValid() (bool, error) {
	if len(k) < minByteKeySize {
		return false, errorByteKeyTooShort
	}
	return true, nil
}

func (bytesKey) GetDictionaryChars() dictionaryChars {
	return dictionaryChars("abcdefghijk")
}

func (bytesKey) GetDictionary() dictionary {
	return dictionary{
		decoders: []decodeRef{
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

// EncodeBytes encodes a slice of bytes against a key which is a slice of bytes.
func EncodeBytes(input []byte, key []byte) ([]byte, error) {
	return encode(
		input, bytesKey(key), findBytePattern)
}

// EncodeBytesStream encodes a byte stream against a key which is a slice of bytes.
func EncodeBytesStream(input io.Reader, key []byte) (*io.PipeReader, error) {
	return encodeStream(
		input, bytesKey(key), findBytePattern)
}

// EncodeBytesStreamPartial encodes a byte stream partially against a key which is a slice of bytes.
// Arguments take and skip are used to determine how many bytes to take, and skip along a stream.
func EncodeBytesStreamPartial(input io.Reader, key []byte, take int, skip int) (*io.PipeReader, error) {
	return encodePartialStream(
		input, bytesKey(key), take, skip, findBytePattern)
}

// DecodeBytes decodes a slice of bytes against a key which is a slice of bytes.
func DecodeBytes(input []byte, key []byte) ([]byte, error) {
	return decode(
		input, bytesKey(key), 2, getByteDefs)
}

// DecodeBytesStream decodes a byte stream against a key which is a slice of bytes.
func DecodeBytesStream(input io.Reader, key []byte) (*io.PipeReader, error) {
	return decodeStream(
		input, bytesKey(key), 2, getByteDefs)
}

// DecodeBytesStreamPartial decodes a byte stream with delimiters
// against a key which is a slice of bytes.
func DecodeBytesStreamPartial(input io.Reader, key []byte) (*io.PipeReader, error) {
	return decodePartialStream(
		input, bytesKey(key), 2, getByteDefs)
}

// AnalyzeBytesKey takes a slice of bytes and analyzes its scale of usefulness at encoding.
func AnalyzeBytesKey(key []byte) (scale int) {
	dict := bytesKey(key).GetDictionary()
	found := 0.0
	for i := 0; i < 255; i++ {
		perByte := 0.0
		for j := 0; j < matchFindRetriesByte; j++ {
			randByte := key[rand.Intn(len(key))]
			for k := 0; k < len(key); k++ {
				success, _, _ := checkByteMatch(byte(i), randByte, key[k], dict)
				if success {
					perByte++
				}
				continue
			}
		}
		found += perByte / float64(matchFindRetriesByte)
	}
	return int(float64(found) / 255.0)
}

func getByteDefs(key encodingKey, group decodeGroup) (byte, error) {
	if len(group.place) < 2 {
		return 0, errors.New("decode group missing locations")
	}
	bytes, ok := key.(bytesKey)
	if !ok {
		return 0, errors.New("failed to cast key")
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
			if len(bytes) >= loc1 {
				change1 = bytes[loc1] + g.amount
			} else {
				return 0, errors.New("decode error")
			}
		}
	}
	for _, g := range dict.decoders {
		if g.character == group.kind[1] {
			if len(bytes) >= loc2 {
				change2 = bytes[loc2] + g.amount
			} else {
				return 0, errors.New("decode error")
			}
		}
	}
	return change2 - change1, nil
}

func findBytePattern(char byte, key encodingKey) ([]byte, error) {
	bytesKey, ok := key.(bytesKey)
	if !ok {
		return nil, errorKeyCastFailed
	}
	pattern := make([]byte, 0)
	var err error

	for i := 0; i < matchFindRetriesByte; i++ {
		pattern, err = getBytePattern(char, bytesKey)
		if err == nil {
			return pattern, nil
		}
	}

	return nil, err
}

func getBytePattern(char byte, key bytesKey) ([]byte, error) {
	bounds := len(key)
	current := rand.Intn(bounds)
	startFinding := rand.Intn(bounds)
	dictionary := key.GetDictionary()

	var pattern []byte
	var err error

	if startFinding > bounds/2 {
		for x := startFinding; x >= 0; x-- {
			pattern, err = findBytePartner(current, x, char, key, dictionary)
			if err == nil {
				return pattern, nil
			}
		}
	} else {
		for x := startFinding; x < bounds; x++ {
			pattern, err = findBytePartner(current, x, char, key, dictionary)
			if err == nil {
				return pattern, nil
			}
		}
	}

	return nil, err
}

func findBytePartner(
	current int,
	checked int,
	difference byte,
	bytes []byte,
	dict dictionary) ([]byte, error) {
	if match, firstType, secondType := checkByteMatch(
		difference, bytes[current], bytes[checked], dict); match {
		return []byte(fmt.Sprintf(
			"%s%v%s%v",
			string(firstType), current,
			string(secondType), checked)), nil
	}

	return nil, errorMatchNotFound
}

func checkByteMatch(
	diff byte,
	current byte,
	checked byte,
	dict dictionary) (bool, uint8, uint8) {
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
