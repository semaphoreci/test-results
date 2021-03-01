package parsers

import (
	"fmt"

	types "github.com/semaphoreci/test-results/pkg/types"
)

// ParserExUnit ...
type ParserExUnit struct {
	Parser
}

// NewParserExUnit ...
// func NewParserExUnit() Parser {
// 	return &parserExUnit{
// 		Parser: NewParser(),
// 	}
// }

// Applicable ...
func (me *ParserExUnit) Applicable(testsuites types.XMLTestSuites) bool {
	return true
}

// ParseSuites ...
func (me *ParserExUnit) ParseSuites(testsuites types.XMLTestSuites) types.Suites {
	return me.Parser.ParseSuites(testsuites)
}

// ParseSuite ...
func (me *ParserExUnit) ParseSuite(suites types.Suites, testcase types.XMLTestSuite) types.Suite {
	return me.Parser.ParseSuite(suites, testcase)
}

// ParseTest ...
func (me *ParserExUnit) ParseTest(suites types.Suites, suite types.Suite, testcase types.XMLTestCase) types.Test {
	fmt.Printf("TESTS")
	test := me.Parser.ParseTest(suites, suite, testcase)

	test.Classname = "Elixir"

	return test
}