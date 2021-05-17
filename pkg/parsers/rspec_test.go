package parsers

import (
	"bytes"
	"testing"
	"time"

	"github.com/semaphoreci/test-results/pkg/fileloader"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func Test_RSpec_ParseTestSuite(t *testing.T) {
	reader := bytes.NewReader([]byte(`
		<?xml version="1.0"?>
			<testsuite name="foo" id="1234">
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
				<testcase name="bar">
				</testcase>
				<testcase id="1" classname="Foo" name="foo bar" file="foo/bar.o" time="0.1234">
				</testcase>
			</testsuite>
	`))

	path := fileloader.Ensure(reader)

	p := NewRSpec()
	testResults := p.Parse(path)
	assert.Equal(t, "Rspec Suite", testResults.Name)
	assert.Equal(t, "rspec", testResults.Framework)
	assert.Equal(t, "a9497520-97f8-30ac-a6f4-6e3f69dc69a5", testResults.ID)
	assert.Equal(t, parser.StatusSuccess, testResults.Status)
	assert.Equal(t, "", testResults.StatusMessage)
	assert.Equal(t, 2, len(testResults.Suites))

	assert.Equal(t, "foo", testResults.Suites[0].Name)
	assert.Equal(t, "061933f1-470b-30d8-b5cd-a254d9e99a5d", testResults.Suites[0].ID)
	assert.Equal(t, 3, len(testResults.Suites[0].Tests))
	assert.Equal(t, 3, testResults.Suites[0].Summary.Total)
	assert.Equal(t, 3, testResults.Suites[0].Summary.Passed)
	assert.Equal(t, 0, testResults.Suites[0].Summary.Disabled)
	assert.Equal(t, 0, testResults.Suites[0].Summary.Failed)
	assert.Equal(t, 0, testResults.Suites[0].Summary.Error)
	assert.Equal(t, 0, testResults.Suites[0].Summary.Skipped)
	assert.Equal(t, time.Duration(0), testResults.Suites[0].Summary.Duration)

	assert.Equal(t, "bar", testResults.Suites[0].Tests[0].Name)
	assert.Equal(t, "baz", testResults.Suites[0].Tests[1].Name)
	assert.Equal(t, "bar", testResults.Suites[0].Tests[2].Name)

	assert.Equal(t, "foo bar", testResults.Suites[1].Tests[0].Name)
	assert.Equal(t, "cf2b566c-25e1-3260-b62a-c8e6dd934d2d", testResults.Suites[1].Tests[0].ID)
	assert.Equal(t, "Foo", testResults.Suites[1].Tests[0].Classname)
	assert.Equal(t, "foo/bar.o", testResults.Suites[1].Tests[0].File)
	assert.Equal(t, time.Duration(123400000), testResults.Suites[1].Tests[0].Duration)

	assert.Equal(t, "5dbd5f2f-cce8-3bad-82c0-acef31c3e60d", testResults.Suites[0].Tests[0].ID)
	assert.Equal(t, "5944a4e5-6a2f-3e47-8a46-061f291a2deb", testResults.Suites[0].Tests[1].ID)
	assert.Equal(t, "5dbd5f2f-cce8-3bad-82c0-acef31c3e60d", testResults.Suites[0].Tests[2].ID)

}

func Test_RSpec_ParseTestSuites(t *testing.T) {
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
			ID:   "e1be58d1-a5b8-3aad-8f87-ef128cc35750",
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
			ID:   "88ea00ac-433b-307a-81d5-11d7c3afa199",
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

	path := fileloader.Ensure(reader)

	p := NewRSpec()
	testResults := p.Parse(path)
	assert.Equal(t, "ff", testResults.Name)
	assert.Equal(t, "rspec", testResults.Framework)
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

func Test_RSpec_ParseInvalidRoot(t *testing.T) {
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

	p := NewRSpec()
	testResults := p.Parse(path)
	assert.Equal(t, parser.StatusError, testResults.Status)
	assert.NotEmpty(t, testResults.StatusMessage)
}
