package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type TranscoderInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var TranscodersList []TranscoderInfo

func init() {
	TranscodersList = make([]TranscoderInfo, 0)
	file, err := os.Open("versions.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&TranscodersList)
	if err != nil {
		panic(err)
	}
}

func GetTranscoderMeta(name string) (string, error) {
	for i := range TranscodersList {
		if TranscodersList[i].Name == name {
			return fmt.Sprintf(
				"[dcplt-%s-%s]",
				TranscodersList[i].Name,
				TranscodersList[i].Version,
			), nil
		}
	}
	return "", errors.New("invalid transcoder metadata")
}

func CheckTranscoder(
	transcoderType TranscoderType,
	message *[]byte) error {
	meta, err := GetTranscoderMeta(string(transcoderType))
	if err != nil {
		return err
	}
	metaSize := len([]byte(meta))
	if bytes.Equal((*message)[:metaSize], []byte(meta)) {
		*message = (*message)[metaSize:]
		return nil
	}
	return errors.New("transcoder version does not match")
}

func WriteVersion(transcoderType TranscoderType) ([]byte, error) {
	meta, err := GetTranscoderMeta(string(transcoderType))
	if err != nil {
		return nil, err
	}
	return []byte(meta), nil
}
