package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type JobMetric struct {
	Timestamp  string
	CPU        float64
	Memory     float64
	SystemDisk float64
	DockerDisk float64
}

// combineCmd represents the combine command
var jobMetricsCmd = &cobra.Command{
	Use:   "job-metrics",
	Short: "TBD",
	Long:  `TBD"`,
	Args:  cobra.PositionalArgs(nil),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcFile, err := cmd.Flags().GetString("src")
		if err != nil {
			return fmt.Errorf("src cannot be parsed: %w", err)
		}

		file, err := os.Open(srcFile)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		defer file.Close()

		metricLineRegex := regexp.MustCompile(`^(.*?) \|  cpu:(.*)%,  mem:\s*(.*)%,  system_disk:\s*(.*)%,  docker_disk:\s*(.*)%,(.*)$`)

		var metrics []JobMetric
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			matches := metricLineRegex.FindStringSubmatch(line)
			if len(matches) != 7 {
				continue
			}
			var m JobMetric
			m.Timestamp = matches[1]
			fmt.Sscanf(matches[2], "%f", &m.CPU)
			fmt.Sscanf(matches[3], "%f", &m.Memory)
			fmt.Sscanf(matches[4], "%f", &m.SystemDisk)
			fmt.Sscanf(matches[5], "%f", &m.DockerDisk)
			metrics = append(metrics, m)
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		if len(metrics) == 0 {
			return fmt.Errorf("no valid data found")
		}

		step := 1
		if len(metrics) > 100 {
			step = len(metrics) / 100
		}

		var (
			xLabels          []string
			cpuSeries        []string
			memSeries        []string
			sysDiskSeries    []string
			dockerDiskSeries []string
		)

		// Parse the first timestamp to use as a reference for relative times
		layout := "Mon 02 Jan 2006 03:04:05 PM MST"
		startTime, err := time.Parse(layout, metrics[0].Timestamp)
		if err != nil {
			return fmt.Errorf("failed to parse start time: %w", err)
		}

		min := func(f1, f2 float64) float64 {
			if f1 < f2 {
				return f1
			}
			return f2
		}

		max := func(f1, f2 float64) float64 {
			if f1 > f2 {
				return f1
			}
			return f2
		}

		cpuMin, cpuMax := metrics[0].CPU, metrics[0].CPU
		memMin, memMax := metrics[0].Memory, metrics[0].Memory
		diskMin, diskMax := metrics[0].SystemDisk, metrics[0].SystemDisk
		dockerMin, dockerMax := metrics[0].DockerDisk, metrics[0].DockerDisk

		for i := 0; i < len(metrics); i += step {
			m := metrics[i]
			cpuMin = min(cpuMin, m.CPU)
			cpuMax = max(cpuMax, m.CPU)
			memMin = min(memMin, m.Memory)
			memMax = max(memMax, m.Memory)
			diskMin = min(diskMin, m.SystemDisk)
			diskMax = max(diskMax, m.SystemDisk)
			dockerMin = min(dockerMin, m.DockerDisk)
			dockerMax = max(dockerMax, m.DockerDisk)

			t, err := time.Parse(layout, m.Timestamp)
			if err != nil {
				xLabels = append(xLabels, "\"??:??\"")
			} else {
				duration := t.Sub(startTime)
				seconds := int(duration.Seconds())
				xLabels = append(xLabels, fmt.Sprintf("\"%02d:%02d\"", seconds/60, seconds%60))
			}
			cpuSeries = append(cpuSeries, fmt.Sprintf("%.2f", m.CPU))
			memSeries = append(memSeries, fmt.Sprintf("%.2f", m.Memory))
			sysDiskSeries = append(sysDiskSeries, fmt.Sprintf("%.2f", m.SystemDisk))
			dockerDiskSeries = append(dockerDiskSeries, fmt.Sprintf("%.2f", m.DockerDisk))
		}

		out := "## 🎯 System Metrics Summary\n\n"
		out += "```mermaid\n"
		out += "gantt\n"
		out += "    title 📊 Metrics Overview Timeline\n"
		out += "    dateFormat HH:mm\n"
		out += "    axisFormat %H:%M\n"
		out += "    section CPU\n"
		out += "    Min CPU :done, cpuMin, 00:00, 1min\n"
		out += "    Max CPU :crit, cpuMax, after cpuMin, 1min\n"
		out += "    section Memory\n"
		out += "    Min Mem :done, memMin, 00:02, 1min\n"
		out += "    Max Mem :crit, memMax, after memMin, 1min\n"
		out += "    section Disk\n"
		out += "    Sys Min :done, diskMin, 00:04, 1min\n"
		out += "    Docker Max :crit, dockerMax, after diskMin, 1min\n"
		out += "```\n\n"

		out += fmt.Sprintf("**Total datapoints:** `%d`  \n", len(metrics))
		out += fmt.Sprintf("**🕒 Time Range:** `%s` → `%s`  \n\n", metrics[0].Timestamp, metrics[len(metrics)-1].Timestamp)
		out += fmt.Sprintf("- **🔥 CPU:** `min: %.2f%%`, `max: %.2f%%`  \n", cpuMin, cpuMax)
		out += fmt.Sprintf("- **🧠 Memory:** `min: %.2f%%`, `max: %.2f%%`  \n", memMin, memMax)
		out += fmt.Sprintf("- **💽 System Disk:** `min: %.2f%%`, `max: %.2f%%`  \n", diskMin, diskMax)
		out += fmt.Sprintf("- **🐳 Docker Disk:** `min: %.2f%%`, `max: %.2f%%`\n\n", dockerMin, dockerMax)
		out += "---\n\n"

		out += "```mermaid\n"
		out += "xychart-beta\n"
		out += "title \"CPU and Memory Usage\"\n"
		out += fmt.Sprintf("x-axis [%s]\n", strings.Join(xLabels, ", "))
		out += "y-axis \"Usage (%)\"\n"
		out += fmt.Sprintf("bar [%s]\n", strings.Join(cpuSeries, ", "))
		out += fmt.Sprintf("line [%s]\n", strings.Join(memSeries, ", "))
		out += "```\n\n"

		out += "```mermaid\n"
		out += "xychart-beta\n"
		out += "title \"Disk Usage\"\n"
		out += fmt.Sprintf("x-axis [%s]\n", strings.Join(xLabels, ", "))
		out += "y-axis \"Disk Usage (%)\"\n"
		out += fmt.Sprintf("bar [%s]\n", strings.Join(sysDiskSeries, ", "))
		out += fmt.Sprintf("line [%s]\n", strings.Join(dockerDiskSeries, ", "))
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
	jobMetricsCmd.Flags().String("src", "/tmp/system-metrics", "source file to read system metrics from")
	rootCmd.AddCommand(jobMetricsCmd)
}
