package fileloader

import (
	"bytes"
	"io"
)

var readers map[string]bytes.Reader = make(map[string]bytes.Reader)

// Load reader from internal buffer or create new one
func Load(path string, reader *bytes.Reader) (bytes.Reader, bool) {
	return decode(path, reader)
}

func decode(path string, reader *bytes.Reader) (bytes.Reader, bool) {
	foundReader, exists := readers[path]
	if exists {
		return foundReader, true
	}
	readers[path] = *reader

	(*reader).Seek(0, io.SeekStart)

	return *reader, false
}
