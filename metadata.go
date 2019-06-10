package decouplet

import (
	"bytes"
	"errors"
	"fmt"
)

type encoderInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (i encoderInfo) getEncoderString() (string, error) {
	return fmt.Sprintf(
		"[dcplt-%s-%s]",
		i.Name,
		i.Version,
	), nil
}

func (i encoderInfo) checkEncoder(message *[]byte) error {
	meta, err := i.getEncoderString()
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

func (i encoderInfo) writeVersion() ([]byte, error) {
	meta, err := i.getEncoderString()
	if err != nil {
		return nil, err
	}
	return []byte(meta), nil
}
