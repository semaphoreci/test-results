package parsers

import (
	"fmt"

	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/semaphoreci/test-results/pkg/parser"
)

var availableParsers = []parser.Parser{
	NewRSpec(),
	NewExUnit(),
	NewMocha(),
	NewGeneric(),
}

// FindParser ...
func FindParser(name string, path string) (parser.Parser, error) {
	fields := logger.Fields{"app": "parser", "name": name, "path": path}

	if name != "auto" {
		logger.Info(fields, "Looking for parser")
		for _, p := range availableParsers {
			if p.GetName() == name {
				logger.Info(fields, "Found parser")
				return p, nil
			}
		}
		logger.Info(fields, "Parser not found")
	}

	for _, p := range availableParsers {
		isApplicable := p.IsApplicable(path)
		logger.Debug(fields, "Looking for applicable parser, checking %s -> %b", p.GetName(), isApplicable)
		if isApplicable {
			logger.Trace(fields, "Found applicable parser: %s", p.GetName())
			return p, nil
		}
	}
	logger.Error(fields, "No applicable parsers found")

	return nil, fmt.Errorf("Parser not found")
}
