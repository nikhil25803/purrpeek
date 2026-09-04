package internal

import (
	"errors"
	"strings"
	"testing"
)

func TestCollectionWarningIsSingleLine(t *testing.T) {
	warning := collectionWarning(errors.Join(errors.New("cpu: unavailable"), errors.New("disk: unavailable")))
	if strings.Count(warning, "\n") != 0 || !strings.Contains(warning, "cpu: unavailable; disk: unavailable") {
		t.Fatalf("collectionWarning() = %q", warning)
	}
}
