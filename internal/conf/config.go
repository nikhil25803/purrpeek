package conf

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const fileName = "purrpeek-conf.yaml"

//go:embed purrpeek-conf.yaml
var defaultYAML []byte

type Config struct {
	Images []string     `yaml:"images"`
	Render RenderConfig `yaml:"render"`
}

type Field struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Enabled     bool   `yaml:"enabled"`
}

type RenderConfig struct {
	OS        OSRender        `yaml:"os"`
	Uptime    UptimeRender    `yaml:"uptime"`
	Time      TimeRender      `yaml:"time"`
	CPU       CPURender       `yaml:"cpu"`
	GPUs      GPUsRender      `yaml:"gpus"`
	Memory    MemoryRender    `yaml:"memory"`
	Disk      DiskRender      `yaml:"disk"`
	Network   NetworkRender   `yaml:"network"`
	Batteries BatteriesRender `yaml:"batteries"`
	Shell     ShellRender     `yaml:"shell"`
	Terminal  TerminalRender  `yaml:"terminal"`
}

type OSRender struct {
	Username      Field `yaml:"username"`
	Hostname      Field `yaml:"hostname"`
	Name          Field `yaml:"name"`
	Version       Field `yaml:"version"`
	Architecture  Field `yaml:"architecture"`
	KernelVersion Field `yaml:"kernelVersion"`
}

type UptimeRender struct {
	Duration Field `yaml:"duration"`
	BootTime Field `yaml:"bootTime"`
}

type TimeRender struct {
	CurrentTime Field `yaml:"currentTime"`
	TimeZone    Field `yaml:"timeZone"`
	UTCOffset   Field `yaml:"utcOffset"`
}

type CPURender struct {
	Model         Field `yaml:"model"`
	PhysicalCores Field `yaml:"physicalCores"`
	LogicalCores  Field `yaml:"logicalCores"`
	UsagePercent  Field `yaml:"usagePercent"`
	FrequencyMHz  Field `yaml:"frequencyMHz"`
}

type GPUsRender struct {
	Models Field `yaml:"models"`
}

type MemoryRender struct {
	Used        Field `yaml:"used"`
	Available   Field `yaml:"available"`
	Total       Field `yaml:"total"`
	UsedPercent Field `yaml:"usedPercent"`
}

type DiskRender struct {
	HomeUsage DiskUsageRender `yaml:"homeUsage"`
	Volumes   VolumeRender    `yaml:"volumes"`
}

type DiskUsageRender struct {
	Total       Field `yaml:"total"`
	Used        Field `yaml:"used"`
	Free        Field `yaml:"free"`
	UsedPercent Field `yaml:"usedPercent"`
}

type VolumeRender struct {
	FileSystem  Field `yaml:"fileSystem"`
	MountPoint  Field `yaml:"mountPoint"`
	Total       Field `yaml:"total"`
	Used        Field `yaml:"used"`
	Free        Field `yaml:"free"`
	UsedPercent Field `yaml:"usedPercent"`
}

type NetworkRender struct {
	Hostname         Field            `yaml:"hostname"`
	PrimaryInterface Field            `yaml:"primaryInterface"`
	LocalIPv4        Field            `yaml:"localIPv4"`
	LocalIPv6        Field            `yaml:"localIPv6"`
	MACAddress       Field            `yaml:"macAddress"`
	Interfaces       InterfacesRender `yaml:"interfaces"`
}

type InterfacesRender struct {
	Name       Field `yaml:"name"`
	Addresses  Field `yaml:"addresses"`
	MACAddress Field `yaml:"macAddress"`
	MTU        Field `yaml:"mtu"`
}

type BatteriesRender struct {
	Names       Field `yaml:"names"`
	Percentages Field `yaml:"percentages"`
}

type ShellRender struct {
	Summary Field `yaml:"summary"`
	Name    Field `yaml:"name"`
	Version Field `yaml:"version"`
	Path    Field `yaml:"path"`
}

type TerminalRender struct {
	Summary   Field `yaml:"summary"`
	Name      Field `yaml:"name"`
	Version   Field `yaml:"version"`
	Term      Field `yaml:"term"`
	ColorTerm Field `yaml:"colorTerm"`
	Width     Field `yaml:"width"`
	Height    Field `yaml:"height"`
}

func Load() (Config, error) {
	return load(os.UserConfigDir, os.ReadFile)
}

func load(userConfigDir func() (string, error), readFile func(string) ([]byte, error)) (Config, error) {
	directory, err := userConfigDir()
	if err != nil {
		return Config{}, fmt.Errorf("locate user config: %w", err)
	}

	path := filepath.Join(directory, "purrpeek", fileName)
	data, err := readFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return parse(defaultYAML)
	}
	if err != nil {
		return Config{}, fmt.Errorf("read user config %q: %w", path, err)
	}

	config, err := parseMerged(defaultYAML, data)
	if err == nil {
		return config, nil
	}
	return Config{}, fmt.Errorf("parse user config %q: %w", path, err)
}

func parseMerged(defaults, override []byte) (Config, error) {
	var base, patch yaml.Node
	if err := decodeNode(defaults, &base); err != nil {
		return Config{}, err
	}
	if err := decodeNode(override, &patch); err != nil {
		return Config{}, err
	}
	mergeMappings(base.Content[0], patch.Content[0])
	merged, err := yaml.Marshal(base.Content[0])
	if err != nil {
		return Config{}, err
	}
	return parse(merged)
}

func decodeNode(data []byte, node *yaml.Node) error {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(node); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			err = errors.New("multiple YAML documents are not supported")
		}
		return err
	}
	if len(node.Content) != 1 || node.Content[0].Kind != yaml.MappingNode {
		return errors.New("configuration must be a YAML mapping")
	}
	return nil
}

func mergeMappings(base, patch *yaml.Node) {
	for index := 0; index < len(patch.Content); index += 2 {
		key, value := patch.Content[index], patch.Content[index+1]
		matched := false
		for baseIndex := 0; baseIndex < len(base.Content); baseIndex += 2 {
			if base.Content[baseIndex].Value != key.Value {
				continue
			}
			baseValue := base.Content[baseIndex+1]
			if baseValue.Kind == yaml.MappingNode && value.Kind == yaml.MappingNode {
				mergeMappings(baseValue, value)
			} else {
				base.Content[baseIndex+1] = value
			}
			matched = true
			break
		}
		if !matched {
			base.Content = append(base.Content, key, value)
		}
	}
}

func parse(data []byte) (Config, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	var config Config
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			err = errors.New("multiple YAML documents are not supported")
		}
		return Config{}, err
	}

	seen := make(map[string]struct{}, len(config.Images))
	images := make([]string, 0, len(config.Images))
	for _, name := range config.Images {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		images = append(images, name)
	}
	config.Images = images
	return config, nil
}
