package parser

// Parser interface defines the methods that a parser must implement.
type Parser interface {
	// Parse return a TestResults struct containing the results of the parsing file at path
	Parse(path string) TestResults
	// IsApplicable returns true if the parser is applicable to the file at path
	IsApplicable(path string) bool
	// GetName returns the name of the parser
	GetName() string
}
