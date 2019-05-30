package decouplet

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestEncodeBytes(t *testing.T) {
	newMessage, err := EncodeBytes([]byte("Test"), []byte("tEst Key3#$T234Alklgn"))
	if err != nil {
		t.Error(err)
	}
	t.Log(string(newMessage))
}

func TestDecoderBytes(t *testing.T) {
	message, err := DecodeBytes(
		[]byte("[dcplt-byteec-0.2]e12i16d17k11k17g11k12e4"), []byte("tEst Key3#$"))
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
	newMessage, err := EncodeBytes(fileBytes, key)
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
		t.Fatal(err)
	}
	newReader, err := DecodeBytesStream(reader, key)
	if err != nil {
		t.Fatal(err)
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
		t.FailNow()
	}
	newReader, err := DecodeBytesStreamPartial(reader, key)
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

func TestAnalyzeBytesKey(t *testing.T) {
	badKey := []byte("badkey")
	scale := AnalyzeBytesKey(badKey)
	t.Log("bad key analysis:", scale)
	if scale > 10 {
		t.Log("small, insufficient keys usually register under 10")
		t.Fail()
	}
	goodKey := []byte("This is a Key$%@#$@^^%$&$%%^*{})([p[]Should _-!`~")
	scale = AnalyzeBytesKey(goodKey)
	t.Log("good key analysis:", scale)
	if scale < 10 {
		t.Log("good keys should be 10 or over")
		t.Fail()
	}
	greatKey := []byte(
		"GREAFgolanVMb elefwoejgitoiqwaz12353445789870-0=)" +
			"(_#@$^#$&$%&$*$&$0238959_=2340+=12!@#$%^&*(()")
	scale = AnalyzeBytesKey(greatKey)
	if scale < 20 {
		t.Log("great keys should be 20 or over(not really a hard number)")
		t.Fail()
	}
	t.Log("great key analysis:", scale)
}

// Examples

func ExampleEncodeBytes() {
	message, err := EncodeBytes([]byte("Test"), []byte("tEst Key3#$T234"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("encoded:", string(message))
}

func ExampleEncodeBytesStream() {
	key := []byte("tEst Key3#$!@&*()[]:;")
	msg := []byte("Test this message and see it stream")
	input := bytes.NewReader(msg)
	reader, err := EncodeBytesStream(input, key)
	if err != nil {
		fmt.Println(err)
	}
	message, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("output:", string(message))
}

func ExampleEncodeBytesStreamPartial() {
	key := []byte("tEst Key3#$!@&*()[]:;")
	msg := []byte("Test this message and see it stream")
	take := 4
	skip := 10
	input := bytes.NewReader(msg)
	reader, err := EncodeBytesStreamPartial(input, key, take, skip)
	if err != nil {
		fmt.Println(err)
	}
	message, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("output:", string(message))
}

func ExampleDecodeBytes() {
	message, err := DecodeBytes(
		[]byte("[dcplt-byteec-0.2]a9c0e8j4j8d4j8c9"), []byte("tEst Key3#$"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("output:", string(message))
}

func ExampleDecodeBytesStream() {
	key := []byte("tEst Key3#$!@&*()[]:;")
	msg := []byte(
		"c18j2a4j1a11i0a1j7a6j1k2c6a18k1b14i0g9k4a12j4e11" +
			"i0d5j0k18h12a13i7a9h0k17c12a4j1b17f0k1e4k17d1k6d1" +
			"c8h1k18h12g1j0a4j1b18g0a5j0a1j7e4b1f14j6j3f4k2i1k18c12d13j5j0b4")
	input := bytes.NewReader(msg)
	reader, err := DecodeBytesStream(input, key)
	if err != nil {
		fmt.Println(err)
	}
	message, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("output:", string(message))

}

func ExampleDecodeBytesStreamPartial() {
	key := []byte("tEst Key3#$!@&*()[]:;")
	msg := []byte(
		";[&e19i0d5j0d16k4f10k4&]; this mess;[&a10j1j7a4i3k7c9a1&];" +
			"and see it;[&c9a1c10i7a4i0k17f1&];eam")
	input := bytes.NewReader(msg)
	reader, err := DecodeBytesStreamPartial(input, key)
	if err != nil {
		fmt.Println(err)
	}
	message, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("output:", string(message))
}

func TestImageBytePPM(t *testing.T) {
	file, err := os.Open("images/body.bin")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	key := make([]byte, 256)
	_, err = rand.Read(key)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	key2 := make([]byte, 256)
	_, err = rand.Read(key2)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	reader, err := EncodeBytesStream(file, key)
	if err != nil {
		fmt.Println(err)
	}
	decoded, err := DecodeBytesStream(reader, key2)
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

func TestImageImagePPM(t *testing.T) {
	file, err := os.Open("images/body.bin")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	key, err := LoadImage("images/gopher.png")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	reader, err := EncodeImageStream(file, key)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	key2, err := LoadImage("images/test.png")
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
