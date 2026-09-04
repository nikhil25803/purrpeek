package asset

import (
	"bytes"
	"testing"
)

func TestBundledImages(t *testing.T) {
	for _, name := range []string{
		"mongo_no_bg.png",
		"mongo_purrpeek.png",
		"snow_no_bg.png",
		"snow_purrpeek.png",
	} {
		data, err := Load(name)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")) {
			t.Fatalf("%s does not contain an embedded PNG", name)
		}
	}
}

func TestLoadRejectsInvalidNames(t *testing.T) {
	for _, name := range []string{"", "../mongo_no_bg.png", `folder\mongo_no_bg.png`, "mongo_no_bg.jpg", "missing.png"} {
		if _, err := Load(name); err == nil {
			t.Fatalf("Load(%q) succeeded", name)
		}
	}
}

func TestSelectImage(t *testing.T) {
	name, data, err := selectImage([]string{"missing.png", "mongo_no_bg.png", "mongo_no_bg.png", "snow_no_bg.png"}, func(count int) int {
		if count != 2 {
			t.Fatalf("candidate count = %d, want 2", count)
		}
		return 1
	})
	if name != "snow_no_bg.png" || len(data) == 0 {
		t.Fatalf("selected %q with %d bytes", name, len(data))
	}
	if err == nil {
		t.Fatal("invalid configured image was not reported")
	}
}

func TestSelectImageFallsBackToDefault(t *testing.T) {
	name, data, err := selectImage(nil, func(int) int { return 0 })
	if name != DefaultImage || len(data) == 0 {
		t.Fatalf("selected %q with %d bytes", name, len(data))
	}
	if err == nil {
		t.Fatal("empty image list was not reported")
	}
}
