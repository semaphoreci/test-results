package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/semaphoreci/test-results/pkg/parsers"
	"github.com/semaphoreci/test-results/pkg/types"
)

func main() {
	start := types.XMLTestSuites{}

	xmlFile, err := os.Open("test/exunit.xml")
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()


	byteValue, _ := ioutil.ReadAll(xmlFile)

	if err := xml.Unmarshal(byteValue, &start); err != nil {
		log.Fatal(err)
	}

	parser := &parsers.ParserExUnit{}
	results := parser.Parse(start)

	file, err := json.Marshal(results)
	if err != nil {
		fmt.Println(err)
	}

 _ = ioutil.WriteFile("test/exunit.json", file, 0644)


	// fmt.Printf("%#v", b)
}
