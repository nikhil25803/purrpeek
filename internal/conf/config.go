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
	Images []string `yaml:"images"`
}

func Load() (Config, error) {
	return load(os.UserConfigDir, os.ReadFile)
}

func load(userConfigDir func() (string, error), readFile func(string) ([]byte, error)) (Config, error) {
	directory, err := userConfigDir()
	if err != nil {
		config, defaultErr := parse(defaultYAML)
		return config, errors.Join(defaultErr, fmt.Errorf("locate user config: %w", err))
	}

	path := filepath.Join(directory, "purrpeek", fileName)
	data, err := readFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return parse(defaultYAML)
	}
	if err != nil {
		config, defaultErr := parse(defaultYAML)
		return config, errors.Join(defaultErr, fmt.Errorf("read user config %q: %w", path, err))
	}

	config, err := parse(data)
	if err == nil {
		return config, nil
	}
	defaultConfig, defaultErr := parse(defaultYAML)
	return defaultConfig, errors.Join(fmt.Errorf("parse user config %q: %w", path, err), defaultErr)
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
