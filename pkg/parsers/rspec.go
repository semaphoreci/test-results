package parsers

import (
	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/semaphoreci/test-results/pkg/parser"
)

// RSpec ...
type RSpec struct {
	logFields logger.Fields
}

// NewRSpec ...
func NewRSpec() RSpec {
	return RSpec{logger.Fields{"app": "parser.rspec"}}
}

// GetName ...
func (me RSpec) GetName() string {
	return "rspec"
}

// IsApplicable ...
func (me RSpec) IsApplicable(path string) bool {
	return true
}

// Parse ...
func (me RSpec) Parse(path string) parser.TestResults {
	return parser.NewTestResults()
}
