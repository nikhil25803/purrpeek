package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/nikhil25803/purrpeek/internal/asset"
	"github.com/nikhil25803/purrpeek/internal/conf"
	"github.com/nikhil25803/purrpeek/internal/render"
	"github.com/nikhil25803/purrpeek/internal/system"
	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	verbose    bool
)

const (
	imageColumns = 48
	imageRows    = 24
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
		if jsonOutput {
			if verbose && collectionErr != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), diagnosticWarning(collectionErr))
			}
			return printJSON(cmd.OutOrStdout(), render.JSON(info))
		}

		config, configErr := conf.Load()
		_, image, imageErr := asset.Select(config.Images)
		if verbose {
			if err := errors.Join(collectionErr, configErr, imageErr); err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), diagnosticWarning(err))
			}
		}
		return renderArtwork(cmd.OutOrStdout(), image, os.Getenv)
	},
}

func renderArtwork(output io.Writer, data []byte, getenv func(string) string) error {
	var prepared bytes.Buffer
	var err error
	switch graphicsProtocol(getenv) {
	case "kitty":
		err = render.KittyImage(&prepared, data, imageColumns, imageRows)
	case "iterm":
		err = render.ITermImage(&prepared, data, imageColumns, imageRows)
	default:
		return render.BrailleImage(output, data, imageColumns, imageRows)
	}
	if err != nil {
		return render.BrailleImage(output, data, imageColumns, imageRows)
	}
	if _, err := prepared.WriteTo(output); err != nil {
		return fmt.Errorf("render artwork: %w", err)
	}
	return nil
}

func graphicsProtocol(getenv func(string) string) string {
	program := strings.ToLower(getenv("TERM_PROGRAM"))
	term := strings.ToLower(getenv("TERM"))
	if getenv("KITTY_WINDOW_ID") != "" || getenv("WEZTERM_PANE") != "" ||
		program == "ghostty" || program == "wezterm" || strings.Contains(term, "kitty") {
		return "kitty"
	}
	if program == "iterm.app" || strings.EqualFold(getenv("LC_TERMINAL"), "iTerm2") {
		return "iterm"
	}
	return ""
}

func diagnosticWarning(err error) string {
	return "warning: partial output: " + strings.ReplaceAll(err.Error(), "\n", "; ")
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
