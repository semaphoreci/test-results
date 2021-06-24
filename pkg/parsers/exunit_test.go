package parsers

import (
	"bytes"
	"testing"

	"github.com/semaphoreci/test-results/pkg/fileloader"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func Test_ExUnit_ParseTestSuite(t *testing.T) {
	reader := bytes.NewReader([]byte(`
		<?xml version="1.0"?>
			<testsuite name="foo" id="1234">
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
				<testcase name="bar">
				</testcase>
			</testsuite>
	`))

	path := fileloader.Ensure(reader)

	p := NewExUnit()
	testResults := p.Parse(path)
	assert.Equal(t, "Exunit Suite", testResults.Name)
	assert.Equal(t, "exunit", testResults.Framework)
	assert.Equal(t, "b37fd0fa-7cfa-3b6b-a992-67b61d24b79f", testResults.ID)
	assert.Equal(t, parser.StatusSuccess, testResults.Status)
	assert.Equal(t, "", testResults.StatusMessage)

	assert.Equal(t, "foo", testResults.Suites[0].Name)
	assert.Equal(t, "3e3fd52e-60da-3a34-ace5-e69db05bc15e", testResults.Suites[0].ID)

	assert.Equal(t, "bar", testResults.Suites[0].Tests[0].Name)
	assert.Equal(t, "baz", testResults.Suites[0].Tests[1].Name)
	assert.Equal(t, "bar", testResults.Suites[0].Tests[2].Name)

	assert.Equal(t, "43ad972b-a4fe-3835-8e9d-a9c765261a64", testResults.Suites[0].Tests[0].ID)
	assert.Equal(t, "bfb9eae0-3a41-312d-b429-a32e76be493f", testResults.Suites[0].Tests[1].ID)
	assert.Equal(t, "43ad972b-a4fe-3835-8e9d-a9c765261a64", testResults.Suites[0].Tests[2].ID)

}

func Test_ExUnit_ParseTestSuites(t *testing.T) {
	reader := bytes.NewReader([]byte(`
		<?xml version="1.0"?>
		<testsuites name="ff">
			<testsuite name="foo" id="1234">
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
			</testsuite>
			<testsuite name="1234">
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
			</testsuite>
			<testsuite id="1234">
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
			</testsuite>
			<testsuite name="1235">
				<testcase name="bar" file="foo/bar:123">
				</testcase>
				<testcase name="baz" file="foo/baz">
				</testcase>
			</testsuite>
		</testsuites>
	`))
	type test struct {
		ID   string
		Name string
		File string
	}

	var fixtures = []struct {
		ID    string
		Name  string
		Tests []test
	}{
		{
			ID:   "46c015b1-a078-3514-8beb-a2b7f4c74cd8",
			Name: "foo",
			Tests: []test{
				{
					ID:   "43dde248-0845-32a2-897f-410b5e614c25",
					Name: "bar",
				},
				{
					ID:   "37eed9e1-1ca1-33df-9b21-b8a406f7081c",
					Name: "baz",
				},
			},
		},
		{
			ID:   "46c015b1-a078-3514-8beb-a2b7f4c74cd8",
			Name: "1234",
			Tests: []test{
				{
					ID:   "43dde248-0845-32a2-897f-410b5e614c25",
					Name: "bar",
				},
				{
					ID:   "37eed9e1-1ca1-33df-9b21-b8a406f7081c",
					Name: "baz",
				},
			},
		},
		{
			ID:   "46c015b1-a078-3514-8beb-a2b7f4c74cd8",
			Name: "",
			Tests: []test{
				{
					ID:   "43dde248-0845-32a2-897f-410b5e614c25",
					Name: "bar",
				},
				{
					ID:   "37eed9e1-1ca1-33df-9b21-b8a406f7081c",
					Name: "baz",
				},
			},
		},
		{
			ID:   "b2ca06e8-bc26-3b78-bd36-80c7889331eb",
			Name: "1235",
			Tests: []test{
				{
					ID:   "2ac47d0b-27b9-395a-8b35-62191b8258cc",
					Name: "bar",
					File: "foo/bar:123",
				},
				{
					ID:   "f419224e-a9db-30ca-8bf3-700ae80629d8",
					Name: "baz",
					File: "foo/baz",
				},
			},
		},
	}

	path := fileloader.Ensure(reader)

	p := NewExUnit()
	testResults := p.Parse(path)
	assert.Equal(t, "ff", testResults.Name)
	assert.Equal(t, "exunit", testResults.Framework)
	assert.Equal(t, "09cccd60-46cb-31a4-846f-c8329a1164f9", testResults.ID)
	assert.Equal(t, parser.StatusSuccess, testResults.Status)
	assert.Equal(t, "", testResults.StatusMessage)

	for i, suite := range testResults.Suites {
		fixture := fixtures[i]
		assert.Equal(t, fixture.ID, suite.ID)
		assert.Equal(t, fixture.Name, suite.Name)
		for i := range fixture.Tests {
			assert.Equal(t, fixture.Tests[i].Name, suite.Tests[i].Name)
			assert.Equal(t, fixture.Tests[i].ID, suite.Tests[i].ID)
			assert.Equal(t, fixture.Tests[i].File, suite.Tests[i].File)
		}
	}
}

func Test_ExUnit_ParseInvalidRoot(t *testing.T) {
	reader := bytes.NewReader([]byte(`
		<?xml version="1.0"?>
		<nontestsuites name="ff">
			<testsuite name="foo" id="1234">
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
			</testsuite>
		</nontestsuites>
	`))

	path := fileloader.Ensure(reader)

	p := NewExUnit()
	testResults := p.Parse(path)
	assert.Equal(t, parser.StatusError, testResults.Status)
	assert.NotEmpty(t, testResults.StatusMessage)
}
