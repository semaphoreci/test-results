package cmd

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
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cli.SetLogLevel(cmd)
		if err != nil {
			return err
		}

		var dir string
		removeDir := true

		pipelineID, found := os.LookupEnv("SEMAPHORE_PIPELINE_ID")
		if !found {
			logger.Error("SEMAPHORE_PIPELINE_ID env is missing")
			return err
		}

		if len(args) == 0 {
			dir, err = ioutil.TempDir("/tmp", "test-results")
			if err != nil {
				logger.Error("Creating temporary directory failed %v", err)
				return err
			}

			dir, err = cli.PullArtifacts("workflow", path.Join("test-results", pipelineID), dir, cmd)
			if err != nil {
				return err
			}
		} else {
			dir = args[0]
			removeDir = false
		}

		result, err := cli.MergeFiles(dir, cmd)
		if err != nil {
			return err
		}

		if removeDir {
			defer os.Remove(dir)
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			logger.Error("Marshaling results failed with: %v", err)
			return err
		}

		fileName, err := cli.WriteToTmpFile(jsonData)
		if err != nil {
			return err
		}

		_, err = cli.PushArtifacts("workflow", fileName, path.Join("test-results", pipelineID+".json"), cmd)
		if err != nil {
			return err
		}

		defer os.Remove(fileName)

		return nil
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
