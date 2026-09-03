package system

import (
	"github.com/nikhil25803/purrpeek/internal/system/compute"
	"github.com/nikhil25803/purrpeek/internal/system/platform"
)

type SystemInfo struct {
	OS      *platform.OSInformation
	Uptime  *platform.UptimeInformation
	CPUInfo *compute.CPUInfo
	GPUs    []compute.GPUInfo
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

	cpuInfo, err := compute.GetCPUInformation()
	if err != nil {
		return nil, err
	}

	gpus := compute.GetGPUInformation()

	return &SystemInfo{
		OS:      osInfo,
		Uptime:  uptimeInfo,
		CPUInfo: cpuInfo,
		GPUs:    gpus,
	}, nil
}
