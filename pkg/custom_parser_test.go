package custom_parser_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/semaphoreci/test-results/cmd"
	"github.com/semaphoreci/test-results/pkg/cli"
	"github.com/stretchr/testify/assert"

	jd "github.com/josephburnett/jd/lib"
)

type Assertion struct {
	name    string
	outFile *os.File
}

func NewAssertion(path string) *Assertion {
	file, _ := ioutil.TempFile("", path)

	return &Assertion{
		name:    path,
		outFile: file,
	}
}

func (me *Assertion) Name() string {
	return me.name
}

func (me *Assertion) XMLFile() string {
	return "asserts/" + me.name + ".xml"
}

func (me *Assertion) JSONFile() string {
	return "asserts/" + me.name + ".json"
}

type CustomParserTest struct {
	name       string
	assertions []*Assertion
}

func NewCustomParserTest() *CustomParserTest {
	return &CustomParserTest{
		name:       "",
		assertions: []*Assertion{},
	}
}

func (me *CustomParserTest) AddAssertion(fileName string) {
	path := filepath.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))])
	me.assertions = append(me.assertions, NewAssertion(path))
}

func (me *CustomParserTest) BasePath() string {
	return "../parsers/" + me.name + "/"
}

func (me *CustomParserTest) TemplatePath() string {
	return me.BasePath() + "/main.tpl"
}
func Test_CustomParsers(t *testing.T) {
	dirs, err := ioutil.ReadDir("../parsers")
	if err != nil {
		t.Error(err)
	}

	for _, parserDir := range dirs {
		parserTest := NewCustomParserTest()
		parserTest.name = parserDir.Name()
		files, err := cli.LoadFiles([]string{parserTest.BasePath()}, ".xml")
		if err != nil {
			t.Error(err)
		}

		for _, file := range files {
			parserTest.AddAssertion(file)
		}

		t.Run(parserTest.name, func(t *testing.T) {
			for _, assertion := range parserTest.assertions {
				t.Run(assertion.Name(), func(t *testing.T) {
					compileCmd := cmd.NewCompileCmd()
					xmlPath := filepath.Join(parserTest.BasePath(), assertion.XMLFile())
					outJsonPath := filepath.Join(assertion.outFile.Name() + ".json")
					assertJsonPath := filepath.Join(parserTest.BasePath(), assertion.JSONFile())
					compileCmd.SetArgs([]string{"--template", parserTest.TemplatePath(), xmlPath, outJsonPath})
					compileCmd.Execute()

					outJson, err := ioutil.ReadFile(outJsonPath)

					if err != nil {
						t.Error(err)
					}
					assertJson, err := ioutil.ReadFile(assertJsonPath)

					if err != nil {
						t.Error(err)
					}

					a, _ := jd.ReadJsonString(string(outJson))
					b, _ := jd.ReadJsonString(string(assertJson))

					c := a.Diff(b).Render()
					assert.Equal(t, "", c)
				})
			}
		})

	}
}
