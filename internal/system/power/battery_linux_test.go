//go:build linux

package power

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestGetLinuxBatteries(t *testing.T) {
	root := t.TempDir()
	for name, capacity := range map[string]string{"BAT1": "80", "BAT0": "40"} {
		directory := filepath.Join(root, name)
		if err := os.Mkdir(directory, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(directory, "type"), []byte("Battery\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(directory, "capacity"), []byte(capacity), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	batteries, err := getLinuxBatteries(context.Background(), root)
	if err != nil || len(batteries) != 2 || batteries[0].Name != "BAT0" || batteries[1].Name != "BAT1" {
		t.Fatalf("getLinuxBatteries() = %#v, %v", batteries, err)
	}
}

func TestGetLinuxBatteriesInvalidCapacity(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, "BAT0")
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "type"), []byte("Battery"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "capacity"), []byte("invalid"), 0o644); err != nil {
		t.Fatal(err)
	}
	batteries, err := getLinuxBatteries(context.Background(), root)
	if err == nil || len(batteries) != 0 {
		t.Fatalf("getLinuxBatteries() = %#v, %v", batteries, err)
	}
}
