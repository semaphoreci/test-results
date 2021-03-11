package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fooParser struct {
	Parser
	name string
}

func (me fooParser) Parse(path string) TestResults {
	return NewTestResults()
}

func (me fooParser) IsApplicable(path string) bool {
	return path == "foo"
}

func (me fooParser) GetName() string {
	return me.name
}

func newFooParser() fooParser {
	parser := fooParser{}
	parser.name = "foo"
	return parser
}

type barParser struct {
	Parser
	name string
}

func (me barParser) Parse(path string) TestResults {
	return NewTestResults()
}

func (me barParser) IsApplicable(path string) bool {
	return path == "bar"
}

func (me barParser) GetName() string {
	return me.name
}

func newBarParser() barParser {
	parser := barParser{}
	parser.name = "bar"
	return parser
}

func TestFindParser(t *testing.T) {
	availableParsers := []Parser{
		newFooParser(),
		newBarParser(),
	}

	parser, _ := FindParser("foo", "path", availableParsers)
	assert.IsType(t, fooParser{}, parser, "Should return correct parser")

	parser, _ = FindParser("bar", "path", availableParsers)
	assert.IsType(t, barParser{}, parser, "Should return correct parser")

	_, err := FindParser("baz", "path", availableParsers)
	assert.Error(t, err, "Should return error")

	parser, _ = FindParser("", "bar", availableParsers)
	assert.IsType(t, barParser{}, parser, "Should return correct parser")
}
