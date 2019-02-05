package example

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func aesEncrypt(message []byte, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	encoded := aesgcm.Seal(nil, nonce, message, nil)
	return encoded, nonce, nil
}

func aesDecrypt(input []byte, key []byte, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plaintext, err := aesgcm.Open(nil, nonce, input, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
