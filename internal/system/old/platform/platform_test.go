package platform

import (
	"testing"
)

func TestCurrent(t *testing.T) {

	current := Current()

	switch current {
	case Windows, Linux, MacOS:
		// OK
	default:
		t.Fatalf("unsupported platform: %v", current)
	}
}

func TestDefaultDisk(t *testing.T) {

	expected := map[Platform]string{
		Windows: "C:",
		Linux:   "/",
		MacOS:   "/",
		Unknown: ".",
	}

	got := DefaultDisk()
	want := expected[Current()]

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
