package parsers

import (
	"bytes"
	"testing"

	"github.com/semaphoreci/test-results/pkg/fileloader"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	reader := bytes.NewReader([]byte(`
		<?xml version="1.0"?>
			<testsuite name="foo" id="1234">
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
			</testsuite>
	`))

	fileloader.Load("/path", reader)

	p := NewGeneric()
	testResults, err := p.Parse("/path")
	assert.Equal(t, nil, err)
	assert.Equal(t, "Generic Parser", testResults.Name)

	assert.Equal(t, "foo", testResults.Suites[0].Name)
	assert.Equal(t, "1234", testResults.Suites[0].ID)

	assert.Equal(t, "bar", testResults.Suites[0].Tests[0].Name)
	assert.Equal(t, "baz", testResults.Suites[0].Tests[1].Name)
}

func Test_newTestResults(t *testing.T) {
	element := parser.NewXMLElement()

	reader := bytes.NewReader([]byte(`
		<?xml version="1.0"?>
			<testsuite name="foo" id="1234">
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
			</testsuite>
	`))

	element.Parse(reader)
}
