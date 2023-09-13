package parsers

import (
	"bytes"
	"github.com/semaphoreci/test-results/pkg/fileloader"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Embedded_ParseTestSuites(t *testing.T) {
	reader := bytes.NewReader([]byte(`
		<?xml version="1.0"?>
	    <testsuites>
			<testsuite name="foo" id="1234">
				<testcase name="bar">
				</testcase>
				<testsuite name="zap" id="4321">
					<testcase name="baz">
					</testcase>
				</testsuite>
				<testsuite name="zup" id="54321">
					<testcase name="bar">
					</testcase>
				</testsuite>
			</testsuite>
		</testsuites>
	`))

	path := fileloader.Ensure(reader)

	e := NewEmbedded()
	testResults := e.Parse(path)
	assert.Equal(t, "Suite", testResults.Name)
	assert.Equal(t, "embedded", testResults.Framework)
	assert.Equal(t, "c5bec5ae-e57f-3dac-98fa-825a5a2cfd55", testResults.ID)
	assert.Equal(t, parser.StatusSuccess, testResults.Status)
	assert.Equal(t, "", testResults.StatusMessage)

	require.Len(t, testResults.Suites, 3)

	assert.Equal(t, "foo", testResults.Suites[0].Name)
	assert.Equal(t, "a4f80929-16cf-38db-a9b0-e7fa3c11d398", testResults.Suites[0].ID)
	require.Len(t, testResults.Suites[0].Tests, 1)
	assert.Equal(t, "bar", testResults.Suites[0].Tests[0].Name)

	assert.Equal(t, "foo\\zap", testResults.Suites[1].Name)
	assert.Equal(t, "3c1b98c1-1db8-3c73-b33e-fa37e44ba2ab", testResults.Suites[1].ID)

	assert.Equal(t, "foo\\zup", testResults.Suites[2].Name)
	assert.Equal(t, "3afdab2c-33ab-30a4-85bf-f9e8872cf38b", testResults.Suites[2].ID)

	require.Len(t, testResults.Suites[0].Tests, 1)
	assert.Equal(t, "6787a9cb-b2c2-3c28-82e0-ff4505555d10", testResults.Suites[0].Tests[0].ID)
	assert.Equal(t, "bar", testResults.Suites[0].Tests[0].Name)

	require.Len(t, testResults.Suites[1].Tests, 1)
	assert.Equal(t, "174855db-b79b-3bda-9404-4874d2eadd10", testResults.Suites[1].Tests[0].ID)
	assert.Equal(t, "baz", testResults.Suites[1].Tests[0].Name)

	require.Len(t, testResults.Suites[2].Tests, 1)
	assert.Equal(t, "b42d9487-cefc-359d-bffd-10ba38879fe1", testResults.Suites[2].Tests[0].ID)
	assert.Equal(t, "bar", testResults.Suites[2].Tests[0].Name)
}

func Test_Embedded_ParseInvalidRoot(t *testing.T) {
	reader := bytes.NewReader([]byte(`
		<?xml version="1.0"?>
		<nontestsuites name="em">
			<testsuite name="foo" id="1234">
				<testcase name="bar">
				</testcase>
				<testsuite name="zap" id="4321">
					<testcase name="baz">
					</testcase>
				</testsuite>
				<testsuite name="zup" id="54321">
					<testcase name="bar">
					</testcase>
				</testsuite>
			</testsuite>
		</nontestsuites>
	`))

	path := fileloader.Ensure(reader)

	p := NewEmbedded()
	testResults := p.Parse(path)
	assert.Equal(t, parser.StatusError, testResults.Status)
	assert.NotEmpty(t, testResults.StatusMessage)
}
