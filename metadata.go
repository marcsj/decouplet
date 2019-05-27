package decouplet

import (
	"bytes"
	"errors"
	"fmt"
)

type EncoderInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var EncodersList []EncoderInfo

func (i EncoderInfo) GetEncoderString() (string, error) {
	return fmt.Sprintf(
		"[dcplt-%s-%s]",
		i.Name,
		i.Version,
	), nil
}

func (i EncoderInfo) CheckEncoder(message *[]byte) error {
	meta, err := i.GetEncoderString()
	if err != nil {
		return err
	}
	metaBytes := []byte(meta)
	if bytes.HasPrefix(*message, metaBytes) {
		*message = (*message)[len(metaBytes):]
		return nil
	}
	return errors.New("encoder version does not match")
}

func (i EncoderInfo) WriteVersion() ([]byte, error) {
	meta, err := i.GetEncoderString()
	if err != nil {
		return nil, err
	}
	return []byte(meta), nil
}
