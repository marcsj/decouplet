package main

import (
	"testing"
)

func TestTranscodeBytes(t *testing.T) {
	newMessage, err := TranscodeBytes([]byte("Test"), []byte("Test this! $%@#&*"))
	if err != nil {
		t.Error(err)
	}
	t.Log(string(newMessage))
}