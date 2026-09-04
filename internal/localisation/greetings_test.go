package localisation

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/nikhil25803/purrpeek/internal/asset"
)

func TestEmbeddedGreetings(t *testing.T) {
	catalog, err := parse(asset.Greetings())
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog) != 8 || len(catalog["hi"].Morning) != 2 {
		t.Fatalf("embedded catalog = %+v", catalog)
	}
}

func TestLoadMergesUserGreetings(t *testing.T) {
	defaults := []byte(`{"en":{"morning":["Good morning"],"night":["Good night"]}}`)
	catalog, err := load(defaults,
		func() (string, error) { return "/config", nil },
		func(path string) ([]byte, error) {
			if want := filepath.Join("/config", "purrpeek", fileName); path != want {
				t.Fatalf("greetings path = %q, want %q", path, want)
			}
			return []byte(`{"en":{"morning":[" Hi ","Hi",""]},"de":{"morning":["Guten Morgen"]}}`), nil
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(catalog["en"].Morning, []string{"Hi"}) ||
		!reflect.DeepEqual(catalog["en"].Night, []string{"Good night"}) ||
		!reflect.DeepEqual(catalog["de"].Morning, []string{"Guten Morgen"}) {
		t.Fatalf("merged catalog = %+v", catalog)
	}
}

func TestLoadUsesEmbeddedGreetingsWhenUserFileIsMissing(t *testing.T) {
	catalog, err := load([]byte(`{"en":{"morning":["Morning"]}}`),
		func() (string, error) { return "/config", nil },
		func(string) ([]byte, error) { return nil, os.ErrNotExist },
	)
	if err != nil || catalog["en"].Morning[0] != "Morning" {
		t.Fatalf("catalog = %+v, error = %v", catalog, err)
	}
}

func TestLoadRejectsGreetingErrors(t *testing.T) {
	valid := []byte(`{"en":{"morning":["Morning"]}}`)
	tests := []struct {
		name     string
		defaults []byte
		dir      func() (string, error)
		read     func(string) ([]byte, error)
	}{
		{"invalid embedded", []byte(`{`), func() (string, error) { return "/config", nil }, func(string) ([]byte, error) { return nil, os.ErrNotExist }},
		{"unavailable directory", valid, func() (string, error) { return "", errors.New("unavailable") }, func(string) ([]byte, error) { return nil, nil }},
		{"unreadable override", valid, func() (string, error) { return "/config", nil }, func(string) ([]byte, error) { return nil, errors.New("denied") }},
		{"malformed override", valid, func() (string, error) { return "/config", nil }, func(string) ([]byte, error) { return []byte(`{`), nil }},
		{"multiple values", valid, func() (string, error) { return "/config", nil }, func(string) ([]byte, error) { return []byte(`{} {}`), nil }},
		{"unknown period", valid, func() (string, error) { return "/config", nil }, func(string) ([]byte, error) { return []byte(`{"en":{"dawn":["Hi"]}}`), nil }},
		{"invalid value", valid, func() (string, error) { return "/config", nil }, func(string) ([]byte, error) { return []byte(`{"en":{"morning":"Hi"}}`), nil }},
		{"control characters", valid, func() (string, error) { return "/config", nil }, func(string) ([]byte, error) { return []byte(`{"en":{"morning":["Hi\u001b"]}}`), nil }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if catalog, err := load(test.defaults, test.dir, test.read); err == nil || catalog != nil {
				t.Fatalf("catalog = %+v, error = %v", catalog, err)
			}
		})
	}
}

func TestGreetingPeriods(t *testing.T) {
	catalog := Catalog{"en": {
		Morning: []string{"Morning"}, Afternoon: []string{"Afternoon"},
		Evening: []string{"Evening"}, Night: []string{"Night"},
	}}
	tests := map[string]string{
		"2026-09-04T04:59:00+05:30": "Night, nikhil",
		"2026-09-04T05:00:00+05:30": "Morning, nikhil",
		"2026-09-04T11:59:00+05:30": "Morning, nikhil",
		"2026-09-04T12:00:00+05:30": "Afternoon, nikhil",
		"2026-09-04T17:00:00+05:30": "Evening, nikhil",
		"2026-09-04T21:00:00+05:30": "Night, nikhil",
	}
	for timestamp, want := range tests {
		if got := greeting(catalog, timestamp, "nikhil", func(int) int { return 0 }); got != want {
			t.Errorf("greeting(%q) = %q, want %q", timestamp, got, want)
		}
	}
}

func TestGreetingRandomSelectionAndFallback(t *testing.T) {
	catalog := Catalog{
		"en": {Morning: []string{"Good morning"}},
		"hi": {Morning: []string{"सुप्रभात", "शुभ प्रभात"}},
	}
	choices := []int{1, 1}
	got := greeting(catalog, "2026-09-04T08:00:00+05:30", "nikhil", func(int) int {
		choice := choices[0]
		choices = choices[1:]
		return choice
	})
	if got != "शुभ प्रभात, nikhil" {
		t.Fatalf("random greeting = %q", got)
	}
	if got := greeting(nil, "invalid", "nikhil", func(int) int { return 0 }); got != "Hello, nikhil" {
		t.Fatalf("fallback greeting = %q", got)
	}
	if got := greeting(nil, "2026-09-04T08:00:00+05:30", "nikhil", func(int) int { return 0 }); got != "Hello, nikhil" {
		t.Fatalf("empty-catalog greeting = %q", got)
	}
}
