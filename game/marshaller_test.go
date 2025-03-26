package game

import (
	"errors"
	"log"
	"os"
	"testing"
)

type Success struct {
}

func (d *Success) Marshal() ([]byte, error) {
	return []byte("SUCCESS"), nil
}

func (d *Success) Unmarshal(data []byte) error {
	if string(data) == "SUCCESS" {
		return nil
	}
	return errors.New("invalid data")
}

func TestMarshaller_FileMarshal(t *testing.T) {
	var success Success
	defer func() {
		_ = os.Remove(".success.txt")
	}()
	err := FileMarshal(".success.txt", &success)
	if err != nil {
		log.Fatal(err)
	}
}

func TestMarshaller_FileUnmarshal(t *testing.T) {
	var success Success
	defer func() {
		_ = os.Remove(".success.txt")
	}()
	err := FileMarshal(".success.txt", &success)
	if err != nil {
		log.Fatal(err)
	}
	err = FileUnmarshal(".success.txt", &success)
	if err != nil {
		log.Fatal(err)
	}
}
