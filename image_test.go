package main

import (
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
		"r54,238r3,243r842,140r51,338b823,470b228,193r314,478r114,111",
		), image)
	if err != nil {
		t.Error(err)
	}
	t.Log("Message: ", string(message))
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
	t.Log("Message: ", string(message))
}