package compute

import (
	"encoding/json"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jaypipes/ghw"
)

const unknownGPUModel = "Unknown GPU"

type GPUInfo struct {
	GPUModel string
}

func GetGPUInformation() []GPUInfo {
	if runtime.GOOS == "darwin" {
		output, err := exec.Command(
			"system_profiler",
			"SPDisplaysDataType",
			"-detailLevel", "mini",
			"-json",
		).Output()
		if err != nil {
			return []GPUInfo{}
		}

		gpus, err := parseMacOSGPUInformation(output)
		if err != nil {
			return []GPUInfo{}
		}
		return gpus
	}

	if runtime.GOOS != "linux" && runtime.GOOS != "windows" {
		return []GPUInfo{}
	}

	info, err := ghw.GPU(ghw.WithDisableWarnings())
	if err != nil {
		return []GPUInfo{}
	}
	return gpuInformationFromCards(info.GraphicsCards)
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
		gpus = append(gpus, GPUInfo{GPUModel: model})
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
		gpus = append(gpus, GPUInfo{GPUModel: model})
	}
	return gpus
}
