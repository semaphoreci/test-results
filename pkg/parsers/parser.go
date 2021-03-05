package parsers

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"os"

	types "github.com/semaphoreci/test-results/pkg/types"
)

// Parser ...
type Parser interface {
	NewSuites(types.XMLTestSuites) types.Suites
	NewSuite(types.Suites, types.XMLTestSuite) types.Suite
	NewTest(types.Suites, types.Suite, types.XMLTestCase) types.Test
	Applicable(types.XMLTestSuites) bool // IsApplicable
	Name() string
}

var availableParsers = map[string]Parser{
	"exunit":  ExUnit{},
	"mocha":   Mocha{},
	"generic": Generic{},
}

// MergeResults merges multiple types.Suites together.
// Use ID of suite to determine if its unique or not.
func MergeResults(suites []types.Suites) []types.Suites {
	return suites
}

// FindParser ...
func FindParser(name string, startElement types.XMLTestSuites) Parser {
	var selectedParser Parser

	selectedParser, found := availableParsers[name]

	if found {
		return selectedParser
	}

	for _, parser := range availableParsers {
		if parser.Applicable(startElement) {
			return parser
		}
	}

	return availableParsers["generic"]
}

// LoadXML ...
func LoadXML(inFile string) types.XMLTestSuites {
	start := "testsuite"

	xmlFile, err := os.Open(inFile)
	if err != nil {
		log.Fatal(err)
	}

	decoder := xml.NewDecoder(xmlFile)

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "testsuites" {
				start = "testsuites"
				break
			}
		}
	}
	xmlFile.Seek(0, io.SeekStart)
	byteValue, _ := ioutil.ReadAll(xmlFile)

	defer xmlFile.Close()
	switch start {
	case "testsuites":
		element := types.XMLTestSuites{}
		if err := xml.Unmarshal(byteValue, &element); err != nil {
			log.Fatal(err)
		}

		return element
	default:
		element := []types.XMLTestSuite{}
		if err := xml.Unmarshal(byteValue, &element); err != nil {
			log.Fatal(err)
		}

		return types.XMLTestSuites{TestSuites: element}
	}
}

// Parse ...
func Parse(parser Parser, xmltestsuites types.XMLTestSuites) types.Suites {
	suites := parser.NewSuites(xmltestsuites)

	for _, xmltestsuite := range xmltestsuites.TestSuites {
		suite := parser.NewSuite(suites, xmltestsuite)

		for _, xmltestcase := range xmltestsuite.TestCases {
			testcase := parser.NewTest(suites, suite, xmltestcase)
			suite.Tests = append(suite.Tests, testcase)
		}

		suites.Suites = append(suites.Suites, suite)
	}

	suites.Aggregate()

	return suites
}
