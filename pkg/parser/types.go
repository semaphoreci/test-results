package parser

import (
	"fmt"
	"os"
	"sort"
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

// Result ...
type Result struct {
	TestResults []TestResults `json:"testResults"`
}

// NewResult ...
func NewResult() Result {
	return Result{
		TestResults: []TestResults{},
	}
}

// Combine test results that are part of result
func (me *Result) Combine(other Result) {
	for i := range other.TestResults {
		testResult := other.TestResults[i]
		testResult.Flatten()
		foundTestResultsIdx, found := me.hasTestResults(testResult)
		if found {
			me.TestResults[foundTestResultsIdx].Combine(testResult)
			me.TestResults[foundTestResultsIdx].Aggregate()
		} else {
			me.TestResults = append(me.TestResults, testResult)
		}
	}

	sort.SliceStable(me.TestResults, func(i, j int) bool { return me.TestResults[i].ID < me.TestResults[j].ID })

	for i := range me.TestResults {
		me.TestResults[i].Aggregate()
	}
}

// Flatten makes sure we don't have duplicated suites in test results
func (me *TestResults) Flatten() {
	testResults := NewTestResults()

	for i := range me.Suites {
		foundSuiteIdx, found := testResults.hasSuite(me.Suites[i])
		if found {
			testResults.Suites[foundSuiteIdx].Combine(me.Suites[i])
			testResults.Suites[foundSuiteIdx].Aggregate()
		} else {
			testResults.Suites = append(testResults.Suites, me.Suites[i])
		}
	}
	me.Suites = testResults.Suites

	sort.SliceStable(me.Suites, func(i, j int) bool {
		return me.Suites[i].ID < me.Suites[j].ID
	})
}

func (me *Result) hasTestResults(testResults TestResults) (int, bool) {
	for i := range me.TestResults {
		if me.TestResults[i].ID == testResults.ID {
			return i, true
		}
	}
	return -1, false
}

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

// Combine ...
func (me *TestResults) Combine(other TestResults) {
	if me.ID == other.ID {
		for i := range other.Suites {
			foundSuiteIdx, found := me.hasSuite(other.Suites[i])
			if found {
				me.Suites[foundSuiteIdx].Combine(other.Suites[i])
				me.Suites[foundSuiteIdx].Aggregate()
			} else {
				me.Suites = append(me.Suites, other.Suites[i])
			}
		}

		sort.SliceStable(me.Suites, func(i, j int) bool {
			return me.Suites[i].ID < me.Suites[j].ID
		})
	}
}

func (me *TestResults) hasSuite(suite Suite) (int, bool) {
	for i := range me.Suites {
		if me.Suites[i].ID == suite.ID {
			return i, true
		}
	}
	return -1, false
}

// ArrangeSuitesByTestFile ...
func (me *TestResults) ArrangeSuitesByTestFile() {
	newSuites := []Suite{}

	for _, suite := range me.Suites {
		for _, test := range suite.Tests {
			var (
				idx        int
				foundSuite *Suite
			)
			if test.File != "" {
				idx, foundSuite = EnsureSuiteByName(newSuites, test.File)
			} else {
				idx, foundSuite = EnsureSuiteByName(newSuites, suite.Name)
			}

			foundSuite.Tests = append(foundSuite.Tests, test)
			foundSuite.Aggregate()

			if idx == -1 {
				foundSuite.EnsureID(*me)
				newSuites = append(newSuites, *foundSuite)
			}
		}
	}

	me.Suites = newSuites
	me.Aggregate()
}

// EnsureSuiteByName ...
func EnsureSuiteByName(suites []Suite, name string) (int, *Suite) {
	for i := range suites {
		if suites[i].Name == name {
			return i, &suites[i]
		}
	}
	suite := NewSuite()
	suite.Name = name

	return -1, &suite
}

// EnsureID ...
func (me *TestResults) EnsureID() {
	if me.ID == "" {
		me.ID = me.Name
	}

	if me.Framework != "" {
		me.ID = fmt.Sprintf("%s%s", me.ID, me.Framework)
	}

	me.ID = UUID(uuid.Nil, me.ID).String()
}

// RegenerateID ...
func (me *TestResults) RegenerateID() {
	me.ID = ""
	me.EnsureID()
	for suiteIdx := range me.Suites {
		me.Suites[suiteIdx].ID = ""
		me.Suites[suiteIdx].EnsureID(*me)
		for testIdx := range me.Suites[suiteIdx].Tests {
			me.Suites[suiteIdx].Tests[testIdx].ID = ""
			me.Suites[suiteIdx].Tests[testIdx].EnsureID(me.Suites[suiteIdx])
		}
	}
}

// Aggregate all test suite summaries
func (me *TestResults) Aggregate() {
	summary := Summary{}

	for i := range me.Suites {
		summary.Duration += me.Suites[i].Summary.Duration
		summary.Skipped += me.Suites[i].Summary.Skipped
		summary.Error += me.Suites[i].Summary.Error
		summary.Total += me.Suites[i].Summary.Total
		summary.Failed += me.Suites[i].Summary.Failed
		summary.Passed += me.Suites[i].Summary.Passed
		summary.Disabled += me.Suites[i].Summary.Disabled
	}

	me.Summary = summary
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
	return Suite{Tests: []Test{}}
}

// Combine ...
func (me *Suite) Combine(other Suite) {
	if me.ID == other.ID {
		for i := range other.Tests {
			if !me.hasTest(other.Tests[i]) {
				me.Tests = append(me.Tests, other.Tests[i])
			}

			shouldReplace, indexToReplace := me.shouldReplaceTest(other.Tests[i])

			if shouldReplace && indexToReplace != -1 {
				me.Tests[indexToReplace] = other.Tests[i]
			}

		}

		sort.SliceStable(me.Tests, func(i, j int) bool {
			return me.Tests[i].ID < me.Tests[j].ID
		})
	}
}
func (me *Suite) shouldReplaceTest(test Test) (shouldReplace bool, foundIndex int) {
	foundIndex = -1
	shouldReplace = false
	for i := range me.Tests {
		if me.Tests[i].ID == test.ID {
			foundIndex = i
			break
		}
	}

	if foundIndex == -1 {
		return
	} else {
		foundTest := me.Tests[foundIndex]

		if foundTest.State == StateSkipped {
			shouldReplace = true
			return
		}
		if foundTest.State == StatePassed && test.State == StateFailed || test.State == StateError {
			shouldReplace = true
			return
		}

		return
	}
}
func (me *Suite) hasTest(test Test) bool {
	for i := range me.Tests {
		if me.Tests[i].ID == test.ID {
			return true
		}
	}
	return false
}

// Aggregate all tests in suite
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

	// If current duration is not zero and current duration is bigger than calculated duration, use it
	if me.Summary.Duration > 0 && me.Summary.Duration > summary.Duration {
		summary.Duration = me.Summary.Duration
	}

	me.Summary = summary
}

// EnsureID ...
func (me *Suite) EnsureID(tr TestResults) {
	if me.ID == "" {
		me.ID = me.Name
	}

	oldID, err := uuid.Parse(tr.ID)
	if err != nil {
		oldID = uuid.Nil
	}

	me.ID = UUID(oldID, me.ID).String()
}

// AppendTest ...
func (me *Suite) AppendTest(test Test) {
	me.Tests = append(me.Tests, test)
	me.Aggregate()
}

type SemEnv struct {
	IP         string `json:"ip"`
	PipelineId string `json:"pipeline_id"`
	WorkflowId string `json:"workflow_id"`

	JobName string `json:"name"`
	JobId   string `json:"id"`

	AgentType    string `json:"agent_type"`
	AgentOsImage string `json:"agent_os_image"`

	GitRefType string `json:"git_ref_type"`
	GitRefName string `json:"git_ref_name"`
	GitRefSha  string `json:"ref_sha"`
}

func NewSemEnv() *SemEnv {
	refName := ""
	refSha := ""
	switch os.Getenv("SEMAPHORE_GIT_REF_TYPE") {
	case "branch":
		refName = os.Getenv("SEMAPHORE_GIT_BRANCH")
		refSha = os.Getenv("SEMAPHORE_GIT_SHA")
	case "tag":
		refName = os.Getenv("SEMAPHORE_GIT_BRANCH")
		refSha = os.Getenv("SEMAPHORE_GIT_SHA")
	case "pull-request":
		refName = os.Getenv("SEMAPHORE_GIT_PR_BRANCH")
		refSha = os.Getenv("SEMAPHORE_GIT_PR_SHA")
	}

	return &SemEnv{
		IP:           os.Getenv("IP"),
		PipelineId:   os.Getenv("SEMAPHORE_PIPELINE_ID"),
		WorkflowId:   os.Getenv("SEMAPHORE_WORKFLOW_ID"),
		JobName:      os.Getenv("SEMAPHORE_JOB_NAME"),
		JobId:        os.Getenv("SEMAPHORE_JOB_ID"),
		AgentType:    os.Getenv("SEMAPHORE_AGENT_MACHINE_TYPE"),
		AgentOsImage: os.Getenv("SEMAPHORE_AGENT_MACHINE_OS_IMAGE"),
		GitRefType:   os.Getenv("SEMAPHORE_GIT_REF_TYPE"),
		GitRefName:   refName,
		GitRefSha:    refSha,
	}
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
	SemEnv    *SemEnv       `json:"semaphore_env"`
}

// NewTest ...
func NewTest() Test {
	return Test{
		State:  StatePassed,
		SemEnv: NewSemEnv(),
	}
}

// EnsureID ...
func (me *Test) EnsureID(s Suite) {
	if me.ID == "" {
		me.ID = me.Name
	}

	if me.Classname != "" {
		me.ID = fmt.Sprintf("%s.%s", me.Classname, me.ID)
	}

	if me.Failure != nil {
		me.ID = fmt.Sprintf("%s.%s", "Failure", me.ID)
	}

	if me.Error != nil {
		me.ID = fmt.Sprintf("%s.%s", "Error", me.ID)
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

//Merge merges two summaries together summing each field
func (s *Summary) Merge(withSummary *Summary) {
	s.Total += withSummary.Total
	s.Passed += withSummary.Passed
	s.Skipped += withSummary.Skipped
	s.Error += withSummary.Error
	s.Failed += withSummary.Failed
	s.Disabled += withSummary.Disabled
	s.Duration += withSummary.Duration
}

// UUID ...
func UUID(id uuid.UUID, str string) uuid.UUID {
	return uuid.NewMD5(id, []byte(str))
}
