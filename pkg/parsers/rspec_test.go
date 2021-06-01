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
	assert.Equal(t, "99ec6b78-8d28-33bb-9c4b-e38fd0000bf4", testResults.ID)
	assert.Equal(t, parser.StatusSuccess, testResults.Status)
	assert.Equal(t, "", testResults.StatusMessage)
	assert.Equal(t, 2, len(testResults.Suites))

	assert.Equal(t, "foo", testResults.Suites[0].Name)
	assert.Equal(t, "961b9fe2-d1d3-3f8a-9d14-2adc16701583", testResults.Suites[0].ID)
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
	assert.Equal(t, "977c3162-ea09-3747-8c87-6f8c78585f51", testResults.Suites[1].Tests[0].ID)
	assert.Equal(t, "Foo", testResults.Suites[1].Tests[0].Classname)
	assert.Equal(t, "foo/bar.o", testResults.Suites[1].Tests[0].File)
	assert.Equal(t, time.Duration(123400000), testResults.Suites[1].Tests[0].Duration)

	assert.Equal(t, "2bd6df3a-319b-30bf-8193-b2fdeb1deabd", testResults.Suites[0].Tests[0].ID)
	assert.Equal(t, "c8ffb89f-c186-3104-aa72-c0c18452f5b5", testResults.Suites[0].Tests[1].ID)
	assert.Equal(t, "2bd6df3a-319b-30bf-8193-b2fdeb1deabd", testResults.Suites[0].Tests[2].ID)

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
			ID:   "03999e27-03e1-37e7-adbf-e712e5f35d67",
			Name: "foo",
			Tests: []test{
				{
					ID:   "5830975e-54e3-3951-8b40-c786de5131e6",
					Name: "bar",
				},
				{
					ID:   "42680dd4-916d-3895-978f-c56dec5bb1f0",
					Name: "baz",
				},
			},
		},
		{
			ID:   "7a9dd0d0-961d-36b7-af47-94deab34e474",
			Name: "1234",
			Tests: []test{
				{
					ID:   "5830975e-54e3-3951-8b40-c786de5131e6",
					Name: "bar",
				},
				{
					ID:   "42680dd4-916d-3895-978f-c56dec5bb1f0",
					Name: "baz",
				},
			},
		},
		{
			ID:   "f2385b0c-5155-3ead-ac47-64e58b31546f",
			Name: "",
			Tests: []test{
				{
					ID:   "5830975e-54e3-3951-8b40-c786de5131e6",
					Name: "bar",
				},
				{
					ID:   "42680dd4-916d-3895-978f-c56dec5bb1f0",
					Name: "baz",
				},
			},
		},
		{
			ID:   "5a0a09d5-7020-3585-badc-44dd2f681c73",
			Name: "1235",
			Tests: []test{
				{
					ID:   "6d356aa1-05b8-3015-a446-ed898d92f9e0",
					Name: "bar",
				},
				{
					ID:   "7e0510b2-d80f-32b2-a880-097ca189b7aa",
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
	assert.Equal(t, "dda9b4d2-9e8d-3547-9fd0-24bd78148a7a", testResults.ID)
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
