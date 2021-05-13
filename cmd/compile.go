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
	"github.com/semaphoreci/test-results/pkg/cli"
	"github.com/spf13/cobra"
)

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile [xml-file] [json-file]",
	Short: "parses xml file to well defined json schema",
	Long:  `Parses xml file to well defined json schema`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := cli.SetLogLevel(cmd)
		if err != nil {
			return
		}

		inFile, err := cli.CheckFile(args[0])
		if err != nil {
			return
		}

		parser, err := cli.FindParser(inFile, cmd)
		if err != nil {
			return
		}

		testResults, err := cli.Parse(parser, inFile, cmd)
		if err != nil {
			return
		}

		jsonData, err := cli.Marshal(testResults)
		if err != nil {
			return
		}

		_, err = cli.WriteToFile(jsonData, args[1])
		if err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
