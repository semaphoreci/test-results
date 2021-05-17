package fileloader

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/semaphoreci/test-results/pkg/logger"
)

var readers map[string]*bytes.Reader = make(map[string]*bytes.Reader)

// Load reader from internal buffer or create new one
func Load(path string, reader *bytes.Reader) (*bytes.Reader, bool) {
	return decode(path, reader)
}

// Ensure file exists at given path with contents read from reader
func Ensure(reader *bytes.Reader) string {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}

	reader.WriteTo(file)
	defer file.Close()

	return file.Name()
}

func decode(path string, reader *bytes.Reader) (*bytes.Reader, bool) {
	foundReader, exists := readers[path]
	if exists && foundReader != nil && foundReader.Size() == reader.Size() {
		logger.Debug("Path read from cache")
		foundReader.Seek(0, io.SeekStart)
		return foundReader, true
	}
	readers[path] = reader
	logger.Debug("No path in cache")
	return reader, false
}
