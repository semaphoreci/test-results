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

// genPipelineReportCmd represents the publish command
var genPipelineReportCmd = &cobra.Command{
	Use:   "gen-pipeline-report [<path>...]",
	Short: "fetches workflow level JUnit reports and combines them together",
	Long: `fetches workflow level junit reports and combines them

	When <path>s are provided it recursively traverses through path structure and
	combines all .json files into one JSON schema file.
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		err := cli.SetLogLevel(cmd)
		if err != nil {
			return
		}

		var dir string

		if len(args) == 0 {

			pipelineID, found := os.LookupEnv("SEMAPHORE_PIPELINE_ID")
			if !found {
				logger.Error("SEMAPHORE_PIPELINE_ID env is missing")
				return
			}
			dir, err = ioutil.TempDir("/tmp", "test-results")
			if err != nil {
				logger.Error("Creating temporary directory failed %v", err)
				return
			}

			dir, err = cli.PullArtifacts("workflow", path.Join("test-results", pipelineID), dir, cmd)
			if err != nil {
				return
			}
		} else {
			dir = args[0]
		}

		result, err := cli.MergeFiles(dir, cmd)
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

		_, err = cli.PushArtifacts("workflow", fileName, path.Join("test-results", "junit.json"), cmd)
		if err != nil {
			return
		}

	},
}

func init() {
	desc := `Removes the files after the given amount of time.
Nd for N days
Nw for N weeks
Nm for N months
Ny for N years
`
	genPipelineReportCmd.Flags().StringP("expire-in", "", "", desc)
	genPipelineReportCmd.Flags().BoolP("force", "f", false, "force artifact push, passes -f flag to artifact CLI")
	rootCmd.AddCommand(genPipelineReportCmd)
}
