package parsers

import (
	"strings"

	types "github.com/semaphoreci/test-results/pkg/types"
)

// ExUnit ...
type ExUnit struct {
	Parser
}

// NewSuites interface for generic parser
func (me ExUnit) NewSuites(testsuites types.XMLTestSuites) types.Suites {
	var parser Generic
	suites := parser.NewSuites(testsuites)

	suites.Name = me.Name()

	return suites
}

// NewSuite ...
func (me ExUnit) NewSuite(suites types.Suites, testsuite types.XMLTestSuite) types.Suite {
	var parser Generic
	return parser.NewSuite(suites, testsuite)
}

// NewTest ...
func (me ExUnit) NewTest(suites types.Suites, suite types.Suite, testcase types.XMLTestCase) types.Test {
	var parser Generic
	test := parser.NewTest(suites, suite, testcase)

	test.Classname = testcase.Classname
	test.Package = testcase.Classname

	return test
}

// Name interface for generic parser
func (me ExUnit) Name() string {
	return "ExUnit"
}

// Applicable interface for generic parser
func (me ExUnit) Applicable(xmltestsuites types.XMLTestSuites) bool {
	for _, xmlTestSuite := range xmltestsuites.TestSuites {
		if strings.Contains(strings.ToLower(xmlTestSuite.Name), "elixir") {
			return true
		}
	}
	return false
}
