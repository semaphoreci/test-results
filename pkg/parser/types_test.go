package parser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTestResults(t *testing.T) {
	testResults := NewTestResults()

	assert.IsType(t, testResults, TestResults{})
}

func TestTestResults_Aggregate(t *testing.T) {
	testResults := TestResults{}

	testResults.Aggregate()
	assert.Equal(t, testResults.Summary, Summary{})

	suite := NewSuite()
	suite.Summary.Total = 5
	suite.Summary.Passed = 1
	suite.Summary.Skipped = 2
	suite.Summary.Error = 2
	suite.Summary.Failed = 1
	suite.Summary.Duration = time.Duration(1)

	testResults.Suites = append(testResults.Suites, suite)

	testResults.Aggregate()
	assert.Equal(t, testResults.Summary, Summary{5, 1, 2, 2, 1, 1})

	suite = NewSuite()
	suite.Summary.Total = 10
	suite.Summary.Passed = 2
	suite.Summary.Skipped = 4
	suite.Summary.Error = 2
	suite.Summary.Failed = 2
	suite.Summary.Duration = time.Duration(10)

	testResults.Suites = append(testResults.Suites, suite)
	testResults.Aggregate()

	assert.Equal(t, testResults.Summary, Summary{15, 3, 6, 4, 3, 11})
}

func TestNewSuite(t *testing.T) {
	suite := NewSuite()

	assert.IsType(t, suite, Suite{})
}

func TestSuite_Aggregate(t *testing.T) {
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
}

func TestNewTest(t *testing.T) {
	test := NewTest()

	assert.Equal(t, test.State, StatePassed, "is in passed state by default")
	assert.IsType(t, test, Test{})
}

func TestNewError(t *testing.T) {
	obj := NewError()

	assert.IsType(t, obj, Error{})
}

func TestNewFailure(t *testing.T) {
	obj := NewFailure()

	assert.IsType(t, obj, Failure{})
}
