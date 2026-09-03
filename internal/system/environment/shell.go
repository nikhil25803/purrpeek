package environment

import (
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strings"
)

var versionPattern = regexp.MustCompile(`\d+(?:\.\d+)+`)

type ShellInfo struct {
	Name    string
	Version string
	Path    string
}

func GetShellInformation() *ShellInfo {
	shellPath := configuredShell(runtime.GOOS, os.Getenv)
	if shellPath != "" {
		if resolvedPath, err := exec.LookPath(shellPath); err == nil {
			shellPath = resolvedPath
		}
	}

	name := shellName(shellPath)
	return &ShellInfo{
		Name:    name,
		Version: shellVersion(shellPath, name),
		Path:    shellPath,
	}
}

func configuredShell(goos string, getenv func(string) string) string {
	if shell := getenv("SHELL"); shell != "" {
		return shell
	}
	if goos == "windows" {
		return getenv("ComSpec")
	}
	return ""
}

func shellName(shellPath string) string {
	if shellPath == "" {
		return ""
	}
	name := path.Base(strings.ReplaceAll(shellPath, `\`, "/"))
	if strings.EqualFold(path.Ext(name), ".exe") {
		name = strings.TrimSuffix(name, path.Ext(name))
	}
	return name
}

func shellVersion(shellPath, name string) string {
	if shellPath == "" {
		return ""
	}

	output, err := exec.Command(shellPath, shellVersionArgs(name)...).CombinedOutput()
	if err != nil {
		return ""
	}
	return extractVersion(string(output))
}

func shellVersionArgs(name string) []string {
	switch strings.ToLower(name) {
	case "powershell", "pwsh":
		return []string{"-NoProfile", "-Command", "$PSVersionTable.PSVersion.ToString()"}
	case "cmd":
		return []string{"/C", "ver"}
	default:
		return []string{"--version"}
	}
}

func extractVersion(output string) string {
	return versionPattern.FindString(output)
}
