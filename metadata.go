package decouplet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
)

type EncoderInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var EncodersList []EncoderInfo

func init() {
	EncodersList = make([]EncoderInfo, 0)

	_, runFile, _, _ := runtime.Caller(0)
	file, err := os.Open(path.Dir(runFile) + "/versions.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&EncodersList)
	if err != nil {
		panic(err)
	}
}

func GetEncoderMeta(name string) (string, error) {
	for i := range EncodersList {
		if EncodersList[i].Name == name {
			return fmt.Sprintf(
				"[dcplt-%s-%s]",
				EncodersList[i].Name,
				EncodersList[i].Version,
			), nil
		}
	}
	return "", errors.New("invalid Encoder metadata")
}

func CheckEncoder(
	EncoderType encoderType,
	message *[]byte) error {
	meta, err := GetEncoderMeta(string(EncoderType))
	if err != nil {
		return err
	}
	metaSize := len([]byte(meta))
	if bytes.Equal((*message)[:metaSize], []byte(meta)) {
		*message = (*message)[metaSize:]
		return nil
	}
	return errors.New("Encoder version does not match")
}

func WriteVersion(EncoderType encoderType) ([]byte, error) {
	meta, err := GetEncoderMeta(string(EncoderType))
	if err != nil {
		return nil, err
	}
	return []byte(meta), nil
}
