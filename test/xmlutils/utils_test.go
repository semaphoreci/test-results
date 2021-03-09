package xmlutils

import (
	"bytes"
	"testing"

	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/semaphoreci/test-results/pkg/xmlutils"
)

func TestLoad(t *testing.T) {

	reader := bytes.NewReader([]byte(`
		<?xml version="1.0"?>
			<testsuite name="foo" id="1234">
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
			</testsuite>
	`))
	decoder, err := xmlutils.Load("sample.xml", *reader)
	if err != nil {
		t.Error(err)
	}
	xml := parser.LoadXML(decoder)
	xmlutils.Log(xml)

	// xml, err = xmlutils.Load("sample.xml", reader)
	// if err != nil {
	// 	t.Error(err)
	// }
	// t.Log(xml)
	// parser, err := parser.FindParser(xml)
	// if err != nil {
	// 	t.Error(err)
	// }
	// testsuite, err := parser.Parse(xml)
	// if err != nil {
	// 	t.Error(err)
	// }

	// t.Log(file)
	// t.Log(err)
	// t.Log("TEST")
	// got := -1
	// if got != 1 {
	// 	t.Errorf("Abs(-1) = %d; want 1", got)
	// }
	t.Fail()
}

func TestLoad1(t *testing.T) {

	reader := bytes.NewReader([]byte(`
		<?xml version="1.0"?>
		<testsuites id="123">
			<testsuite name="foo" id="123">
				<testcase name="foo1" />
			</testsuite>
			<testsuite name="bar" id="124"></testsuite>
		</testsuites>
	`))
	decoder, err := xmlutils.Load("sample1.xml", *reader)
	if err != nil {
		t.Error(err)
	}

	_ = parser.LoadXML(decoder)

	t.Fail()
}
