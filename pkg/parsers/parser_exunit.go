package parsers

import (
	"strings"

	types "github.com/semaphoreci/test-results/pkg/types"
)

// ParserExUnit ...
type ParserExUnit struct {
	Parser
}

// NewSuites interface for generic parser
func (me ParserExUnit) NewSuites(testsuites types.XMLTestSuites) types.Suites {
	var parser ParserGeneric
	suites := parser.NewSuites(testsuites)

	suites.Name = me.Name()

	return suites
}

// NewSuite ...
func (me ParserExUnit) NewSuite(suites types.Suites, testsuite types.XMLTestSuite) types.Suite {
	var parser ParserGeneric
	return parser.NewSuite(suites, testsuite)
}

// NewTest ...
func (me ParserExUnit) NewTest(suites types.Suites, suite types.Suite, testcase types.XMLTestCase) types.Test {
	var parser ParserGeneric
	test := parser.NewTest(suites, suite, testcase)

	test.Classname = testcase.Classname
	test.Package = testcase.Classname

	return test
}

// Name interface for generic parser
func (me ParserExUnit) Name() string {
	return "ExUnit"
}

// Applicable interface for generic parser
func (me ParserExUnit) Applicable(xmltestsuites types.XMLTestSuites) bool {
	for _, xmlTestSuite := range xmltestsuites.TestSuites {
		if strings.Contains(strings.ToLower(xmlTestSuite.Name), "elixir") {
			return true
		}
	}
	return false
}