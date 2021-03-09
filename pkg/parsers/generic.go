package parsers

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/semaphoreci/test-results/pkg/fileloader"
	"github.com/semaphoreci/test-results/pkg/parser"
)

// Generic ...
type Generic struct {
	parser.Parser
}

// NewGeneric ...
func NewGeneric() Generic {
	return Generic{}
}

// Parse ...
func (p Generic) Parse(path string) (*parser.TestResults, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	bytes := bytes.NewReader(file)

	reader, err := fileloader.Load(path, *bytes)
	if err != nil {
		return nil, err
	}

	xmlElement := parser.NewXMLElement()

	err = xmlElement.Parse(reader)
	if err != nil {
		return nil, err
	}

	var results parser.TestResults

	switch xmlElement.Tag() {
	case "testsuites":
		results = newTestResults(xmlElement)
	case "testsuite":
		results = parser.NewTestResults()
		results.Name = "Generic Parser"
		results.Suites = []parser.Suite{newSuite(xmlElement)}
	}

	results.Aggregate()

	return &results, nil
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
		return 0
	}

	return d
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func parseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return b
}
