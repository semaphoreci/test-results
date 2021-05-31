package cli_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/semaphoreci/test-results/pkg/cli"
	"github.com/stretchr/testify/assert"
)

func Test_LoadFiles(t *testing.T) {

	t.Run("with invalid path to file", func(t *testing.T) {
		filePath := generateFile(t)
		paths, err := cli.LoadFiles([]string{fmt.Sprintf("%s1", filePath)})

		assert.Len(t, paths, 0, "should return correct number of files")
		assert.NotNil(t, err, "should throw error")
		os.RemoveAll(filePath)
	})

	t.Run("with single file", func(t *testing.T) {
		filePath := generateFile(t)
		paths, err := cli.LoadFiles([]string{filePath})

		assert.Equal(t, filePath, paths[0], "should contain correct file path")
		assert.Len(t, paths, 1, "should return correct number of files")
		assert.Nil(t, err, "should not throw error")
		os.RemoveAll(filePath)
	})

	t.Run("with directory", func(t *testing.T) {
		dirPath := generateDir(t)
		assert.NotEqual(t, "", dirPath)

		paths, err := cli.LoadFiles([]string{dirPath})
		assert.Len(t, paths, 5, "should return correct number of files")
		assert.Nil(t, err, "should not throw error")

		os.RemoveAll(dirPath)
	})
}

func generateFile(t *testing.T) string {
	filePath, err := ioutil.TempFile("", "file-*.xml")
	if err != nil {
		t.Errorf("Failed to create temporary file: %v", err)
	}

	return filePath.Name()
}

func generateDir(t *testing.T) string {
	dirPath, err := ioutil.TempDir("", "")
	assert.Nil(t, err)

	nestedDir, err := ioutil.TempDir(dirPath, "xml-*")
	assert.Nil(t, err)

	for i := 0; i < 5; i++ {
		_, err = ioutil.TempFile(nestedDir, "file-*.xml")
		assert.Nil(t, err)
	}

	nestedDir, _ = ioutil.TempDir(dirPath, "json-*")
	for i := 0; i < 3; i++ {
		_, err := ioutil.TempFile(nestedDir, "file-*.json")
		assert.Nil(t, err)
	}

	return dirPath
}
