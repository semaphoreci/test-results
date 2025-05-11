package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// commandMetricsCmd represents the job-metrics command
var commandMetricsCmd = &cobra.Command{
	Use:   "command-metrics",
	Short: "Show job metrics and flowchart",
	Long:  `Show job metrics and flowchart breakdown from job logs.`,
	Args:  cobra.PositionalArgs(nil),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Placeholder for metrics parsing logic
		out := ""
		// [metrics parsing logic would go here]

		srcFile, err := cmd.Flags().GetString("src")
		if err != nil {
			return fmt.Errorf("src cannot be parsed: %w", err)
		}

		// Find the most recent /tmp/job_log_*.json file
		matches, err := filepath.Glob(srcFile)
		if err != nil || len(matches) == 0 {
			return fmt.Errorf("failed to find job log file: %w", err)
		}

		type CmdFinished struct {
			Event      string `json:"event"`
			Directive  string `json:"directive"`
			StartedAt  int64  `json:"started_at"`
			FinishedAt int64  `json:"finished_at"`
		}

		lines, err := os.ReadFile(matches[0])
		if err != nil {
			return fmt.Errorf("could not read job log: %w", err)
		}

		var flowNodes []CmdFinished
		for _, raw := range strings.Split(string(lines), "\n") {
			if strings.TrimSpace(raw) == "" {
				continue
			}
			var entry CmdFinished
			if err := json.Unmarshal([]byte(raw), &entry); err != nil {
				continue
			}
			if entry.Event == "cmd_finished" {
				flowNodes = append(flowNodes, entry)
			}
		}

		out += "## 🧭 Job Timeline\n\n```mermaid\ngantt\n    title Job Command Timeline\n    dateFormat X\n"
		for i, node := range flowNodes {
			duration := node.FinishedAt - node.StartedAt
			if duration < 1 {
				duration = 1 // ensure minimum duration for visibility
			}
			out += fmt.Sprintf("    %s :step%d, %d, %ds\n", node.Directive, i, node.StartedAt, duration)
		}
		out += "```\n"

		if len(args) != 1 {
			return fmt.Errorf("please provide the output file path as a positional argument")
		}
		if err := os.WriteFile(args[0], []byte(out), 0644); err != nil {
			return fmt.Errorf("failed to write output: %w", err)
		}
		return nil
	},
}

func init() {
	commandMetricsCmd.Flags().String("src", "/tmp/job_log_*.json", "source file to read system metrics from")
	rootCmd.AddCommand(commandMetricsCmd)
}
