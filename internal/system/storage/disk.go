package storage

import (
	"fmt"
	"os"
	"slices"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskUsage struct {
	TotalSpace      string
	UsedSpace       string
	FreeSpace       string
	UsagePercentage string
}

type DiskInfo struct {
	FileSystem string
	MountPoint string
	DiskUsage
}

type DiskInformation struct {
	Summary *DiskUsage
	Volumes []DiskInfo
}

func GetDiskInformation() (*DiskInformation, error) {
	info := &DiskInformation{Volumes: []DiskInfo{}}

	if home, err := os.UserHomeDir(); err == nil {
		if usage, err := disk.Usage(home); err == nil {
			summary := formatDiskUsage(usage)
			info.Summary = &summary
		}
	}

	partitions, _ := disk.Partitions(false)
	for _, partition := range partitions {
		if slices.Contains(partition.Opts, "nobrowse") {
			continue
		}

		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil || usage.Total == 0 {
			continue
		}

		info.Volumes = append(info.Volumes, DiskInfo{
			FileSystem: partition.Fstype,
			MountPoint: partition.Mountpoint,
			DiskUsage:  formatDiskUsage(usage),
		})
	}

	return info, nil
}

func formatDiskUsage(usage *disk.UsageStat) DiskUsage {
	const GB = 1024 * 1024 * 1024

	return DiskUsage{
		TotalSpace:      fmt.Sprintf("%.2f GB", float64(usage.Total)/GB),
		UsedSpace:       fmt.Sprintf("%.2f GB", float64(usage.Used)/GB),
		FreeSpace:       fmt.Sprintf("%.2f GB", float64(usage.Free)/GB),
		UsagePercentage: fmt.Sprintf("%.2f%%", usage.UsedPercent),
	}
}
