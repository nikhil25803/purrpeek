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
}

func TestLoadFallsBackToEmbeddedConfig(t *testing.T) {
	tests := []struct {
		name      string
		directory func() (string, error)
		readFile  func(string) ([]byte, error)
		wantError bool
	}{
		{
			name:      "missing user config",
			directory: func() (string, error) { return "/config", nil },
			readFile:  func(string) ([]byte, error) { return nil, os.ErrNotExist },
		},
		{
			name:      "malformed user config",
			directory: func() (string, error) { return "/config", nil },
			readFile:  func(string) ([]byte, error) { return []byte("images: ["), nil },
			wantError: true,
		},
		{
			name:      "unavailable config directory",
			directory: func() (string, error) { return "", errors.New("unavailable") },
			readFile:  func(string) ([]byte, error) { return nil, nil },
			wantError: true,
		},
	}

	want := []string{"mongo_no_bg.png", "snow_no_bg.png"}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, err := load(test.directory, test.readFile)
			if (err != nil) != test.wantError {
				t.Fatalf("error = %v, wantError = %t", err, test.wantError)
			}
			if !reflect.DeepEqual(config.Images, want) {
				t.Fatalf("images = %v, want %v", config.Images, want)
			}
		})
	}
}
