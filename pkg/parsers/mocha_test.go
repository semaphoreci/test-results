package parsers

import (
	"bytes"
	"testing"

	"github.com/semaphoreci/test-results/pkg/fileloader"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func Test_Mocha_ParseTestSuite(t *testing.T) {
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

	fileloader.Load("/path1", reader)

	p := NewMocha()
	testResults := p.Parse("/path1")
	assert.Equal(t, "Mocha Suite", testResults.Name)
	assert.Equal(t, "mocha", testResults.Framework)
	assert.Equal(t, "3351c17a-881f-3a48-9a3a-5d62681955ed", testResults.ID)
	assert.Equal(t, parser.StatusSuccess, testResults.Status)
	assert.Equal(t, "", testResults.StatusMessage)

	assert.Equal(t, "foo", testResults.Suites[0].Name)
	assert.Equal(t, "52ce7b29-afd4-3689-a6cc-7d3a80df2ad3", testResults.Suites[0].ID)

	assert.Equal(t, "bar", testResults.Suites[0].Tests[0].Name)
	assert.Equal(t, "baz", testResults.Suites[0].Tests[1].Name)
	assert.Equal(t, "bar", testResults.Suites[0].Tests[2].Name)

	assert.Equal(t, "b8a9e778-6b95-3b7e-9643-532875abbfad", testResults.Suites[0].Tests[0].ID)
	assert.Equal(t, "66727304-bdcc-30bf-b76e-6f3ee9374aa8", testResults.Suites[0].Tests[1].ID)
	assert.Equal(t, "b8a9e778-6b95-3b7e-9643-532875abbfad", testResults.Suites[0].Tests[2].ID)

}

func Test_Mocha_ParseTestSuites(t *testing.T) {
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
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
			</testsuite>
		</testsuites>
	`))
	type test struct {
		ID   string
		Name string
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
				},
				{
					ID:   "9f6f64c2-9b45-3f43-a380-6358cc960cca",
					Name: "baz",
				},
			},
		},
	}

	fileloader.Load("/path2", reader)

	p := NewMocha()
	testResults := p.Parse("/path2")
	assert.Equal(t, "ff", testResults.Name)
	assert.Equal(t, "mocha", testResults.Framework)
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
		}
	}
}

func Test_Mocha_ParseInvalidRoot(t *testing.T) {
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

	fileloader.Load("/path3", reader)

	p := NewMocha()
	testResults := p.Parse("/path3")
	assert.Equal(t, parser.StatusError, testResults.Status)
	assert.NotEmpty(t, testResults.StatusMessage)
}
