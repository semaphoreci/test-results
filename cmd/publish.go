package cmd

import (
	"encoding/json"
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		inputs := args
		err := cli.SetLogLevel(cmd)
		if err != nil {
			return err
		}

		paths, err := cli.LoadFiles(inputs, ".xml")
		if err != nil {
			return err
		}

		dirPath, err := ioutil.TempDir("", "test-results-*")

		if err != nil {
			return err
		}

		for _, path := range paths {
			parser, err := cli.FindParser(path, cmd)
			if err != nil {
				return err
			}

			testResults, err := cli.Parse(parser, path, cmd)
			if err != nil {
				return err
			}

			jsonData, err := cli.Marshal(testResults)
			if err != nil {
				return err
			}

			tmpFile, err := ioutil.TempFile(dirPath, "result-*.json")
			if err != nil {
				return err
			}

			_, err = cli.WriteToFile(jsonData, tmpFile.Name())
			if err != nil {
				return err
			}
		}

		result, err := cli.MergeFiles(dirPath, cmd)
		if err != nil {
			return err
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

		_, err = cli.PushArtifacts("job", fileName, path.Join("test-results", "junit.json"), cmd)
		if err != nil {
			return err
		}

		pipelineID, found := os.LookupEnv("SEMAPHORE_PIPELINE_ID")
		if !found {
			logger.Error("SEMAPHORE_PIPELINE_ID env is missing")
			return err
		}

		jobID, found := os.LookupEnv("SEMAPHORE_JOB_ID")
		if !found {
			logger.Error("SEMAPHORE_JOB_ID env is missing")
			return err
		}

		_, err = cli.PushArtifacts("workflow", fileName, path.Join("test-results", pipelineID, jobID+".json"), cmd)
		if err != nil {
			return err
		}

		noRaw, err := cmd.Flags().GetBool("no-raw")
		if err != nil {
			logger.Error("Reading flag error: %v", err)
			return err
		}

		if !noRaw {
			singlePath := true
			if len(paths) > 1 {
				singlePath = false
			}

			for idx, rawFilePath := range paths {
				outPath := path.Join("test-results", "junit.xml")
				if !singlePath {
					outPath = path.Join("test-results", fmt.Sprintf("junit-%d.xml", idx))
				}

				_, err = cli.PushArtifacts("job", rawFilePath, outPath, cmd)
				if err != nil {
					return err
				}
			}
		}

		defer os.Remove(fileName)

		return nil
	},
}

func init() {

	desc := `Skips uploading raw XML files`
	publishCmd.Flags().BoolP("no-raw", "", false, desc)
	publishCmd.Flags().BoolP("force", "f", false, "force artifact push, passes -f flag to artifact CLI")

	desc = `Removes the files after the given amount of time.
Nd for N days
Nw for N weeks
Nm for N months
Ny for N years
`
	publishCmd.Flags().StringP("expire-in", "", "", desc)

	rootCmd.AddCommand(publishCmd)
}
