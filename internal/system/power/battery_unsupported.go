//go:build !darwin && !linux && !windows

package power

import "context"

func GetBatteryInformation(context.Context) ([]BatteryInfo, error) {
	return []BatteryInfo{}, nil
}
