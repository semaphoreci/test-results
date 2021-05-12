package parser

import (
	"time"

	"github.com/google/uuid"
)

// Properties maps additional parameters for test suites
type Properties map[string]string

// State indicates state of specific test
type State string

const (
	// StatePassed indicates that test was successful
	StatePassed State = "passed"
	// StateError indicates that test errored due to unexpected behaviour when running test i.e. exception
	StateError State = "error"
	// StateFailed indicates that test failed due to invalid test result
	StateFailed State = "failed"
	// StateSkipped indicates that test was skipped
	StateSkipped State = "skipped"
	// StateDisabled indicates that test was disabled
	StateDisabled State = "disabled"
)

// Status stores information about parsing results
type Status string

const (
	// StatusSuccess indicates that parsing was successful
	StatusSuccess Status = "success"

	// StatusError indicates that parsing failed due to error
	StatusError Status = "error"
)

// TestResults ...
type TestResults struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Framework     string  `json:"framework"`
	IsDisabled    bool    `json:"isDisabled"`
	Suites        []Suite `json:"suites"`
	Summary       Summary `json:"summary"`
	Status        Status  `json:"status"`
	StatusMessage string  `json:"statusMessage"`
}

// NewTestResults ...
func NewTestResults() TestResults {
	return TestResults{
		Suites:        []Suite{},
		Status:        StatusSuccess,
		StatusMessage: "",
	}
}

// EnsureID ...
func (me *TestResults) EnsureID() {
	if me.ID == "" {
		me.ID = me.Name
	}

	me.ID = UUID(uuid.Nil, me.ID).String()
}

// Aggregate all test suite summaries
func (me *TestResults) Aggregate() {
	summary := Summary{}

	for _, suite := range me.Suites {
		summary.Duration += suite.Summary.Duration
		summary.Skipped += suite.Summary.Skipped
		summary.Error += suite.Summary.Error
		summary.Total += suite.Summary.Total
		summary.Failed += suite.Summary.Failed
		summary.Passed += suite.Summary.Passed
		summary.Disabled += suite.Summary.Disabled
	}

	(*me).Summary = summary
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
func (me *Suite) Aggregate() {
	summary := Summary{}

	for _, test := range me.Tests {
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
		case StateDisabled:
			summary.Disabled++
		}
	}

	(*me).Summary = summary
}

// EnsureID ...
func (me *Suite) EnsureID(tr TestResults) {
	if me.ID == "" {
		me.ID = me.Name
	}

	me.ID = UUID(uuid.MustParse(tr.ID), me.ID).String()
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

// EnsureID ...
func (me *Test) EnsureID(s Suite) {
	if me.ID == "" {
		me.ID = me.Name
	}

	me.ID = UUID(uuid.MustParse(s.ID), me.ID).String()
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
	Disabled int           `json:"disabled"`
	Duration time.Duration `json:"duration"`
}

// UUID ...
func UUID(id uuid.UUID, str string) uuid.UUID {
	return uuid.NewMD5(id, []byte(str))
}
