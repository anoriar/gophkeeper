package reader

import (
	"encoding/json"
	"os"

	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

// EntryFileReader missing godoc.
type EntryFileReader struct {
	file    *os.File
	decoder *json.Decoder
}

// NewEntryFileReader missing godoc.
func NewEntryFileReader(filename string) (*EntryFileReader, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &EntryFileReader{file: file, decoder: json.NewDecoder(file)}, nil
}

// ReadEntry missing godoc.
func (c *EntryFileReader) ReadEntry() (*entity.Entry, error) {
	entry := &entity.Entry{}
	err := c.decoder.Decode(entry)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// Close missing godoc.
func (c *EntryFileReader) Close() error {
	return c.file.Close()
}
