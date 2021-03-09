package fileloader

import (
	"bytes"
	"fmt"
	"log"
)

var readers map[string]bytes.Reader = make(map[string]bytes.Reader)

// Load reader from internal buffer or create new one
func Load(path string, reader bytes.Reader) (*bytes.Reader, error) {
	decoder, err := decode(path, reader)

	if err != nil {
		return nil, err
	}

	return &decoder, nil
}

func decode(path string, reader bytes.Reader) (bytes.Reader, error) {
	for key, reader := range readers {
		if key == path {
			log.Printf("Path from cache: %s\n", path)
			return reader, nil
		}
	}

	readers[path] = reader

	log.Printf("New path in cache: %s\n", path)
	return reader, nil
}

// Log ...
func Log(any interface{}) {
	fmt.Printf("%#v\n", any)
}
