package parsers

import (
	"strings"

	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/semaphoreci/test-results/pkg/parser"
)

// ExUnit ...
type ExUnit struct {
}

// NewExUnit ...
func NewExUnit() ExUnit {
	return ExUnit{}
}

// GetName ...
func (me ExUnit) GetName() string {
	return "exunit"
}

// IsApplicable ...
func (me ExUnit) IsApplicable(path string) bool {
	xmlElement, err := LoadXML(path)
	logger.Trace("Checking applicability")

	if err != nil {
		logger.Error("Loading XML failed: %v", err)
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
	var results parser.TestResults

	xmlElement, err := LoadXML(path)

	if err != nil {
		logger.Error("Loading XML failed: %v", err)
		return results
	}

	switch xmlElement.Tag() {
	case "testsuites":
		logger.Debug("Root <testsuites> element found")
		results = me.newTestResults(*xmlElement)
	case "testsuite":
		logger.Debug("No root <testsuites> element found")
		results = parser.NewTestResults()
		results.Name = "ExUnit Parser"
		results.Suites = append(results.Suites, me.newSuite(*xmlElement))
	}

	results.Aggregate()

	return results
}

func (me ExUnit) newTestResults(xml parser.XMLElement) parser.TestResults {
	testResults := parser.NewTestResults()

	for _, node := range xml.Children {
		switch node.Tag() {
		case "testsuite":
			testResults.Suites = append(testResults.Suites, me.newSuite(node))
		}
	}

	for attr, value := range xml.Attributes {
		switch attr {
		case "name":
			testResults.Name = value
		case "time":
			testResults.Summary.Duration = parseTime(value)
		case "tests":
			testResults.Summary.Total = parseInt(value)
		case "failures":
			testResults.Summary.Failed = parseInt(value)
		case "errors":
			testResults.Summary.Error = parseInt(value)
		case "disabled":
			testResults.IsDisabled = parseBool(value)
		}
	}
	testResults.Summary.Passed = testResults.Summary.Total - testResults.Summary.Error - testResults.Summary.Failed

	return testResults
}

func (me ExUnit) newSuite(xml parser.XMLElement) parser.Suite {
	suite := parser.NewSuite()

	for _, node := range xml.Children {
		switch node.Tag() {
		case "properties":
			suite.Properties = me.parseProperties(node)
		case "system-out":
			suite.SystemOut = string(node.Contents)
		case "system-err":
			suite.SystemErr = string(node.Contents)
		case "testcase":
			suite.Tests = append(suite.Tests, me.newTest(node))
		}
	}

	for attr, value := range xml.Attributes {
		switch attr {
		case "name":
			suite.Name = strings.Trim(value, "Elixir.")
		case "tests":
			suite.Summary.Total = parseInt(value)
		case "failures":
			suite.Summary.Failed = parseInt(value)
		case "errors":
			suite.Summary.Error = parseInt(value)
		case "time":
			suite.Summary.Duration = parseTime(value)
		case "disabled":
			suite.IsDisabled = parseBool(value)
		case "skipped":
			suite.IsSkipped = parseBool(value)
		case "timestamp":
			suite.Timestamp = value
		case "hostname":
			suite.Hostname = value
		case "id":
			suite.ID = value
		case "package":
			suite.Package = value
		}
	}

	suite.Aggregate()

	return suite
}

func (me ExUnit) newTest(xml parser.XMLElement) parser.Test {
	test := parser.NewTest()

	for _, node := range xml.Children {
		switch node.Tag() {
		case "failure":
			test.State = parser.StateFailed
			test.Failure = me.parseFailure(node)
		case "error":
			test.State = parser.StateError
			test.Error = me.parseError(node)
		case "skipped":
			test.State = parser.StateSkipped
		case "system-out":
			test.SystemOut = string(node.Contents)
		case "system-err":
			test.SystemErr = string(node.Contents)
		}
	}

	for attr, value := range xml.Attributes {
		switch attr {
		case "name":
			test.Name = value
		case "time":
			test.Duration = parseTime(value)
		case "classname":
			test.Classname = value
		}
	}

	return test
}

func (me ExUnit) parseProperties(xml parser.XMLElement) parser.Properties {
	properties := make(map[string]string)
	for _, node := range xml.Children {
		properties[node.Attr("name")] = node.Attr("value")
	}

	return properties
}

func (me ExUnit) parseFailure(xml parser.XMLElement) *parser.Failure {
	failure := parser.NewFailure()

	failure.Body = string(xml.Contents)
	failure.Message = xml.Attr("message")
	failure.Type = xml.Attr("type")

	return &failure
}

func (me ExUnit) parseError(xml parser.XMLElement) *parser.Error {
	err := parser.NewError()

	err.Body = string(xml.Contents)
	err.Message = xml.Attr("message")
	err.Type = xml.Attr("type")

	return &err
}
