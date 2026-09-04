package conf

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestEmbeddedConfig(t *testing.T) {
	config, err := parse(defaultYAML)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"mongo_no_bg.png", "snow_no_bg.png"}
	if !reflect.DeepEqual(config.Images, want) {
		t.Fatalf("images = %v, want %v", config.Images, want)
	}
	if !config.Render.OS.Name.Enabled || config.Render.OS.Name.Name != "OS" ||
		!config.Render.Terminal.Summary.Enabled {
		t.Fatalf("embedded render defaults were not loaded: %+v", config.Render)
	}
	if config.Render.OS.Hostname.Enabled || config.Render.Uptime.BootTime.Enabled || config.Render.CPU.UsagePercent.Enabled ||
		config.Render.Memory.Used.Enabled || config.Render.Disk.Volumes.MountPoint.Enabled ||
		config.Render.Network.Interfaces.Addresses.Enabled || config.Render.Shell.Path.Enabled ||
		config.Render.Terminal.Term.Enabled {
		t.Fatal("additional render fields must be disabled by default")
	}
	assertFieldsConfigured(t, reflect.ValueOf(config.Render), "render")
}

func assertFieldsConfigured(t *testing.T, value reflect.Value, path string) {
	t.Helper()
	fieldType := reflect.TypeFor[Field]()
	for index := range value.NumField() {
		field := value.Field(index)
		name := value.Type().Field(index).Tag.Get("yaml")
		if field.Type() == fieldType {
			configured := field.Interface().(Field)
			if configured.Name == "" || configured.Description == "" {
				t.Errorf("%s.%s is missing name or description", path, name)
			}
			continue
		}
		assertFieldsConfigured(t, field, path+"."+name)
	}
}

func TestParseNormalizesImages(t *testing.T) {
	config, err := parse([]byte("images:\n  - ' mongo_no_bg.png '\n  - mongo_no_bg.png\n  - ''\n"))
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"mongo_no_bg.png"}; !reflect.DeepEqual(config.Images, want) {
		t.Fatalf("images = %v, want %v", config.Images, want)
	}
	if _, err := parse([]byte("unknown: true\n")); err == nil {
		t.Fatal("unknown YAML field was accepted")
	}
	if _, err := parse([]byte("images: []\n---\nimages: []\n")); err == nil {
		t.Fatal("multiple YAML documents were accepted")
	}
}

func TestLoadUsesUserConfig(t *testing.T) {
	config, err := load(
		func() (string, error) { return "/config", nil },
		func(path string) ([]byte, error) {
			if want := filepath.Join("/config", "purrpeek", fileName); path != want {
				t.Fatalf("config path = %q", path)
			}
			return []byte("images:\n  - snow_purrpeek.png\n"), nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"snow_purrpeek.png"}; !reflect.DeepEqual(config.Images, want) {
		t.Fatalf("images = %v, want %v", config.Images, want)
	}
	if !config.Render.OS.Name.Enabled {
		t.Fatal("image-only user config did not inherit render defaults")
	}
}

func TestUserConfigOverridesRenderFields(t *testing.T) {
	config, err := parseMerged(defaultYAML, []byte("render:\n  os:\n    name:\n      name: Platform\n      enabled: false\n"))
	if err != nil {
		t.Fatal(err)
	}
	if config.Render.OS.Name.Name != "Platform" || config.Render.OS.Name.Enabled {
		t.Fatalf("OS name override = %+v", config.Render.OS.Name)
	}
	if config.Render.OS.Name.Description != "Operating system" || !config.Render.CPU.Model.Enabled {
		t.Fatal("partial override did not retain embedded defaults")
	}
	if _, err := parseMerged(defaultYAML, []byte("render:\n  os:\n    madeUp: {enabled: true}\n")); err == nil {
		t.Fatal("unknown render field was accepted")
	}
}

func TestLoadUsesEmbeddedConfigWhenUserConfigIsMissing(t *testing.T) {
	config, err := load(
		func() (string, error) { return "/config", nil },
		func(string) ([]byte, error) { return nil, os.ErrNotExist },
	)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"mongo_no_bg.png", "snow_no_bg.png"}; !reflect.DeepEqual(config.Images, want) {
		t.Fatalf("images = %v, want %v", config.Images, want)
	}
}

func TestLoadRejectsConfigurationErrors(t *testing.T) {
	tests := []struct {
		name      string
		directory func() (string, error)
		readFile  func(string) ([]byte, error)
	}{
		{
			name:      "malformed user config",
			directory: func() (string, error) { return "/config", nil },
			readFile:  func(string) ([]byte, error) { return []byte("images: ["), nil },
		},
		{
			name:      "unknown user field",
			directory: func() (string, error) { return "/config", nil },
			readFile:  func(string) ([]byte, error) { return []byte("render:\n  unknown: true\n"), nil },
		},
		{
			name:      "multiple documents",
			directory: func() (string, error) { return "/config", nil },
			readFile:  func(string) ([]byte, error) { return []byte("images: []\n---\nimages: []\n"), nil },
		},
		{
			name:      "unreadable user config",
			directory: func() (string, error) { return "/config", nil },
			readFile:  func(string) ([]byte, error) { return nil, errors.New("permission denied") },
		},
		{
			name:      "unavailable config directory",
			directory: func() (string, error) { return "", errors.New("unavailable") },
			readFile:  func(string) ([]byte, error) { return nil, nil },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, err := load(test.directory, test.readFile)
			if err == nil {
				t.Fatal("load() accepted invalid configuration")
			}
			if !reflect.DeepEqual(config, Config{}) {
				t.Fatalf("config = %+v, want zero value", config)
			}
		})
	}
}
