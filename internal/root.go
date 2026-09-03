package internal

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nikhil25803/purrpeek/internal/system"
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
		info, err := system.GetSystemInformation()
		if err != nil {
			return err
		}

		if jsonOutput {
			return printJSON(info)
		}

		return nil
	},
}

func printJSON(data any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")

	return encoder.Encode(data)
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
