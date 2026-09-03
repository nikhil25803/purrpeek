package internal

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
)

var rootCmd = &cobra.Command{
	Use:   "purrpeek",
	Short: "System information, from a Cat's persepective! Purr Approved.",
	Long:  `Purrpeek! System information, from a Cat's persepective! Purr Approved.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("purrpeek")

		if jsonOutput {
			fmt.Println("JSON Output")
		}

		return nil
	},
}

func init() {

	rootCmd.Flags().BoolVar(
		&jsonOutput,
		"json",
		false,
		"Output in JSON format",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
