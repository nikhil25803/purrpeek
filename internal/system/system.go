package system

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikhil25803/purrpeek/internal/system/compute"
	"github.com/nikhil25803/purrpeek/internal/system/environment"
	"github.com/nikhil25803/purrpeek/internal/system/memory"
	"github.com/nikhil25803/purrpeek/internal/system/network"
	"github.com/nikhil25803/purrpeek/internal/system/platform"
	"github.com/nikhil25803/purrpeek/internal/system/power"
	"github.com/nikhil25803/purrpeek/internal/system/storage"
)

type SystemInfo struct {
	OS        *platform.OSInformation     `json:"os,omitempty"`
	Uptime    *platform.UptimeInformation `json:"uptime,omitempty"`
	Time      *platform.TimeInfo          `json:"time,omitempty"`
	CPU       *compute.CPUInfo            `json:"cpu,omitempty"`
	GPUs      []compute.GPUInfo           `json:"gpus"`
	Memory    *memory.MemoryInfo          `json:"memory,omitempty"`
	Disk      *storage.DiskInformation    `json:"disk,omitempty"`
	Network   *network.NetworkInfo        `json:"network,omitempty"`
	Batteries []power.BatteryInfo         `json:"batteries"`
	Shell     *environment.ShellInfo      `json:"shell,omitempty"`
	Terminal  *environment.TerminalInfo   `json:"terminal,omitempty"`
}

func GetSystemInformation(ctx context.Context) (*SystemInfo, error) {
	info := &SystemInfo{
		Time:      platform.GetTimeInformation(),
		GPUs:      []compute.GPUInfo{},
		Batteries: []power.BatteryInfo{},
		Shell:     environment.GetShellInformation(ctx),
		Terminal:  environment.GetTerminalInformation(),
	}
	var errs []error
	collect := func(component string, err error) {
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", component, err))
		}
	}

	var err error
	info.OS, err = platform.GetOSInformation(ctx)
	collect("os", err)
	info.Uptime, err = platform.GetUptimeInformation(ctx)
	collect("uptime", err)
	info.CPU, err = compute.GetCPUInformation(ctx)
	collect("cpu", err)
	info.GPUs, err = compute.GetGPUInformation(ctx)
	collect("gpus", err)
	info.Memory, err = memory.GetMemoryInformation(ctx)
	collect("memory", err)
	info.Disk, err = storage.GetDiskInformation(ctx)
	collect("disk", err)
	info.Network, err = network.GetNetworkInformation()
	collect("network", err)
	info.Batteries, err = power.GetBatteryInformation(ctx)
	collect("batteries", err)

	return info, errors.Join(errs...)
}
