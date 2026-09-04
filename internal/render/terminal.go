package render

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/nikhil25803/purrpeek/internal/conf"
	"github.com/nikhil25803/purrpeek/internal/localisation"
	"github.com/nikhil25803/purrpeek/internal/system"
)

const (
	yellow = "\x1b[33m"
	white  = "\x1b[37m"
	reset  = "\x1b[0m"
)

type TerminalPanel struct {
	Lines []string
	Width int
}

func SystemPanel(info *system.SystemInfo, config conf.RenderConfig, greetings localisation.Catalog) TerminalPanel {
	if info == nil {
		return TerminalPanel{}
	}

	panel := TerminalPanel{}
	addPlain := func(text string) {
		if text == "" {
			return
		}
		panel.Width = max(panel.Width, runewidth.StringWidth(text))
		panel.Lines = append(panel.Lines, yellow+text+reset)
	}
	add := func(field conf.Field, value string) {
		name, value := strings.TrimSpace(field.Name), strings.TrimSpace(value)
		if !field.Enabled || name == "" || value == "" {
			return
		}
		plain := name + ": " + value
		panel.Width = max(panel.Width, runewidth.StringWidth(plain))
		panel.Lines = append(panel.Lines, yellow+name+":"+reset+" "+white+value+reset)
	}

	if info.OS != nil {
		if config.OS.Username.Enabled {
			currentTime := ""
			if info.Time != nil {
				currentTime = info.Time.CurrentTime
			}
			text := localisation.Greeting(greetings, currentTime, info.OS.Username)
			addPlain(text)
			addPlain(strings.Repeat("-", runewidth.StringWidth(text)))
		}
		add(config.OS.Hostname, info.OS.Hostname)
		add(config.OS.Name, info.OS.Name)
		add(config.OS.Version, info.OS.Version)
		add(config.OS.Architecture, info.OS.Architecture)
		add(config.OS.KernelVersion, info.OS.KernelVersion)
	}
	if info.Uptime != nil {
		if info.Uptime.DurationSeconds > 0 {
			add(config.Uptime.Duration, formatDuration(info.Uptime.DurationSeconds))
		}
		add(config.Uptime.BootTime, readableTime(info.Uptime.BootTime))
	}
	if info.Time != nil {
		add(config.Time.CurrentTime, readableTime(info.Time.CurrentTime))
		add(config.Time.TimeZone, info.Time.TimeZone)
		add(config.Time.UTCOffset, info.Time.UTCOffset)
	}
	if info.CPU != nil {
		add(config.CPU.Model, info.CPU.Model)
		if info.CPU.PhysicalCores > 0 {
			add(config.CPU.PhysicalCores, strconv.Itoa(info.CPU.PhysicalCores))
		}
		if info.CPU.LogicalCores > 0 {
			add(config.CPU.LogicalCores, strconv.Itoa(info.CPU.LogicalCores))
		}
		add(config.CPU.UsagePercent, formatPercent(info.CPU.UsagePercent))
		if info.CPU.FrequencyMHz != nil {
			add(config.CPU.FrequencyMHz, fmt.Sprintf("%.0f MHz", *info.CPU.FrequencyMHz))
		}
	}
	models := make([]string, 0, len(info.GPUs))
	for _, gpu := range info.GPUs {
		if model := strings.TrimSpace(gpu.Model); model != "" {
			models = append(models, model)
		}
	}
	add(config.GPUs.Models, strings.Join(models, ", "))
	if info.Memory != nil {
		if info.Memory.TotalBytes > 0 {
			add(config.Memory.Total, formatBytes(info.Memory.TotalBytes))
		}
		add(config.Memory.Used, formatBytes(info.Memory.UsedBytes))
		add(config.Memory.Available, formatBytes(info.Memory.AvailableBytes))
		add(config.Memory.UsedPercent, formatPercent(info.Memory.UsedPercent))
	}
	if info.Disk != nil && info.Disk.HomeUsage != nil {
		if info.Disk.HomeUsage.TotalBytes > 0 {
			add(config.Disk.HomeUsage.Total, formatBytes(info.Disk.HomeUsage.TotalBytes))
		}
		add(config.Disk.HomeUsage.Used, formatBytes(info.Disk.HomeUsage.UsedBytes))
		add(config.Disk.HomeUsage.Free, formatBytes(info.Disk.HomeUsage.FreeBytes))
		add(config.Disk.HomeUsage.UsedPercent, formatPercent(info.Disk.HomeUsage.UsedPercent))
	}
	if info.Disk != nil {
		fileSystems, mounts, totals, used, free, percentages := []string{}, []string{}, []string{}, []string{}, []string{}, []string{}
		for _, volume := range info.Disk.Volumes {
			fileSystems = appendNonEmpty(fileSystems, volume.FileSystem)
			mounts = appendNonEmpty(mounts, volume.MountPoint)
			totals = append(totals, formatBytes(volume.TotalBytes))
			used = append(used, formatBytes(volume.UsedBytes))
			free = append(free, formatBytes(volume.FreeBytes))
			percentages = append(percentages, formatPercent(volume.UsedPercent))
		}
		add(config.Disk.Volumes.FileSystem, strings.Join(fileSystems, ", "))
		add(config.Disk.Volumes.MountPoint, strings.Join(mounts, ", "))
		add(config.Disk.Volumes.Total, strings.Join(totals, ", "))
		add(config.Disk.Volumes.Used, strings.Join(used, ", "))
		add(config.Disk.Volumes.Free, strings.Join(free, ", "))
		add(config.Disk.Volumes.UsedPercent, strings.Join(percentages, ", "))
	}
	if info.Network != nil {
		add(config.Network.Hostname, info.Network.Hostname)
		add(config.Network.PrimaryInterface, info.Network.PrimaryInterface)
		add(config.Network.LocalIPv4, info.Network.LocalIPv4)
		add(config.Network.LocalIPv6, info.Network.LocalIPv6)
		add(config.Network.MACAddress, info.Network.MACAddress)
		names, addresses, macAddresses, mtus := []string{}, []string{}, []string{}, []string{}
		for _, networkInterface := range info.Network.Interfaces {
			names = appendNonEmpty(names, networkInterface.Name)
			addresses = append(addresses, networkInterface.Addresses...)
			macAddresses = appendNonEmpty(macAddresses, networkInterface.MACAddress)
			if networkInterface.MTU > 0 {
				mtus = append(mtus, strconv.Itoa(networkInterface.MTU))
			}
		}
		add(config.Network.Interfaces.Name, strings.Join(names, ", "))
		add(config.Network.Interfaces.Addresses, strings.Join(addresses, ", "))
		add(config.Network.Interfaces.MACAddress, strings.Join(macAddresses, ", "))
		add(config.Network.Interfaces.MTU, strings.Join(mtus, ", "))
	}
	names := make([]string, 0, len(info.Batteries))
	percentages := make([]string, 0, len(info.Batteries))
	for _, battery := range info.Batteries {
		names = appendNonEmpty(names, battery.Name)
		percentages = append(percentages, strconv.FormatFloat(battery.Percentage, 'f', -1, 64)+"%")
	}
	add(config.Batteries.Names, strings.Join(names, ", "))
	add(config.Batteries.Percentages, strings.Join(percentages, ", "))
	if info.Shell != nil {
		add(config.Shell.Summary, joinSummary(info.Shell.Name, info.Shell.Version))
		add(config.Shell.Name, info.Shell.Name)
		add(config.Shell.Version, info.Shell.Version)
		add(config.Shell.Path, info.Shell.Path)
	}
	if info.Terminal != nil {
		add(config.Terminal.Summary, joinSummary(info.Terminal.Name, info.Terminal.Version))
		add(config.Terminal.Name, info.Terminal.Name)
		add(config.Terminal.Version, info.Terminal.Version)
		add(config.Terminal.Term, info.Terminal.Term)
		add(config.Terminal.ColorTerm, info.Terminal.ColorTerm)
		if info.Terminal.Width > 0 {
			add(config.Terminal.Width, strconv.Itoa(info.Terminal.Width))
		}
		if info.Terminal.Height > 0 {
			add(config.Terminal.Height, strconv.Itoa(info.Terminal.Height))
		}
	}
	return panel
}

func readableTime(value string) string {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return value
	}
	return parsed.Format("2006-01-02 15:04:05")
}

func joinSummary(name, version string) string {
	return strings.TrimSpace(strings.TrimSpace(name) + " " + strings.TrimSpace(version))
}

func appendNonEmpty(values []string, value string) []string {
	if value = strings.TrimSpace(value); value != "" {
		return append(values, value)
	}
	return values
}

func BraillePanel(output io.Writer, data []byte, columns, rows, sectionColumns, sectionRows int, panel TerminalPanel) error {
	art, err := BrailleLines(data, columns, rows)
	if err != nil {
		return err
	}
	if _, err := io.WriteString(output, "\n"); err != nil {
		return fmt.Errorf("render panel: %w", err)
	}
	lineCount := max(sectionRows, len(panel.Lines))
	artTop := (sectionRows - rows) / 2
	artLeft := (sectionColumns - columns) / 2
	for index := range lineCount {
		line := ""
		artIndex := index - artTop
		if artIndex >= 0 && artIndex < len(art) {
			line = strings.Repeat(" ", artLeft) + art[artIndex]
		}
		if index < len(panel.Lines) {
			line += strings.Repeat(" ", sectionColumns-runewidth.StringWidth(line)+4) + panel.Lines[index]
		}
		if _, err := fmt.Fprintln(output, line); err != nil {
			return fmt.Errorf("render panel: %w", err)
		}
	}
	if _, err := io.WriteString(output, reset+"\n"); err != nil {
		return fmt.Errorf("render panel: %w", err)
	}
	return nil
}

func ImagePanel(output io.Writer, image []byte, columns, rows, sectionColumns, sectionRows int, panel TerminalPanel) error {
	lineCount := max(sectionRows, len(panel.Lines))
	artTop := (sectionRows - rows) / 2
	artLeft := (sectionColumns - columns) / 2
	if _, err := fmt.Fprintf(output, "\r\n\r%s\x1b[%dA", strings.Repeat("\r\n", lineCount), lineCount); err != nil {
		return fmt.Errorf("render panel: %w", err)
	}
	if artTop > 0 {
		if _, err := fmt.Fprintf(output, "\x1b[%dB", artTop); err != nil {
			return fmt.Errorf("render panel: %w", err)
		}
	}
	if artLeft > 0 {
		if _, err := fmt.Fprintf(output, "\x1b[%dC", artLeft); err != nil {
			return fmt.Errorf("render panel: %w", err)
		}
	}
	if _, err := output.Write(image); err != nil {
		return fmt.Errorf("render artwork: %w", err)
	}
	if _, err := io.WriteString(output, "\r"); err != nil {
		return fmt.Errorf("render panel: %w", err)
	}
	if artTop > 0 {
		if _, err := fmt.Fprintf(output, "\x1b[%dA", artTop); err != nil {
			return fmt.Errorf("render panel: %w", err)
		}
	}
	for _, line := range panel.Lines {
		if _, err := fmt.Fprintf(output, "\r\x1b[%dC%s\r\n", sectionColumns+4, line); err != nil {
			return fmt.Errorf("render panel: %w", err)
		}
	}
	if remaining := lineCount - len(panel.Lines); remaining > 0 {
		if _, err := io.WriteString(output, strings.Repeat("\r\n", remaining)); err != nil {
			return fmt.Errorf("render panel: %w", err)
		}
	}
	if _, err := io.WriteString(output, reset+"\r\n"); err != nil {
		return fmt.Errorf("render panel: %w", err)
	}
	return nil
}
