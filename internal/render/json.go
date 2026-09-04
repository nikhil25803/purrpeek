package render

import (
	"fmt"

	"github.com/nikhil25803/purrpeek/internal/system"
	"github.com/nikhil25803/purrpeek/internal/system/compute"
	"github.com/nikhil25803/purrpeek/internal/system/environment"
	"github.com/nikhil25803/purrpeek/internal/system/network"
	"github.com/nikhil25803/purrpeek/internal/system/platform"
	"github.com/nikhil25803/purrpeek/internal/system/power"
)

const gibibyte = 1024 * 1024 * 1024

type jsonReport struct {
	OS        *platform.OSInformation   `json:"os,omitempty"`
	Uptime    *jsonUptime               `json:"uptime,omitempty"`
	Time      *platform.TimeInfo        `json:"time,omitempty"`
	CPU       *jsonCPU                  `json:"cpu,omitempty"`
	GPUs      []compute.GPUInfo         `json:"gpus"`
	Memory    *jsonMemory               `json:"memory,omitempty"`
	Disk      *jsonDisk                 `json:"disk,omitempty"`
	Network   *network.NetworkInfo      `json:"network,omitempty"`
	Batteries []power.BatteryInfo       `json:"batteries"`
	Shell     *environment.ShellInfo    `json:"shell,omitempty"`
	Terminal  *environment.TerminalInfo `json:"terminal,omitempty"`
}

type jsonUptime struct {
	Duration string `json:"duration,omitempty"`
	BootTime string `json:"bootTime,omitempty"`
}

type jsonCPU struct {
	Model         string   `json:"model,omitempty"`
	PhysicalCores int      `json:"physicalCores"`
	LogicalCores  int      `json:"logicalCores"`
	UsagePercent  string   `json:"usagePercent"`
	FrequencyMHz  *float64 `json:"frequencyMHz,omitempty"`
}

type jsonMemory struct {
	Total       string `json:"total"`
	Used        string `json:"used"`
	Available   string `json:"available"`
	UsedPercent string `json:"usedPercent"`
}

type jsonDisk struct {
	HomeUsage *jsonDiskUsage `json:"homeUsage,omitempty"`
	Volumes   []jsonVolume   `json:"volumes"`
}

type jsonDiskUsage struct {
	Total       string `json:"total"`
	Used        string `json:"used"`
	Free        string `json:"free"`
	UsedPercent string `json:"usedPercent"`
}

type jsonVolume struct {
	FileSystem string `json:"fileSystem,omitempty"`
	MountPoint string `json:"mountPoint"`
	jsonDiskUsage
}

func JSON(info *system.SystemInfo) any {
	if info == nil {
		return nil
	}

	report := &jsonReport{
		OS:        info.OS,
		Time:      info.Time,
		GPUs:      info.GPUs,
		Network:   info.Network,
		Batteries: info.Batteries,
		Shell:     info.Shell,
		Terminal:  info.Terminal,
	}
	if info.Uptime != nil && (info.Uptime.DurationSeconds > 0 || info.Uptime.BootTime != "") {
		report.Uptime = &jsonUptime{BootTime: info.Uptime.BootTime}
		if info.Uptime.DurationSeconds > 0 {
			report.Uptime.Duration = formatDuration(info.Uptime.DurationSeconds)
		}
	}
	if info.CPU != nil {
		report.CPU = &jsonCPU{
			Model:         info.CPU.Model,
			PhysicalCores: info.CPU.PhysicalCores,
			LogicalCores:  info.CPU.LogicalCores,
			UsagePercent:  formatPercent(info.CPU.UsagePercent),
			FrequencyMHz:  info.CPU.FrequencyMHz,
		}
	}
	if info.Memory != nil {
		report.Memory = &jsonMemory{
			Total:       formatBytes(info.Memory.TotalBytes),
			Used:        formatBytes(info.Memory.UsedBytes),
			Available:   formatBytes(info.Memory.AvailableBytes),
			UsedPercent: formatPercent(info.Memory.UsedPercent),
		}
	}
	if info.Disk != nil {
		report.Disk = &jsonDisk{Volumes: make([]jsonVolume, len(info.Disk.Volumes))}
		if info.Disk.HomeUsage != nil {
			homeUsage := diskUsage(info.Disk.HomeUsage.TotalBytes, info.Disk.HomeUsage.UsedBytes, info.Disk.HomeUsage.FreeBytes, info.Disk.HomeUsage.UsedPercent)
			report.Disk.HomeUsage = &homeUsage
		}
		for index, volume := range info.Disk.Volumes {
			report.Disk.Volumes[index] = jsonVolume{
				FileSystem:    volume.FileSystem,
				MountPoint:    volume.MountPoint,
				jsonDiskUsage: diskUsage(volume.TotalBytes, volume.UsedBytes, volume.FreeBytes, volume.UsedPercent),
			}
		}
	}
	return report
}

func diskUsage(total, used, free uint64, percent float64) jsonDiskUsage {
	return jsonDiskUsage{
		Total:       formatBytes(total),
		Used:        formatBytes(used),
		Free:        formatBytes(free),
		UsedPercent: formatPercent(percent),
	}
}

func formatBytes(bytes uint64) string {
	return fmt.Sprintf("%.2f GiB", float64(bytes)/gibibyte)
}

func formatPercent(percent float64) string {
	return fmt.Sprintf("%.2f%%", percent)
}

func formatDuration(seconds uint64) string {
	days := seconds / 86400
	seconds %= 86400
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60
	return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
}
