package parsers

import (
	"bytes"
	"testing"
	"time"

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
			<testsuite id="1234" name="foo" time="0.1234" tests="10" failures="5" errors="1">
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
			</testsuite>
	`))

	err := element.Parse(reader)
	assert.Nil(t, err)

	testResults := newTestResults(element)

	duration, _ := time.ParseDuration("0.1234s")
	assert.Equal(t, "foo", testResults.Name)
	assert.Equal(t, duration, testResults.Summary.Duration)
	assert.Equal(t, 10, testResults.Summary.Total)
	assert.Equal(t, 5, testResults.Summary.Failed)
	assert.Equal(t, 1, testResults.Summary.Error)
	assert.Equal(t, 4, testResults.Summary.Passed)
	assert.Equal(t, false, testResults.IsDisabled)
	assert.Equal(t, 1, len(testResults.Suites))
}
