package parsers

import (
	"strings"

	types "github.com/semaphoreci/test-results/pkg/types"
)

// Mocha ...
type Mocha struct {
	Parser
}

// NewSuites interface for generic parser
func (me Mocha) NewSuites(testsuites types.XMLTestSuites) types.Suites {
	var parser Generic
	suites := parser.NewSuites(testsuites)

	suites.Name = me.Name()

	return suites
}

// NewSuite ...
func (me Mocha) NewSuite(suites types.Suites, testsuite types.XMLTestSuite) types.Suite {
	var parser Generic
	return parser.NewSuite(suites, testsuite)
}

// NewTest ...
func (me Mocha) NewTest(suites types.Suites, suite types.Suite, testcase types.XMLTestCase) types.Test {
	var parser Generic
	return parser.NewTest(suites, suite, testcase)
}

// Name interface for generic parser
func (me Mocha) Name() string {
	return "Mocha"
}

// Applicable interface for generic parser
func (me Mocha) Applicable(xmltestsuites types.XMLTestSuites) bool {
	for _, xmlTestSuite := range xmltestsuites.TestSuites {
		if strings.Contains(strings.ToLower(xmlTestSuite.Name), "mocha") {
			return true
		}
	}
	return false
}
