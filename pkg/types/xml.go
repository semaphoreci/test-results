package types

// XMLTestSuites ...
type XMLTestSuites struct {
	ID string `xml:"id,attr"`
	TestSuites []XMLTestSuite `xml:"testsuite"`
}

// XMLTestSuite ...
type XMLTestSuite struct {
	Name string `xml:"name,attr"`
	Errors int `xml:"errors,attr"`
	Failures int `xml:"failures,attr"`
	Tests int `xml:"tests,attr"`
	TestsCases []XMLTestCase `xml:"testcase"`
	Time float64 `xml:"time,attr"`
}

// XMLTestCase ...
type XMLTestCase struct {
	Name string `xml:"name,attr"`
	File string `xml:"file,attr"`
	ClassName string `xml:"classname,attr"`
	Time float64 `xml:"time,attr"`
	Failure *XMLFailure `xml:"failure"`
}

// XMLFailure ...
type XMLFailure struct {
	Message string `xml:"message,attr"`
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}
