package cmd

/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"github.com/semaphoreci/test-results/pkg/cli"
	"github.com/semaphoreci/test-results/pkg/logger"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/spf13/cobra"
)

// combineCmd represents the combine command
func NewCombineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "combine <json-file-path>... <json-file>]",
		Short: "combines multiples json summary files into one",
		Long:  `Combines multiples json summary files into one"`,
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputs := args[:len(args)-1]
			output := args[len(args)-1]

			err := cli.SetLogLevel(cmd)
			if err != nil {
				return err
			}

			paths, err := cli.LoadFiles(inputs, ".json")
			if err != nil {
				return err
			}

			result := parser.NewResult()
			for _, path := range paths {
				inFile, err := cli.CheckFile(path)
				if err != nil {
					logger.Error(err.Error())
					return err
				}

				newResult, err := cli.Load(inFile)

				if err != nil {
					logger.Error(err.Error())
					return err
				}
				result.Combine(*newResult)
			}

			jsonData, err := cli.Marshal(result)
			if err != nil {
				return err
			}

			_, err = cli.WriteToFile(jsonData, output)
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.test-results.yaml)")
	cmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	cmd.PersistentFlags().BoolP("trace", "", false, "trace output")
	cmd.PersistentFlags().StringP("name", "N", "", "name of the suite")
	cmd.PersistentFlags().StringP("suite-prefix", "S", "", "prefix for each suite")
	cmd.PersistentFlags().StringP("parser", "p", "auto", "override parser to be used")
	return cmd
}

func init() {
	combineCmd := NewCombineCmd()
	rootCmd.AddCommand(combineCmd)
}
