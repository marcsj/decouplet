package decouplet

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

func decode(
	input []byte,
	key encodingKey,
	groups int,
	decodeFunc func(encodingKey, decodeGroup) (byte, error),
) (output []byte, err error) {
	err = CheckEncoder(key.GetType(), &input)
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

func decodeStream(
	input io.Reader,
	key encodingKey,
	groups int,
	decodeFunc func(encodingKey, decodeGroup) (byte, error),
) (output *io.PipeReader, err error) {
	reader, writer := io.Pipe()

	go func() {
		chars := key.GetDictionaryChars()
		defer writer.Close()

		charSplit := splitInfo{chars: chars, groups: groups}
		scanner := bufio.NewScanner(input)
		scanner.Split(charSplit.scanDecodeSplit)

		for scanner.Scan() {
			err := writeDecodeBuffer(decodeFunc, scanner.Bytes(), groups, key, writer)
			if err != nil {
				writer.CloseWithError(err)
			}
		}

		if err := scanner.Err(); err != nil {
			writer.CloseWithError(err)
		}
	}()

	return reader, nil
}

func (t splitInfo) scanDecodeSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	charsFound := 0
	for i := range data {
		if t.chars.checkIn(data[i]) {
			charsFound++
			if charsFound == t.groups+1 {
				return i, data[:i], nil
			}
		}
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func scanDecodeStream(
	input io.Reader,
	key encodingKey,
	groups int,
	decodeFunc func(encodingKey, decodeGroup) (byte, error),
	writer *io.PipeWriter,
) {
	defer writer.Close()

	scanner := bufio.NewScanner(input)
	scanner.Split(scanPartialSplit)

	for scanner.Scan() {
		splitBytes := bytes.SplitAfter(scanner.Bytes(), partialStartBytes)
		var skipped []byte
		if len(splitBytes[0]) < len(partialStartBytes) {
			skipped = splitBytes[0][0:len(splitBytes[0])]
		} else {
			skipped = splitBytes[0][0 : len(splitBytes[0])-len(partialStartBytes)]
		}
		_, err := writer.Write(skipped)
		if err != nil {
			writer.CloseWithError(err)
		}
		if len(splitBytes) > 1 {
			if len(splitBytes[1]) > 0 {
				readBytes := splitBytes[1]
				err := writeDecodeBuffer(decodeFunc, readBytes, groups, key, writer)
				if err != nil {
					if err != io.EOF {
						writer.CloseWithError(err)
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		writer.CloseWithError(err)
	}
}

func scanPartialSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, partialEndBytes); i >= 0 {
		return i + len(partialStartBytes), data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func decodePartialStream(
	input io.Reader,
	key encodingKey,
	groups int,
	decodeFunc func(encodingKey, decodeGroup) (byte, error),
) (output *io.PipeReader, err error) {
	reader, writer := io.Pipe()

	go scanDecodeStream(input, key, groups, decodeFunc, writer)

	return reader, nil
}

func writeDecodeBuffer(
	decodeFunc func(encodingKey, decodeGroup) (byte, error),
	buffer []byte,
	groups int,
	key encodingKey,
	writer io.Writer,
) error {
	decodeGroups, err := findDecodeGroups(buffer, key.GetDictionaryChars(), groups)
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
	characters dictionaryChars,
	numGroups int,
) (decodeGroups []decodeGroup, err error) {
	if !characters.checkIn(input[0]) {
		return decodeGroups, errors.New(
			"no decode characters found")
	}
	decode := decodeGroup{
		kind:  []uint8{},
		place: []string{},
	}
	buffer := make([]uint8, 0)
	numberAdded := 0

	for i := range input {
		if characters.checkIn(input[i]) {
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
	key encodingKey,
	decodeGroups []decodeGroup,
	decodeFunc func(encodingKey, decodeGroup) (byte, error),
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
