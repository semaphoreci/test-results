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

var name string
var parser string

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile [xml-file] [json-file]",
	Short: "parses xml file to well defined json schema",
	Long:  `Parses xml file to well defined json schema`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		inFile := args[0]
		outFile := args[1]

		_, err := os.Stat(inFile)
		if err != nil {
			logger.Error("compile-cmd", "Input file read failed: %v", err)
		}

		parser := parsers.NewGeneric()

		testResults, err := parser.Parse(inFile)

		if err != nil {
			logger.Error("compile-cmd", "Parsing failed: %v", err)
		}

		file, err := json.Marshal(testResults)
		if err != nil {
			logger.Error("compile-cmd", "JSON marshaling failed: %v", err)
		}

		// Todo: Check if file can be created at location
		err = ioutil.WriteFile(outFile, file, 0644)
		if err != nil {
			logger.Error("compile-cmd", "Output file write failed: %v", err)
		}
	},
}

func init() {
	compileCmd.Flags().StringVarP(&name, "name", "N", "suite", "name of the suite")
	compileCmd.Flags().StringVarP(&parser, "parser", "p", "auto", "override parser to be used")
	rootCmd.AddCommand(compileCmd)
}
