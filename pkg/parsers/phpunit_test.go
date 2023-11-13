package parsers

import (
	"testing"

	"github.com/semaphoreci/test-results/pkg/parser"
)

func Test_PHPUnit_CommonParse(t *testing.T) {
	parserWants := map[string]parser.TestResults{
		"empty": {
			ID:         "fbebf2b6-a680-36d2-974c-0bed1f2db373",
			Name:       "PHPUnit Suite",
			Framework:  "phpunit",
			IsDisabled: false,
			Summary: parser.Summary{
				Total:    0,
				Passed:   0,
				Skipped:  0,
				Error:    0,
				Failed:   0,
				Disabled: 0,
				Duration: 0,
			},
			Status:        "error",
			StatusMessage: "EOF",
			Suites:        []parser.Suite{},
		},
		"basic": {
			ID:         "fbebf2b6-a680-36d2-974c-0bed1f2db373",
			Name:       "PHPUnit Suite",
			Framework:  "phpunit",
			IsDisabled: false,
			Summary: parser.Summary{
				Total:    0,
				Passed:   0,
				Skipped:  0,
				Error:    0,
				Failed:   0,
				Disabled: 0,
				Duration: 0,
			},
			Status:        "success",
			StatusMessage: "",
			Suites: []parser.Suite{
				{
					ID:         "4bb3c7c8-483f-3294-9c83-1ea4a103be84",
					Name:       "foo",
					IsSkipped:  false,
					IsDisabled: false,
					Timestamp:  "",
					Hostname:   "",
					Package:    "",
					Properties: parser.Properties(nil),
					Summary: parser.Summary{
						Total:    0,
						Passed:   0,
						Skipped:  0,
						Error:    0,
						Failed:   0,
						Disabled: 0,
						Duration: 0,
					},
					SystemOut: "",
					SystemErr: "",
					Tests:     []parser.Test{},
				},
			},
		},
		"multi-suite": {
			ID:         "c1388f1b-b9b5-39ea-8f93-49f4a41ad528",
			Name:       "ff",
			Framework:  "phpunit",
			IsDisabled: false,
			Summary: parser.Summary{
				Total:    10,
				Passed:   10,
				Skipped:  0,
				Error:    0,
				Failed:   0,
				Disabled: 0,
				Duration: 0,
			},
			Status:        "success",
			StatusMessage: "",
			Suites: []parser.Suite{
				{
					ID:         "bb72528b-47d4-3cfa-9734-dd98e9e22280",
					Name:       "foo",
					IsSkipped:  false,
					IsDisabled: false,
					Timestamp:  "",
					Hostname:   "",
					Package:    "",
					Properties: parser.Properties(nil),
					Summary: parser.Summary{
						Total:    2,
						Passed:   2,
						Skipped:  0,
						Error:    0,
						Failed:   0,
						Disabled: 0,
						Duration: 0,
					},
					SystemOut: "",
					SystemErr: "",
					Tests: []parser.Test{
						{
							ID:        "70b27675-3433-3bc5-8f38-7cce102a9304",
							File:      "",
							Classname: "",
							Package:   "",
							Name:      "bar",
							Duration:  0,
							State:     "passed",
							Failure:   (*parser.Failure)(nil),
							Error:     (*parser.Error)(nil),
							SystemOut: "",
							SystemErr: "",
							SemEnv: parser.SemEnv{
								ProjectId:    "project-id",
								PipelineId:   "ppl-id",
								WorkflowId:   "wf-id",
								JobStartedAt: "job-creation-time",
								JobName:      "job-name",
								JobId:        "job-id",
								AgentType:    "agent-machine-type",
								AgentOsImage: "agent-machine-os-image",
								GitRefType:   "git-ref-type",
								GitRefName:   "",
								GitRefSha:    "",
							},
						},
						{
							ID:        "3118f840-3afe-370b-a80f-5996ec01df73",
							File:      "",
							Classname: "",
							Package:   "",
							Name:      "baz",
							Duration:  0,
							State:     "passed",
							Failure:   (*parser.Failure)(nil),
							Error:     (*parser.Error)(nil),
							SystemOut: "",
							SystemErr: "",
							SemEnv: parser.SemEnv{
								ProjectId:    "project-id",
								PipelineId:   "ppl-id",
								WorkflowId:   "wf-id",
								JobStartedAt: "job-creation-time",
								JobName:      "job-name",
								JobId:        "job-id",
								AgentType:    "agent-machine-type",
								AgentOsImage: "agent-machine-os-image",
								GitRefType:   "git-ref-type",
								GitRefName:   "",
								GitRefSha:    "",
							},
						},
					},
				},
				{
					ID:         "bb72528b-47d4-3cfa-9734-dd98e9e22280",
					Name:       "1234",
					IsSkipped:  false,
					IsDisabled: false,
					Timestamp:  "",
					Hostname:   "",
					Package:    "",
					Properties: parser.Properties(nil),
					Summary: parser.Summary{
						Total:    2,
						Passed:   2,
						Skipped:  0,
						Error:    0,
						Failed:   0,
						Disabled: 0,
						Duration: 0,
					},
					SystemOut: "",
					SystemErr: "",
					Tests: []parser.Test{
						{
							ID:        "70b27675-3433-3bc5-8f38-7cce102a9304",
							File:      "",
							Classname: "",
							Package:   "",
							Name:      "bar",
							Duration:  0,
							State:     "passed",
							Failure:   (*parser.Failure)(nil),
							Error:     (*parser.Error)(nil),
							SystemOut: "",
							SystemErr: "",
							SemEnv: parser.SemEnv{
								ProjectId:    "project-id",
								PipelineId:   "ppl-id",
								WorkflowId:   "wf-id",
								JobStartedAt: "job-creation-time",
								JobName:      "job-name",
								JobId:        "job-id",
								AgentType:    "agent-machine-type",
								AgentOsImage: "agent-machine-os-image",
								GitRefType:   "git-ref-type",
								GitRefName:   "",
								GitRefSha:    "",
							},
						},
						{
							ID:        "3118f840-3afe-370b-a80f-5996ec01df73",
							File:      "",
							Classname: "",
							Package:   "",
							Name:      "baz",
							Duration:  0,
							State:     "passed",
							Failure:   (*parser.Failure)(nil),
							Error:     (*parser.Error)(nil),
							SystemOut: "",
							SystemErr: "",
							SemEnv: parser.SemEnv{
								ProjectId:    "project-id",
								PipelineId:   "ppl-id",
								WorkflowId:   "wf-id",
								JobStartedAt: "job-creation-time",
								JobName:      "job-name",
								JobId:        "job-id",
								AgentType:    "agent-machine-type",
								AgentOsImage: "agent-machine-os-image",
								GitRefType:   "git-ref-type",
								GitRefName:   "",
								GitRefSha:    "",
							},
						},
					},
				},
				{
					ID:         "bb72528b-47d4-3cfa-9734-dd98e9e22280",
					Name:       "",
					IsSkipped:  false,
					IsDisabled: false,
					Timestamp:  "",
					Hostname:   "",
					Package:    "",
					Properties: parser.Properties(nil),
					Summary: parser.Summary{
						Total:    2,
						Passed:   2,
						Skipped:  0,
						Error:    0,
						Failed:   0,
						Disabled: 0,
						Duration: 0,
					},
					SystemOut: "",
					SystemErr: "",
					Tests: []parser.Test{
						{
							ID:        "70b27675-3433-3bc5-8f38-7cce102a9304",
							File:      "",
							Classname: "",
							Package:   "",
							Name:      "bar",
							Duration:  0,
							State:     "passed",
							Failure:   (*parser.Failure)(nil),
							Error:     (*parser.Error)(nil),
							SystemOut: "",
							SystemErr: "",
							SemEnv: parser.SemEnv{
								ProjectId:    "project-id",
								PipelineId:   "ppl-id",
								WorkflowId:   "wf-id",
								JobStartedAt: "job-creation-time",
								JobName:      "job-name",
								JobId:        "job-id",
								AgentType:    "agent-machine-type",
								AgentOsImage: "agent-machine-os-image",
								GitRefType:   "git-ref-type",
								GitRefName:   "",
								GitRefSha:    "",
							},
						},
						{
							ID:        "3118f840-3afe-370b-a80f-5996ec01df73",
							File:      "",
							Classname: "",
							Package:   "",
							Name:      "baz",
							Duration:  0,
							State:     "passed",
							Failure:   (*parser.Failure)(nil),
							Error:     (*parser.Error)(nil),
							SystemOut: "",
							SystemErr: "",
							SemEnv: parser.SemEnv{
								ProjectId:    "project-id",
								PipelineId:   "ppl-id",
								WorkflowId:   "wf-id",
								JobStartedAt: "job-creation-time",
								JobName:      "job-name",
								JobId:        "job-id",
								AgentType:    "agent-machine-type",
								AgentOsImage: "agent-machine-os-image",
								GitRefType:   "git-ref-type",
								GitRefName:   "",
								GitRefSha:    "",
							},
						},
					},
				},
				{
					ID:         "6ab72d38-571f-38e8-bc1a-f2f0a892d94f",
					Name:       "1235",
					IsSkipped:  false,
					IsDisabled: false,
					Timestamp:  "",
					Hostname:   "",
					Package:    "",
					Properties: parser.Properties(nil),
					Summary: parser.Summary{
						Total:    2,
						Passed:   2,
						Skipped:  0,
						Error:    0,
						Failed:   0,
						Disabled: 0,
						Duration: 0,
					},
					SystemOut: "",
					SystemErr: "",
					Tests: []parser.Test{
						{
							ID:        "961a8a17-78f1-3335-b272-e89b78e2d223",
							File:      "foo/bar:123",
							Classname: "",
							Package:   "",
							Name:      "bar",
							Duration:  0,
							State:     "passed",
							Failure:   (*parser.Failure)(nil),
							Error:     (*parser.Error)(nil),
							SystemOut: "",
							SystemErr: "",
							SemEnv: parser.SemEnv{
								ProjectId:    "project-id",
								PipelineId:   "ppl-id",
								WorkflowId:   "wf-id",
								JobStartedAt: "job-creation-time",
								JobName:      "job-name",
								JobId:        "job-id",
								AgentType:    "agent-machine-type",
								AgentOsImage: "agent-machine-os-image",
								GitRefType:   "git-ref-type",
								GitRefName:   "",
								GitRefSha:    "",
							},
						},
						{
							ID:        "731ae751-cd7f-3696-b913-2dbbdaa66772",
							File:      "foo/baz",
							Classname: "",
							Package:   "",
							Name:      "baz",
							Duration:  0,
							State:     "passed",
							Failure:   (*parser.Failure)(nil),
							Error:     (*parser.Error)(nil),
							SystemOut: "",
							SystemErr: "",
							SemEnv: parser.SemEnv{
								ProjectId:    "project-id",
								PipelineId:   "ppl-id",
								WorkflowId:   "wf-id",
								JobStartedAt: "job-creation-time",
								JobName:      "job-name",
								JobId:        "job-id",
								AgentType:    "agent-machine-type",
								AgentOsImage: "agent-machine-os-image",
								GitRefType:   "git-ref-type",
								GitRefName:   "",
								GitRefSha:    "",
							},
						},
					},
				},
				{
					ID:         "09db2b35-5e0d-3560-be28-fe9252a73a37",
					Name:       "diff by classname",
					IsSkipped:  false,
					IsDisabled: false,
					Timestamp:  "",
					Hostname:   "",
					Package:    "",
					Properties: parser.Properties(nil),
					Summary: parser.Summary{
						Total:    2,
						Passed:   2,
						Skipped:  0,
						Error:    0,
						Failed:   0,
						Disabled: 0,
						Duration: 0,
					},
					SystemOut: "",
					SystemErr: "",
					Tests: []parser.Test{
						{
							ID:        "4f478561-3125-36f7-850f-c1d19985d412",
							File:      "foo/bar",
							Classname: "foo",
							Package:   "",
							Name:      "bar",
							Duration:  0,
							State:     "passed",
							Failure:   (*parser.Failure)(nil),
							Error:     (*parser.Error)(nil),
							SystemOut: "",
							SystemErr: "",
							SemEnv: parser.SemEnv{
								ProjectId:    "project-id",
								PipelineId:   "ppl-id",
								WorkflowId:   "wf-id",
								JobStartedAt: "job-creation-time",
								JobName:      "job-name",
								JobId:        "job-id",
								AgentType:    "agent-machine-type",
								AgentOsImage: "agent-machine-os-image",
								GitRefType:   "git-ref-type",
								GitRefName:   "",
								GitRefSha:    "",
							},
						},
						{
							ID:        "51038177-419f-32d9-b0d4-438f7a898efc",
							File:      "foo/bar",
							Classname: "bar",
							Package:   "",
							Name:      "bar",
							Duration:  0,
							State:     "passed",
							Failure:   (*parser.Failure)(nil),
							Error:     (*parser.Error)(nil),
							SystemOut: "",
							SystemErr: "",
							SemEnv: parser.SemEnv{
								ProjectId:    "project-id",
								PipelineId:   "ppl-id",
								WorkflowId:   "wf-id",
								JobStartedAt: "job-creation-time",
								JobName:      "job-name",
								JobId:        "job-id",
								AgentType:    "agent-machine-type",
								AgentOsImage: "agent-machine-os-image",
								GitRefType:   "git-ref-type",
								GitRefName:   "",
								GitRefSha:    "",
							},
						},
					},
				},
			},
		},
		"invalid-root": {
			ID:         "fbebf2b6-a680-36d2-974c-0bed1f2db373",
			Name:       "PHPUnit Suite",
			Framework:  "phpunit",
			IsDisabled: false,
			Summary: parser.Summary{
				Total:    0,
				Passed:   0,
				Skipped:  0,
				Error:    0,
				Failed:   0,
				Disabled: 0,
				Duration: 0,
			},
			Status:        "error",
			StatusMessage: "Invalid root element found: <nontestsuites>, must be one of <testsuites>, <testsuite>",
			Suites:        []parser.Suite{},
		},
	}

	testCases := buildParserTestCases(commonParserTestCases, parserWants)
	runParserTests(t, NewPHPUnit(), testCases)
}

func Test_PHPUnit_SpecificParse(t *testing.T) {
	specificParserTestCases := map[string]string{}
	parserWants := map[string]parser.TestResults{}

	testCases := buildParserTestCases(specificParserTestCases, parserWants)
	runParserTests(t, NewPHPUnit(), testCases)
}
