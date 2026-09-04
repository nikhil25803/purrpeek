//go:build darwin

package power

import (
	"context"
	"os/exec"
	"regexp"
	"strconv"
)

var pmsetBatteryPattern = regexp.MustCompile(`(?m)^\s*-([^(]+?)\s+\([^)]*\)\s+([0-9]+(?:\.[0-9]+)?)%;`)

func GetBatteryInformation(ctx context.Context) ([]BatteryInfo, error) {
	output, err := exec.CommandContext(ctx, "pmset", "-g", "batt").Output()
	if err != nil {
		return []BatteryInfo{}, err
	}
	return parsePMSetBatteries(output), nil
}

func parsePMSetBatteries(output []byte) []BatteryInfo {
	batteries := []BatteryInfo{}
	for _, match := range pmsetBatteryPattern.FindAllSubmatch(output, -1) {
		percentage, err := strconv.ParseFloat(string(match[2]), 64)
		if err == nil {
			batteries = append(batteries, BatteryInfo{Name: string(match[1]), Percentage: percentage})
		}
	}
	return normalizeBatteries(batteries)
}
