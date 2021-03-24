package parsers

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/semaphoreci/test-results/pkg/fileloader"
	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/semaphoreci/test-results/pkg/parser"
)

// Generic ...
type Generic struct {
	logFields logger.Fields
}

// NewGeneric ...
func NewGeneric() Generic {
	return Generic{logger.Fields{"app": "parser.generic"}}
}

// IsApplicable ...
func (me Generic) IsApplicable(path string) bool {
	return true
}

// GetName ...
func (me Generic) GetName() string {
	return "generic"
}

// Parse ...
func (me Generic) Parse(path string) parser.TestResults {
	me.logFields["path"] = path

	var results parser.TestResults
	var reader *bytes.Reader
	// Preload path with loader. If nothing is found in file cache - load it up from path.
	reader, found := fileloader.Load(path, nil)

	if found == false {
		file, err := ioutil.ReadFile(path)

		if err != nil {
			logger.Error(me.logFields, "Reading file failed: %v", err)
			// TODO: Add status with reading file failure
			return results
		}

		b := bytes.NewReader(file)
		reader, _ = fileloader.Load(path, b)
	}

	xmlElement := parser.NewXMLElement()

	err := xmlElement.Parse(reader)
	if err != nil {
		logger.Error(me.logFields, "Parsing XML failed: %v", err)
		// TODO: Add status with parsing XML failure
		return results
	}

	switch xmlElement.Tag() {
	case "testsuites":
		logger.Debug(me.logFields, "Root <testsuites> element found")
		results = newTestResults(xmlElement)
	case "testsuite":
		logger.Debug(me.logFields, "No root <testsuites> element found")
		results = parser.NewTestResults()
		results.Name = "Generic Parser"
		results.Suites = []parser.Suite{newSuite(xmlElement)}
	}

	results.Aggregate()

	return results
}

func newTestResults(xml parser.XMLElement) parser.TestResults {
	testResults := parser.NewTestResults()

	for _, node := range xml.Children {
		switch node.Tag() {
		case "testsuite":
			testResults.Suites = append(testResults.Suites, newSuite(node))
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

func newSuites(elements []parser.XMLElement) []parser.Suite {
	var suites []parser.Suite

	for _, element := range elements {
		suites = append(suites, newSuite(element))
	}

	return suites
}

func newSuite(xml parser.XMLElement) parser.Suite {
	suite := parser.NewSuite()

	for _, node := range xml.Children {
		switch node.Tag() {
		case "properties":
			properties := newProperties(node)
			suite.Properties = properties
		case "system-out":
			suite.SystemOut = string(node.Contents)
		case "system-err":
			suite.SystemErr = string(node.Contents)
		case "testcase":
			suite.Tests = append(suite.Tests, newTest(node))
		}
	}

	for attr, value := range xml.Attributes {
		switch attr {
		case "name":
			suite.Name = value
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

func newTest(xml parser.XMLElement) parser.Test {
	test := parser.NewTest()

	for _, node := range xml.Children {
		switch node.Tag() {
		case "failure":
			test.State = parser.StateFailed
			test.Failure = parseFailure(node)
		case "error":
			test.State = parser.StateError
			test.Error = parseError(node)
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

func parseFailure(xml parser.XMLElement) *parser.Failure {
	failure := parser.NewFailure()

	failure.Body = string(xml.Contents)
	failure.Message = xml.Attr("message")
	failure.Type = xml.Attr("type")

	return &failure
}

func parseError(xml parser.XMLElement) *parser.Error {
	err := parser.NewError()

	err.Body = string(xml.Contents)
	err.Message = xml.Attr("message")
	err.Type = xml.Attr("type")

	return &err
}

func newProperties(xml parser.XMLElement) parser.Properties {
	properties := make(map[string]string)
	for _, node := range xml.Children {
		properties[node.Attr("name")] = node.Attr("value")
	}

	return properties
}

func parseTime(s string) time.Duration {
	// append 's' to end of input to use `time` built in duration parser
	d, err := time.ParseDuration(s + "s")
	if err != nil {
		logger.Warn(logger.Fields{}, "Duration parsing failed: %v", err)
		return 0
	}

	return d
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		logger.Warn(logger.Fields{}, "Integer parsing failed: %v", err)
		return 0
	}
	return i
}

func parseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		logger.Warn(logger.Fields{}, "Boolean parsing failed: %v", err)
		return false
	}
	return b
}
