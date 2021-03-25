package parsers

import (
	"strings"

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
	me.logFields["fun"] = "IsApplicable"
	xmlElement, err := LoadXML(path)
	logger.Trace(me.logFields, "Checking applicability")

	if err != nil {
		logger.Error(me.logFields, "Loading XML failed: %v", err)
		return false
	}

	switch xmlElement.Tag() {
	case "testsuites":
		testsuites := xmlElement.Children

		for _, testsuite := range testsuites {
			switch testsuite.Tag() {
			case "testsuite":
				for attr, value := range testsuite.Attributes {
					switch attr {
					case "name":
						if strings.HasPrefix(value, "Elixir.") {
							return true
						}
					}
				}
			}
		}

	case "testsuite":
		for attr, value := range xmlElement.Attributes {
			switch attr {
			case "name":
				if strings.HasPrefix(value, "Elixir.") {
					return true
				}
			}
		}
	}

	return false
}

// Parse ...
func (me ExUnit) Parse(path string) parser.TestResults {
	parser := NewGeneric()
	parser.logFields = me.logFields

	return parser.Parse(path)
}
