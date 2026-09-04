package render

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/nikhil25803/purrpeek/internal/system"
	"github.com/nikhil25803/purrpeek/internal/system/compute"
	"github.com/nikhil25803/purrpeek/internal/system/memory"
	"github.com/nikhil25803/purrpeek/internal/system/platform"
	"github.com/nikhil25803/purrpeek/internal/system/power"
	"github.com/nikhil25803/purrpeek/internal/system/storage"
)

func TestFormatters(t *testing.T) {
	for _, test := range []struct {
		name string
		got  string
		want string
	}{
		{"zero bytes", formatBytes(0), "0.00 GiB"},
		{"bytes", formatBytes(17179869184), "16.00 GiB"},
		{"percent", formatPercent(6.341463409758692), "6.34%"},
		{"zero duration", formatDuration(0), "0d 0h 0m 0s"},
		{"duration", formatDuration(266881), "3d 2h 8m 1s"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("got %q, want %q", test.got, test.want)
			}
		})
	}
}

func TestJSONPresentation(t *testing.T) {
	data, err := json.Marshal(JSON(&system.SystemInfo{
		Uptime: &platform.UptimeInformation{DurationSeconds: 266881, BootTime: "2026-09-01T07:34:23+05:30"},
		CPU:    &compute.CPUInfo{Model: "Example", UsagePercent: 6.341463409758692},
		GPUs:   []compute.GPUInfo{},
		Memory: &memory.MemoryInfo{
			TotalBytes:     17179869184,
			UsedBytes:      11177099264,
			AvailableBytes: 6002769920,
			UsedPercent:    65.05928039550781,
		},
		Disk: &storage.DiskInformation{
			HomeUsage: &storage.DiskUsage{TotalBytes: 17179869184, UsedBytes: 10737418240, FreeBytes: 6442450944, UsedPercent: 62.5},
			Volumes:   []storage.DiskInfo{},
		},
		Batteries: []power.BatteryInfo{},
	}))
	if err != nil {
		t.Fatal(err)
	}
	jsonText := string(data)
	for _, value := range []string{
		`"duration":"3d 2h 8m 1s"`,
		`"usagePercent":"6.34%"`,
		`"total":"16.00 GiB"`,
		`"usedPercent":"65.06%"`,
		`"gpus":[]`,
		`"batteries":[]`,
		`"volumes":[]`,
	} {
		if !strings.Contains(jsonText, value) {
			t.Errorf("JSON %s missing %s", jsonText, value)
		}
	}
	for _, oldKey := range []string{`"durationSeconds"`, `"totalBytes"`, `"OS"`} {
		if strings.Contains(jsonText, oldKey) {
			t.Errorf("JSON %s contains raw key %s", jsonText, oldKey)
		}
	}
}

func TestJSONNilAndPartial(t *testing.T) {
	if JSON(nil) != nil {
		t.Fatal("JSON(nil) should return nil")
	}
	data, err := json.Marshal(JSON(&system.SystemInfo{
		Uptime:    &platform.UptimeInformation{},
		GPUs:      []compute.GPUInfo{},
		Batteries: []power.BatteryInfo{},
	}))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(data), `{"gpus":[],"batteries":[]}`; got != want {
		t.Fatalf("partial JSON = %s, want %s", got, want)
	}
}
