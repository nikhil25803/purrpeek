package compute

import (
	"math"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUInfo struct {
	CPUModel      string
	PhysicalCores int
	LogicalCores  int
	CPUUsage      float64
	CPUFrequency  float64
}

func GetCPUInformation() (*CPUInfo, error) {

	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	logicalCores, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}

	physicalCores, err := cpu.Counts(false)
	if err != nil {
		return nil, err
	}

	cpuPercentages, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	cpuFrequency, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	return &CPUInfo{
		CPUModel:      cpuInfo[0].ModelName,
		PhysicalCores: physicalCores,
		LogicalCores:  logicalCores,
		CPUUsage:      math.Round(cpuPercentages[0]*100) / 100,
		CPUFrequency:  cpuFrequency[0].Mhz,
	}, nil
}
