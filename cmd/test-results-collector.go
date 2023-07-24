package cmd

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/semaphoreci/test-results/pkg/cli"
	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/spf13/cobra"
)

var genCSVReport = &cobra.Command{
	Use:   "test-results-collector",
	Short: "fetches workflow report and generates CSV file",
	Long: `fetches workflow report and generates CSV file for test results ingestion
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := cli.SetLogLevel(cmd)
		if err != nil {
			return err
		}

		testSuiteId, _ := genTestSuite(cmd, args)
		return genTest(testSuiteId, cmd, args)
	}}

func genTestSuite(cmd *cobra.Command, args []string) (string, error) {
	projectID, found := os.LookupEnv("SEMAPHORE_PROJECT_ID")
	if !found {
		logger.Error("SEMAPHORE_PROJECT_ID env is missing")
		return "", errors.New("SEMAPHORE_PROJECT_ID env is missing")
	}

	gitBranch, found := os.LookupEnv("SEMAPHORE_GIT_BRANCH")
	if !found {
		logger.Error("SEMAPHORE_PIPELINE_ID env is missing")
		return "", errors.New("SEMAPHORE_GIT_BRANCH env is missing")
	}

	testSuiteID := uuid.NewSHA1(uuid.Nil, []byte(projectID+gitBranch))

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err := writer.WriteAll([][]string{{"test_suite_id", "project_id", "branch_name"},
		{testSuiteID.String(), projectID, gitBranch}})
	if err != nil {
		logger.Error("Writing to CSV failed %v", err)
		return "", err
	}
	fileName, err := cli.WriteToTmpFile(buf.Bytes())
	if err != nil {
		logger.Error("Writing to CSV failed %v", err)
		return "", err
	}
	_, err = cli.PushArtifacts("workflow", fileName, path.Join("test-results-collector", "test-suite.csv"), cmd)
	return testSuiteID.String(), err
}

func genTest(testSuiteId string, cmd *cobra.Command, args []string) error {
	pipelineID, found := os.LookupEnv("SEMAPHORE_PIPELINE_ID")
	if !found {
		logger.Error("SEMAPHORE_PIPELINE_ID env is missing")
		return errors.New("SEMAPHORE_PIPELINE_ID env is missing")
	}

	dir, err := ioutil.TempDir("", "test-results")
	if err != nil {
		logger.Error("Creating temporary directory failed %v", err)
		return err
	}
	defer os.Remove(dir)

	dir, err = cli.PullArtifacts("workflow", path.Join("test-results", pipelineID+".json"), path.Join(dir, pipelineID+".json"), cmd)
	if err != nil {
		return err
	}

	results, err := cli.Load(dir)
	if err != nil {
		return err
	}

	// iterate over test results, generate test identity and test result CSV
	// test identity contains test suite id, test id name file name and framework
	// test result contains test id, git sha, duration, job id and state
	logger.Info("Generating test identity and test result CSV")

	testIdentitys := [][]string{{"test_suite_id", "test_id", "test_name", "file_name", "runner_name"}}
	testResults := [][]string{{"test_id", "git_sha", "duration", "job_id", "state"}}
	for _, result := range results.TestResults {
		for _, suite := range result.Suites {
			for _, test := range suite.Tests {
				fileName := test.File
				if test.File == "" {
					fileName = test.Classname
				}

				runnerName := result.Framework
				if result.Framework == "" {
					runnerName = suite.Package
				}
				testIdentity := parser.TestIdentity{TestSuiteId: testSuiteId, TestId: test.ID, TestName: test.Name, FileName: fileName, RunnerName: runnerName}
				testResult := parser.TestResult{TestId: test.ID, GitSha: test.SemEnv.GitRefSha, Duration: test.Duration, JobId: test.SemEnv.JobId, State: test.State}

				testIdentitys = append(testIdentitys, testIdentity.String())
				testResults = append(testResults, testResult.String())
			}
		}
	}

	testIdFile, err := ioutil.TempFile("", "test-identity")
	if err != nil {
		logger.Error("Creating temporary file failed %v", err)
		return err
	}
	err = csv.NewWriter(testIdFile).WriteAll(testIdentitys)
	if err != nil {
		logger.Error("Writing test identitys to CSV failed %v", err)
		return err
	}


	testResFile, err := ioutil.TempFile("", "test-identity")
	if err != nil {
		logger.Error("Creating temporary file failed %v", err)
		return err
	}

	err = csv.NewWriter(testResFile).WriteAll(testResults)
	if err != nil {
		logger.Error("Writing test results to CSV failed %v", err)
		return err
	}

	_, err = cli.PushArtifacts("workflow", testIdFile.Name(), path.Join("test-results-collector", "test-identity.csv"), cmd)
	if err != nil {
		return err
	}
	_, err = cli.PushArtifacts("workflow", testResFile.Name(), path.Join("test-results-collector", "test-result.csv"), cmd)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	genCSVReport.Flags().BoolP("force", "f", false, "force artifact push, passes -f flag to artifact CLI")
	rootCmd.AddCommand(genCSVReport)
}
