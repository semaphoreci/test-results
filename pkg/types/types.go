package types

import (
	"time"
)

// State ...
type State string

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

// Summary ...
type Summary struct {
	Total int `json:"total"`
	Passed int `json:"passed"`
	Skipped int `json:"skipped"`
	Error int `json:"error"`
	Failed int `json:"failed"`
	Duration time.Duration `json:"duration"`
}

// Failure ...
type Failure struct {
	Message string `json:"message"`
	Type string `json:"type"`
	Body string `json:"body"`
}

// Test ...
type Test struct {
	ID string `json:"id"`
	File string `json:"file"`
	Classname string `json:"classname"`
	Package string `json:"package"`
	Name string `json:"name"`
	Duration time.Duration `json:"duration"`
	State State `json:"state"`
	Failure *Failure `json:"failure"`
	SystemOut string `json:"systemOut"`
	SystemErr string `json:"systemErr"`
}

// Suite ...
type Suite struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Tests []Test `json:"tests"`
	Properties map[string]string `json:"properties"`
	Summary Summary `json:"summary"`
}

// Suites ...
type Suites struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Suites []Suite `json:"suites"`
	Summary Summary `json:"summary"`
}

// Aggregate all tests in suite
func (suite *Suite) Aggregate() {
	summary := Summary{}

	for _, test := range suite.Tests {
		summary.Duration += test.Duration
		summary.Total++
		switch(test.State) {
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
	suite.Summary = summary
}

// Aggregate all test suites
func (suites *Suites) Aggregate() {
	summary := Summary{}

	for _, suite := range suites.Suites {
		summary.Duration += suite.Summary.Duration
		summary.Skipped += suite.Summary.Skipped
		summary.Error += suite.Summary.Error
		summary.Total += suite.Summary.Total
		summary.Failed += suite.Summary.Failed
		summary.Passed += suite.Summary.Passed
	}

	suites.Summary = summary
}
