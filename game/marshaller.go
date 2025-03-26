package game

import (
	"io"
	"log"
	"os"
)

type Marshal interface {
	Marshal() ([]byte, error)
}

type Unmarshal interface {
	Unmarshal([]byte) error
}

func FileMarshal(file string, marshal Marshal) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)
	marshaled, err := marshal.Marshal()
	if err != nil {
		return err
	}
	if _, err := f.Write(marshaled); err != nil {
		return err
	}
	return nil
}

func FileUnmarshal(file string, unmarshal Unmarshal) error {
	f, err := os.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)
	bytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	return unmarshal.Unmarshal(bytes)
}
