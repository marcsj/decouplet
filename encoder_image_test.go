package decouplet

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestEncodeImage(t *testing.T) {
	image, err := LoadImage("images/test.png")
	if err != nil {
		t.Error(err)
	}
	newMessage, err := EncodeImage([]byte("Test"), image)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(newMessage))
}

func TestDecoderImage(t *testing.T) {
	image, err := LoadImage("images/test.png")
	if err != nil {
		t.Error(err)
	}
	message, err := DecodeImage([]byte(
		"[dcplt-imgec-0.2]a182145r90241r590295k137282r6777k139200r460987c138337",
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
		"!!**_-+Test THIS bigger message with More Symbols" +
			"@$_()#$%^#@!~#2364###$%! *(#$%)^@#%$@"
	newMessage, err := EncodeImage(
		[]byte(originalMessage), image)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(newMessage))
	message, err := DecodeImage(newMessage, image)
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
	newMessage, err := EncodeImage(fileBytes, image)
	if err != nil {
		t.Error(err)
	}
	message, err := DecodeImage(newMessage, image)
	if err != nil {
		t.Error(err)
	}
	t.Log("Length of finished: ", len(message))
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

func TestEncodeImageConcurrent(t *testing.T) {
	image, err := LoadImage("images/test.jpg")
	if err != nil {
		t.Error(err)
	}
	msg := []byte("Test this message and see it stream")
	input := bytes.NewReader(msg)
	reader, err := EncodeImageStream(input, image)
	if err != nil {
		t.Error(err)
	}
	newReader, err := DecodeImageStream(reader, image)
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

func TestEncodeImageConcurrentPartial(t *testing.T) {
	image, err := LoadImage("images/test.jpg")
	if err != nil {
		t.Error(err)
	}
	take := 1
	skip := 3
	msg := []byte("Test this message and see it stream, using partial encoding.")
	input := bytes.NewReader(msg)
	reader, err := EncodeImageStreamPartial(input, image, take, skip)
	if err != nil {
		fmt.Println(err)
	}
	newReader, err := DecodeImageStreamPartial(reader, image)
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

func TestImageImagePPM(t *testing.T) {
	file, err := os.Open("images/body.bin")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	key, err := LoadImage("images/test.png")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	reader, err := EncodeImageStream(file, key)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	key2, err := LoadImage("images/image.jpg")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	decoded, err := DecodeImageStream(reader, key2)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	data, err := ioutil.ReadAll(decoded)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = ioutil.WriteFile("images/body.ecb.bin", data, 0644)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

}
