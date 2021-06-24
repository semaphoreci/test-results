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

	path := fileloader.Ensure(reader)

	p := NewMocha()
	testResults := p.Parse(path)
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
			ID:   "eb15cdbd-e1c4-3544-b31b-fa61427ec0ba",
			Name: "foo",
			Tests: []test{
				{
					ID:   "04a92157-8b5e-3d53-9209-bb849750f267",
					Name: "bar",
				},
				{
					ID:   "20a228ae-ff3c-3f5d-8798-ea7eaecbbc32",
					Name: "baz",
				},
			},
		},
		{
			ID:   "eb15cdbd-e1c4-3544-b31b-fa61427ec0ba",
			Name: "1234",
			Tests: []test{
				{
					ID:   "04a92157-8b5e-3d53-9209-bb849750f267",
					Name: "bar",
				},
				{
					ID:   "20a228ae-ff3c-3f5d-8798-ea7eaecbbc32",
					Name: "baz",
				},
			},
		},
		{
			ID:   "eb15cdbd-e1c4-3544-b31b-fa61427ec0ba",
			Name: "",
			Tests: []test{
				{
					ID:   "04a92157-8b5e-3d53-9209-bb849750f267",
					Name: "bar",
				},
				{
					ID:   "20a228ae-ff3c-3f5d-8798-ea7eaecbbc32",
					Name: "baz",
				},
			},
		},
		{
			ID:   "fb5253d7-944e-3833-958c-118b9a8bb7c3",
			Name: "1235",
			Tests: []test{
				{
					ID:   "5a8bd9c8-02b4-35f8-af6a-16608c914743",
					Name: "bar",
					File: "foo/bar:123",
				},
				{
					ID:   "e1015998-a7e4-37fc-b750-7a1d1fad8e25",
					Name: "baz",
					File: "foo/baz",
				},
			},
		},
	}

	path := fileloader.Ensure(reader)

	p := NewMocha()
	testResults := p.Parse(path)
	assert.Equal(t, "ff", testResults.Name)
	assert.Equal(t, "mocha", testResults.Framework)
	assert.Equal(t, "8037b8cf-de65-31d9-bc40-d0dfcc6ad82c", testResults.ID)
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

	path := fileloader.Ensure(reader)

	p := NewMocha()
	testResults := p.Parse(path)
	assert.Equal(t, parser.StatusError, testResults.Status)
	assert.NotEmpty(t, testResults.StatusMessage)
}
