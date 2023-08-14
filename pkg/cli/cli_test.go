package cli_test

import (
	"encoding/json"
	"fmt"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/semaphoreci/test-results/pkg/cli"
	"github.com/stretchr/testify/assert"
)

func Test_LoadFiles(t *testing.T) {

	t.Run("with invalid path to file", func(t *testing.T) {
		filePath := generateFile(t)
		paths, err := cli.LoadFiles([]string{fmt.Sprintf("%s1", filePath)}, ".xml")

		assert.Len(t, paths, 0, "should return correct number of files")
		assert.NotNil(t, err, "should throw error")
		os.RemoveAll(filePath)
	})

	t.Run("with single file", func(t *testing.T) {
		filePath := generateFile(t)
		paths, err := cli.LoadFiles([]string{filePath}, ".xml")

		assert.Equal(t, filePath, paths[0], "should contain correct file path")
		assert.Len(t, paths, 1, "should return correct number of files")
		assert.Nil(t, err, "should not throw error")
		os.RemoveAll(filePath)
	})

	t.Run("with directory", func(t *testing.T) {
		dirPath := generateDir(t)
		assert.NotEqual(t, "", dirPath)

		paths, err := cli.LoadFiles([]string{dirPath}, ".xml")
		assert.Len(t, paths, 5, "should return correct number of files")
		assert.Nil(t, err, "should not throw error")

		os.RemoveAll(dirPath)
	})

	t.Run("with big directory", func(t *testing.T) {
		dirPath := generateDirWithFilesAndNestedDir(t, 2600, 3)
		assert.NotEmpty(t, dirPath)

		paths, err := cli.LoadFiles([]string{dirPath}, ".xml")
		assert.Len(t, paths, 2600, "should return correct number of files")
		assert.Nil(t, err, "should not throw error")

		os.RemoveAll(dirPath)
	})

}

func generateFile(t *testing.T) string {
	filePath, err := os.CreateTemp("", "file-*.xml")
	if err != nil {
		t.Errorf("Failed to create temporary file: %v", err)
	}

	return filePath.Name()
}

func generateDir(t *testing.T) string {
	return generateDirWithFilesAndNestedDir(t, 5, 3)
}

func generateDirWithFilesAndNestedDir(t *testing.T, fNumber, dirNumber int) string {
	dirPath, err := os.MkdirTemp("", "")
	assert.Nil(t, err)

	nestedDir, err := os.MkdirTemp(dirPath, "xml-*")
	assert.Nil(t, err)

	for i := 0; i < fNumber; i++ {
		_, err = os.CreateTemp(nestedDir, "file-*.xml")
		assert.Nil(t, err)
	}

	nestedDir, _ = os.MkdirTemp(dirPath, "json-*")
	for i := 0; i < dirNumber; i++ {
		_, err := os.MkdirTemp(nestedDir, "file-*.json")
		assert.Nil(t, err)
	}

	return dirPath
}

func TestWriteToTmpFile(t *testing.T) {
	tr := parser.TestResults{
		ID:         "1234",
		Name:       "Test",
		Framework:  "JUnit",
		IsDisabled: false,
		Suites:     nil,
		Summary: parser.Summary{
			Total:    10,
			Passed:   5,
			Skipped:  0,
			Error:    0,
			Failed:   5,
			Disabled: 0,
			Duration: 360,
		},
		Status:        "OK",
		StatusMessage: "Test",
	}
	result := parser.Result{TestResults: []parser.TestResults{tr}}
	jsonData, _ := json.Marshal(&result)

	t.Run("Write to one tmp file", func(t *testing.T) {
		file, err := cli.WriteToTmpFile(jsonData)
		assert.NoError(t, err)
		os.Remove(file)
	})

	t.Run("Write to three thousand tmp files", func(t *testing.T) {
		fileNumber := 3000
		files := make([]string, 0, fileNumber)

		for i := 0; i < fileNumber; i++ {
			file, err := cli.WriteToTmpFile(jsonData)
			assert.NoError(t, err)

			files = append(files, file)
		}

		for _, file := range files {
			os.Remove(file)
		}
	})
}

func TestWriteToFile(t *testing.T) {
	tr := parser.TestResults{
		ID:         "1234",
		Name:       "Test",
		Framework:  "JUnit",
		IsDisabled: false,
		Suites:     nil,
		Summary: parser.Summary{
			Total:    10,
			Passed:   5,
			Skipped:  0,
			Error:    0,
			Failed:   5,
			Disabled: 0,
			Duration: 360,
		},
		Status:        "OK",
		StatusMessage: "Test",
	}
	result := parser.Result{TestResults: []parser.TestResults{tr}}
	jsonData, _ := json.Marshal(&result)

	t.Run("Write to one file", func(t *testing.T) {
		file, err := cli.WriteToFile(jsonData, "out")
		assert.NoError(t, err)
		os.Remove(file)
	})

	t.Run("Write to three thousand files", func(t *testing.T) {
		fileNumber := 3000
		dirPath, err := os.MkdirTemp("", "test-results-*")
		require.NoError(t, err)

		defer os.RemoveAll(dirPath)

		for i := 0; i < fileNumber; i++ {
			tmpFile, err := os.CreateTemp(dirPath, "result-*.json")
			require.NoError(t, err)

			_, err = cli.WriteToFile(jsonData, tmpFile.Name())
			require.NoError(t, err)
		}
	})
}
