package parsers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/semaphoreci/test-results/pkg/parsers"
)

func TestParse(t *testing.T) {
	var parser parsers.Generic

	testResults, err := parser.Parse("../priv/exunit.xml")
	if err != nil {
		t.Error(err)
	}

	file, err := json.Marshal(testResults)
	if err != nil {
		log.Fatal(err)
	}

	_ = ioutil.WriteFile("/tmp/results.json", file, 0644)

	t.Fail()
}
