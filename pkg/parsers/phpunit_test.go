package parsers

import (
	"bytes"
	"testing"

	"github.com/semaphoreci/test-results/pkg/fileloader"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func Test_PHPUnit_ParseTestSuite(t *testing.T) {
	reader := bytes.NewReader([]byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<testsuites>
			<testsuite name="" tests="5" assertions="5" errors="0" warnings="0" failures="0" skipped="0" time="0.186278">
				<testsuite name="test1" tests="2" assertions="2" errors="0" warnings="0" failures="0" skipped="0" time="0.185621">
					<testsuite name="Tests\Tests1\FirstTest" file="/app/tests1/FirstTest.php" tests="2" assertions="2" errors="0" warnings="0" failures="0" skipped="0" time="0.184274">
						<testcase name="testFakeFirst" class="Tests\Tests1\FirstTest" classname="Tests.Tests1.FirstTest" file="/app/tests1/FirstTest.php" line="9" assertions="1" time="0.184274"/>
						<testcase name="secondTakeOnFirstTest" class="Tests\Tests1\FirstTest" classname="Tests.Tests1.FirstTest" file="/app/tests1/FirstTest.php" line="11" assertions="1" time="0"/>
					</testsuite>
					<testsuite name="Tests\Tests1\SecondTest" file="/app/tests1/SecondTest.php" tests="1" assertions="1" errors="0" warnings="0" failures="0" skipped="0" time="0.001347">
						<testcase name="testFakeSecond" class="Tests\Tests1\SecondTest" classname="Tests.Tests1.SecondTest" file="/app/tests1/SecondTest.php" line="9" assertions="1" time="0.001347"/>
					</testsuite>
				</testsuite>
				<testsuite name="test2" tests="5" assertions="5" errors="0" warnings="0" failures="0" skipped="0" time="0.000657">
					<testsuite name="Tests\Tests2\FirstTest" file="/app/tests2/FirstTest.php" tests="4" assertions="1" errors="0" warnings="0" failures="0" skipped="0" time="0.000146">
						<testcase name="testFakeFirst1" class="Tests\Tests2\FirstTest" classname="Tests.Tests2.FirstTest" file="/app/tests2/FirstTest.php" line="9" assertions="1" time="0.000146"/>
						<testcase name="testFakeFirst2" class="Tests\Tests2\FirstTest" classname="Tests.Tests2.FirstTest" file="/app/tests2/FirstTest.php" line="10" assertions="1" time="0.000146"/>
						<testcase name="testFakeFirst3" class="Tests\Tests2\FirstTest" classname="Tests.Tests2.FirstTest" file="/app/tests2/FirstTest.php" line="11" assertions="1" time="0.000146"/>
						<testcase name="testFakeFirst4" class="Tests\Tests2\FirstTest" classname="Tests.Tests2.FirstTest" file="/app/tests2/FirstTest.php" line="12" assertions="1" time="0.000146"/>
					</testsuite>
					<testsuite name="Tests\Tests2\SecondTest" file="/app/tests2/SecondTest.php" tests="1" assertions="1" errors="0" warnings="0" failures="0" skipped="0" time="0.000510">
						<testcase name="testFakeSecond" class="Tests\Tests2\SecondTest" classname="Tests.Tests2.SecondTest" file="/app/tests2/SecondTest.php" line="9" assertions="1" time="0.000510"/>
					</testsuite>
				</testsuite>
			</testsuite>
		</testsuites>
	`))

	path := fileloader.Ensure(reader)

	p := NewPHPUnit()
	testResults := p.Parse(path)
	assert.Equal(t, "PHPUnit Suite", testResults.Name)
	assert.Equal(t, "phpunit", testResults.Framework)
	assert.Equal(t, parser.StatusSuccess, testResults.Status)
	assert.Equal(t, "", testResults.StatusMessage)

	assert.Equal(t, "test1\\Tests\\Tests1\\FirstTest", testResults.Suites[0].Name)
	assert.Equal(t, "test1\\Tests\\Tests1\\SecondTest", testResults.Suites[1].Name)
	assert.Equal(t, "test2\\Tests\\Tests2\\FirstTest", testResults.Suites[2].Name)
	assert.Equal(t, "test2\\Tests\\Tests2\\SecondTest", testResults.Suites[3].Name)

	assert.Equal(t, 4, len(testResults.Suites))
	assert.Equal(t, 8, testResults.Summary.Total)
}

func Test_PHPUnit_flattenTestSuite(t *testing.T) {
	path := fileloader.Ensure(bytes.NewReader([]byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<testsuites>
			<testsuite name="foo">
				<testsuite name="1">
					<testcase></testcase>
				</testsuite>
				<testsuite name="2">
					<testcase></testcase>
				</testsuite>
				<testsuite name="3">
					<testcase></testcase>
				</testsuite>
			</testsuite>
			<testsuite name="bar">
				<testsuite name="1">
					<testcase></testcase>
				</testsuite>
			</testsuite>
			<testsuite name="baz">
				<testsuite name="1">
					<testsuite name="test">
						<testcase></testcase>
					</testsuite>
				</testsuite>
			</testsuite>
			<testsuite name="im empty">
				<testsuite name="with no cases">
					<testsuite name="so im not present in report">
					</testsuite>
				</testsuite>
			</testsuite>
		</testsuites>
	`)))

	xml, err := LoadXML(path)
	flattenTestSuites(xml)
	assert.Nil(t, err)
	assert.Equal(t, 5, len(xml.Children))
	assert.Equal(t, "foo\\1", xml.Children[0].Attr("name"))
	assert.Equal(t, "foo\\2", xml.Children[1].Attr("name"))
	assert.Equal(t, "foo\\3", xml.Children[2].Attr("name"))
	assert.Equal(t, "bar\\1", xml.Children[3].Attr("name"))
	assert.Equal(t, "baz\\1\\test", xml.Children[4].Attr("name"))
}

func Test_PHPUnit_prefixSuiteName(t *testing.T) {
	assertions := []struct {
		prefix         string
		name           string
		expectedResult string
		desc           string
	}{
		{prefix: "", name: "foo", expectedResult: "foo", desc: "no prefix"},
		{prefix: "", name: "", expectedResult: "", desc: "empty"},
		{prefix: "bar", name: "", expectedResult: "bar", desc: "no name"},
		{prefix: "bar", name: "foo", expectedResult: "bar\\foo", desc: "with prefix #1"},
		{prefix: "bar\\baz", name: "foo", expectedResult: "bar\\baz\\foo", desc: "with prefix #2"},
	}

	for _, tt := range assertions {
		t.Run(tt.desc, func(t *testing.T) {
			actualResult := prefixSuiteName(tt.name, tt.prefix)

			assert.Equal(t, tt.expectedResult, actualResult)
		})
	}
}
