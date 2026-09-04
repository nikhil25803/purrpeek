package platform

import "testing"

func TestDisplayOSName(t *testing.T) {
	for _, test := range []struct{ goos, platform, want string }{
		{"darwin", "darwin", "macOS"},
		{"linux", "ubuntu", "ubuntu"},
		{"windows", "", "windows"},
	} {
		if got := displayOSName(test.goos, test.platform); got != test.want {
			t.Errorf("displayOSName(%q, %q) = %q, want %q", test.goos, test.platform, got, test.want)
		}
	}
}
