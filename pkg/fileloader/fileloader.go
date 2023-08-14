package fileloader

import (
	"bytes"
	"io"
	"os"

	"github.com/semaphoreci/test-results/pkg/logger"
)

var readers map[string]*bytes.Reader = make(map[string]*bytes.Reader)

// Load reader from internal buffer if path was already loaded or create new one if not
func Load(path string, reader *bytes.Reader) (*bytes.Reader, bool) {
	return decode(path, reader)
}

// Ensure puts reader data into temporary created file.
func Ensure(reader *bytes.Reader) (fileName string) {
	file, err := os.CreateTemp("", "")
	defer file.Close() // #nosec
	if err != nil {
		panic(err)
	}

	fileName = file.Name()

	_, err = reader.WriteTo(file)
	if err != nil {
		panic(err)
	}

	return
}

func decode(path string, reader *bytes.Reader) (*bytes.Reader, bool) {
	foundReader, exists := readers[path]
	if exists && foundReader != nil && foundReader.Size() == reader.Size() {
		logger.Debug("Path read from cache")
		_, err := foundReader.Seek(0, io.SeekStart)
		if err != nil {
			logger.Error("Cannot seek to start of reader: %v", err)
		}

		return foundReader, true
	}
	readers[path] = reader
	logger.Debug("No path in cache")
	return reader, false
}
