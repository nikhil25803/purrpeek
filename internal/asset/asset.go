package asset

import (
	"embed"
	"errors"
	"fmt"
	"math/rand/v2"
	"path"
	"strings"
)

const DefaultImage = "mongo_no_bg.png"

//go:embed images/*.png
var images embed.FS

func Load(name string) ([]byte, error) {
	name = strings.TrimSpace(name)
	if name == "" || strings.ContainsAny(name, `/\`) || path.Base(name) != name || path.Ext(name) != ".png" {
		return nil, fmt.Errorf("invalid bundled image name %q", name)
	}
	data, err := images.ReadFile(path.Join("images", name))
	if err != nil {
		return nil, fmt.Errorf("load bundled image %q: %w", name, err)
	}
	return data, nil
}

func Select(names []string) (string, []byte, error) {
	return selectImage(names, rand.IntN)
}

func selectImage(names []string, choose func(int) int) (string, []byte, error) {
	type candidate struct {
		name string
		data []byte
	}

	seen := make(map[string]struct{}, len(names))
	valid := make([]candidate, 0, len(names))
	var problems []error
	for _, name := range names {
		name = strings.TrimSpace(name)
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		data, err := Load(name)
		if err != nil {
			problems = append(problems, err)
			continue
		}
		valid = append(valid, candidate{name: name, data: data})
	}

	if len(valid) > 0 {
		selected := valid[choose(len(valid))]
		return selected.name, selected.data, errors.Join(problems...)
	}
	problems = append(problems, errors.New("no valid configured images; using default"))
	data, err := Load(DefaultImage)
	if err != nil {
		return "", nil, errors.Join(append(problems, err)...)
	}
	return DefaultImage, data, errors.Join(problems...)
}
