package fileloader

import (
	"bytes"
	"io"

	"github.com/semaphoreci/test-results/pkg/logger"
)

var readers map[string]*bytes.Reader = make(map[string]*bytes.Reader)

// Load reader from internal buffer or create new one
func Load(path string, reader *bytes.Reader) (*bytes.Reader, bool) {
	return decode(path, reader)
}

func decode(path string, reader *bytes.Reader) (*bytes.Reader, bool) {
	foundReader, exists := readers[path]
	if exists && foundReader != nil {
		logger.Debug("Path read from cache")
		foundReader.Seek(0, io.SeekStart)
		return foundReader, true
	}
	readers[path] = reader
	logger.Debug("No path in cache")
	return reader, false
}
