package main

import (
	"testing"
)

func TestTranscodeBytes(t *testing.T) {
	newMessage, err := TranscodeBytes([]byte("This is a test message! Very secret."), []byte("Test this! $%@#&*"))
	if err != nil {
		t.Error(err)
	}
	t.Log(string(newMessage))
}

func TestTransdecodeBytes(t *testing.T) {
	message, err := TransdecodeBytes([]byte("[dcplt-bytetc-0.1]a11d2d12h2b7" +
		"j0j3b1a2i0j2e0a4i0h4i4h6j7h4i4g1j7h14i6j1c0h13j0g15b0c5j1h10i1i9j0" +
		"j3b1c4h1j2d0f11h2f6h2f7h2a11e2j2c0a13i2e11i1h10i4e16i1j2c0f15h2c8" +
		"j7b9h6h0j6c14a0"), []byte("Test this! $%@#&*"))
	if err != nil {
		t.Error(err)
	}
	t.Log(string(message))
}