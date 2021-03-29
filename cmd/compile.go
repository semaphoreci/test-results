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
			logger.SetLevel(logger.TraceLevel)
		} else if verbose {
			logger.SetLevel(logger.DebugLevel)
		}

		inFile := args[0]
		outFile := args[1]

		_, err := os.Stat(inFile)

		logger.Info("Reading %s", inFile)
		if err != nil {
			logger.Error("Input file read failed: %v", err)
			return
		}

		parser, err := parsers.FindParser(parser, inFile)
		if err != nil {
			logger.Error("Could not find parser: %v", err)
			return
		}
		logger.Info("Using %s parser", parser.GetName())

		testResults := parser.Parse(inFile)
		if name != "" {
			logger.Debug("Overriding test results name to %s", name)
			testResults.Name = name
		}

		testResults.Framework = parser.GetName()

		file, err := json.Marshal(testResults)
		if err != nil {
			logger.Error("Marshaling results failed with: %v", err)
		}

		err = ioutil.WriteFile(outFile, file, 0644)
		if err != nil {
			logger.Error("Output file write failed: %v", err)
		}
		logger.Info("Saving results to %s", outFile)
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
