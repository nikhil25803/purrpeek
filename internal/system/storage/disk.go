package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskUsage struct {
	TotalBytes  uint64  `json:"totalBytes"`
	UsedBytes   uint64  `json:"usedBytes"`
	FreeBytes   uint64  `json:"freeBytes"`
	UsedPercent float64 `json:"usedPercent"`
}

type DiskInfo struct {
	FileSystem string `json:"fileSystem,omitempty"`
	MountPoint string `json:"mountPoint"`
	DiskUsage
}

type DiskInformation struct {
	HomeUsage *DiskUsage `json:"homeUsage,omitempty"`
	Volumes   []DiskInfo `json:"volumes"`
}

func GetDiskInformation(ctx context.Context) (*DiskInformation, error) {
	info := &DiskInformation{Volumes: []DiskInfo{}}
	var errs []error

	if home, err := os.UserHomeDir(); err != nil {
		errs = append(errs, fmt.Errorf("home directory: %w", err))
	} else if usage, err := disk.UsageWithContext(ctx, home); err != nil {
		errs = append(errs, fmt.Errorf("home usage: %w", err))
	} else {
		homeUsage := diskUsage(usage)
		info.HomeUsage = &homeUsage
	}

	partitions, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		errs = append(errs, fmt.Errorf("partitions: %w", err))
		return info, errors.Join(errs...)
	}

	skipped := 0
	for _, partition := range partitions {
		if slices.Contains(partition.Opts, "nobrowse") {
			continue
		}
		usage, err := disk.UsageWithContext(ctx, partition.Mountpoint)
		if err != nil {
			skipped++
			continue
		}
		if usage.Total == 0 {
			continue
		}
		info.Volumes = append(info.Volumes, DiskInfo{
			FileSystem: partition.Fstype,
			MountPoint: partition.Mountpoint,
			DiskUsage:  diskUsage(usage),
		})
	}
	if skipped > 0 {
		errs = append(errs, fmt.Errorf("%d volume(s) unavailable", skipped))
	}
	slices.SortFunc(info.Volumes, func(a, b DiskInfo) int {
		return strings.Compare(a.MountPoint, b.MountPoint)
	})
	return info, errors.Join(errs...)
}

func diskUsage(usage *disk.UsageStat) DiskUsage {
	return DiskUsage{
		TotalBytes:  usage.Total,
		UsedBytes:   usage.Used,
		FreeBytes:   usage.Free,
		UsedPercent: usage.UsedPercent,
	}
}
