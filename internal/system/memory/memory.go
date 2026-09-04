package memory

import (
	"context"

	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryInfo struct {
	TotalBytes     uint64  `json:"totalBytes"`
	UsedBytes      uint64  `json:"usedBytes"`
	AvailableBytes uint64  `json:"availableBytes"`
	UsedPercent    float64 `json:"usedPercent"`
}

func GetMemoryInformation(ctx context.Context) (*MemoryInfo, error) {
	stat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, err
	}

	return &MemoryInfo{
		TotalBytes:     stat.Total,
		UsedBytes:      stat.Used,
		AvailableBytes: stat.Available,
		UsedPercent:    stat.UsedPercent,
	}, nil
}
