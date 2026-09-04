package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/nikhil25803/purrpeek/internal/render"
	"github.com/nikhil25803/purrpeek/internal/system"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:           "purrpeek",
	Short:         "System information, from a cat's perspective. Purr approved.",
	Long:          "Purrpeek! System information, from a cat's perspective. Purr approved.",
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(cmd.Context(), 5*time.Second)
		defer cancel()

		info, collectionErr := system.GetSystemInformation(ctx)
		if verbose && collectionErr != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), collectionWarning(collectionErr))
		}
		if jsonOutput {
			return printJSON(cmd.OutOrStdout(), render.JSON(info))
		}
		return nil
	},
}

func collectionWarning(err error) string {
	return "warning: partial system information: " + strings.ReplaceAll(err.Error(), "\n", "; ")
}

func printJSON(output io.Writer, data any) error {
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func init() {
	rootCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show collection warnings")
}

func Execute() error {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	return rootCmd.Execute()
}
