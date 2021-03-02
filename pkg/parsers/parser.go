package parsers

import (
	types "github.com/semaphoreci/test-results/pkg/types"
)

// Parser ...
type Parser interface {
	NewSuites(types.XMLTestSuites) types.Suites
	NewSuite(types.Suites, types.XMLTestSuite) types.Suite
	NewTest(types.Suites, types.Suite, types.XMLTestCase) types.Test
	Applicable(types.XMLTestSuites) bool
	Name() string
}

