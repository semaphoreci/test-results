package fileloader

import (
	"bytes"

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
		logger.Log("fileloader", "FileLoader: Path %s read from cache", path)
		return foundReader, true
	}
	readers[path] = reader
	logger.Debug("fileloader", "FileLoader: No path %s in cache", path)
	return reader, false
}
