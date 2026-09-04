package compute

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

const cpuSampleDuration = 200 * time.Millisecond

type CPUInfo struct {
	Model         string   `json:"model,omitempty"`
	PhysicalCores int      `json:"physicalCores"`
	LogicalCores  int      `json:"logicalCores"`
	UsagePercent  float64  `json:"usagePercent"`
	FrequencyMHz  *float64 `json:"frequencyMHz,omitempty"`
}

func GetCPUInformation(ctx context.Context) (*CPUInfo, error) {
	info := &CPUInfo{}
	var errs []error

	stats, err := cpu.InfoWithContext(ctx)
	if err != nil {
		errs = append(errs, fmt.Errorf("hardware information: %w", err))
	} else {
		model, frequency, err := cpuHardware(stats)
		if err != nil {
			errs = append(errs, fmt.Errorf("hardware information: %w", err))
		} else {
			info.Model = model
			info.FrequencyMHz = frequency
		}
	}

	if info.LogicalCores, err = cpu.CountsWithContext(ctx, true); err != nil {
		errs = append(errs, fmt.Errorf("logical cores: %w", err))
	}
	if info.PhysicalCores, err = cpu.CountsWithContext(ctx, false); err != nil {
		errs = append(errs, fmt.Errorf("physical cores: %w", err))
	}
	percentages, err := cpu.PercentWithContext(ctx, cpuSampleDuration, false)
	if err != nil {
		errs = append(errs, fmt.Errorf("usage: %w", err))
	} else if len(percentages) == 0 {
		errs = append(errs, errors.New("usage: no CPU returned"))
	} else {
		info.UsagePercent = percentages[0]
	}

	return info, errors.Join(errs...)
}

func cpuHardware(stats []cpu.InfoStat) (string, *float64, error) {
	if len(stats) == 0 {
		return "", nil, errors.New("no CPU returned")
	}
	if stats[0].Mhz < 100 {
		return stats[0].ModelName, nil, nil
	}
	frequency := stats[0].Mhz
	return stats[0].ModelName, &frequency, nil
}
