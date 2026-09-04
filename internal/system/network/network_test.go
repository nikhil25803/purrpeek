package network

import "testing"

func TestFindPrimaryInterfaceFallback(t *testing.T) {
	interfaces := []NetworkInterface{{Name: "first"}, {Name: "second"}}
	if got := findPrimaryInterface(interfaces); got == nil || got.Name != "first" {
		t.Fatalf("findPrimaryInterface() = %#v", got)
	}
	if got := findPrimaryInterface(nil); got != nil {
		t.Fatalf("findPrimaryInterface(nil) = %#v", got)
	}
}
