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

	"github.com/semaphoreci/test-results/pkg/cli"
	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/spf13/cobra"
)

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile <xml-file-path>... <json-file>]",
	Short: "parses xml files to well defined json schema",
	Long: `Parses xml file to well defined json schema

	It traverses through directory sturcture specificed by <xml-file-path> and compiles
	every .xml file.
	`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		inputs := args[:len(args)-1]
		output := args[len(args)-1]

		err := cli.SetLogLevel(cmd)
		if err != nil {
			return
		}

		paths, err := cli.LoadFiles(inputs)
		if err != nil {
			return
		}

		dirPath, err := ioutil.TempDir("", "test-results-*")
		for _, path := range paths {
			parser, err := cli.FindParser(path, cmd)
			if err != nil {
				return
			}

			testResults, err := cli.Parse(parser, path, cmd)
			if err != nil {
				return
			}

			jsonData, err := cli.Marshal(testResults)
			if err != nil {
				return
			}

			tmpFile, err := ioutil.TempFile(dirPath, "result-*.json")

			_, err = cli.WriteToFile(jsonData, tmpFile.Name())
			if err != nil {
				return
			}
		}

		result, err := cli.MergeFiles(dirPath, cmd)
		if err != nil {
			return
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			logger.Error("Marshaling results failed with: %v", err)
			return
		}

		_, err = cli.WriteToFile(jsonData, output)
		if err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
