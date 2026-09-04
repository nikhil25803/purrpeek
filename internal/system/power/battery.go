package power

import (
	"slices"
	"strings"
)

type BatteryInfo struct {
	Name       string  `json:"name,omitempty"`
	Percentage float64 `json:"percentage"`
}

func normalizeBatteries(batteries []BatteryInfo) []BatteryInfo {
	valid := batteries[:0]
	for _, battery := range batteries {
		battery.Name = strings.TrimSpace(battery.Name)
		if battery.Percentage >= 0 && battery.Percentage <= 100 {
			valid = append(valid, battery)
		}
	}
	slices.SortFunc(valid, func(a, b BatteryInfo) int { return strings.Compare(a.Name, b.Name) })
	return valid
}
