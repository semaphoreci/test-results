package parser

import (
	"time"
)

// State ...
type State string

// Properties ...
type Properties map[string]string

const (
	// StatePassed ...
	StatePassed State = "passed"
	// StateError ...
	StateError State = "error"
	// StateFailed ...
	StateFailed State = "failed"
	// StateSkipped ...
	StateSkipped State = "skipped"
)

// TestResults ...
type TestResults struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Framework  string  `json:"framework"`
	IsDisabled bool    `json:"isDisabled"`
	Suites     []Suite `json:"suites"`
	Summary    Summary `json:"summary"`
}

// NewTestResults ...
func NewTestResults() TestResults {
	return TestResults{}
}

// Aggregate all test suites
func (testResults *TestResults) Aggregate() {
	summary := Summary{}

	for _, suite := range testResults.Suites {
		summary.Duration += suite.Summary.Duration
		summary.Skipped += suite.Summary.Skipped
		summary.Error += suite.Summary.Error
		summary.Total += suite.Summary.Total
		summary.Failed += suite.Summary.Failed
		summary.Passed += suite.Summary.Passed
	}

	(*testResults).Summary = summary
}

// Suite ...
type Suite struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	IsSkipped  bool       `json:"isSkipped"`
	IsDisabled bool       `json:"isDisabled"`
	Timestamp  string     `json:"timestamp"`
	Hostname   string     `json:"hostname"`
	Package    string     `json:"package"`
	Tests      []Test     `json:"tests"`
	Properties Properties `json:"properties"`
	Summary    Summary    `json:"summary"`
	SystemOut  string     `json:"systemOut"`
	SystemErr  string     `json:"systemErr"`
}

// NewSuite ...
func NewSuite() Suite {
	return Suite{}
}

// Aggregate all tests in suite
// TODO: add flag to skip aggregating already present data
func (suite *Suite) Aggregate() {
	summary := Summary{}

	for _, test := range suite.Tests {
		summary.Duration += test.Duration
		summary.Total++
		switch test.State {
		case StateSkipped:
			summary.Skipped++
		case StateFailed:
			summary.Failed++
		case StateError:
			summary.Error++
		case StatePassed:
			summary.Passed++
		}
	}

	(*suite).Summary = summary
}

// Test ...
type Test struct {
	ID        string        `json:"id"`
	File      string        `json:"file"`
	Classname string        `json:"classname"`
	Package   string        `json:"package"`
	Name      string        `json:"name"`
	Duration  time.Duration `json:"duration"`
	State     State         `json:"state"`
	Failure   *Failure      `json:"failure"`
	Error     *Error        `json:"error"`
	SystemOut string        `json:"systemOut"`
	SystemErr string        `json:"systemErr"`
}

// NewTest ...
func NewTest() Test {
	return Test{
		State: StatePassed,
	}
}

type err struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Body    string `json:"body"`
}

// Failure ...
type Failure err

// NewFailure ...
func NewFailure() Failure {
	return Failure{}
}

// Error ...
type Error err

// NewError ...
func NewError() Error {
	return Error{}
}

// Summary ...
type Summary struct {
	Total    int           `json:"total"`
	Passed   int           `json:"passed"`
	Skipped  int           `json:"skipped"`
	Error    int           `json:"error"`
	Failed   int           `json:"failed"`
	Duration time.Duration `json:"duration"`
}
