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
	"path"

	"github.com/semaphoreci/test-results/pkg/cli"
	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish <xml-file-path>...",
	Short: "parses xml file to well defined json schema and publishes results to artifacts storage",
	Long: `Parses xml file to well defined json schema and publishes results to artifacts storage

	It traverses through directory sturcture specificed by <xml-file-path>, compiles
	every .xml file and publishes it as one artifact.
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		inputs := args
		err := cli.SetLogLevel(cmd)
		if err != nil {
			return
		}

		paths, err := cli.LoadFiles(inputs, ".xml")
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

		fileName, err := cli.WriteToTmpFile(jsonData)
		if err != nil {
			return
		}

		_, err = cli.PushArtifacts("job", fileName, path.Join("test-results", "junit.json"), cmd)
		if err != nil {
			return
		}

		pipelineID, found := os.LookupEnv("SEMAPHORE_PIPELINE_ID")
		if !found {
			logger.Error("SEMAPHORE_PIPELINE_ID env is missing")
			return
		}

		jobID, found := os.LookupEnv("SEMAPHORE_JOB_ID")
		if !found {
			logger.Error("SEMAPHORE_JOB_ID env is missing")
			return
		}

		_, err = cli.PushArtifacts("workflow", fileName, path.Join("test-results", pipelineID, jobID+".json"), cmd)
		if err != nil {
			return
		}

		noRaw, err := cmd.Flags().GetBool("no-raw")
		if err != nil {
			logger.Error("Reading flag error: %v", err)
			return
		}
		if !noRaw {
			for _, rawFilePath := range paths {
				_, err = cli.PushArtifacts("job", rawFilePath, path.Join("test-results/raw", rawFilePath), cmd)
				if err != nil {
					return
				}
			}
		}
	},
}

func init() {

	desc := `Skips uploading raw XML files`
	publishCmd.Flags().BoolP("no-raw", "", false, desc)

	desc = `Removes the files after the given amount of time.
Nd for N days
Nw for N weeks
Nm for N months
Ny for N years
`
	publishCmd.Flags().StringP("expire-in", "", "", desc)

	rootCmd.AddCommand(publishCmd)
}
