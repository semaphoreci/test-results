package parsers

import (
	"bytes"
	"testing"

	"github.com/semaphoreci/test-results/pkg/fileloader"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func Test_GoLang_ParseTestSuite(t *testing.T) {
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

	p := NewGoLang()
	testResults := p.Parse(path)
	assert.Equal(t, "Golang Suite", testResults.Name)
	assert.Equal(t, "golang", testResults.Framework)
	assert.Equal(t, "69cd6757-6b3d-30ca-bb19-0b892b4f399e", testResults.ID)
	assert.Equal(t, parser.StatusSuccess, testResults.Status)
	assert.Equal(t, "", testResults.StatusMessage)

	assert.Equal(t, "foo", testResults.Suites[0].Name)
	assert.Equal(t, "57959178-09f2-30bd-ad59-bef48adce2bb", testResults.Suites[0].ID)

	assert.Equal(t, "bar", testResults.Suites[0].Tests[0].Name)
	assert.Equal(t, "baz", testResults.Suites[0].Tests[1].Name)
	assert.Equal(t, "bar", testResults.Suites[0].Tests[2].Name)

	assert.Equal(t, "58846687-ae5e-387c-b172-1f3b89fa6a67", testResults.Suites[0].Tests[0].ID)
	assert.Equal(t, "1464edb1-8e78-3e35-8194-7daf93eda8c7", testResults.Suites[0].Tests[1].ID)
	assert.Equal(t, "58846687-ae5e-387c-b172-1f3b89fa6a67", testResults.Suites[0].Tests[2].ID)

}

func Test_GoLang_ParseTestSuites(t *testing.T) {
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
			ID:   "618eb6d2-8166-317a-a050-6ad11e07ca06",
			Name: "foo",
			Tests: []test{
				{
					ID:   "20ab6d5e-8e9e-3e53-bec5-88b8daf74faa",
					Name: "bar",
				},
				{
					ID:   "8c783fed-516b-3ad8-9e70-87694988d12e",
					Name: "baz",
				},
			},
		},
		{
			ID:   "618eb6d2-8166-317a-a050-6ad11e07ca06",
			Name: "1234",
			Tests: []test{
				{
					ID:   "20ab6d5e-8e9e-3e53-bec5-88b8daf74faa",
					Name: "bar",
				},
				{
					ID:   "8c783fed-516b-3ad8-9e70-87694988d12e",
					Name: "baz",
				},
			},
		},
		{
			ID:   "618eb6d2-8166-317a-a050-6ad11e07ca06",
			Name: "",
			Tests: []test{
				{
					ID:   "20ab6d5e-8e9e-3e53-bec5-88b8daf74faa",
					Name: "bar",
				},
				{
					ID:   "8c783fed-516b-3ad8-9e70-87694988d12e",
					Name: "baz",
				},
			},
		},
		{
			ID:   "4d4f9483-afc4-3e47-97f7-e216ec50f225",
			Name: "1235",
			Tests: []test{
				{
					ID:   "8c7b771b-3ec2-33db-b8d8-2c09af0c01cc",
					Name: "bar",
				},
				{
					ID:   "9cb3f829-ba01-3d32-b864-5cf7148b6f98",
					Name: "baz",
				},
			},
		},
	}

	path := fileloader.Ensure(reader)

	p := NewGoLang()
	testResults := p.Parse(path)
	assert.Equal(t, "ff", testResults.Name)
	assert.Equal(t, "golang", testResults.Framework)
	assert.Equal(t, "488d730c-9521-3018-acd2-ce18a75a7077", testResults.ID)
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

func Test_GoLang_ParseInvalidRoot(t *testing.T) {
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

	p := NewGoLang()
	testResults := p.Parse(path)
	assert.Equal(t, parser.StatusError, testResults.Status)
	assert.NotEmpty(t, testResults.StatusMessage)
}
