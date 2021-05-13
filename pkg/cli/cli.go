package cli

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/semaphoreci/test-results/pkg/parsers"
	"github.com/spf13/cobra"
)

// CheckFile checks if file exists and can be `stat`ed at given `path`
func CheckFile(path string) (string, error) {
	_, err := os.Stat(path)
	if err != nil {
		logger.Error("Input file read failed: %v", err)
		return path, err
	}

	return path, nil
}

// FindParser finds parser according to file type or flag specified by user
func FindParser(path string, cmd *cobra.Command) (parser.Parser, error) {

	parserName, err := cmd.Flags().GetString("parser")
	if err != nil {
		logger.Error("Reading flag error: %v", err)
		return nil, err
	}

	parser, err := parsers.FindParser(parserName, path)
	if err != nil {
		logger.Error("Could not find parser: %v", err)
		return nil, err
	}
	logger.Info("Using %s parser", parser.GetName())
	return parser, nil
}

// Parse parses file at `path` with given `parser`
func Parse(parser parser.Parser, path string, cmd *cobra.Command) (parser.TestResults, error) {
	testResults := parser.Parse(path)

	testResultsName, err := cmd.Flags().GetString("name")
	if err != nil {
		logger.Error("Reading flag error: %v", err)
		return testResults, err
	}

	if testResultsName != "" {
		logger.Debug("Overriding test results name to %s", testResultsName)
		testResults.Name = testResultsName
	}

	return testResults, nil
}

// Marshal provides json output for given test results
func Marshal(testResults parser.TestResults) ([]byte, error) {
	jsonData, err := json.Marshal(testResults)
	if err != nil {
		logger.Error("Marshaling results failed with: %v", err)
		return nil, err
	}
	return jsonData, nil
}

// WriteToFile saves data to given file
func WriteToFile(data []byte, path string) (string, error) {
	file, err := os.Create(path)

	if err != nil {
		logger.Error("Opening file %s: %v", path, err)
		return "", err
	}
	return writeToFile(data, file)
}

// WriteToTmpFile saves data to temporary file
func WriteToTmpFile(data []byte) (string, error) {
	file, err := ioutil.TempFile("/tmp", "test-results")

	if err != nil {
		logger.Error("Opening file %s: %v", file.Name(), err)
		return "", err
	}
	return writeToFile(data, file)
}

func writeToFile(data []byte, file *os.File) (string, error) {
	logger.Info("Saving results to %s", file.Name())

	_, err := file.Write(data)
	if err != nil {
		logger.Error("Output file write failed: %v", err)
		return "", err
	}

	return file.Name(), nil
}

// PushArtifacts publishes artifacts to semaphore artifact storage
func PushArtifacts(level string, file string, destination string, cmd *cobra.Command) error {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		logger.Error("Reading flag error: %v", err)
		return err
	}

	expireIn, err := cmd.Flags().GetString("expire-in")
	if err != nil {
		logger.Error("Reading flag error: %v", err)
		return err
	}

	artifactsPush := exec.Command("artifact")
	artifactsPush.Args = append(artifactsPush.Args, "push", level, file, "-d", destination)
	if verbose {
		artifactsPush.Args = append(artifactsPush.Args, "-v")
	}

	if expireIn != "" {
		artifactsPush.Args = append(artifactsPush.Args, "--expire-in", expireIn)
	}

	output, err := artifactsPush.CombinedOutput()

	logger.Info("Pushing json artifacts:\n > %s", artifactsPush.String())

	if err != nil {
		logger.Error("Pushing artifacts failed: %v\n%s", err, string(output))
		return err
	}
	return nil
}

// SetLogLevel sets log level according to flags
func SetLogLevel(cmd *cobra.Command) error {
	trace, err := cmd.Flags().GetBool("trace")
	if err != nil {
		logger.Error("Reading flag error: %v", err)
		return err
	}

	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		logger.Error("Reading flag error: %v", err)
		return err
	}

	if trace {
		logger.SetLevel(logger.TraceLevel)
	} else if verbose {
		logger.SetLevel(logger.DebugLevel)
	}
	return nil
}
