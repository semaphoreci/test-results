# Custom parsers

Generic parser parses JUnit XML documents according to [JUnit XML Schema](https://github.com/windyroad/JUnit-Schema/blob/master/JUnit.xsd). In some situations, it might be necessary to write a  custom parser. This document will guide you through the process.

> **Note:** This guide assumes you have a basic understanding of [golang](https://go.dev/).

## Creating a new parser

Every parser provides an implementation of the [`Parser`](https://pkg.go.dev/github.com/semaphoreci/test-results@v0.4.13/pkg/parser#Parser) interface, in particular:

- `GetName() string` - returns the name of the parser that can then be used in the CLI as a `--parser` option

- `IsApplicable(path string) bool` - should return true if the parser is applicable to the given file at `path`

- `Parse(path string) parser.TestResults` - parses the file at `path` and returns a `parser.TestResults` struct. This struct has a well-defined format and can be serialized to JSON.

[Generic parser implementation](https://github.com/semaphoreci/test-results/blob/master/pkg/parsers/generic.go) is a good starting template for creating a custom parser.

After the parser is implemented, it has to be added to the list of [available parsers](https://github.com/semaphoreci/test-results/blob/master/pkg/parsers/parsers.go#L10).

## Good parser qualities

- The parser is as generic as possible.

    Currently, custom parsers need to be compiled, thus merged to the main repository. Each test runner should have at most one parser.

- Parsing is idempotent.

    If you parse the same file twice, the results should be the same. It's is highly crucial for test IDs. Please refer to [`ID generation guide`](id-generation.md) for more details.

- It's tested.

    If the parser is lacking tests, it will most likely be rejected.
