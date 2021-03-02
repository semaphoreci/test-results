package parsers

import (
	"strings"

	types "github.com/semaphoreci/test-results/pkg/types"
)

// ParserMocha ...
type ParserMocha struct {
	Parser
}

// NewSuites interface for generic parser
func (me ParserMocha) NewSuites(testsuites types.XMLTestSuites) types.Suites {
	var parser ParserGeneric
	suites := parser.NewSuites(testsuites)

	suites.Name = me.Name()

	return suites
}

// NewSuite ...
func (me ParserMocha) NewSuite(suites types.Suites, testsuite types.XMLTestSuite) types.Suite {
	var parser ParserGeneric
	return parser.NewSuite(suites, testsuite)
}

// NewTest ...
func (me ParserMocha) NewTest(suites types.Suites, suite types.Suite, testcase types.XMLTestCase) types.Test {
	var parser ParserGeneric
	return parser.NewTest(suites, suite, testcase)
}

// Name interface for generic parser
func (me ParserMocha) Name() string {
	return "Mocha"
}

// Applicable interface for generic parser
func (me ParserMocha) Applicable(xmltestsuites types.XMLTestSuites) bool {
	for _, xmlTestSuite := range xmltestsuites.TestSuites {
		if strings.Contains(strings.ToLower(xmlTestSuite.Name), "mocha") {
			return true
		}
	}
	return false
}