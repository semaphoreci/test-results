package parser

import (
	"fmt"

	"github.com/semaphoreci/test-results/pkg/logger"
)

// Parser ...
type Parser interface {
	Parse(string) TestResults
	IsApplicable(string) bool
	GetName() string
}

// FindParser ...
func FindParser(name string, path string, availableParsers []Parser) (Parser, error) {
	logger.Log("parser", "Looking for parser %s", name)
	for _, p := range availableParsers {
		if p.GetName() == name {
			logger.Log("parser", "Using parser %s", name)
			return p, nil
		}
	}

	for _, p := range availableParsers {
		if p.IsApplicable(path) {
			logger.Log("parser", "Using applicable parser %s", p.GetName())
			return p, nil
		}
	}

	logger.Error("parser", "No applicable parsers found")

	return nil, fmt.Errorf("Parser not found")
}
