package parsers

import (
	"fmt"
	"time"

	types "github.com/semaphoreci/test-results/pkg/types"
)

// ParserGeneric ...
type ParserGeneric struct {
	Parser
}

// NewSuites interface for generic parser
func (me ParserGeneric) NewSuites(testsuites types.XMLTestSuites) types.Suites {
	suites := types.Suites{ID: testsuites.ID, Name: me.Name()}

	return suites
}

// NewSuite ...
func (me ParserGeneric) NewSuite(suites types.Suites, testsuite types.XMLTestSuite) types.Suite {
	suite := types.Suite{Name: testsuite.Name, ID: testsuite.Name}

	return suite
}

// NewTest ...
func (me ParserGeneric) NewTest(suites types.Suites, suite types.Suite, testcase types.XMLTestCase) types.Test {
	duration, _ := time.ParseDuration(fmt.Sprintf("%fs", testcase.Time))

	state := types.StatePassed
	var failure *types.Failure

	if testcase.Failure != nil {
		state = types.StateFailed
		failure = &types.Failure{
			Message: testcase.Failure.Message,
			Type: testcase.Failure.Type,
			Body: testcase.Failure.Body,
		}
	}

	test := types.Test{
		File: testcase.File,
		Name: testcase.Name,
		Package: suites.Name,
		Duration: duration,
		State: state,
		Failure: failure,
	}

	return test
}

// Name interface for generic parser
func (me ParserGeneric) Name() string {
	return "Generic"
}

// Applicable interface for generic parser
func (me ParserGeneric) Applicable(types.XMLTestSuites) bool {
	return true
}