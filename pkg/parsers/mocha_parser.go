package parsers

// MochaParser ...
type MochaParser struct {
	Parser
}

// // Parse interface for mocha parser
// func (parser MochaParser) Parse(testsuites types.XMLTestSuites) types.Suites {

// 	suites := types.Suites{ID: testsuites.ID}

// 	for _, xmlTestSuite := range testsuites.TestSuites {
// 	}
// 	suites.Aggregate()

// 	return suites
// }

// // Applicable interface for mocha parser
// func (parser MochaParser) Applicable(types.XMLTestSuites) bool {
// 	return true
// }

// // Name interface for mocha parser
// func (parser MochaParser) Name() string {
// 	return "mocha"
// }