package writer

import (
	"encoding/json"
	"os"

	"github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

// EntryFileWriter missing godoc.
type EntryFileWriter struct {
	file    *os.File
	encoder *json.Encoder
}

// NewEntryFileWriter missing godoc.
func NewEntryFileWriter(filename string) (*EntryFileWriter, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &EntryFileWriter{file: file, encoder: json.NewEncoder(file)}, nil
}

// NewEntryFileEmptyWriter missing godoc.
func NewEntryFileEmptyWriter(filename string) (*EntryFileWriter, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}

	return &EntryFileWriter{file: file, encoder: json.NewEncoder(file)}, nil
}

// WriteEntry missing godoc.
func (w *EntryFileWriter) WriteEntry(entry entity.Entry) error {
	err := w.encoder.Encode(entry)
	if err != nil {
		return err
	}
	return nil
}

// Close missing godoc.
func (w *EntryFileWriter) Close() error {
	return w.file.Close()
}
