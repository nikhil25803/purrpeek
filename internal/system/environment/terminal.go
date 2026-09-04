package environment

import (
	"os"

	"golang.org/x/term"
)

type TerminalInfo struct {
	Name      string `json:"name,omitempty"`
	Version   string `json:"version,omitempty"`
	Term      string `json:"term,omitempty"`
	ColorTerm string `json:"colorTerm,omitempty"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

func GetTerminalInformation() *TerminalInfo {
	name, version := terminalIdentity(os.Getenv)
	width, height := terminalSize()

	return &TerminalInfo{
		Name:      name,
		Version:   version,
		Term:      os.Getenv("TERM"),
		ColorTerm: os.Getenv("COLORTERM"),
		Width:     width,
		Height:    height,
	}
}

func terminalIdentity(getenv func(string) string) (string, string) {
	if name := getenv("TERM_PROGRAM"); name != "" {
		return name, getenv("TERM_PROGRAM_VERSION")
	}
	if name := getenv("LC_TERMINAL"); name != "" {
		return name, getenv("LC_TERMINAL_VERSION")
	}
	if getenv("WT_SESSION") != "" {
		return "Windows Terminal", ""
	}
	if getenv("ConEmuPID") != "" || getenv("ConEmuANSI") != "" {
		return "ConEmu", getenv("ConEmuBuild")
	}
	if getenv("ANSICON") != "" {
		return "ANSICON", ""
	}
	return "", ""
}

func terminalSize() (int, int) {
	for _, file := range []*os.File{os.Stdout, os.Stderr, os.Stdin} {
		width, height, err := term.GetSize(int(file.Fd()))
		if err == nil && width > 0 && height > 0 {
			return width, height
		}
	}
	return 0, 0
}
