package parser

import (
	"fmt"
)

// Parser ...
type Parser interface {
	Parse(string) TestResults
	IsApplicable(string) bool
	GetName() string
}

// FindParser ...
func FindParser(name string, path string, availableParsers []Parser) (Parser, error) {
	for _, p := range availableParsers {
		if p.GetName() == name {
			return p, nil
		}
	}

	for _, p := range availableParsers {
		if p.IsApplicable(path) {
			return p, nil
		}
	}

	return nil, fmt.Errorf("Parser not found")
}
