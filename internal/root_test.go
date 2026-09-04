package internal

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/nikhil25803/purrpeek/internal/conf"
	"github.com/nikhil25803/purrpeek/internal/localisation"
	"github.com/spf13/cobra"
)

func TestCollectionWarningIsSingleLine(t *testing.T) {
	warning := diagnosticWarning(errors.Join(errors.New("cpu: unavailable"), errors.New("disk: unavailable")))
	if strings.Count(warning, "\n") != 0 || !strings.Contains(warning, "cpu: unavailable; disk: unavailable") {
		t.Fatalf("collectionWarning() = %q", warning)
	}
}

func TestPrintJSONDoesNotRenderImage(t *testing.T) {
	var output bytes.Buffer
	if err := printJSON(&output, map[string]string{"status": "ok"}); err != nil {
		t.Fatal(err)
	}
	if strings.ContainsRune(output.String(), '\x1b') {
		t.Fatal("JSON output contains a terminal graphics command")
	}
}

func TestGraphicsProtocol(t *testing.T) {
	tests := []struct {
		environment map[string]string
		want        string
	}{
		{map[string]string{"TERM_PROGRAM": "ghostty"}, "kitty"},
		{map[string]string{"KITTY_WINDOW_ID": "1", "TERM_PROGRAM": "Apple_Terminal"}, "kitty"},
		{map[string]string{"WEZTERM_PANE": "1"}, "kitty"},
		{map[string]string{"TERM_PROGRAM": "iTerm.app"}, "iterm"},
		{map[string]string{"LC_TERMINAL": "iTerm2"}, "iterm"},
		{map[string]string{"TERM_PROGRAM": "Apple_Terminal"}, ""},
		{map[string]string{"TERM_PROGRAM": "vscode"}, ""},
		{map[string]string{"WT_SESSION": "1"}, ""},
		{map[string]string{"TERM_PROGRAM": "unknown"}, ""},
	}

	for _, test := range tests {
		got := graphicsProtocol(func(key string) string { return test.environment[key] })
		if got != test.want {
			t.Fatalf("graphicsProtocol(%v) = %q, want %q", test.environment, got, test.want)
		}
	}
}

func TestArtworkSize(t *testing.T) {
	tests := []struct {
		name                      string
		terminalWidth, panelWidth int
		wantColumns, wantRows     int
	}{
		{name: "wide", terminalWidth: 120, panelWidth: 50, wantColumns: 36, wantRows: 18},
		{name: "constrained", terminalWidth: 80, panelWidth: 50, wantColumns: 26, wantRows: 13},
		{name: "very narrow", terminalWidth: 40, panelWidth: 50, wantColumns: 8, wantRows: 4},
		{name: "unknown", terminalWidth: 0, panelWidth: 50, wantColumns: 36, wantRows: 18},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			columns, rows := artworkSize(test.terminalWidth, test.panelWidth)
			if columns != test.wantColumns || rows != test.wantRows {
				t.Fatalf("artworkSize(%d, %d) = %dx%d, want %dx%d", test.terminalWidth, test.panelWidth, columns, rows, test.wantColumns, test.wantRows)
			}
		})
	}
}

func TestPlainOutputRejectsInvalidConfiguration(t *testing.T) {
	originalLoader, originalJSON := loadConfig, jsonOutput
	t.Cleanup(func() {
		loadConfig = originalLoader
		jsonOutput = originalJSON
	})
	loadConfig = func() (conf.Config, error) { return conf.Config{}, errors.New("invalid YAML") }
	jsonOutput = false

	err := rootCmd.RunE(&cobra.Command{}, nil)
	if err == nil || !strings.Contains(err.Error(), "load configuration: invalid YAML") {
		t.Fatalf("plain output error = %v", err)
	}
}

func TestPlainOutputRejectsInvalidGreetings(t *testing.T) {
	originalConfig, originalGreetings, originalJSON := loadConfig, loadGreetings, jsonOutput
	t.Cleanup(func() {
		loadConfig = originalConfig
		loadGreetings = originalGreetings
		jsonOutput = originalJSON
	})
	loadConfig = func() (conf.Config, error) { return conf.Config{}, nil }
	loadGreetings = func() (localisation.Catalog, error) { return nil, errors.New("invalid JSON") }
	jsonOutput = false

	err := rootCmd.RunE(&cobra.Command{}, nil)
	if err == nil || !strings.Contains(err.Error(), "load greetings: invalid JSON") {
		t.Fatalf("plain output error = %v", err)
	}
}

func TestJSONDoesNotLoadConfiguration(t *testing.T) {
	originalLoader, originalGreetings, originalJSON := loadConfig, loadGreetings, jsonOutput
	t.Cleanup(func() {
		loadConfig = originalLoader
		loadGreetings = originalGreetings
		jsonOutput = originalJSON
	})
	called := false
	loadConfig = func() (conf.Config, error) {
		called = true
		return conf.Config{}, errors.New("invalid YAML")
	}
	loadGreetings = func() (localisation.Catalog, error) {
		called = true
		return nil, errors.New("invalid greetings")
	}
	jsonOutput = true
	var output bytes.Buffer
	command := &cobra.Command{}
	command.SetContext(context.Background())
	command.SetOut(&output)

	if err := rootCmd.RunE(command, nil); err != nil {
		t.Fatal(err)
	}
	if called {
		t.Fatal("JSON output loaded artwork configuration")
	}
	if !strings.HasPrefix(output.String(), "{") {
		t.Fatalf("JSON output = %q", output.String())
	}
}
