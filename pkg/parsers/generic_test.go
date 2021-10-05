package parsers

import (
	"bytes"
	"testing"

	"github.com/semaphoreci/test-results/pkg/fileloader"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func Test_Generic_ParseTestSuite(t *testing.T) {
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

	p := NewGeneric()
	testResults := p.Parse(path)
	assert.Equal(t, "Generic Suite", testResults.Name)
	assert.Equal(t, "17990af8-cb17-371c-9a8e-215e0e201902", testResults.ID)
	assert.Equal(t, parser.StatusSuccess, testResults.Status)
	assert.Equal(t, "", testResults.StatusMessage)

	assert.Equal(t, "foo", testResults.Suites[0].Name)
	assert.Equal(t, "b088d75b-c907-3ac1-9f80-489cfacb1619", testResults.Suites[0].ID)

	assert.Equal(t, "bar", testResults.Suites[0].Tests[0].Name)
	assert.Equal(t, "baz", testResults.Suites[0].Tests[1].Name)
	assert.Equal(t, "bar", testResults.Suites[0].Tests[2].Name)

	assert.Equal(t, "a4e4268d-208f-398d-baa5-f7fd5f904216", testResults.Suites[0].Tests[0].ID)
	assert.Equal(t, "237dc38f-c2ec-3ee4-af95-3b9d003da11c", testResults.Suites[0].Tests[1].ID)
	assert.Equal(t, "a4e4268d-208f-398d-baa5-f7fd5f904216", testResults.Suites[0].Tests[2].ID)

}

func Test_Generic_ParseTestSuites(t *testing.T) {
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
			<testsuite name="diff by classname">
				<testcase name="bar" file="foo/bar" classname="foo">
				</testcase>
				<testcase name="bar" file="foo/bar" classname="bar">
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
			ID:   "b3d01848-18ba-3046-b1c9-8806a10742d6",
			Name: "foo",
			Tests: []test{
				{
					ID:   "a718e5fd-8ba4-3f28-a837-b89ec9e6e2e9",
					Name: "bar",
				},
				{
					ID:   "7e336b3f-0878-3a72-b5cd-69c794251a2a",
					Name: "baz",
				},
			},
		},
		{
			ID:   "b3d01848-18ba-3046-b1c9-8806a10742d6",
			Name: "1234",
			Tests: []test{
				{
					ID:   "a718e5fd-8ba4-3f28-a837-b89ec9e6e2e9",
					Name: "bar",
				},
				{
					ID:   "7e336b3f-0878-3a72-b5cd-69c794251a2a",
					Name: "baz",
				},
			},
		},
		{
			ID:   "b3d01848-18ba-3046-b1c9-8806a10742d6",
			Name: "",
			Tests: []test{
				{
					ID:   "a718e5fd-8ba4-3f28-a837-b89ec9e6e2e9",
					Name: "bar",
				},
				{
					ID:   "7e336b3f-0878-3a72-b5cd-69c794251a2a",
					Name: "baz",
				},
			},
		},
		{
			ID:   "376e439f-2352-341a-8400-ff142935ddda",
			Name: "1235",
			Tests: []test{
				{
					ID:   "1ccf0a01-e251-3fad-8dbc-4b4112a5902a",
					Name: "bar",
					File: "foo/bar:123",
				},
				{
					ID:   "9f6f64c2-9b45-3f43-a380-6358cc960cca",
					Name: "baz",
					File: "foo/baz",
				},
			},
		},
		{
			ID:   "db748ed5-c393-3a56-85b1-12fbedbf9c46",
			Name: "diff by classname",
			Tests: []test{
				{
					ID:   "5f5244e4-8c37-3e5f-8d50-e22bb46391d7",
					Name: "bar",
					File: "foo/bar",
				},
				{
					ID:   "aaffd1c8-90ce-33b9-898f-359e203b551b",
					Name: "bar",
					File: "foo/bar",
				},
			},
		},
	}

	path := fileloader.Ensure(reader)

	p := NewGeneric()
	testResults := p.Parse(path)
	assert.Equal(t, "ff", testResults.Name)
	assert.Equal(t, "cd9c81c6-06c6-3623-b337-6819885fbfe8", testResults.ID)
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

func Test_Generic_ParseInvalidRoot(t *testing.T) {
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

	p := NewGeneric()
	testResults := p.Parse(path)
	assert.Equal(t, parser.StatusError, testResults.Status)
	assert.NotEmpty(t, testResults.StatusMessage)
}
