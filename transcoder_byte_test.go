package decouplet

import (
	"bytes"
	"os"
	"testing"
)

func TestTranscodeBytes(t *testing.T) {
	newMessage, err := TranscodeBytes([]byte("Test"), []byte("tEst Key3#$"))
	if err != nil {
		t.Error(err)
	}
	t.Log(string(newMessage))
}

func TestTransdecodeBytes(t *testing.T) {
	message, err := TransdecodeBytes(
		[]byte("[dcplt-bytetc-0.1]a9c0e8j4j8d4j8c9"), []byte("tEst Key3#$"))
	if err != nil {
		t.Error(err)
	}
	t.Log(string(message))
}

func TestByteMessage(t *testing.T) {
	originalMessage :=
		"!!**_-+Test THIS bigger message with More Symbols" +
			"@$_()#$%^#@!~#2364###$%! *(#$%)^@#%$@"
	newMessage, err := TranscodeBytes(
		[]byte(originalMessage), []byte("Test Key!@# $"))
	if err != nil {
		t.Error(err)
	}
	t.Log(string(newMessage))
	message, err := TransdecodeBytes(newMessage, []byte("Test Key!@# $"))
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
			"^@#%@#)^Test byte Key!@#$" +
			"^GEWg gwefwgwef _#$%@#$%L",
	)
	t.Log("Length of original:", len(fileBytes))
	newMessage, err := TranscodeBytesConcurrent(fileBytes, key)
	if err != nil {
		t.Error(err)
	}
	message, err := TransdecodeBytes(newMessage, key)
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
