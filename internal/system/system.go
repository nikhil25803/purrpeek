package system

import (
	"github.com/nikhil25803/purrpeek/internal/system/platform"
)

type SystemInfo struct {
	OS     *platform.OSInformation
	Uptime *platform.UptimeInformation
}

func GetSystemInformation() (*SystemInfo, error) {
	osInfo, err := platform.GetOSInformation()
	if err != nil {
		return nil, err
	}

	uptimeInfo, err := platform.GetUptimeInformation()
	if err != nil {
		return nil, err
	}

	return &SystemInfo{
		OS:     osInfo,
		Uptime: uptimeInfo,
	}, nil
}
