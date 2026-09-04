//go:build windows

package power

import (
	"context"
	"encoding/json"
	"os/exec"
)

const windowsBatteryQuery = `ConvertTo-Json -Compress -InputObject @(Get-CimInstance Win32_Battery | Select-Object Name, EstimatedChargeRemaining)`

func GetBatteryInformation(ctx context.Context) ([]BatteryInfo, error) {
	output, err := exec.CommandContext(ctx, "powershell.exe", "-NoProfile", "-NonInteractive", "-Command", windowsBatteryQuery).Output()
	if err != nil {
		return []BatteryInfo{}, err
	}
	return parseWindowsBatteries(output)
}

func parseWindowsBatteries(output []byte) ([]BatteryInfo, error) {
	var values []struct {
		Name                     string  `json:"Name"`
		EstimatedChargeRemaining float64 `json:"EstimatedChargeRemaining"`
	}
	if err := json.Unmarshal(output, &values); err != nil {
		return []BatteryInfo{}, err
	}
	batteries := make([]BatteryInfo, 0, len(values))
	for _, value := range values {
		batteries = append(batteries, BatteryInfo{Name: value.Name, Percentage: value.EstimatedChargeRemaining})
	}
	return normalizeBatteries(batteries), nil
}
