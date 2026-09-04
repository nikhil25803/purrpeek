package platform

import (
	"testing"
	"time"
)

func TestTimeInformation(t *testing.T) {
	tests := []struct {
		name       string
		offset     int
		wantOffset string
	}{
		{"positive half-hour", 5*60*60 + 30*60, "+05:30"},
		{"negative", -7 * 60 * 60, "-07:00"},
		{"UTC", 0, "+00:00"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			now := time.Date(2026, 9, 4, 9, 30, 0, 0, time.FixedZone("TEST", test.offset))
			info := timeInformation(now)
			if info.CurrentTime != now.Format(time.RFC3339) || info.TimeZone != "TEST" || info.UTCOffset != test.wantOffset {
				t.Fatalf("timeInformation() = %#v", info)
			}
		})
	}
}
