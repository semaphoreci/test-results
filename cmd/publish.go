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
	"os"
	"path"

	"github.com/semaphoreci/test-results/pkg/cli"
	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish [xml-file]",
	Short: "parses xml file to well defined json schema and publishes results to artifacts storage",
	Long:  `Parses xml file to well defined json schema and publishes results to artifacts storage`,
	Args:  cobra.MinimumNArgs(1),
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

		fileName, err := cli.WriteToTmpFile(jsonData)
		if err != nil {
			return
		}

		err = cli.PushArtifacts("job", fileName, path.Join("test-results", "junit.json"), cmd)
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

		err = cli.PushArtifacts("workflow", fileName, path.Join("test-results", pipelineID, jobID+".json"), cmd)
		if err != nil {
			return
		}

		noRaw, err := cmd.Flags().GetBool("no-raw")
		if err != nil {
			logger.Error("Reading flag error: %v", err)
			return
		}
		if !noRaw {
			err = cli.PushArtifacts("job", inFile, "test-results/junit.xml", cmd)
			if err != nil {
				return
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
