package internal

import (
	"bytes"
	"errors"
	"strings"
	"testing"
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
