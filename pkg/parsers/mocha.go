package parsers

import (
	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/semaphoreci/test-results/pkg/parser"
)

// Mocha ...
type Mocha struct {
	logFields logger.Fields
}

// NewMocha ...
func NewMocha() Mocha {
	return Mocha{logger.Fields{"app": "parser.mocha"}}
}

// GetName ...
func (me Mocha) GetName() string {
	return "mocha"
}

// IsApplicable ...
func (me Mocha) IsApplicable(path string) bool {
	return true
}

// Parse ...
func (me Mocha) Parse(path string) parser.TestResults {
	return parser.NewTestResults()
}
