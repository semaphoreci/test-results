package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"

	"github.com/semaphoreci/test-results/pkg/parsers"
	"github.com/semaphoreci/test-results/pkg/types"
)

// Parse ...
func Parse(parser parsers.Parser, xmltestsuites types.XMLTestSuites) types.Suites {
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

func main() {
	start := types.XMLTestSuites{}

	xmlFile, err := os.Open("test/exunit.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()


	byteValue, _ := ioutil.ReadAll(xmlFile)

	if err := xml.Unmarshal(byteValue, &start); err != nil {
		log.Fatal(err)
	}



	var myParsers []parsers.Parser
	myParsers = append(myParsers, parsers.ParserExUnit{})
	myParsers = append(myParsers, parsers.ParserMocha{})
	myParsers = append(myParsers, parsers.ParserGeneric{})

	var selectedParser parsers.Parser

	for _, parser := range myParsers {
		if(parser.Applicable(start)) {
			selectedParser = parser
			break;
		}
	}
	results := Parse(selectedParser, start)

	file, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}

 _ = ioutil.WriteFile("test/exunit.json", file, 0644)

}
