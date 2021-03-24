package parsers

import (
	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/semaphoreci/test-results/pkg/parser"
)

// ExUnit ...
type ExUnit struct {
	logFields logger.Fields
}

// NewExUnit ...
func NewExUnit() ExUnit {
	return ExUnit{logger.Fields{"app": "parser.exunit"}}
}

// GetName ...
func (me ExUnit) GetName() string {
	return "exunit"
}

// IsApplicable ...
func (me ExUnit) IsApplicable(path string) bool {
	return true
}

// Parse ...
func (me ExUnit) Parse(path string) parser.TestResults {
	return parser.NewTestResults()
}
