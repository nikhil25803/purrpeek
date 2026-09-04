package system

import (
	"context"
	"testing"
)

func TestCanceledCollectionReturnsPartialInformation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	info, err := GetSystemInformation(ctx)
	if info == nil || info.Time == nil || info.GPUs == nil || info.Batteries == nil {
		t.Fatalf("partial collection returned %#v", info)
	}
	if err == nil {
		t.Fatal("canceled collection returned no error")
	}
}
