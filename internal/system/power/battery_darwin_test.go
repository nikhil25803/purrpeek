//go:build darwin

package power

import (
	"context"
	"testing"
)

func TestParsePMSetBatteries(t *testing.T) {
	output := []byte("Now drawing from 'Battery Power'\n -InternalBattery-0 (id=1)\t78%; discharging; 4:00 remaining present: true\n")
	batteries := parsePMSetBatteries(output)
	if len(batteries) != 1 || batteries[0].Name != "InternalBattery-0" || batteries[0].Percentage != 78 {
		t.Fatalf("parsePMSetBatteries() = %#v", batteries)
	}
}

func TestCanceledPMSet(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := GetBatteryInformation(ctx); err == nil {
		t.Fatal("GetBatteryInformation(canceled context) returned no error")
	}
}
