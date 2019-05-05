package decouplet

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestEncodeBytes(t *testing.T) {
	newMessage, err := EncodeBytes([]byte("Test"), []byte("tEst Key3#$T234"))
	if err != nil {
		t.Error(err)
	}
	t.Log(string(newMessage))
}

func TestDecoderBytes(t *testing.T) {
	message, err := DecodeBytes(
		[]byte("[dcplt-byteec-0.2]a9c0e8j4j8d4j8c9"), []byte("tEst Key3#$"))
	if err != nil {
		t.Error(err)
	}
	t.Log(string(message))
}

func TestByteMessage(t *testing.T) {
	originalMessage :=
		"!!**_-+Test THIS bigger message with More Symbols" +
			"@$_()#$%^#@!~#2364###$%! *(#$%)^@#%$@"
	newMessage, err := EncodeBytes(
		[]byte(originalMessage), []byte("Test encodingKey!@# $"))
	if err != nil {
		t.Error(err)
	}
	t.Log(string(newMessage))
	message, err := DecodeBytes(newMessage, []byte("Test encodingKey!@# $"))
	if err != nil {
		t.Error(err)
	}
	if originalMessage != string(message) {
		t.Fail()
	}
	t.Log("Message:", string(message))
}

func TestByteMessage_Byte(t *testing.T) {
	imageFile, err := os.Open("images/test.jpg")
	if err != nil {
		t.Error(err)
	}
	defer imageFile.Close()
	fileInfo, err := imageFile.Stat()
	if err != nil {
		t.Error(err)
	}
	fileBytes := make([]byte, fileInfo.Size())
	_, err = imageFile.Read(fileBytes)
	if err != nil {
		t.Error(err)
	}
	key := []byte(
		"$#%#%@#$@$)*^_#@$*^)@$)@#" +
			"^@#%@#)^Test byte encodingKey!@#$" +
			"^GEWg gwefwgwef _#$%@#$%L",
	)
	t.Log("Length of original:", len(fileBytes))
	newMessage, err := EncodeBytesConcurrent(fileBytes, key)
	if err != nil {
		t.Error(err)
	}
	message, err := DecodeBytes(newMessage, key)
	if err != nil {
		t.Error(err)
	}
	t.Log("Length of finished:", len(message))
	if len(message) != len(fileBytes) {
		t.Log("sizes are not the same:",
			len(message), len(fileBytes))
		t.Fail()
	}
	if !bytes.Equal(fileBytes, message) {
		t.Log("bytes are not equal")
		t.Fail()
	}
}

func TestEncodeBytesConcurrent(t *testing.T) {
	key := []byte("tEst Key3#$!@&*()[]:;")
	msg := []byte("Test this message and see it stream")
	input := bytes.NewReader(msg)
	reader, err := EncodeBytesStream(input, key)
	if err != nil {
		t.Error(err)
	}
	newReader, err := DecodeBytesStream(reader, key)
	if err != nil {
		t.Error(err)
	}
	b, err := ioutil.ReadAll(newReader)
	t.Log(string(b))
	if !bytes.Equal(msg, b) {
		t.Log("bytes are not equal")
		t.Fail()
	}
}

func TestEncodeBytesConcurrentPartial(t *testing.T) {
	key := []byte("tEst Key3#$!@&*()[]:;")
	msg := []byte("Test this message and see it stream and be partially encoded! here")
	take := 1
	skip := 3
	input := bytes.NewReader(msg)
	reader, err := EncodeBytesStreamPartial(input, key, take, skip)
	if err != nil {
		t.Error(err)
	}
	newReader, err := DecodeByteStreamPartial(reader, key)
	if err != nil {
		t.Error(err)
	}
	b, err := ioutil.ReadAll(newReader)
	t.Log(string(b))
	if !bytes.Equal(msg, b) {
		t.Log("bytes are not equal")
		t.Fail()
	}
}
