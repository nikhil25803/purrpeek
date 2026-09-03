package memory

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryInfo struct {
	RAMTotal     string
	RAMUsed      string
	RAMAvailable string
	RAMFree      string
	RAMUsage     string
}

func GetMemoryInformation() (*MemoryInfo, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	const GB = 1024 * 1024 * 1024

	return &MemoryInfo{
		RAMTotal:     fmt.Sprintf("%.2f GB", float64(vmStat.Total)/GB),
		RAMUsed:      fmt.Sprintf("%.2f GB", float64(vmStat.Used)/GB),
		RAMAvailable: fmt.Sprintf("%.2f GB", float64(vmStat.Available)/GB),
		RAMFree:      fmt.Sprintf("%.2f GB", float64(vmStat.Free)/GB),
		RAMUsage:     fmt.Sprintf("%.2f%%", vmStat.UsedPercent),
	}, nil
}
