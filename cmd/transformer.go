package cmd

import (
	"github.com/semaphoreci/test-results/pkg/cli"
	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/semaphoreci/test-results/pkg/transformer"
	"github.com/spf13/cobra"
)

// transformCmd represents the transform command
var transformCmd = &cobra.Command{
	Use:   "transform",
	Short: "transform transforms a JUnit XML file into another XML file according to template",
	Long:  `transform transforms a JUnit XML file into another XML file according to template`,

	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		tplPath, err := cmd.Flags().GetString("tpl")
		if err != nil {
			logger.Error("Failed to get template path: %v", err)
			return
		}

		template, err := transformer.LoadTemplate(tplPath)
		if err != nil {
			return
		}

		xml, err := transformer.LoadXML(args[0])
		if err != nil {
			return
		}

		output, err := transformer.Transform(template, xml)
		if err != nil {
			return
		}

		_, err = cli.WriteToFile([]byte(output), args[1])
		if err != nil {
			return
		}

		return
	},
}

func init() {
	desc := `path to template file`
	transformCmd.Flags().StringP("tpl", "t", "", desc)
	transformCmd.Flags().StringP("out", "o", "", desc)

	rootCmd.AddCommand(transformCmd)
}
