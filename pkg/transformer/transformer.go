package transformer

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/semaphoreci/test-results/pkg/parsers"
)

var toXMLTemplate = `
	<{{ tag . }} {{ attributes .}}>
		{{ range $child := children . }}
			{{ template "toXML" $child }}
		{{ end }}
		{{ text . }}
	</{{ tag .}}>
`

func Transform(template template.Template, data map[string]interface{}) (output string, err error) {
	var buf bytes.Buffer

	err = template.ExecuteTemplate(&buf, "main", data)

	if err != nil {
		logger.Error("Failed to execute template: %v", err)
	}

	return buf.String(), err
}

func LoadTemplate(path string) (tpl template.Template, err error) {
	rawTemplate, err := os.ReadFile(path) // #nosec
	if err != nil {
		logger.Error("Failed to read template file: %v", err)
		return
	}

	fieldFunc := func(node map[string]interface{}, field string) string {
		value, ok := node[field].(string)
		if !ok {
			logger.Error("Failed to get field %s from node %v", field, node)
			return ""
		}
		return value
	}

	helpers := template.FuncMap{
		"field": func(node map[string]interface{}, field string) string {
			return fieldFunc(node, fmt.Sprintf("@%v", field))
		},
		"text": func(node map[string]interface{}) string {
			return fieldFunc(node, "#text")
		},
		"tag": func(node map[string]interface{}) string {
			return fieldFunc(node, "#tag")
		},
		"attributes": func(node map[string]interface{}) string {
			attributes := make([]string, 0)
			for keyName, keyValue := range node {
				if keyName[0] == '@' {
					attributes = append(attributes, fmt.Sprintf("%s=\"%s\"", keyName[1:], keyValue))
				}
			}
			return strings.Join(attributes, " ")
		},
		"children": func(node map[string]interface{}) (values []map[string]interface{}) {
			for k, v := range node {
				if k[0] != '@' && k[0] != '#' {
					values = v.([]map[string]interface{})
				}
			}
			return
		},
	}

	temp, err := template.
		New("toXML").
		Funcs(helpers).
		Parse(toXMLTemplate)

	if err != nil {
		logger.Error("Failed to parse template: %v", err)
		return
	}

	temp, err = temp.
		New("main").
		Funcs(helpers).
		Parse(string(rawTemplate))

	if err != nil {
		logger.Error("Failed to parse template: %v", err)
		return
	}

	tpl = *temp

	return
}

func LoadXML(path string) (result map[string]interface{}, err error) {
	xmlElement, err := parsers.LoadXML(path)

	if err != nil {
		logger.Error("Failed to parse xml file: %v", err)
		return
	}

	result = walk(*xmlElement, func(node parser.XMLElement) map[string]interface{} {
		el := make(map[string]interface{})
		for name, value := range node.Attributes {
			el[fmt.Sprintf("@%v", name)] = value
		}

		el["#tag"] = string(node.Tag())
		el["#text"] = string(node.Contents)

		return el
	})

	return
}

func walk(node parser.XMLElement, fn func(parser.XMLElement) map[string]interface{}) map[string]interface{} {
	result := fn(node)

	for _, child := range node.Children {
		collectionName := child.Tag()
		if _, ok := result[collectionName]; !ok {
			result[collectionName] = make([]map[string]interface{}, 0)
		}
		c := walk(child, fn)
		result[collectionName] = append(result[collectionName].([]map[string]interface{}), c)
	}

	return result
}
