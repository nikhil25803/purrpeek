package compute

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"runtime"
	"slices"
	"strings"

	"github.com/jaypipes/ghw"
)

const unknownGPUModel = "Unknown GPU"

type GPUInfo struct {
	Model string `json:"model"`
}

func GetGPUInformation(ctx context.Context) ([]GPUInfo, error) {
	if runtime.GOOS == "darwin" {
		output, err := exec.CommandContext(
			ctx,
			"system_profiler",
			"SPDisplaysDataType",
			"-detailLevel", "mini",
			"-json",
		).Output()
		if err != nil {
			return []GPUInfo{}, err
		}

		gpus, err := parseMacOSGPUInformation(output)
		if err != nil {
			return []GPUInfo{}, err
		}
		sortGPUs(gpus)
		return gpus, nil
	}

	if runtime.GOOS != "linux" && runtime.GOOS != "windows" {
		return []GPUInfo{}, nil
	}

	info, err := ghw.GPU(ghw.WithDisableWarnings())
	if err != nil {
		return []GPUInfo{}, err
	}
	if info == nil {
		return []GPUInfo{}, errors.New("no GPU information returned")
	}
	gpus := gpuInformationFromCards(info.GraphicsCards)
	sortGPUs(gpus)
	return gpus, nil
}

func parseMacOSGPUInformation(data []byte) ([]GPUInfo, error) {
	var report struct {
		Displays []struct {
			Name  string `json:"_name"`
			Model string `json:"sppci_model"`
		} `json:"SPDisplaysDataType"`
	}
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, err
	}

	gpus := make([]GPUInfo, 0, len(report.Displays))
	for _, display := range report.Displays {
		model := strings.TrimSpace(display.Model)
		if model == "" {
			model = strings.TrimSpace(display.Name)
		}
		if model == "" {
			model = unknownGPUModel
		}
		gpus = append(gpus, GPUInfo{Model: model})
	}
	return gpus, nil
}

func gpuInformationFromCards(cards []*ghw.GraphicsCard) []GPUInfo {
	gpus := make([]GPUInfo, 0, len(cards))
	for _, card := range cards {
		model := ""
		if card != nil && card.DeviceInfo != nil && card.DeviceInfo.Product != nil {
			model = strings.TrimSpace(card.DeviceInfo.Product.Name)
		}
		if model == "" {
			model = unknownGPUModel
		}
		gpus = append(gpus, GPUInfo{Model: model})
	}
	return gpus
}

func sortGPUs(gpus []GPUInfo) {
	slices.SortFunc(gpus, func(a, b GPUInfo) int { return strings.Compare(a.Model, b.Model) })
}
