package parsers

import (
	"fmt"
	"time"

	types "github.com/semaphoreci/test-results/pkg/types"
)

// ParserInt ...
type ParserInt interface {
	Parse(types.XMLTestSuites) types.Suites
	ParseSuites(types.XMLTestSuites) types.Suites
	ParseSuite(types.Suites, types.XMLTestSuite) types.Suite
	ParseTest(types.Suites, types.Suite, types.XMLTestCase) types.Test
	Applicable(types.XMLTestSuites) bool
}

// Parser ...
type Parser struct {
	ParserInt
	Name string
}

// Parse interface for generic parser
func (me *Parser) Parse(testsuites types.XMLTestSuites) types.Suites {
	suites := me.ParseSuites(testsuites)

	suites.Aggregate()
	return suites
}

// ParseSuites ...
func (me *Parser) ParseSuites(testsuites types.XMLTestSuites) types.Suites {
	suites := types.Suites{ID: testsuites.ID, Name: me.Name}

	for _, testsuite := range testsuites.TestSuites {
		suites.Suites = append(suites.Suites, me.ParseSuite(suites, testsuite))
	}

	return suites
}

// ParseSuite ...
func (me *Parser) ParseSuite(suites types.Suites, testsuite types.XMLTestSuite) types.Suite {
	suite := types.Suite{Name: testsuite.Name, ID: testsuite.Name}

	for _, testcase := range testsuite.TestsCases {
		suite.Tests = append(suite.Tests, me.ParseTest(suites, suite, testcase))
	}

	suite.Aggregate()

	return suite
}

// ParseTest ...
func (me *Parser) ParseTest(suites types.Suites, suite types.Suite, testcase types.XMLTestCase) types.Test {
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
// func (me *Parser) Name() string {
// 	return me.name
// }

// Applicable interface for generic parser
func (me *Parser) Applicable(types.XMLTestSuites) bool {
	return true
}