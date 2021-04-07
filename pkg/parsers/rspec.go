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
	xmlElement, err := LoadXML(path)
	logger.Debug("Checking applicability of %s parser", me.GetName())

	if err != nil {
		logger.Error("Loading XML failed: %v", err)
		return false
	}

	switch xmlElement.Tag() {
	case "testsuite":
		for attr, value := range xmlElement.Attributes {
			switch attr {
			case "name":
				if value == "rspec" {
					return true
				}
			}
		}
	}
	return false
}

// Parse ...
func (me RSpec) Parse(path string) parser.TestResults {
	parser := NewGeneric()

	return parser.Parse(path)
}
