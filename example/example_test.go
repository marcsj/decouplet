package example

import (
	"bytes"
	"crypto/rand"
	"github.com/marcsj/decouplet"
	"io"
	"log"
	"testing"
)

func Test_AESExample(t *testing.T) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		t.Error(err)
	}
	unencrypted := []byte("This is an example message to encode.")
	encryptedBytes, nonce, err := aesEncrypt(unencrypted, key)
	if err != nil {
		t.Error(err)
	}
	log.Println("Encrypted message:", string(encryptedBytes))
	transcodeKey := []byte("This is an example key!@#$%^&*()1234567890")
	output, err := decouplet.TranscodeBytes(encryptedBytes, transcodeKey)
	if err != nil {
		t.Error(err)
	}
	log.Println("Transcoded message:", string(output))
	transdecoded, err := decouplet.TransdecodeBytes(output, transcodeKey)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(transdecoded, encryptedBytes) {
		t.Log("Transdecoded bytes do not equal encrypted bytes")
		t.Fail()
	}
	decrypted, err := aesDecrypt(transdecoded, key, nonce)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(decrypted, unencrypted) {
		t.Log("Decrypted text does not equal original text")
		t.Fail()
	}
	log.Println("Decrypted text:", string(decrypted))
}
