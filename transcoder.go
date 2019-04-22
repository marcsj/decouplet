package decouplet

import (
	"errors"
	"io"
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
	errorCh := make(chan error, len(input))
	wg := &sync.WaitGroup{}
	wg.Add(len(input))

	for i := range input {
		bytesRelay(
			i, input, byteGroups, key, encoder, errorCh, wg)
	}

	select {
	case err := <-errorCh:
		return nil, err
	default:
		break
	}

	for _, byteGroup := range byteGroups {
		for _, b := range byteGroup.bytes {
			bytes = append(bytes, b)
		}
	}
	return bytes, nil
}

func TranscodeStream(
	input io.Reader,
	key Key,
	encoder func(byte, Key) ([]byte, error),
) (reader io.Reader, err error) {
	reader, writer := io.Pipe()
	go func() {
		for {
			b := make([]byte, 1)
			_, err := input.Read(b)
			if err != nil {
				if err == io.EOF {
					writer.Close()
					return
				}
				writer.CloseWithError(err)
				return
			}
			m, err := encoder(b[0], key)
			if err != nil {
				writer.CloseWithError(err)
			}
			_, err = writer.Write(m)
			if err != nil {
				writer.CloseWithError(err)
			}
		}
	}()
	return reader, nil
}

func TranscodeStreamPartial(
	input io.Reader,
	key Key,
	take int,
	skip int,
	encoder func(byte, Key) ([]byte, error),
) (reader io.Reader, err error) {
	reader, writer := io.Pipe()
	go func() {
		taken := 0
		skipped := 0
		taking := true
		for {
			b := make([]byte, 1)
			_, err := input.Read(b)
			if err != nil {
				if err == io.EOF {
					writer.Close()
					return
				}
				writer.CloseWithError(err)
				return
			}
			if taking {
				m, err := encoder(b[0], key)
				if err != nil {
					writer.CloseWithError(err)
				}
				_, err = writer.Write(m)
				if err != nil {
					writer.CloseWithError(err)
				}
				taken++
				if taken >= take {
					taken = 0
					taking = false
				}
			} else {
				_, err = writer.Write(b)
				if err != nil {
					writer.CloseWithError(err)
				}
				skipped++
				if skipped >= skip {
					skipped = 0
					taking = true
				}
			}
		}
	}()
	return reader, nil
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
	errorCh := make(chan error, len(input))
	wg := &sync.WaitGroup{}
	wg.Add(len(input))

	for i := range input {
		go bytesRelay(
			i, input, byteGroups, key, encoder, errorCh, wg)
	}
	wg.Wait()

	select {
	case err := <-errorCh:
		return nil, err
	default:
		break
	}

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
	errorCh chan error,
	wg *sync.WaitGroup) {
	defer wg.Done()
	byteGroup := ByteGroup{
		bytes: make([]byte, 0),
	}
	msg, err := encoder(input[index], key)
	if err != nil {
		errorCh <- err
		return
	}
	for _, b := range msg {
		byteGroup.bytes = append(byteGroup.bytes, b)
	}
	bytes[index] = byteGroup
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
	decodeGroups, err := findDecodeGroups(
		input, key.GetDictionaryChars(), groups)
	if err != nil {
		return nil, err
	}
	decoded, err := decodeBytes(key, decodeGroups, decodeFunc)
	return decoded, err
}

func TransdecodeStream(
	input io.Reader,
	key Key,
	groups int,
	decodeFunc func(Key, DecodeGroup) (byte, error),
) (output io.Reader, err error) {
	chars := key.GetDictionaryChars()
	groupsFound := 0
	buffer := make([]byte, 0)
	reader, writer := io.Pipe()
	go func() {
		for {
			b, err := readTranscodedStream(
				input, writer, buffer, key, groups, decodeFunc)
			if err != nil {
				writer.CloseWithError(err)
			}
			if chars.CheckIn(b) {
				if groupsFound == groups {
					err = writeDecodeBuffer(
						decodeFunc, buffer, chars, groups, key, writer)
					if err != nil {
						writer.CloseWithError(err)
						return
					}
					buffer = make([]byte, 0)
					groupsFound = 0
				}
				groupsFound++
			}
			buffer = append(buffer, b)
		}
	}()

	return reader, nil
}

func readTranscodedStream(
	input io.Reader,
	writer *io.PipeWriter,
	buffer []byte,
	key Key,
	groups int,
	decodeFunc func(Key, DecodeGroup) (byte, error),
) (byte, error) {
	chars := key.GetDictionaryChars()
	b := make([]byte, 1)
	_, err := input.Read(b)
	if err != nil {
		if err == io.EOF {
			err = writeDecodeBuffer(
				decodeFunc, buffer, chars, groups, key, writer)
			if err != nil {
				return b[0], err
			}
			writer.Close()
		} else {
			return b[0], err
		}
	}
	return b[0], nil
}

func writeDecodeBuffer(
	decodeFunc func(Key, DecodeGroup) (byte, error),
	buffer []byte,
	chars DictionaryChars,
	groups int,
	key Key,
	writer io.Writer,
) error {
	decodeGroups, err := findDecodeGroups(buffer, chars, groups)
	if err != nil {
		return err
	}
	decoded, err := decodeBytes(key, decodeGroups, decodeFunc)
	if err != nil {
		return err
	}
	_, err = writer.Write(decoded)
	if err != nil {
		return err
	}
	return nil
}

func findDecodeGroups(
	input []byte,
	characters DictionaryChars,
	numGroups int,
) (decodeGroups []DecodeGroup, err error) {
	if !characters.CheckIn(input[0]) {
		return decodeGroups, errors.New(
			"no decode characters found")
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
