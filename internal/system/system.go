package system

import (
	"github.com/nikhil25803/purrpeek/internal/system/compute"
	"github.com/nikhil25803/purrpeek/internal/system/memory"
	"github.com/nikhil25803/purrpeek/internal/system/platform"
	"github.com/nikhil25803/purrpeek/internal/system/storage"
)

type SystemInfo struct {
	OS         *platform.OSInformation
	Uptime     *platform.UptimeInformation
	CPUInfo    *compute.CPUInfo
	GPUs       []compute.GPUInfo
	MemoryInfo *memory.MemoryInfo
	DiskInfo   *storage.DiskInformation
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

	memoryInfo, err := memory.GetMemoryInformation()
	if err != nil {
		return nil, err
	}

	diskInfo, err := storage.GetDiskInformation()
	if err != nil {
		return nil, err
	}

	return &SystemInfo{
		OS:         osInfo,
		Uptime:     uptimeInfo,
		CPUInfo:    cpuInfo,
		GPUs:       gpus,
		MemoryInfo: memoryInfo,
		DiskInfo:   diskInfo,
	}, nil
}
