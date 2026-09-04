//go:build windows

package power

import "testing"

func TestParseWindowsBatteries(t *testing.T) {
	batteries, err := parseWindowsBatteries([]byte(`[{"Name":"Battery","EstimatedChargeRemaining":55}]`))
	if err != nil || len(batteries) != 1 || batteries[0].Percentage != 55 {
		t.Fatalf("parseWindowsBatteries() = %#v, %v", batteries, err)
	}
}
