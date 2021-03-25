package cmd

/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/semaphoreci/test-results/pkg/parsers"
	"github.com/spf13/cobra"
)

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile [xml-file] [json-file]",
	Short: "parses xml file to well defined json schema",
	Long:  `Parses xml file to well defined json schema`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		if trace {
			logger.LogEntry.SetLevel(logger.TraceLevel)
		} else if verbose {
			logger.LogEntry.SetLevel(logger.DebugLevel)
		}

		var logFields = logger.Fields{"app": "compile"}
		inFile := args[0]
		outFile := args[1]

		_, err := os.Stat(inFile)

		if err != nil {
			logger.Error(logFields, "Input file read failed: %v", err)
		} else {
			logger.Info(logFields, "File successfuly read: %s", inFile)
		}

		parser, err := parsers.FindParser(parser, inFile)
		if err != nil {
			logger.Error(logFields, "Could not find parser: %v", err)
		} else {
			logger.Info(logFields, "Parser found: %s", parser.GetName())
		}

		testResults := parser.Parse(inFile)

		file, err := json.Marshal(testResults)
		if err != nil {
			logger.Error(logFields, "JSON marshaling failed: %v", err)
		} else {
			logger.Info(logFields, "JSON marshaling completed: %s", inFile)
		}

		// Todo: Check if file can be created at location
		err = ioutil.WriteFile(outFile, file, 0644)
		if err != nil {
			logger.Error(logFields, "Output file write failed: %v", err)
		} else {
			logger.Info(logFields, "File saved to: %s", outFile)
		}
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
