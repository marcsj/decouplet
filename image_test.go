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
		"[dcplt-imgtc-0.2]a182145r90241r590295k137282r6777k139200r460987c138337",
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
	fileBytes := make([]byte, fileInfo.Size())
	_, err = imageFile.Read(fileBytes)
	if err != nil {
		t.Error(err)
	}
	t.Log("Length of original: ", len(fileBytes))
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
	if len(message) != len(fileBytes) {
		t.Fail()
	}
	if string(fileBytes) != string(message) {
		t.Log("fileBytes:", fileBytes)
		t.Log("message:", message)
		t.Fail()
	}
}