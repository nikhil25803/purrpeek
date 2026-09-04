//go:build linux

package power

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func GetBatteryInformation(ctx context.Context) ([]BatteryInfo, error) {
	return getLinuxBatteries(ctx, "/sys/class/power_supply")
}

func getLinuxBatteries(ctx context.Context, root string) ([]BatteryInfo, error) {
	entries, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		return []BatteryInfo{}, nil
	}
	if err != nil {
		return []BatteryInfo{}, err
	}

	batteries := []BatteryInfo{}
	var errs []error
	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return normalizeBatteries(batteries), err
		}
		directory := filepath.Join(root, entry.Name())
		supplyType, err := os.ReadFile(filepath.Join(directory, "type"))
		if err != nil || strings.TrimSpace(string(supplyType)) != "Battery" {
			continue
		}
		capacity, err := os.ReadFile(filepath.Join(directory, "capacity"))
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", entry.Name(), err))
			continue
		}
		percentage, err := strconv.ParseFloat(strings.TrimSpace(string(capacity)), 64)
		if err != nil || percentage < 0 || percentage > 100 {
			errs = append(errs, fmt.Errorf("%s: invalid capacity", entry.Name()))
			continue
		}
		batteries = append(batteries, BatteryInfo{Name: entry.Name(), Percentage: percentage})
	}
	return normalizeBatteries(batteries), errors.Join(errs...)
}
