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
	"os/exec"

	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/semaphoreci/test-results/pkg/parsers"
	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish [xml-file]",
	Short: "parses xml file to well defined json schema and publishes results to artifacts storage",
	Long:  `Parses xml file to well defined json schema and publishes results to artifacts storage`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inFile := args[0]

		_, err := os.Stat(inFile)
		if err != nil {
			logger.Error("publish-cmd", "Input file read failed: %v", err)
		}

		parser := parsers.NewGeneric()

		testResults, err := parser.Parse(inFile)

		if err != nil {
			logger.Error("publish-cmd", "Parsing failed: %v", err)
		}

		file, err := json.Marshal(testResults)
		if err != nil {
			logger.Error("publish-cmd", "JSON marshaling failed: %v", err)
		}

		tmpFile, err := ioutil.TempFile("/tmp", "test-results")

		// Todo: Check if file can be created at location
		_, err = tmpFile.Write(file)
		if err != nil {
			logger.Error("publish-cmd", "Output file write failed: %v", err)
		}

		artifactsPush := exec.Command("artifact", "push", "job", tmpFile.Name(), "-d", "test-results/junit.json")
		err = artifactsPush.Run()
		if err != nil {
			logger.Error("publish-cmd", "Pushing artifacts failed: %v", err)
		}

		artifactsPush = exec.Command("artifact", "push", "job", inFile, "-d", "test-results/junit.xml")
		err = artifactsPush.Run()
		if err != nil {
			logger.Error("publish-cmd", "Pushing artifacts failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
}
