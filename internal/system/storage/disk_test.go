package storage

import (
	"testing"

	"github.com/shirou/gopsutil/v4/disk"
)

func TestDiskUsageUsesRawValues(t *testing.T) {
	got := diskUsage(&disk.UsageStat{Total: 100, Used: 60, Free: 40, UsedPercent: 60})
	if got.TotalBytes != 100 || got.UsedBytes != 60 || got.FreeBytes != 40 || got.UsedPercent != 60 {
		t.Fatalf("diskUsage() = %#v", got)
	}
}
