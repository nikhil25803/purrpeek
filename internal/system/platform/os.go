package platform

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/user"
	"runtime"

	"github.com/shirou/gopsutil/v4/host"
)

type OSInformation struct {
	Username      string `json:"username,omitempty"`
	Hostname      string `json:"hostname,omitempty"`
	Name          string `json:"name,omitempty"`
	Version       string `json:"version,omitempty"`
	KernelVersion string `json:"kernelVersion,omitempty"`
	Architecture  string `json:"architecture,omitempty"`
}

func GetOSInformation(ctx context.Context) (*OSInformation, error) {
	info := &OSInformation{}
	var errs []error

	currentUser, err := user.Current()
	if err != nil {
		errs = append(errs, fmt.Errorf("username: %w", err))
	} else {
		info.Username = currentUser.Username
	}
	if info.Hostname, err = os.Hostname(); err != nil {
		errs = append(errs, fmt.Errorf("hostname: %w", err))
	}

	platform, _, version, err := host.PlatformInformationWithContext(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("platform: %w", err))
	} else {
		info.Name = displayOSName(runtime.GOOS, platform)
		info.Version = version
	}
	if info.KernelVersion, err = host.KernelVersionWithContext(ctx); err != nil {
		errs = append(errs, fmt.Errorf("kernel: %w", err))
	}
	if info.Architecture, err = host.KernelArch(); err != nil {
		info.Architecture = runtime.GOARCH
		errs = append(errs, fmt.Errorf("architecture: %w", err))
	}

	return info, errors.Join(errs...)
}

func displayOSName(goos, platform string) string {
	if goos == "darwin" {
		return "macOS"
	}
	if platform != "" {
		return platform
	}
	return goos
}
