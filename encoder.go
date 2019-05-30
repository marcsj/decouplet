package decouplet

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
)

type encodingKey interface {
	GetVersion() EncoderInfo
	CheckValid() (bool, error)
	GetDictionaryChars() dictionaryChars
	GetDictionary() dictionary
}

func encode(
	input []byte,
	key encodingKey,
	encoder func(byte, encodingKey) ([]byte, error),
) ([]byte, error) {
	if valid, err := key.CheckValid(); !valid {
		return nil, err
	}

	b, err := key.GetVersion().WriteVersion()
	if err != nil {
		return nil, err
	}
	reader, err := encodeStream(bytes.NewReader(input), key, encoder)
	if err != nil {
		return nil, err
	}

	output, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	err = reader.Close()
	if err != nil {
		return nil, err
	}
	return append(b, output...), nil
}

func encodeStream(
	input io.Reader,
	key encodingKey,
	encoder func(byte, encodingKey) ([]byte, error),
) (*io.PipeReader, error) {
	if valid, err := key.CheckValid(); !valid {
		return nil, err
	}
	reader, writer := io.Pipe()
	go func(
		input io.Reader,
		writer *io.PipeWriter,
		encoder func(byte, encodingKey) ([]byte, error),
		key encodingKey) {

		scanner := bufio.NewScanner(input)
		scanner.Split(bufio.ScanBytes)

		for scanner.Scan() {
			m, err := encoder(scanner.Bytes()[0], key)
			if err != nil {
				writer.CloseWithError(err)
				return
			}
			_, err = writer.Write(m)
			if err != nil {
				writer.CloseWithError(err)
				return
			}
		}

		if err := scanner.Err(); err != nil {
			writer.CloseWithError(err)
			return
		}
		writer.Close()
	}(input, writer, encoder, key)

	return reader, nil
}

func encodePartialStream(
	input io.Reader,
	key encodingKey,
	take int,
	skip int,
	encoder func(byte, encodingKey) ([]byte, error),
) (*io.PipeReader, error) {
	reader, writer := io.Pipe()
	if valid, err := key.CheckValid(); !valid {
		return nil, err
	}

	go func() {
		defer writer.Close()
		for {
			_, err := writer.Write(partialStartBytes)
			if err != nil {
				writer.CloseWithError(err)
				return
			}
			takeR := io.LimitReader(input, int64(take))
			encodedR, err := encodeStream(takeR, key, encoder)
			if err != nil {
				writer.CloseWithError(err)
				return
			}
			_, err = io.Copy(writer, encodedR)
			if err != nil {
				writer.CloseWithError(err)
				return
			}
			_, err = writer.Write(partialEndBytes)
			if err != nil {
				writer.CloseWithError(err)
				return
			}
			_, err = io.CopyN(writer, input, int64(skip))
			if err != nil {
				writer.CloseWithError(err)
				return
			}
		}
	}()

	return reader, nil
}
