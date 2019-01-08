package main

import (
	"os"
	"testing"
)

func TestTranscodeImage(t *testing.T) {
	image, err := LoadImage("images/test.png")
	if err != nil {
		t.Error(err)
	}
	newMessage, err := TranscodeImage([]byte("Test"), image)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(newMessage))
}

func TestTransdecodeImage(t *testing.T) {
	image, err := LoadImage("images/test.png")
	if err != nil {
		t.Error(err)
	}
	message, err := TransdecodeImage([]byte(
		"[dcplt-imgtc-0.1]g196255r90241k554943k551042c551723k138337r565877k138337",
		), image)
	if err != nil {
		t.Error(err)
	}
	t.Log("Message:", string(message))
}

func TestImageMessage(t *testing.T) {
	image, err := LoadImage("images/test.jpg")
	if err != nil {
		t.Error(err)
	}
	originalMessage :=
			"!!**_-+Test THIS bigger message with More Symbols!" +
			"@$_()#$%^#@!~#2364###$%! *(#$%)^@#%$@More and more"
	newMessage, err := TranscodeImage([]byte(originalMessage), image)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(newMessage))
	message, err := TransdecodeImage(newMessage, image)
	if err != nil {
		t.Error(err)
	}
	if originalMessage != string(message) {
		t.Fail()
	}
	t.Log("Message:", string(message))
}

func TestImageMessage_Image(t *testing.T) {
	imageFile, err := os.Open("images/test.jpg")
	if err != nil {
		t.Error(err)
	}
	defer imageFile.Close()
	fileInfo, err := imageFile.Stat()
	if err != nil {
		t.Error(err)
	}
	fileBytes := make([]byte, fileInfo.Size()) //fix this crap
	_, err = imageFile.Read(fileBytes)
	if err != nil {
		t.Error(err)
	}
	image, err := LoadImage("images/test.jpg")
	if err != nil {
		t.Error(err)
	}
	newMessage, err := TranscodeImage(fileBytes, image)
	if err != nil {
		t.Error(err)
	}
	t.Log("Length of message:", len(newMessage))
	message, err := TransdecodeImage(newMessage, image)
	if err != nil {
		t.Error(err)
	}
	if string(fileBytes) != string(message) {
		t.Fail()
	}
}
