package environment

import (
	"reflect"
	"testing"
)

func TestShellHelpers(t *testing.T) {
	if got := configuredShell("windows", func(key string) string {
		if key == "ComSpec" {
			return `C:\Windows\cmd.exe`
		}
		return ""
	}); got != `C:\Windows\cmd.exe` {
		t.Fatalf("configuredShell() = %q", got)
	}
	if got := shellName(`C:\Program Files\pwsh.EXE`); got != "pwsh" {
		t.Fatalf("shellName() = %q", got)
	}
	if got := extractVersion("GNU bash, version 5.2.37(1)-release"); got != "5.2.37" {
		t.Fatalf("extractVersion() = %q", got)
	}
	if got, want := shellVersionArgs("cmd"), []string{"/C", "ver"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("shellVersionArgs() = %q", got)
	}
}

func TestTerminalIdentity(t *testing.T) {
	env := map[string]string{
		"TERM_PROGRAM":         "vscode",
		"TERM_PROGRAM_VERSION": "1.2.3",
		"WT_SESSION":           "ignored",
	}
	name, version := terminalIdentity(func(key string) string { return env[key] })
	if name != "vscode" || version != "1.2.3" {
		t.Fatalf("terminalIdentity() = %q, %q", name, version)
	}
}
