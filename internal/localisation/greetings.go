package localisation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/nikhil25803/purrpeek/internal/asset"
)

const fileName = "greetings.json"

type Periods struct {
	Morning   []string `json:"morning"`
	Afternoon []string `json:"afternoon"`
	Evening   []string `json:"evening"`
	Night     []string `json:"night"`
}

type Catalog map[string]Periods

func Load() (Catalog, error) {
	return load(asset.Greetings(), os.UserConfigDir, os.ReadFile)
}

func load(defaults []byte, userConfigDir func() (string, error), readFile func(string) ([]byte, error)) (Catalog, error) {
	catalog, err := parse(defaults)
	if err != nil {
		return nil, fmt.Errorf("parse bundled greetings: %w", err)
	}
	directory, err := userConfigDir()
	if err != nil {
		return nil, fmt.Errorf("locate user greetings: %w", err)
	}
	path := filepath.Join(directory, "purrpeek", fileName)
	data, err := readFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return catalog, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read user greetings %q: %w", path, err)
	}
	override, err := parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse user greetings %q: %w", path, err)
	}
	merge(catalog, override)
	return catalog, nil
}

func parse(data []byte) (Catalog, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var catalog Catalog
	if err := decoder.Decode(&catalog); err != nil {
		return nil, err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			err = errors.New("multiple JSON values are not supported")
		}
		return nil, err
	}
	if catalog == nil {
		return nil, errors.New("greetings must be a JSON object")
	}

	normalized := make(Catalog, len(catalog))
	for language, periods := range catalog {
		language = strings.TrimSpace(language)
		if language == "" || strings.ContainsFunc(language, unicode.IsControl) {
			return nil, errors.New("language code cannot be empty or contain control characters")
		}
		if _, exists := normalized[language]; exists {
			return nil, fmt.Errorf("duplicate language code %q", language)
		}
		var err error
		if periods.Morning, err = normalize(periods.Morning); err != nil {
			return nil, fmt.Errorf("language %q morning: %w", language, err)
		}
		if periods.Afternoon, err = normalize(periods.Afternoon); err != nil {
			return nil, fmt.Errorf("language %q afternoon: %w", language, err)
		}
		if periods.Evening, err = normalize(periods.Evening); err != nil {
			return nil, fmt.Errorf("language %q evening: %w", language, err)
		}
		if periods.Night, err = normalize(periods.Night); err != nil {
			return nil, fmt.Errorf("language %q night: %w", language, err)
		}
		normalized[language] = periods
	}
	return normalized, nil
}

func normalize(phrases []string) ([]string, error) {
	if phrases == nil {
		return nil, nil
	}
	seen := make(map[string]struct{}, len(phrases))
	result := make([]string, 0, len(phrases))
	for _, phrase := range phrases {
		phrase = strings.TrimSpace(phrase)
		if phrase == "" {
			continue
		}
		if strings.ContainsFunc(phrase, unicode.IsControl) {
			return nil, errors.New("phrase contains control characters")
		}
		if _, exists := seen[phrase]; exists {
			continue
		}
		seen[phrase] = struct{}{}
		result = append(result, phrase)
	}
	return result, nil
}

func merge(catalog, override Catalog) {
	for language, patch := range override {
		periods := catalog[language]
		if patch.Morning != nil {
			periods.Morning = patch.Morning
		}
		if patch.Afternoon != nil {
			periods.Afternoon = patch.Afternoon
		}
		if patch.Evening != nil {
			periods.Evening = patch.Evening
		}
		if patch.Night != nil {
			periods.Night = patch.Night
		}
		catalog[language] = periods
	}
}

func Greeting(catalog Catalog, currentTime, username string) string {
	return greeting(catalog, currentTime, username, rand.IntN)
}

func greeting(catalog Catalog, currentTime, username string, choose func(int) int) string {
	username = strings.TrimSpace(username)
	if username == "" {
		return ""
	}
	parsed, err := time.Parse(time.RFC3339, currentTime)
	if err != nil {
		return "Hello, " + username
	}
	period := periodForHour(parsed.Hour())
	languages := make([]string, 0, len(catalog))
	for language, periods := range catalog {
		if len(period.phrases(periods)) > 0 {
			languages = append(languages, language)
		}
	}
	if len(languages) == 0 {
		return "Hello, " + username
	}
	slices.Sort(languages)
	phrases := period.phrases(catalog[languages[choose(len(languages))]])
	return phrases[choose(len(phrases))] + ", " + username
}

type dayPeriod byte

const (
	morning dayPeriod = iota
	afternoon
	evening
	night
)

func periodForHour(hour int) dayPeriod {
	switch {
	case hour >= 5 && hour < 12:
		return morning
	case hour >= 12 && hour < 17:
		return afternoon
	case hour >= 17 && hour < 21:
		return evening
	default:
		return night
	}
}

func (period dayPeriod) phrases(periods Periods) []string {
	switch period {
	case morning:
		return periods.Morning
	case afternoon:
		return periods.Afternoon
	case evening:
		return periods.Evening
	default:
		return periods.Night
	}
}
