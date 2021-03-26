package parser

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
