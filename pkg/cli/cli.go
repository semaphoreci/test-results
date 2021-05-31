package cli

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/semaphoreci/test-results/pkg/parsers"
	"github.com/spf13/cobra"
)

// LoadFiles checks if path exists and can be `stat`ed at given `path`
func LoadFiles(path string) ([]string, error) {
	file, err := os.Stat(path)
	paths := []string{}

	if err != nil {
		logger.Error("Input file read failed: %v", err)
		return paths, err
	}

	switch file.IsDir() {
	case true:
		filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
			if d.Type().IsRegular() {
				switch filepath.Ext(d.Name()) {
				case ".xml":
					paths = append(paths, path)
				}
			}
			return nil
		})

	case false:
		paths = append(paths, path)
	}

	sort.Strings(paths)

	return paths, nil
}

// CheckFile checks if path exists and can be `stat`ed at given `path`
func CheckFile(path string) (string, error) {
	_, err := os.Stat(path)

	if err != nil {
		logger.Error("Input file read failed: %v", err)
		return "", err
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
func PushArtifacts(level string, file string, destination string, cmd *cobra.Command) (string, error) {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		logger.Error("Reading flag error: %v", err)
		return "", err
	}

	expireIn, err := cmd.Flags().GetString("expire-in")
	if err != nil {
		logger.Error("Reading flag error: %v", err)
		return "", err
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

	logger.Info("Pushing artifacts:\n$ %s", artifactsPush.String())

	if err != nil {
		logger.Error("Pushing artifacts failed: %v\n%s", err, string(output))
		return "", err
	}
	return destination, nil
}

// PullArtifacts fetches artifacts from semaphore artifact storage
func PullArtifacts(level string, remotePath string, localPath string, cmd *cobra.Command) (string, error) {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		logger.Error("Reading flag error: %v", err)
		return "", err
	}

	artifactsPush := exec.Command("artifact")
	artifactsPush.Args = append(artifactsPush.Args, "pull", level, remotePath, "-d", localPath)
	if verbose {
		artifactsPush.Args = append(artifactsPush.Args, "-v")
	}

	output, err := artifactsPush.CombinedOutput()

	logger.Info("Pulling artifacts:\n$ %s", artifactsPush.String())

	if err != nil {
		logger.Error("Pulling artifacts failed: %v\n%s", err, string(output))
		return "", err
	}

	return localPath, nil
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

// MergeFiles merges all json files found in path into one big blob
func MergeFiles(path string, cmd *cobra.Command) (*parser.Result, error) {
	verbose, err := cmd.Flags().GetBool("verbose")

	_, err = CheckFile(path)
	if err != nil {
		logger.Error(err.Error())
	}

	var result *parser.Result

	fun := func(p string, d fs.DirEntry, err error) error {
		if verbose {
			logger.Info("[verbose] Checking file: %s", p)
		}

		if err != nil {
			logger.Info(err.Error())
			return err
		}

		if d.Type().IsDir() {
			return nil
		}

		fs, err := d.Info()
		if err != nil {
			logger.Error(err.Error())
			return err
		}

		if filepath.Ext(fs.Name()) != ".json" {
			return nil
		}

		inFile, err := CheckFile(p)
		if err != nil {
			logger.Error(err.Error())
			return err
		}

		result, err = Load(inFile, result)
		if err != nil {
			logger.Error(err.Error())
			return err
		}

		if verbose {
			logger.Info("[verbose] File loaded: %s", p)
		}

		return nil
	}

	err = filepath.WalkDir(path, fun)
	if err != nil {
		logger.Error("Test results dir listing failed: %v", err)
		return nil, err
	}
	result.Combine()

	return result, nil
}

// Load ...
func Load(path string, result *parser.Result) (*parser.Result, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var testResults parser.TestResults

	json.Unmarshal(bytes, &testResults)

	result.TestResults = append(result.TestResults, testResults)

	return result, nil
}
