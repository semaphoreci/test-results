package parsers

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/semaphoreci/test-results/pkg/fileloader"
	"github.com/semaphoreci/test-results/pkg/parser"
)

type parserTestCase struct {
	Name  string
	Input string
	want  parser.TestResults
}

var commonParserTestCases = map[string]string{
	"empty": ``,
	"basic": `
			<?xml version="1.0"?>
			<testsuite name="foo" id="1234">
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
				<testcase name="bar">
				</testcase>
			</testsuite>
		`,
	"multi-suite": `
		<?xml version="1.0"?>
		<testsuites name="ff">
			<testsuite name="foo" id="1234">
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
			</testsuite>
			<testsuite name="1234">
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
			</testsuite>
			<testsuite id="1234">
				<testcase name="bar">
				</testcase>
				<testcase name="baz">
				</testcase>
			</testsuite>
			<testsuite name="1235">
				<testcase name="bar" file="foo/bar:123">
				</testcase>
				<testcase name="baz" file="foo/baz">
				</testcase>
			</testsuite>
			<testsuite name="diff by classname">
				<testcase name="bar" file="foo/bar" classname="foo">
				</testcase>
				<testcase name="bar" file="foo/bar" classname="bar">
				</testcase>
			</testsuite>
		</testsuites>
		`,
	"invalid-root": `
			<?xml version="1.0"?>
			<nontestsuites name="em">
				<testsuite name="foo" id="1234">
					<testcase name="bar">
					</testcase>
					<testsuite name="zap" id="4321">
						<testcase name="baz">
						</testcase>
					</testsuite>
					<testsuite name="zup" id="54321">
						<testcase name="bar">
						</testcase>
					</testsuite>
				</testsuite>
			</nontestsuites>
			`,
}

func buildParserTestCases(inputs map[string]string, wants map[string]parser.TestResults) []parserTestCase {
	cases := []parserTestCase{}
	for key, input := range inputs {
		want, exists := wants[key]
		if !exists {
			// Handle missing expected results (either ignore or report an error)
			continue
		}
		cases = append(cases, parserTestCase{
			Name:  key,
			Input: input,
			want:  want,
		})
	}
	return cases
}

func runParserTests(t *testing.T, parser parser.Parser, testCases []parserTestCase) {
	t.Setenv("IP", "192.168.0.1")
	t.Setenv("SEMAPHORE_PIPELINE_ID", "ppl-id")
	t.Setenv("SEMAPHORE_WORKFLOW_ID", "wf-id")
	t.Setenv("SEMAPHORE_JOB_NAME", "job-name")
	t.Setenv("SEMAPHORE_JOB_ID", "job-id")
	t.Setenv("SEMAPHORE_PROJECT_ID", "project-id")
	t.Setenv("SEMAPHORE_AGENT_MACHINE_TYPE", "agent-machine-type")
	t.Setenv("SEMAPHORE_AGENT_MACHINE_OS_IMAGE", "agent-machine-os-image")
	t.Setenv("SEMAPHORE_JOB_CREATION_TIME", "job-creation-time")

	// For branch
	t.Setenv("SEMAPHORE_GIT_REF_TYPE", "git-ref-type")
	t.Setenv("SEMAPHORE_GIT_BRANCH", "git-branch")
	t.Setenv("SEMAPHORE_GIT_SHA", "git-sha")

	for _, tc := range testCases {
		xml := bytes.NewReader([]byte(tc.Input))
		path := fileloader.Ensure(xml)
		got := parser.Parse(path)

		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("%s parsing failed for \"%s\" case:\n%s", parser.GetName(), tc.Name, diff)
			t.Errorf("%#v\n\n", got)
		}
	}
}

func Test_FindParser(t *testing.T) {
	// testCases := []struct {
	// 	Name  string
	// 	Input string
	// 	want  parser.TestResults
	// }{
	// 	{
	// 		Name:  "Empty input",
	// 		Input: ``,
	// 		want:  parser.TestResults{ID: "", Name: "", Framework: "", IsDisabled: false, Suites: []parser.Suite{}, Summary: parser.Summary{Total: 0, Passed: 0, Skipped: 0, Error: 0, Failed: 0, Disabled: 0, Duration: 0}, Status: "error", StatusMessage: "EOF"},
	// 	},
	// }

	// for _, tc := range testCases {
	// 	xml := bytes.NewReader([]byte(tc.Input))
	// 	path := fileloader.Ensure(xml)
	// 	got := NewGeneric().Parse(path)

	// 	if diff := cmp.Diff(tc.want, got); diff != "" {
	// 		t.Errorf("Generic.Parse(\"%s\") failed (-want +got):\n%s", tc.Name, diff)
	// 		t.Errorf("Got: %#v", got)
	// 	}
	// }
}

// func TestFindParser(t *testing.T) {
// 	tests := []struct {
// 		desc    string
// 		name    string
// 		path    string
// 		reader  *bytes.Reader
// 		want    parser.Parser
// 		wantErr bool
// 	}{
// 		{
// 			desc: "finds parser automatically",
// 			name: "auto",
// 			path: fileloader.Ensure(bytes.NewReader([]byte(`
// 				<?xml version="1.0"?>
// 					<testsuites name="foo" time="0.1234" tests="10" failures="5" errors="1">
// 						<testsuite>
// 							<testcase name="bar">
// 							</testcase>
// 							<testcase name="baz">
// 							</testcase>
// 						</testsuite>
// 					</testsuites>
// 			`))),
// 			want:    Generic{},
// 			wantErr: false,
// 		},
// 		{
// 			desc: "finds rspec parser automatically",
// 			name: "auto",
// 			path: fileloader.Ensure(bytes.NewReader([]byte(`
// 				<?xml version="1.0"?>
// 					<testsuite name="rspec">
// 						<testcase name="bar">
// 						</testcase>
// 						<testcase name="baz">
// 						</testcase>
// 					</testsuite>
// 			`))),
// 			want:    RSpec{},
// 			wantErr: false,
// 		},
// 		{
// 			desc: "finds rspec parser automatically for the suite with rspec prefix",
// 			name: "auto",
// 			path: fileloader.Ensure(bytes.NewReader([]byte(`
// 				<?xml version="1.0"?>
// 					<testsuite name="rspec1">
// 						<testcase name="bar">
// 						</testcase>
// 						<testcase name="baz">
// 						</testcase>
// 					</testsuite>
// 			`))),
// 			want:    RSpec{},
// 			wantErr: false,
// 		},
// 		{
// 			desc: "finds exunit parser automatically",
// 			name: "auto",
// 			path: fileloader.Ensure(bytes.NewReader([]byte(`
// 				<?xml version="1.0"?>
// 					<testsuites>
// 						<testsuite name="Elixir.bar">
// 							<testcase name="foo">
// 							</testcase>
// 							<testcase name="baz">
// 							</testcase>
// 						</testsuite>
// 					</testsuites>
// 			`))),
// 			want:    ExUnit{},
// 			wantErr: false,
// 		},
// 		{
// 			desc: "finds mocha parser automatically",
// 			name: "auto",
// 			path: fileloader.Ensure(bytes.NewReader([]byte(`
// 				<?xml version="1.0"?>
// 				<testsuites name="Mocha tests">
// 					<testsuite name="rspec">
// 						<testcase name="bar">
// 						</testcase>
// 						<testcase name="baz">
// 						</testcase>
// 					</testsuite>
// 				</testsuites>
// 			`))),
// 			want:    Mocha{},
// 			wantErr: false,
// 		},
// 		{
// 			desc: "finds golang parser automatically",
// 			name: "auto",
// 			path: fileloader.Ensure(bytes.NewReader([]byte(`
// 				<?xml version="1.0"?>
// 				<testsuites name="tests">
// 					<testsuite>
// 						<properties>
// 							<property name="go.version" value="1.15.0"></property>
// 						</properties>
// 						<testcase name="bar">
// 						</testcase>
// 						<testcase name="baz">
// 						</testcase>
// 					</testsuite>
// 				</testsuites>
// 			`))),
// 			want:    GoLang{},
// 			wantErr: false,
// 		},
// 		{
// 			desc: "finds golang parser automatically",
// 			name: "auto",
// 			path: fileloader.Ensure(bytes.NewReader([]byte(`
// 			<?xml version="1.0"?>
// 			<testsuite>
// 				<properties>
// 					<property name="go.version" value="1.15.0"></property>
// 				</properties>
// 				<testcase name="bar">
// 				</testcase>
// 				<testcase name="baz">
// 				</testcase>
// 			</testsuite>
// 		`))),
// 			want:    GoLang{},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.desc, func(t *testing.T) {

// 			fileloader.Load(tt.path, tt.reader)

// 			got, err := FindParser(tt.name, tt.path)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("FindParser() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !assert.IsType(t, tt.want, got) {
// 				t.Errorf("Type of FindParser() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
