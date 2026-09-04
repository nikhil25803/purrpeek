package internal

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/nikhil25803/purrpeek/internal/conf"
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

func TestJSONDoesNotLoadConfiguration(t *testing.T) {
	originalLoader, originalJSON := loadConfig, jsonOutput
	t.Cleanup(func() {
		loadConfig = originalLoader
		jsonOutput = originalJSON
	})
	called := false
	loadConfig = func() (conf.Config, error) {
		called = true
		return conf.Config{}, errors.New("invalid YAML")
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
