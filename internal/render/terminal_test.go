package render

import (
	"bytes"
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
	"github.com/nikhil25803/purrpeek/internal/conf"
	"github.com/nikhil25803/purrpeek/internal/localisation"
	"github.com/nikhil25803/purrpeek/internal/system"
	"github.com/nikhil25803/purrpeek/internal/system/compute"
	"github.com/nikhil25803/purrpeek/internal/system/environment"
	"github.com/nikhil25803/purrpeek/internal/system/memory"
	"github.com/nikhil25803/purrpeek/internal/system/platform"
	"github.com/nikhil25803/purrpeek/internal/system/power"
	"github.com/nikhil25803/purrpeek/internal/system/storage"
)

func TestSystemPanel(t *testing.T) {
	field := func(name string) conf.Field { return conf.Field{Name: name, Enabled: true} }
	config := conf.RenderConfig{
		OS: conf.OSRender{
			Username: field("User"), Hostname: field("Host"), Name: field("OS"),
			Version: field("Version"), Architecture: field("Arch"), KernelVersion: field("Kernel"),
		},
		Uptime: conf.UptimeRender{Duration: field("Uptime")},
		Time:   conf.TimeRender{CurrentTime: field("Time"), TimeZone: field("Timezone")},
		CPU:    conf.CPURender{Model: field("CPU")}, GPUs: conf.GPUsRender{Models: field("GPU")},
		Memory:    conf.MemoryRender{Total: field("Memory"), UsedPercent: field("Memory Usage")},
		Disk:      conf.DiskRender{HomeUsage: conf.DiskUsageRender{Total: field("Disk"), UsedPercent: field("Disk Usage")}},
		Batteries: conf.BatteriesRender{Percentages: field("Battery")},
		Shell:     conf.ShellRender{Summary: field("Shell")}, Terminal: conf.TerminalRender{Summary: field("Terminal")},
	}
	info := &system.SystemInfo{
		OS:     &platform.OSInformation{Username: "nikhil", Hostname: "mac", Name: "macOS", Version: "26.6", Architecture: "arm64"},
		Uptime: &platform.UptimeInformation{DurationSeconds: 90061},
		Time:   &platform.TimeInfo{CurrentTime: "2026-09-04T15:00:35+05:30", TimeZone: "IST", UTCOffset: "+05:30"},
		CPU:    &compute.CPUInfo{Model: "Apple M4"}, GPUs: []compute.GPUInfo{{Model: "Apple M4"}, {Model: "Other"}},
		Memory:    &memory.MemoryInfo{TotalBytes: 16 * gibibyte, UsedPercent: 69.55},
		Disk:      &storage.DiskInformation{HomeUsage: &storage.DiskUsage{TotalBytes: 228 * gibibyte, UsedPercent: 76.83}},
		Batteries: []power.BatteryInfo{{Percentage: 15}, {Percentage: 82.5}},
		Shell:     &environment.ShellInfo{Name: "zsh", Version: "5.9"},
		Terminal:  &environment.TerminalInfo{Name: "ghostty", Version: "1.2"},
	}

	greetings := localisation.Catalog{"en": {Afternoon: []string{"Good afternoon"}}}
	panel := SystemPanel(info, config, greetings)
	output := strings.Join(panel.Lines, "\n")
	for _, want := range []string{
		yellow + "Good afternoon, nikhil" + reset,
		yellow + "----------------------" + reset,
		yellow + "Host:" + reset + " " + white + "mac" + reset,
		yellow + "OS:" + reset + " " + white + "macOS" + reset,
		"1d 1h 1m 1s", "2026-09-04 15:00:35", white + "IST" + reset,
		"Apple M4, Other", "16.00 GiB", "69.55%", "15%, 82.5%", "zsh 5.9", "ghostty 1.2",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("panel does not contain %q:\n%s", want, output)
		}
	}
	if strings.Contains(output, "Kernel:") {
		t.Fatal("panel rendered an unavailable value")
	}
	if strings.Contains(output, "@") {
		t.Fatalf("hostname leaked into greeting: %s", output)
	}
}

func TestSystemPanelHonorsLabelsAndEnabled(t *testing.T) {
	panel := SystemPanel(&system.SystemInfo{
		OS:       &platform.OSInformation{Name: "macOS", Version: "26"},
		CPU:      &compute.CPUInfo{Model: "Apple M4", LogicalCores: 10},
		Terminal: &environment.TerminalInfo{Name: "ghostty", Term: "xterm-256color"},
	}, conf.RenderConfig{
		Terminal: conf.TerminalRender{Term: conf.Field{Name: "TERM Value", Enabled: true}},
	}, nil)
	if len(panel.Lines) != 1 || !strings.Contains(panel.Lines[0], "TERM Value:") ||
		!strings.Contains(panel.Lines[0], "xterm-256color") {
		t.Fatalf("only enabled field should render: %q", panel.Lines)
	}
}

func TestSystemPanelUsesUnicodeDisplayWidth(t *testing.T) {
	info := &system.SystemInfo{
		OS:   &platform.OSInformation{Username: "nikhil"},
		Time: &platform.TimeInfo{CurrentTime: "2026-09-04T08:00:00+05:30"},
	}
	config := conf.RenderConfig{OS: conf.OSRender{Username: conf.Field{Enabled: true}}}
	panel := SystemPanel(info, config, localisation.Catalog{"hi": {Morning: []string{"शुभ प्रभात"}}})
	want := runewidth.StringWidth("शुभ प्रभात, nikhil")
	if panel.Width != want || len(panel.Lines) != 2 {
		t.Fatalf("Unicode panel width = %d, lines = %d, want %d and 2", panel.Width, len(panel.Lines), want)
	}
}

func TestPanelLayouts(t *testing.T) {
	data := encodePNG(t, image.NewUniform(color.White), image.Rect(0, 0, 2, 4))
	panel := TerminalPanel{Lines: []string{"details"}, Width: 7}

	var side bytes.Buffer
	if err := BraillePanel(&side, data, 1, 1, 1, 1, panel); err != nil {
		t.Fatal(err)
	}
	if got := side.String(); got != "\n⣿    details\n\x1b[0m\n" {
		t.Fatalf("side-by-side panel = %q", got)
	}

	var raster bytes.Buffer
	if err := ImagePanel(&raster, []byte("image"), 8, 4, 36, 18, panel); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(raster.String(), "\x1b[7B\x1b[14Cimage") ||
		!strings.Contains(raster.String(), "\x1b[7A\r\x1b[40Cdetails") || strings.Count(raster.String(), "\r\n") != 38 ||
		!strings.HasPrefix(raster.String(), "\r\n") ||
		!strings.HasSuffix(raster.String(), reset+"\r\n") {
		t.Fatalf("raster panel placement = %q", raster.String())
	}
}

func TestBraillePanelCentersColumns(t *testing.T) {
	data := encodePNG(t, image.NewUniform(color.White), image.Rect(0, 0, 2, 4))
	panel := TerminalPanel{Lines: []string{"one", "two", "three"}, Width: 5}
	var output bytes.Buffer
	if err := BraillePanel(&output, data, 1, 1, 5, 3, panel); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(output.String(), "\n")
	if lines[0] != "" || lines[1] != "         one" || lines[2] != "  ⣿      two" || lines[3] != "         three" {
		t.Fatalf("centered lines = %q", lines[:4])
	}
}
