package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

type TranscoderInfo struct {
	Name string 	`json:"name"`
	Version string 	`json:"version"`
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

func CheckTranscoder(name string, message string) (string, error) {
	meta, err := GetTranscoderMeta(name)
	if err != nil {
		return message, err
	}
	if strings.HasPrefix(message, meta) {
		return strings.TrimPrefix(message, meta), nil
	}
	return message, errors.New("transcoder version does not match")
}

func WriteVersion(name string, bytes []byte) ([]byte, error) {
	meta, err := GetTranscoderMeta(name)
	if err != nil {
		return bytes, err
	}
	for i := range meta {
		bytes = append(bytes, byte(meta[i]))
	}
	return bytes, nil
}