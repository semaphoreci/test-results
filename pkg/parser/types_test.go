package parser

import (
	"fmt"
	"testing"
	"time"

	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func Test_Result_Combine(t *testing.T) {
	result := NewResult()

	resultToMerge := NewResult()

	testResult := NewTestResults()
	suite := newSuite("1", "foo")
	newTest(&suite, "1", "foo.1")
	newTest(&suite, "2", "foo.2")
	testResult.Suites = append(testResult.Suites, suite)
	suite = newSuite("1", "foo")
	newTest(&suite, "3", "foo.3")
	newTest(&suite, "4", "foo.4")
	testResult.Suites = append(testResult.Suites, suite)
	resultToMerge.TestResults = append(resultToMerge.TestResults, testResult)
	result.Combine(resultToMerge)

	resultToMerge = NewResult()
	suite = newSuite("1", "foo")
	newTest(&suite, "3", "foo.3")
	newTest(&suite, "4", "foo.4")
	testResult = NewTestResults()
	testResult.Suites = append(testResult.Suites, suite)
	resultToMerge.TestResults = append(resultToMerge.TestResults, testResult)
	result.Combine(resultToMerge)

	resultToMerge = NewResult()
	suite = newSuite("1", "foo")
	newTest(&suite, "5", "foo.51")
	newTest(&suite, "6", "foo.61")
	testResult = NewTestResults()
	testResult.Suites = append(testResult.Suites, suite)
	resultToMerge.TestResults = append(resultToMerge.TestResults, testResult)
	result.Combine(resultToMerge)

	for _, suite := range result.TestResults[0].Suites {
		logger.Info("%+v\n", suite)
	}

	assert.Equal(t, 1, len(result.TestResults))
	assert.Equal(t, 1, len(result.TestResults[0].Suites))
	assert.Equal(t, 6, len(result.TestResults[0].Suites[0].Tests))
}

func Test_TestResults_Combine(t *testing.T) {
	suite := newSuite("1", "foo")
	newTest(&suite, "1", "foo.1")
	newTest(&suite, "2", "foo.2")
	newTest(&suite, "3", "foo.3")

	testResult := NewTestResults()
	testResult.Suites = append(testResult.Suites, suite)

	suite = newSuite("1", "foo")
	newTest(&suite, "1", "foo.1")
	newTest(&suite, "2", "foo.2")

	testResultToMerge := NewTestResults()
	testResultToMerge.Suites = append(testResultToMerge.Suites, suite)

	testResult.Combine(testResultToMerge)

	assert.Equal(t, 1, len(testResult.Suites))
	assert.Equal(t, 3, len(testResult.Suites[0].Tests))
}

func Test_Suite_Combine(t *testing.T) {
	suite := newSuite("1", "foo")
	newTest(&suite, "1", "foo.1")
	newTest(&suite, "2", "foo.2")
	newTest(&suite, "3", "foo.3")
	suiteToMerge := newSuite("1", "foo")
	newTest(&suiteToMerge, "1", "foo.1")
	newTest(&suiteToMerge, "2", "foo.2")

	suite.Combine(suiteToMerge)

	assert.Equal(t, 3, len(suite.Tests))
}

func Test_NewTest_Results(t *testing.T) {
	testResults := NewTestResults()

	assert.IsType(t, TestResults{}, testResults)
	assert.Equal(t, StatusSuccess, testResults.Status)
	assert.Equal(t, "", testResults.StatusMessage)
}

func Test_TestResults_Aggregate(t *testing.T) {
	testResults := TestResults{}

	testResults.Aggregate()
	assert.Equal(t, testResults.Summary, Summary{})

	suite := NewSuite()
	suite.Summary.Total = 6
	suite.Summary.Passed = 1
	suite.Summary.Skipped = 2
	suite.Summary.Error = 2
	suite.Summary.Failed = 1
	suite.Summary.Disabled = 1
	suite.Summary.Duration = time.Duration(1)

	testResults.Suites = append(testResults.Suites, suite)

	testResults.Aggregate()
	assert.Equal(t, testResults.Summary, Summary{6, 1, 2, 2, 1, 1, 1})

	suite = NewSuite()
	suite.Summary.Total = 12
	suite.Summary.Passed = 2
	suite.Summary.Skipped = 4
	suite.Summary.Error = 2
	suite.Summary.Failed = 2
	suite.Summary.Disabled = 2
	suite.Summary.Duration = time.Duration(10)

	testResults.Suites = append(testResults.Suites, suite)
	testResults.Aggregate()

	assert.Equal(t, testResults.Summary, Summary{18, 3, 6, 4, 3, 3, 11})
}

func Test_TestResults_ArrangeSuitesByTestFile(t *testing.T) {
	testResults := NewTestResults()
	testResults.ID = "1"

	suite := newSuite("1", "test_suite_name/with_special_chars.go")
	newTest(&suite, "1", "foo/foo.go")
	newTest(&suite, "2", "foo/bar.go")
	testResults.Suites = append(testResults.Suites, suite)

	suite = newSuite("2", "golang")
	newTest(&suite, "3", "golang/foo.go")
	newTest(&suite, "4", "golang/bar.go")
	testResults.Suites = append(testResults.Suites, suite)

	suite = newSuite("3", "foo/foo.go")
	newTest(&suite, "5", "foo/foo.go")
	newTest(&suite, "6", "foo/foo.go")
	testResults.Suites = append(testResults.Suites, suite)

	assert.Equal(t, 3, len(testResults.Suites), "test results should have correct number of suites before arrangement")
	for _, suite := range testResults.Suites {
		assert.Equal(t, 2, len(suite.Tests), "suites should have correct number of tests before arrangement")
	}

	testResults.Aggregate()

	testResults.ArrangeSuitesByTestFile()

	assert.Equal(t, 4, len(testResults.Suites), "test results should have correct number of suites after arrangement")

	suite = testResults.Suites[0]
	assert.Equal(t, "foo/foo.go", suite.Name, "suite name should match")
	assert.Equal(t, 3, len(suite.Tests), "should contain correct number of tests")
	assert.Equal(t, "foo/foo.go#1", suite.Tests[0].Name, "test name should match")
	assert.Equal(t, "foo/foo.go#5", suite.Tests[1].Name, "test name should match")
	assert.Equal(t, "foo/foo.go#6", suite.Tests[2].Name, "test name should match")

	suite = testResults.Suites[1]
	assert.Equal(t, "foo/bar.go", suite.Name, "suite name should match")
	assert.Equal(t, 1, len(suite.Tests), "should remove tests from old suite")
	assert.Equal(t, "foo/bar.go#2", suite.Tests[0].Name, "test name should match")

	suite = testResults.Suites[2]
	assert.Equal(t, "golang/foo.go", suite.Name, "suite name should match")
	assert.Equal(t, 1, len(suite.Tests), "should remove tests from old suite")
	assert.Equal(t, "golang/foo.go#3", suite.Tests[0].Name, "test name should match")

	suite = testResults.Suites[3]
	assert.Equal(t, "golang/bar.go", suite.Name, "suite name should match")
	assert.Equal(t, 1, len(suite.Tests), "should remove tests from old suite")
	assert.Equal(t, "golang/bar.go#4", suite.Tests[0].Name, "test name should match")
}

func Test_TestResults_ArrangeSuitesByTestFile_SingleSuite(t *testing.T) {
	testResults := NewTestResults()
	testResults.ID = "1"

	suite := newSuite("1", "test_suite_name/with_special_chars.go")
	newTest(&suite, "1", "foo/foo.go")
	newTest(&suite, "2", "foo/bar.go")
	test := NewTest()
	test.ID = "3"
	test.Name = "Foo"
	suite.Tests = append(suite.Tests, test)
	test.ID = "4"
	test.Name = "Bar"
	suite.Tests = append(suite.Tests, test)

	testResults.Suites = append(testResults.Suites, suite)

	assert.Equal(t, 1, len(testResults.Suites), "test results should have correct number of suites before arrangement")
	for _, suite := range testResults.Suites {
		assert.Equal(t, 4, len(suite.Tests), "suites should have correct number of tests before arrangement")
	}

	testResults.ArrangeSuitesByTestFile()

	assert.Equal(t, 3, len(testResults.Suites), "test results should have correct number of suites after arrangement")

	suite = testResults.Suites[0]
	assert.Equal(t, "foo/foo.go", suite.Name, "suite name should match")
	assert.Equal(t, 1, len(suite.Tests), "should contain correct number of tests")
	assert.Equal(t, "foo/foo.go#1", suite.Tests[0].Name, "test name should match")

	suite = testResults.Suites[1]
	assert.Equal(t, "foo/bar.go", suite.Name, "suite name should match")
	assert.Equal(t, 1, len(suite.Tests), "should remove tests from old suite")
	assert.Equal(t, "foo/bar.go#2", suite.Tests[0].Name, "test name should match")

	suite = testResults.Suites[2]
	assert.Equal(t, "test_suite_name/with_special_chars.go", suite.Name, "suite name should match")
	assert.Equal(t, 2, len(suite.Tests), "should remove tests from old suite")
	assert.Equal(t, "Foo", suite.Tests[0].Name, "test name should match")
	assert.Equal(t, "Bar", suite.Tests[1].Name, "test name should match")

}

func Test_NewSuite(t *testing.T) {
	suite := NewSuite()

	assert.IsType(t, suite, Suite{})
}

func Test_Suite_Aggregate(t *testing.T) {
	suite := NewSuite()

	suite.Aggregate()
	assert.Equal(t, Summary{}, suite.Summary)

	test := NewTest()
	suite.Tests = append(suite.Tests, test)
	suite.Aggregate()

	assert.Equal(t, Summary{Total: 1, Passed: 1}, suite.Summary)

	test = NewTest()
	test.State = StateFailed
	suite.Tests = append(suite.Tests, test)
	suite.Aggregate()

	assert.Equal(t, Summary{Total: 2, Passed: 1, Failed: 1}, suite.Summary)

	test = NewTest()
	test.State = StateSkipped
	test.Duration = time.Duration(10)
	suite.Tests = append(suite.Tests, test)
	suite.Aggregate()

	assert.Equal(t, Summary{Total: 3, Passed: 1, Failed: 1, Skipped: 1, Duration: 10}, suite.Summary)

	test = NewTest()
	test.State = StateError
	test.Duration = time.Duration(50)
	suite.Tests = append(suite.Tests, test)
	suite.Aggregate()

	assert.Equal(t, Summary{Total: 4, Passed: 1, Failed: 1, Skipped: 1, Error: 1, Duration: 60}, suite.Summary)

	test = NewTest()
	test.State = StateDisabled
	test.Duration = time.Duration(50)
	suite.Tests = append(suite.Tests, test)
	suite.Aggregate()

	assert.Equal(t, Summary{Total: 5, Passed: 1, Failed: 1, Skipped: 1, Error: 1, Disabled: 1, Duration: 110}, suite.Summary)
}

func Test_NewTest(t *testing.T) {
	test := NewTest()

	assert.Equal(t, test.State, StatePassed, "is in passed state by default")
	assert.IsType(t, test, Test{})
}

func Test_NewError(t *testing.T) {
	obj := NewError()

	assert.IsType(t, obj, Error{})
}

func Test_NewFailure(t *testing.T) {
	obj := NewFailure()

	assert.IsType(t, obj, Failure{})
}

func newTest(suite *Suite, id string, file string) {
	test := NewTest()
	test.ID = id
	test.File = file
	test.Name = fmt.Sprintf("%s#%s", file, id)
	suite.AppendTest(test)
}

func newSuite(id string, name string) Suite {
	suite := NewSuite()
	suite.ID = id
	suite.Name = name
	return suite
}
