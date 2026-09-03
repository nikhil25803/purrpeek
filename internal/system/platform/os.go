package platform

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v3/host"
)

type OSInformation struct {
	Username      string
	Hostname      string
	OS            string
	KernelVersion string
	Architecture  string
}

func GetOSInformation() (*OSInformation, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	osName := runtime.GOOS

	osDetails, err := GetOSDetails(osName)
	if err != nil {
		return nil, err
	}

	kernelVersion, err := host.KernelVersion()
	if err != nil {
		return nil, err
	}

	architecture := runtime.GOARCH

	return &OSInformation{
		Username:      currentUser.Username,
		Hostname:      hostname,
		OS:            osDetails,
		KernelVersion: kernelVersion,
		Architecture:  architecture,
	}, nil
}

func GetOSDetails(osname string) (string, error) {
	switch osname {
	case "windows":
		return getWindowsVersion()
	case "darwin":
		return getMacVersion()
	case "linux":
		return getLinuxVersion()
	default:
		return "Other", nil
	}
}

func getLinuxVersion() (string, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "", fmt.Errorf("failed to open os-release: %w", err)
	}
	defer file.Close()

	var name, version string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {

			name = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
			return name, nil
		}

		if strings.HasPrefix(line, "NAME=") {
			name = strings.Trim(strings.TrimPrefix(line, "NAME="), "\"")
		} else if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read os-release: %w", err)
	}

	if name != "" {
		return fmt.Sprintf("%s %s", name, version), nil
	}
	return "Unknown Linux", nil
}

func getMacVersion() (string, error) {
	cmd := exec.Command("sw_vers", "-productVersion")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to execute sw_vers: %w", err)
	}

	return "macOS " + strings.TrimSpace(out.String()), nil
}

func getWindowsVersion() (string, error) {

	args := []string{"-NoProfile", "-Command", "(Get-ItemProperty 'HKLM:\\SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion').ProductName"}

	cmd := exec.Command("powershell", args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get windows version: %w", err)
	}

	return strings.TrimSpace(out.String()), nil
}
