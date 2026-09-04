package compute

import (
	"testing"

	"github.com/shirou/gopsutil/v4/cpu"
)

func TestCPUHardware(t *testing.T) {
	if _, _, err := cpuHardware(nil); err == nil {
		t.Fatal("cpuHardware(nil) returned no error")
	}
	model, frequency, err := cpuHardware([]cpu.InfoStat{{ModelName: "Apple M4", Mhz: 4}})
	if err != nil || model != "Apple M4" || frequency != nil {
		t.Fatalf("implausible frequency = %q, %v, %v", model, frequency, err)
	}
	_, frequency, err = cpuHardware([]cpu.InfoStat{{ModelName: "Example", Mhz: 3200}})
	if err != nil || frequency == nil || *frequency != 3200 {
		t.Fatalf("valid frequency = %v, %v", frequency, err)
	}
}
