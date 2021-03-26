package parsers

import (
	"strings"

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
	me.logFields["fun"] = "IsApplicable"
	xmlElement, err := LoadXML(path)
	logger.Trace(me.logFields, "Checking applicability")

	if err != nil {
		logger.Error(me.logFields, "Loading XML failed: %v", err)
		return false
	}

	switch xmlElement.Tag() {
	case "testsuites":
		for attr, value := range xmlElement.Attributes {
			logger.Trace(me.logFields, "%s %s", attr, value)
			switch attr {
			case "name":
				if strings.Contains(strings.ToLower(value), "mocha") {
					return true
				}
			}
		}
	}
	return false
}

// Parse ...
func (me Mocha) Parse(path string) parser.TestResults {
	parser := NewGeneric()

	return parser.Parse(path)
}
