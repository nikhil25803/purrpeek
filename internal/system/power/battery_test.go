package power

import "testing"

func TestNormalizeBatteries(t *testing.T) {
	batteries := normalizeBatteries([]BatteryInfo{
		{Name: " B ", Percentage: 75},
		{Name: "invalid", Percentage: 101},
		{Name: "A", Percentage: 0},
	})
	if len(batteries) != 2 || batteries[0].Name != "A" || batteries[1].Name != "B" {
		t.Fatalf("normalizeBatteries() = %#v", batteries)
	}
}
